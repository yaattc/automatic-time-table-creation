package cmd

import (
	"crypto/rand"
	"crypto/sha1" //nolint
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pkgz/rest"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/uni"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/teacher"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"

	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/user"

	"github.com/go-pkgz/auth/provider"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/service"

	"github.com/go-pkgz/auth/avatar"

	"github.com/go-pkgz/auth/token"

	"github.com/go-pkgz/auth"
	log "github.com/go-pkgz/lgr"
	"github.com/yaattc/automatic-time-table-creation/backend/app/rest/api"
)

// Server runs REST API web server
type Server struct {
	Location string `long:"location" env:"LOCATION" description:"location of all db files" required:"true"`

	ServiceURL string `long:"service_url" env:"SERVICE_URL" description:"http service url" required:"true"`
	Port       int    `long:"service_port" env:"SERVICE_PORT" description:"http server port" default:"8080"`

	Auth struct {
		TTL struct {
			JWT    time.Duration `long:"jwt" env:"JWT" default:"5m" description:"jwt TTL"`
			Cookie time.Duration `long:"cookie" env:"COOKIE" default:"200h" description:"auth cookie TTL"`
		} `group:"ttl" namespace:"ttl" env-namespace:"TTL"`
		Secret     string `long:"secret" env:"SECRET" description:"secret for authentication tokens" required:"true"`
		BCryptCost int    `long:"bcrypt_cost" env:"BCRYPT_COST" description:"bcrypt cost for hashing user password" default:"10"`
	} `group:"auth" namespace:"auth" env-namespace:"AUTH"`
	DBConnStr string `long:"db_conn_str" env:"DB_CONN_STR" required:"true" description:"connection string to db"`

	Admin AdminGroup `group:"admin" namespace:"admin" env-namespace:"ADMIN"`

	CommonOpts
}

// AdminGroup defines options group for admin params
type AdminGroup struct {
	Email    string `long:"email" env:"EMAIL" description:"default admin email" required:"true"`
	Password string `long:"password" env:"PASSWORD" description:"default admin password" required:"true"`
}

// Execute runs http web server
func (s *Server) Execute(_ []string) error {
	pgpool, pgconf, err := preparePostgres(s.DBConnStr)
	if err != nil {
		return err
	}

	// initializing repositories
	ur, err := user.NewPostgres(pgpool, pgconf)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize postgres user repository at %s", s.DBConnStr)
	}

	tr, err := teacher.NewPostgres(pgpool, pgconf)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize postgres teacher repository at %s", s.DBConnStr)
	}

	uor, err := uni.NewPostgres(pgpool, pgconf)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize postgres university organization repository at %s", s.DBConnStr)
	}

	ds := &service.DataStore{
		UserRepository:    ur,
		TeacherRepository: tr,
		UniOrgRepository:  uor,
		BCryptCost:        s.Auth.BCryptCost,
	}

	if _, err = ds.RegisterAdmin(s.Admin.Email, s.Admin.Password); err != nil {
		return errors.Wrapf(err, "failed to register admin %s:%s", s.Admin.Email, s.Admin.Password)
	}

	authenticator := s.makeAuthenticator(ds)
	srv := api.Rest{
		Version:       s.Version,
		Authenticator: authenticator,
		DataStore:     ds,
	}

	srv.Run(s.Port)
	return nil
}

func (s *Server) makeAuthenticator(ds *service.DataStore) *auth.Service {
	authenticator := auth.NewService(auth.Opts{
		SecretReader: token.SecretFunc(func(aud string) (string, error) {
			return s.Auth.Secret, nil
		}),
		ClaimsUpd: token.ClaimsUpdFunc(func(c token.Claims) token.Claims { // set attributes, on new token or refresh
			if c.User == nil {
				return c
			}

			var err error
			c.User.Email, err = ds.GetUserEmail(c.User.ID)
			if err != nil {
				log.Printf("[WARN] can't read email for %s, %v", c.User.ID, err)
			}

			privs, err := ds.GetUserPrivs(c.User.ID)
			if err != nil {
				log.Printf("[WARN] can't get privs for %s, %v ", c.User.ID, err)
			}

			c.User.SetSliceAttr("privileges", store.PrivsToStr(privs))

			return c
		}),
		SecureCookies:  strings.HasPrefix(s.ServiceURL, "https://"),
		TokenDuration:  s.Auth.TTL.JWT,
		CookieDuration: s.Auth.TTL.Cookie,
		DisableXSRF:    true,
		Issuer:         "attc",
		URL:            strings.TrimSuffix(s.ServiceURL, "/"),
		Validator: token.ValidatorFunc(func(token string, claims token.Claims) bool { // check on each auth call (in middleware)
			return claims.User != nil
		}),
		AvatarStore: avatar.NewNoOp(),
		Logger:      log.Default(),
	})
	authenticator.AddCustomHandler(&DirectProvider{
		DirectHandler: &provider.DirectHandler{
			L:            log.Default(),
			ProviderName: "local",
			Issuer:       "attc",
			TokenService: authenticator.TokenService(),
			CredChecker:  provider.CredCheckerFunc(ds.CheckUserCredentials),
			AvatarSaver:  authenticator.AvatarProxy(),
		},
	})
	//authenticator.AddDirectProvider("local", provider.CredCheckerFunc(ds.CheckUserCredentials))
	return authenticator
}

// DirectProvider overrides the LoginHandler to allow to pass the token string inside the response body
// fixme A HUGE KLUDGE
type DirectProvider struct {
	*provider.DirectHandler
}

// LoginHandler just does the same as provider.DirectHandler.LoginHandler but adds the token string
// inside the response body
func (p *DirectProvider) LoginHandler(w http.ResponseWriter, r *http.Request) {
	creds, err := p.getCredentials(w, r)
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusBadRequest, err, "failed to parse credentials")
		return
	}
	sessOnly := r.URL.Query().Get("sess") == "1"
	if p.CredChecker == nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError,
			errors.New("no credential checker"), "no credential checker")
		return
	}
	ok, err := p.CredChecker.Check(creds.User, creds.Password)
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError, err, "failed to check user credentials")
		return
	}
	if !ok {
		rest.SendErrorJSON(w, r, p.L, http.StatusForbidden, nil, "incorrect user or password")
		return
	}
	u := token.User{
		Name: creds.User,
		ID:   p.ProviderName + "_" + token.HashID(sha1.New(), creds.User), //nolint
	}
	u, err = setAvatar(p.AvatarSaver, u, &http.Client{Timeout: 5 * time.Second})
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError, err, "failed to save avatar to proxy")
		return
	}

	cid, err := randToken()
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError, err, "can't make token id")
		return
	}

	claims := token.Claims{
		User: &u,
		StandardClaims: jwt.StandardClaims{
			Id:       cid,
			Issuer:   p.Issuer,
			Audience: creds.Audience,
		},
		SessionOnly: sessOnly,
	}

	ts, ok := p.TokenService.(*token.Service)
	if !ok {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError,
			errors.New("can't get token service"), "can't get token service")
		return
	}

	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Now().Add(ts.TokenDuration).Unix()
	}

	if claims.Issuer == "" {
		claims.Issuer = ts.Issuer
	}

	if !ts.DisableIAT {
		claims.IssuedAt = time.Now().Unix()
	}

	tokenStr, err := ts.Token(claims)
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError,
			err, "can't get token string")
		return
	}

	if _, err = p.TokenService.Set(w, claims); err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError, err, "failed to set token")
		return
	}
	rest.RenderJSON(w, r, rest.JSON{"user": claims.User, "token": tokenStr})
}

// getCredentials extracts user and password from request
func (p *DirectProvider) getCredentials(w http.ResponseWriter, r *http.Request) (credentials, error) {
	if r.Body != nil {
		r.Body = http.MaxBytesReader(w, r.Body, provider.MaxHTTPBodySize)
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "" {
		mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			return credentials{}, err
		}
		contentType = mt
	}

	if contentType == "application/json" {
		var creds credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			return credentials{}, errors.Wrap(err, "failed to parse request body")
		}
		return creds, nil
	}

	if err := r.ParseForm(); err != nil {
		return credentials{}, errors.Wrap(err, "failed to parse request")
	}

	return credentials{
		User:     r.Form.Get("user"),
		Password: r.Form.Get("passwd"),
		Audience: r.Form.Get("aud"),
	}, nil
}

// credentials holds user credentials
type credentials struct {
	User     string `json:"user"`
	Password string `json:"passwd"`
	Audience string `json:"aud"`
}

func randToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "can't get random")
	}
	s := sha1.New() //nolint
	if _, err := s.Write(b); err != nil {
		return "", errors.Wrap(err, "can't write randoms to sha1")
	}
	return fmt.Sprintf("%x", s.Sum(nil)), nil
}

// setAvatar saves avatar and puts proxied URL to u.Picture
func setAvatar(ava provider.AvatarSaver, u token.User, client *http.Client) (token.User, error) {
	if ava != nil {
		avatarURL, e := ava.Put(u, client)
		if e != nil {
			return u, errors.Wrap(e, "failed to save avatar for")
		}
		u.Picture = avatarURL
		return u, nil
	}
	return u, nil // empty AvatarSaver ok, just skipped
}

package cmd

import (
	"strings"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"

	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/engine"

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
		Secret string `long:"secret" env:"SECRET" description:"secret for authentication tokens"`
	} `group:"auth" namespace:"auth" env-namespace:"AUTH"`
	AdminPwd  string `long:"admin_pwd" env:"ADMIN_PWD" default:"" description:"admin basic auth password"`
	DBConnStr string `long:"db_conn_str" env:"DB_CONN_STR" required:"true" description:"connection string to db"`

	Admin AdminGroup `group:"admin" namespace:"admin" env-namespace:"ADMIN"`

	CommonOpts
}

// AdminGroup defines options group for admin params
type AdminGroup struct {
	Type   string `long:"type" env:"TYPE" description:"type of admin store" choice:"shared" choice:"rpc" default:"shared"` //nolint
	Shared struct {
		Admins []string `long:"id" env:"ID" description:"admin(s) ids" env-delim:","`
		Email  string   `long:"email" env:"EMAIL" default:"" description:"admin email"`
	} `group:"shared" namespace:"shared" env-namespace:"SHARED"`
}

// Execute runs http web server
func (s *Server) Execute(_ []string) error {

	pg, err := engine.NewPostgres(s.DBConnStr)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize postgres engine at %s: %v", s.DBConnStr, err)
	}

	ds := &service.DataStore{Engine: pg}

	authenticator := s.makeAuthenticator(ds)
	srv := api.Rest{
		Version:       s.Version,
		Authenticator: authenticator,
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
		JWTQuery:       "jwt",
		Issuer:         "attc",
		URL:            strings.TrimSuffix(s.ServiceURL, "/"),
		Validator: token.ValidatorFunc(func(token string, claims token.Claims) bool { // check on each auth call (in middleware)
			return claims.User != nil
		}),
		AvatarStore: avatar.NewNoOp(),
		AdminPasswd: s.AdminPwd,
		Logger:      log.Default(),
	})
	authenticator.AddDirectProvider("local", provider.CredCheckerFunc(ds.CheckUserCredentials))
	return authenticator
}

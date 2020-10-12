package cmd

import (
	"time"

	"github.com/go-pkgz/auth"
	"github.com/yaattc/automatic-time-table-creation/backend/app/rest/api"
)

// ServeCmd runs REST API web server
type ServeCmd struct {
	Location string `long:"location" env:"LOCATION" description:"location of all db files" required:"true"`
	Port     int    `long:"http_port" env:"HTTP_PORT" description:"http server port" default:"8080"`

	Auth struct {
		TTL struct {
			JWT    time.Duration `long:"jwt" env:"JWT" default:"5m" description:"jwt TTL"`
			Cookie time.Duration `long:"cookie" env:"COOKIE" default:"200h" description:"auth cookie TTL"`
		} `group:"ttl" namespace:"ttl" env-namespace:"TTL"`
	} `group:"auth" namespace:"auth" env-namespace:"AUTH"`
	CommonOpts
}

// Execute runs http web server
func (s *ServeCmd) Execute(_ []string) error {

	authenticator := s.makeAuthenticator()
	srv := api.Rest{
		Version:       s.Version,
		Authenticator: authenticator,
	}

	srv.Run(s.Port)
	return nil
}

func (s *ServeCmd) makeAuthenticator() *auth.Service {
	return nil
}

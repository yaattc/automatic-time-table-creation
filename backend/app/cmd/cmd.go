// Package cmd contains all cli commands, their arguments and tests to them
package cmd

import (
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

// CommonOptionsCommander extends flags.Commander with SetCommon
// All commands should implement this interfaces
type CommonOptionsCommander interface {
	SetCommon(commonOpts CommonOpts)
	Execute(args []string) error
}

// CommonOpts sets externally from main, shared across all commands
type CommonOpts struct {
	Version string
}

// SMTPGroup defines options for SMTP server connection, used in auth and notify modules
type SMTPGroup struct {
	Host     string        `long:"host" env:"HOST" description:"SMTP host"`
	Port     int           `long:"port" env:"PORT" description:"SMTP port"`
	Username string        `long:"username" env:"USERNAME" description:"SMTP username"`
	Password string        `long:"password" env:"PASSWORD" description:"SMTP password"`
	TLS      bool          `long:"tls" env:"TLS" description:"enable TLS"`
	Timeout  time.Duration `long:"timeout" env:"TIMEOUT" default:"10s" description:"SMTP TCP connection timeout"`

	From string `long:"from" env:"FROM" required:"true" description:"from email address"`
}

// SetCommon satisfies CommonOptionsCommander interface and sets common option fields
// The method called by main for each command
func (c *CommonOpts) SetCommon(opts CommonOpts) {
	c.Version = opts.Version
}

func preparePostgres(connStr string) (*pgx.ConnPool, pgx.ConnConfig, error) {
	connConf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		return nil, pgx.ConnConfig{}, errors.Wrapf(err, "failed to parse pg user Store with connstr %s", connStr)
	}

	p, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 5,
		AcquireTimeout: time.Minute,
	})
	if err != nil {
		return nil, pgx.ConnConfig{}, errors.Wrapf(err, "failed to initialize pg user Store with connstr %s", connStr)
	}

	log.Printf("[INFO] initialized postgres connection pool to %s:%d", connConf.Host, connConf.Port)

	return p, connConf, nil
}

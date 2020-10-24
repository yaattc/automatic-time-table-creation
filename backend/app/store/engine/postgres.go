package engine

import (
	"time"

	log "github.com/go-pkgz/lgr"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

// Postgres implements Interface with postgres queries
type Postgres struct {
	connPool *pgx.ConnPool
	connConf pgx.ConnConfig
}

// NewPostgres returns the new instance of Postgres
func NewPostgres(connStr string) (*Postgres, error) {
	connConf, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse pg user Store with connstr %s", connStr)
	}

	p, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 5,
		AfterConnect: func(conn *pgx.Conn) error {
			// todo no-op yet
			return nil
		},
		AcquireTimeout: time.Minute,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize pg user Store with connstr %s", connStr)
	}

	log.Printf("[INFO] initialized postgres connection pool to %s:%d", connConf.Host, connConf.Port)

	return &Postgres{
		connPool: p,
		connConf: connConf,
	}, nil
}

// GetUser from the database by its ID
func (p *Postgres) GetUser(id string) (u store.User, err error) {
	row := p.connPool.QueryRow("SELECT id, email, privileges FROM users WHERE id = $1", id)
	err = row.Scan(&u.ID, &u.Email, &u.Privileges)
	return u, errors.Wrapf(err, "failed to read user with id %s", id)
}

// GetPasswordHash of user by its email
func (p *Postgres) GetPasswordHash(email string) (pwd string, err error) {
	row := p.connPool.QueryRow("SELECT password FROM users WHERE email = $1", email)
	err = row.Scan(&pwd)
	return pwd, errors.Wrapf(err, "failed to read %s user's password", email)
}

package user

import (
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
func NewPostgres(connPool *pgx.ConnPool, connConf pgx.ConnConfig) (*Postgres, error) {
	return &Postgres{
		connPool: connPool,
		connConf: connConf,
	}, nil
}

// GetUser from the database by its ID
func (p *Postgres) GetUser(id string) (u store.User, err error) {
	row := p.connPool.QueryRow(`SELECT id, email, privileges FROM users WHERE id = $1`, id)
	err = row.Scan(&u.ID, &u.Email, &u.Privileges)
	return u, errors.Wrapf(err, "failed to read user with id %s", id)
}

// GetPasswordHash of user by its email
func (p *Postgres) GetPasswordHash(email string) (pwd string, err error) {
	row := p.connPool.QueryRow(`SELECT password FROM users WHERE email = $1`, email)
	err = row.Scan(&pwd)
	return pwd, errors.Wrapf(err, "failed to read %s user's password", email)
}

// AddUser to the database
func (p *Postgres) AddUser(u store.User, pwd string, ignoreIfExists bool) (err error) {
	query := "INSERT INTO users(id, email, password, privileges) VALUES ($1, $2, $3, $4)"
	if ignoreIfExists {
		query += " ON CONFLICT DO NOTHING"
	}
	_, err = p.connPool.Exec(query, u.ID, u.Email, pwd, u.Privileges)

	return errors.Wrapf(err, "failed to add user %s into database", u.ID)
}

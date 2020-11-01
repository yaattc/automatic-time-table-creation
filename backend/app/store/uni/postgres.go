package uni

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/pgh"
)

// Postgres implements Interface with postgres queries
type Postgres struct {
	connPool *pgx.ConnPool
	connConf pgx.ConnConfig //nolint:structcheck,unused
}

// NewPostgres returns the new instance of Postgres
func NewPostgres(connPool *pgx.ConnPool, connConf pgx.ConnConfig) (*Postgres, error) {
	return &Postgres{
		connPool: connPool,
		connConf: connConf,
	}, nil
}

// AddGroup to the database
func (p *Postgres) AddGroup(g store.Group) (id string, err error) {
	_, err = p.connPool.Exec(`INSERT INTO groups("id", "name", "study_year_id") VALUES($1, $2, $3)`,
		g.ID, g.Name, g.StudyYear.ID)
	return g.ID, errors.Wrapf(err, "failed to insert group %s.%s", g.StudyYear.ID, g.Name)
}

// ListGroups registered in the database
func (p *Postgres) ListGroups() (res []store.Group, err error) {
	err = pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {

		return nil
	}))
	return res, err
}

// DeleteGroup from the database
func (p *Postgres) DeleteGroup(id string) error {
	panic("implement me")
}

// AddStudyYear to database
func (p *Postgres) AddStudyYear(sy store.StudyYear) (id string, err error) {
	panic("implement me")
}

// DeleteStudyYear from the database
func (p *Postgres) DeleteStudyYear(studyYearID string) error {
	panic("implement me")
}

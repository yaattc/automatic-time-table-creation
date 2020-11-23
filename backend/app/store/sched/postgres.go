package sched

import (
	"time"

	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/pgh"

	"github.com/jackc/pgx"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

// Postgres implements Interface
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

// ListClasses in the selected period of time for the selected group
func (p *Postgres) ListClasses(from time.Time, till time.Time, groupID string) (res []store.Class, err error) {
	const query = `SELECT id, course_id, group_id, teacher_id, title, location, start_time, duration 
                    FROM classes WHERE start_time >= $1 AND start_time <= $2 AND group_id = $3`

	err = pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {
		rows, err := tx.Query(query, from, till, groupID)
		if err != nil {
			return errors.Wrapf(err, "failed to query classes from %s to %s for group %s",
				from.String(), till.String(), groupID)
		}

		for rows.Next() {
			cl := store.Class{}
			err := rows.Scan(&cl.ID, &cl.Course.ID, &cl.Group.ID, &cl.Teacher.ID,
				&cl.Title, &cl.Location, &cl.Start, &cl.Duration)
			if err != nil {
				return errors.Wrapf(err, "failed to scan classes from %s to %s for group %s", from.String(), till.String(), groupID)
			}
			res = append(res, cl)
		}
		return nil
	}))
	return res, err
}

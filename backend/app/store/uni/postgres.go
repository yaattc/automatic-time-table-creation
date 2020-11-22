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
		rows, err := tx.Query(`SELECT groups.id, groups.name, groups.study_year_id, study_years.name 
						FROM groups 
						LEFT JOIN study_years ON groups.study_year_id = study_years.id`)
		if err != nil {
			return errors.Wrap(err, "failed to query groups")
		}
		defer rows.Close()

		var g store.Group
		for rows.Next() {
			g = store.Group{}
			if err = rows.Scan(&g.ID, &g.Name, &g.StudyYear.ID, &g.StudyYear.Name); err != nil {
				return errors.Wrap(err, "failed to scan group")
			}
			res = append(res, g)
		}
		return nil
	}))
	return res, err
}

// DeleteGroup from the database
func (p *Postgres) DeleteGroup(id string) error {
	_, err := p.connPool.Exec(`DELETE FROM groups WHERE id = $1`, id)
	return errors.Wrapf(err, "failed to remove group with id %s", id)
}

// AddStudyYear to database
func (p *Postgres) AddStudyYear(sy store.StudyYear) (id string, err error) {
	_, err = p.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES($1, $2)`, sy.ID, sy.Name)
	return sy.ID, errors.Wrapf(err, "failed to add study year %+v", sy)
}

// DeleteStudyYear from the database
func (p *Postgres) DeleteStudyYear(studyYearID string) error {
	_, err := p.connPool.Exec(`DELETE FROM study_years WHERE id = $1`, studyYearID)
	return errors.Wrapf(err, "failed to remove study year with id %s", studyYearID)
}

// GetGroup from the database
func (p *Postgres) GetGroup(id string) (g store.Group, err error) {
	row := p.connPool.QueryRow(`SELECT groups.id, groups.name, groups.study_year_id, study_years.name 
										FROM groups
										LEFT JOIN study_years ON groups.study_year_id = study_years.id
										WHERE groups.id = $1`, id)
	err = row.Scan(&g.ID, &g.Name, &g.StudyYear.ID, &g.StudyYear.Name)
	return g, errors.Wrapf(err, "failed to select group %s", id)
}

// GetStudyYear by its id
func (p *Postgres) GetStudyYear(id string) (sy store.StudyYear, err error) {
	row := p.connPool.QueryRow(`SELECT id, name FROM study_years WHERE id = $1`, id)
	err = row.Scan(&sy.ID, &sy.Name)
	return sy, errors.Wrapf(err, "failed to get study year %s", id)
}

// ListStudyYears from the database
func (p *Postgres) ListStudyYears() (res []store.StudyYear, err error) {
	err = pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {
		rows, err := tx.Query(`SELECT id, name FROM study_years`)
		if err != nil {
			return errors.Wrap(err, "failed to query list all study years")
		}
		var s store.StudyYear
		for rows.Next() {
			s = store.StudyYear{}
			if err = rows.Scan(&s.ID, &s.Name); err != nil {
				return errors.Wrap(err, "failed to scan study year")
			}
			res = append(res, s)
		}
		return nil
	}))
	return res, err
}

// AddCourse to the database
func (p *Postgres) AddCourse(course store.Course) (id string, err error) {
	err = pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {
		_, err := tx.Exec(`INSERT INTO courses(id, name, primary_lector_id, assistant_lector_id, edu_program) 
						VALUES ($1, $2, $3, $4, $5)`,
			course.ID, course.Name, course.PrimaryLector.ID,
			course.AssistantLector.ID, course.Program)
		if err != nil {
			return errors.Wrapf(err, "failed to insert course %s", course.Name)
		}
		for _, ta := range course.Assistants {
			_, err := tx.Exec(`INSERT INTO courses_teacher_assistants(course_id, assistant_id) VALUES ($1, $2)`,
				course.ID, ta.ID)
			if err != nil {
				return errors.Wrapf(err, "failed to insert assistant %s to course %s", ta.ID, course.Name)
			}
		}
		return nil
	}))
	return course.ID, err
}

// GetCourseDetails by id
func (p *Postgres) GetCourseDetails(id string) (res store.Course, err error) {
	err = pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {
		row := tx.QueryRow(`SELECT id, name, edu_program, primary_lector_id, assistant_lector_id 
							FROM courses WHERE id = $1`, id)
		if err := row.Scan(&res.ID, &res.Name, &res.Program, &res.PrimaryLector.ID, &res.AssistantLector.ID); err != nil {
			return errors.Wrapf(err, "failed to scan course details for course %s", id)
		}
		rows, err := tx.Query(`SELECT assistant_id FROM courses_teacher_assistants WHERE course_id = $1`, id)
		if err != nil {
			return errors.Wrapf(err, "failed to query TAs for course %s", id)
		}
		for rows.Next() {
			var taID string
			if err := rows.Scan(&taID); err != nil {
				return errors.Wrapf(err, "failed to scan id of TA for course %s", id)
			}
			res.Assistants = append(res.Assistants, store.Teacher{TeacherDetails: store.TeacherDetails{ID: taID}})
		}
		return nil
	}))
	return res, err
}

// ListTimeSlots that are registered in the database
func (p *Postgres) ListTimeSlots() (res []store.TimeSlot, err error) {
	err = pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {
		rows, err := tx.Query(`SELECT id, weekday, start, duration FROM time_slots`)
		if err != nil {
			return errors.Wrap(err, "failed to query list all time slots")
		}
		var ts store.TimeSlot
		for rows.Next() {
			ts = store.TimeSlot{}
			if err = rows.Scan(&ts.ID, &ts.Weekday, &ts.Start, &ts.Duration); err != nil {
				return errors.Wrap(err, "failed to scan time slot")
			}
			res = append(res, ts)
		}
		return nil
	}))
	return res, err
}

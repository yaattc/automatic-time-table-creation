package teacher

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/pgh"
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

// AddTeacher to database
func (p *Postgres) AddTeacher(teacher store.TeacherDetails) (string, error) {
	_, err := p.connPool.Exec(
		`INSERT INTO teachers("id", "name", "surname", "email", "degree", "about") 
					VALUES ($1, $2, $3, $4, $5, $6)`,
		teacher.ID,
		teacher.Name,
		teacher.Surname,
		teacher.Email,
		teacher.Degree,
		teacher.About)
	return teacher.ID, errors.Wrapf(err, "failed to insert teacher %s %s", teacher.Name, teacher.Surname)
}

// DeleteTeacher from the database
func (p *Postgres) DeleteTeacher(teacherID string) error {
	_, err := p.connPool.Exec(`DELETE FROM teachers WHERE id = $1`, teacherID)
	return errors.Wrapf(err, "failed to delete teacher %s", teacherID)
}

// ListTeachers from the database
func (p *Postgres) ListTeachers() ([]store.TeacherDetails, error) {
	var res []store.TeacherDetails
	err := pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {
		// taking teacher IDs
		rows, err := tx.Query(`SELECT id FROM teachers`)
		if err != nil {
			return errors.Wrap(err, "failed to make a query on listing teachers")
		}
		defer rows.Close()

		var ids []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return errors.Wrap(err, "failed to scan teachers")
			}
			ids = append(ids, id)
		}

		// taking teachers' preferences
		for _, id := range ids {
			td, err := p.getTeacherDetails(id, tx)
			if err != nil {
				return errors.Wrapf(err, "failed to load teacher details for %s", id)
			}
			res = append(res, td)
		}
		return nil
	}))
	return res, err
}

// GetTeacherFull gets the full information about the given teacher, including references
// to other teachers, such as preferences in staff
func (p *Postgres) GetTeacherFull(teacherID string) (store.Teacher, error) {
	var tp store.TeacherPreferences
	var td store.TeacherDetails
	err := pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) (err error) {
		tp, err = p.getTeacherPreferences(teacherID, tx)
		if err != nil {
			return errors.Wrapf(err, "failed to load preferences for teacher %s", teacherID)
		}
		td, err = p.getTeacherDetails(teacherID, tx)
		if err != nil {
			return errors.Wrapf(err, "failed to load teacher details for %s", teacherID)
		}
		return nil
	}))
	return store.Teacher{Preferences: tp, TeacherDetails: td}, err
}

// getTeacherPreferences composes teacher preferences from several tables into a single object
func (p *Postgres) getTeacherPreferences(teacherID string, tx *pgx.Tx) (tp store.TeacherPreferences, err error) {
	// loading teaching staff
	if tp.Staff, err = p.getStaffPreferences(teacherID, tx); err != nil {
		return store.TeacherPreferences{}, errors.Wrapf(err, "failed to load teaching staff for %s", teacherID)
	}

	// loading time slots
	if tp.TimeSlots, err = p.getTimeSlotPreferences(teacherID, tx); err != nil {
		return store.TeacherPreferences{}, errors.Wrapf(err, "failed to load time slots for %s", teacherID)
	}

	// loading locations
	if tp.Locations, err = p.getLocationPreferences(teacherID, tx); err != nil {
		return store.TeacherPreferences{}, errors.Wrapf(err, "failed to load location preferences for %s", teacherID)
	}

	return tp, nil
}

// getLocationPreferences loads general teacher's preferences in locations, despite the time slots
func (p *Postgres) getLocationPreferences(teacherID string, tx *pgx.Tx) (locs []store.Location, err error) {
	row := tx.QueryRow(`SELECT locations FROM teacher_preferences WHERE teacher_id = $1`, teacherID)
	err = row.Scan(&locs)
	return locs, errors.Wrapf(err, "failed to scan locations for %s", teacherID)
}

// getTimeSlotPreferences returns the teacher preferences in time slots
func (p *Postgres) getTimeSlotPreferences(teacherID string, tx *pgx.Tx) ([]store.TimeSlot, error) {
	rows, err := tx.Query(`SELECT weekday, start, duration, location 
								FROM teacher_preferences_time_slots WHERE teacher_id = $1`, teacherID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query time slots for %s", teacherID)
	}
	defer rows.Close()

	var tss []store.TimeSlot
	for rows.Next() {
		ts := store.TimeSlot{}
		if err := rows.Scan(&ts.Weekday, &ts.Start, &ts.Duration, &ts.Location); err != nil {
			return nil, errors.Wrapf(err, "failed to scan time slots for %s", teacherID)
		}
		tss = append(tss, ts)
	}

	return tss, nil
}

// getStaffPreferences returns the teacher preferences in staff, loading the details of the staff
func (p *Postgres) getStaffPreferences(teacherID string, tx *pgx.Tx) ([]store.TeacherDetails, error) {
	rows, err := tx.Query(`SELECT staff_id FROM teacher_preferences_staff WHERE teacher_id = $1`, teacherID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query staff preferences for %s", teacherID)
	}
	defer rows.Close()

	var tIDs []string
	for rows.Next() {
		var tID string
		if err = rows.Scan(&tID); err != nil {
			return nil, errors.Wrapf(err, "failed to scan staff for %s", teacherID)
		}
		tIDs = append(tIDs, tID)
	}

	var tds []store.TeacherDetails
	for _, tID := range tIDs {
		t, err := p.getTeacherDetails(tID, tx)
		if err != nil {
			return nil, err
		}
		tds = append(tds, t)
	}
	return tds, nil
}

// getTeacherDetails gets teacher data that relates to only one particular teacher
// without taking links to the others
func (p *Postgres) getTeacherDetails(teacherID string, tx *pgx.Tx) (store.TeacherDetails, error) {
	t := store.TeacherDetails{}
	row := tx.QueryRow(`SELECT id, name, surname, email, degree, about FROM teachers WHERE id = $1`, teacherID)
	if err := row.Scan(&t.ID, &t.Name, &t.Surname, &t.Email, &t.Degree, &t.About); err != nil {
		return store.TeacherDetails{}, errors.Wrapf(err, "failed to scan details for %s", teacherID)
	}
	return t, nil
}

// SetPreferences for the given teacher
func (p *Postgres) SetPreferences(teacherID string, pref store.TeacherPreferences) error {
	err := pgh.Tx(p.connPool, pgh.TxerFunc(func(tx *pgx.Tx) error {
		// setting staff preferences
		for _, t := range pref.Staff {
			_, err := p.connPool.Exec(`INSERT INTO teacher_preferences_staff("teacher_id", "staff_id") VALUES ($1, $2) 
										ON CONFLICT (teacher_id, staff_id) DO UPDATE SET teacher_id = $1, staff_id = $2`,
				teacherID, t.ID)
			if err != nil {
				return errors.Wrapf(err, "failed to insert a staff preference for %s with the staff %s", teacherID, t.ID)
			}
		}

		// setting time slot preferences
		for _, ts := range pref.TimeSlots {
			_, err := p.connPool.Exec(`INSERT INTO teacher_preferences_time_slots("teacher_id", "weekday", "start", "duration", "location") 
			VALUES ($1, $2, $3, $4, $5)`, teacherID, ts.Weekday, ts.Start, ts.Duration, ts.Location)
			if err != nil {
				return errors.Wrapf(err, "failed to insert time slot preference for %s with the time slot %+v", teacherID, ts)
			}
		}

		// setting general location preferences
		_, err := p.connPool.Exec(`INSERT INTO teacher_preferences("teacher_id", "locations") 
		VALUES ($1, $2) ON CONFLICT (teacher_id) DO UPDATE SET locations = $2`, teacherID, pref.Locations)
		if err != nil {
			return errors.Wrapf(err, "failed to insert teacher preference in locations for %s with locs %v", teacherID, pref.Locations)
		}
		return nil
	}))
	return err
}

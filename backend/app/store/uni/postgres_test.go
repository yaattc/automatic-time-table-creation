package uni

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"

	"github.com/stretchr/testify/assert"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
)

func TestPostgres_AddGroup(t *testing.T) {
	srv := preparePgStore(t)

	_, err := srv.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES($1, $2)`,
		"00000000-0000-0000-0000-100000000003", "BS - Year 1 (Computer Science)")
	require.NoError(t, err)

	id, err := srv.AddGroup(store.Group{
		ID:        "00000000-0000-0000-0000-000000000003",
		Name:      "B20-05",
		StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-100000000003"},
	})
	require.NoError(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000003", id)

	var g store.Group
	row := srv.connPool.QueryRow(`SELECT id, name FROM groups`)
	err = row.Scan(&g.ID, &g.Name)
	assert.Equal(t, "00000000-0000-0000-0000-000000000003", g.ID)
	assert.Equal(t, "B20-05", g.Name)
}

func TestPostgres_DeleteGroup(t *testing.T) {
	srv := preparePgStore(t)

	_, err := srv.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES ($1, $2)`,
		"00000000-0000-0000-0000-000000000001", "BS - Year 1 (Computer Science)")
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO groups("id", "study_year_id", "name") VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-100000000001", "00000000-0000-0000-0000-000000000001", "B20-05")
	require.NoError(t, err)

	err = srv.DeleteGroup("00000000-0000-0000-0000-100000000001")
	require.NoError(t, err)

	var cnt int
	row := srv.connPool.QueryRow(`SELECT COUNT(*) FROM groups`)
	err = row.Scan(&cnt)
	require.NoError(t, err)
	assert.Zero(t, cnt)

	var id, name string
	row = srv.connPool.QueryRow(`SELECT id, name FROM study_years`)
	err = row.Scan(&id, &name)
	require.NoError(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", id)
	assert.Equal(t, "BS - Year 1 (Computer Science)", name)
}

func TestPostgres_AddStudyYear(t *testing.T) {
	srv := preparePgStore(t)
	id, err := srv.AddStudyYear(store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"})
	require.NoError(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", id)

	var name string
	row := srv.connPool.QueryRow(`SELECT id, name FROM study_years`)
	err = row.Scan(&id, &name)
	require.NoError(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", id)
	assert.Equal(t, "BS - Year 1 (Computer Science)", name)
}

func TestPostgres_DeleteStudyYear(t *testing.T) {
	srv := preparePgStore(t)
	_, err := srv.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES ($1, $2)`,
		"00000000-0000-0000-0000-000000000001", "BS - Year 1 (Computer Science)")
	require.NoError(t, err)

	err = srv.DeleteStudyYear("00000000-0000-0000-0000-000000000001")
	require.NoError(t, err)

	var cnt int
	row := srv.connPool.QueryRow(`SELECT COUNT(*) FROM study_years`)
	err = row.Scan(&cnt)
	require.NoError(t, err)
	assert.Zero(t, cnt)
}

func TestPostgres_ListStudyYears(t *testing.T) {
	srv := preparePgStore(t)
	expected := []store.StudyYear{
		{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000002", Name: "MS - Year 1 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000003", Name: "BS - Year 2 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000004", Name: "MS - Year 2 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000005", Name: "BS - Year 3 (Computer Science)"},
	}

	addStudyYear := func(sy store.StudyYear) {
		_, err := srv.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES ($1, $2)`,
			sy.ID, sy.Name)
		require.NoError(t, err)
	}

	for _, sy := range expected {
		addStudyYear(sy)
	}

	sys, err := srv.ListStudyYears()
	require.NoError(t, err)

	assert.ElementsMatch(t, expected, sys)
}

func TestPostgres_ListGroups(t *testing.T) {
	srv := preparePgStore(t)
	_, err := srv.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES ($1, $2)`,
		"00000000-0000-0000-0000-000000000001", "BS - Year 1 (Computer Science)")
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO groups("id", "study_year_id", "name") VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-100000000001", "00000000-0000-0000-0000-000000000001", "B20-01")
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO groups("id", "study_year_id", "name") VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-100000000002", "00000000-0000-0000-0000-000000000001", "B20-02")
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO groups("id", "study_year_id", "name") VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-100000000003", "00000000-0000-0000-0000-000000000001", "B20-03")
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO groups("id", "study_year_id", "name") VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-100000000004", "00000000-0000-0000-0000-000000000001", "B20-04")
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO groups("id", "study_year_id", "name") VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-100000000005", "00000000-0000-0000-0000-000000000001", "B20-05")
	require.NoError(t, err)

	gs, err := srv.ListGroups()
	require.NoError(t, err)
	assert.ElementsMatch(t, []store.Group{
		{
			ID:        "00000000-0000-0000-0000-100000000001",
			Name:      "B20-01",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000002",
			Name:      "B20-02",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000003",
			Name:      "B20-03",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000004",
			Name:      "B20-04",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000005",
			Name:      "B20-05",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
	}, gs)
}

func TestPostgres_GetGroup(t *testing.T) {
	srv := preparePgStore(t)
	_, err := srv.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES ($1, $2)`,
		"00000000-0000-0000-0000-000000000001", "BS - Year 1 (Computer Science)")
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO groups("id", "study_year_id", "name") VALUES ($1, $2, $3)`,
		"00000000-0000-0000-0000-100000000001", "00000000-0000-0000-0000-000000000001", "B20-01")
	require.NoError(t, err)

	g, err := srv.GetGroup("00000000-0000-0000-0000-100000000001")
	require.NoError(t, err)
	assert.Equal(t, store.Group{
		ID:   "00000000-0000-0000-0000-100000000001",
		Name: "B20-01",
		StudyYear: store.StudyYear{
			ID:   "00000000-0000-0000-0000-000000000001",
			Name: "BS - Year 1 (Computer Science)",
		},
	}, g)
}

func TestPostgres_GetStudyYear(t *testing.T) {
	srv := preparePgStore(t)
	_, err := srv.connPool.Exec(`INSERT INTO study_years("id", "name") VALUES ($1, $2)`,
		"00000000-0000-0000-0000-000000000001", "BS - Year 1 (Computer Science)")
	require.NoError(t, err)

	sy, err := srv.GetStudyYear("00000000-0000-0000-0000-000000000001")
	require.NoError(t, err)
	assert.Equal(t, store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"}, sy)
}

func TestPostgres_AddCourse(t *testing.T) {
	srv := preparePgStore(t)
	expected := store.Course{
		ID:      "00000000-0000-0000-0000-000000000001",
		Name:    "Operational systems",
		Program: store.Bachelor,
		Assistants: []store.Teacher{
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000001",
				Name:    "somename",
				Surname: "somesurname",
				Email:   "someemail",
				Degree:  "somedegree",
				About:   "something",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000002",
				Name:    "somename2",
				Surname: "somesurname2",
				Email:   "someemail2",
				Degree:  "somedegree2",
				About:   "something2",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000003",
				Name:    "somename3",
				Surname: "somesurname3",
				Email:   "someemail3",
				Degree:  "somedegree3",
				About:   "something3",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000004",
				Name:    "somename4",
				Surname: "somesurname4",
				Email:   "someemail4",
				Degree:  "somedegree4",
				About:   "something4",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000005",
				Name:    "somename5",
				Surname: "somesurname5",
				Email:   "someemail5",
				Degree:  "somedegree5",
				About:   "something5",
			}},
		},
		PrimaryLector: store.Teacher{
			TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-200000000001",
				Name:    "some primary lector name",
				Surname: "some primary lector surname",
				Email:   "some primary lector email",
				Degree:  "some primary lector degree",
				About:   "some primary lector about",
			},
		},
		AssistantLector: store.Teacher{
			TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-200000000002",
				Name:    "some assistant lector name",
				Surname: "some assistant lector surname",
				Email:   "some assistant lector email",
				Degree:  "some assistant lector degree",
				About:   "some assistant lector about",
			},
		},
	}
	tas := []string{
		"00000000-0000-0000-0000-100000000001",
		"00000000-0000-0000-0000-100000000002",
		"00000000-0000-0000-0000-100000000003",
		"00000000-0000-0000-0000-100000000004",
		"00000000-0000-0000-0000-100000000005",
	}

	// filling out dependencies
	for _, ta := range expected.Assistants {
		_, err := srv.connPool.Exec(`INSERT INTO teachers(id, name, surname, email, degree, about)
								VALUES ($1, $2, $3, $4, $5, $6)`, ta.ID, ta.Name, ta.Surname,
			ta.Email, ta.Degree, ta.About)
		require.NoError(t, err)
	}

	_, err := srv.connPool.Exec(`INSERT INTO teachers(id, name, surname, email, degree, about)
								VALUES ($1, $2, $3, $4, $5, $6)`, expected.PrimaryLector.ID, expected.PrimaryLector.Name, expected.PrimaryLector.Surname,
		expected.PrimaryLector.Email, expected.PrimaryLector.Degree, expected.PrimaryLector.About)
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO teachers(id, name, surname, email, degree, about)
								VALUES ($1, $2, $3, $4, $5, $6)`, expected.AssistantLector.ID, expected.AssistantLector.Name, expected.AssistantLector.Surname,
		expected.AssistantLector.Email, expected.AssistantLector.Degree, expected.AssistantLector.About)
	require.NoError(t, err)

	id, err := srv.AddCourse(expected)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, id)

	// checking TAs
	rows, err := srv.connPool.Query(`SELECT course_id, assistant_id FROM courses_teacher_assistants`)
	require.NoError(t, err)

	for rows.Next() {
		var courseID string
		var assistantID string
		err := rows.Scan(&courseID, &assistantID)
		require.NoError(t, err)
		assert.Equal(t, expected.ID, courseID)

		assert.Contains(t, tas, assistantID)
	}

	// checking course itself
	var actual store.Course
	row := srv.connPool.QueryRow(`SELECT name, edu_program, primary_lector_id, assistant_lector_id FROM courses WHERE id = $1`, expected.ID)
	err = row.Scan(&actual.Name, &actual.Program, &actual.PrimaryLector.ID, &actual.AssistantLector.ID)
	require.NoError(t, err)

	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Program, actual.Program)
	assert.Equal(t, expected.PrimaryLector.ID, actual.PrimaryLector.ID)
	assert.Equal(t, expected.AssistantLector.ID, actual.AssistantLector.ID)
}

func TestPostgres_GetCourseDetails(t *testing.T) {
	srv := preparePgStore(t)
	expected := store.Course{
		ID:      "00000000-0000-0000-0000-000000000001",
		Name:    "Operational systems",
		Program: store.Bachelor,
		Assistants: []store.Teacher{
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000001",
				Name:    "somename",
				Surname: "somesurname",
				Email:   "someemail",
				Degree:  "somedegree",
				About:   "something",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000002",
				Name:    "somename2",
				Surname: "somesurname2",
				Email:   "someemail2",
				Degree:  "somedegree2",
				About:   "something2",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000003",
				Name:    "somename3",
				Surname: "somesurname3",
				Email:   "someemail3",
				Degree:  "somedegree3",
				About:   "something3",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000004",
				Name:    "somename4",
				Surname: "somesurname4",
				Email:   "someemail4",
				Degree:  "somedegree4",
				About:   "something4",
			}},
			{TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-100000000005",
				Name:    "somename5",
				Surname: "somesurname5",
				Email:   "someemail5",
				Degree:  "somedegree5",
				About:   "something5",
			}},
		},
		PrimaryLector: store.Teacher{
			TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-200000000001",
				Name:    "some primary lector name",
				Surname: "some primary lector surname",
				Email:   "some primary lector email",
				Degree:  "some primary lector degree",
				About:   "some primary lector about",
			},
		},
		AssistantLector: store.Teacher{
			TeacherDetails: store.TeacherDetails{
				ID:      "00000000-0000-0000-0000-200000000002",
				Name:    "some assistant lector name",
				Surname: "some assistant lector surname",
				Email:   "some assistant lector email",
				Degree:  "some assistant lector degree",
				About:   "some assistant lector about",
			},
		},
	}

	// filling out dependencies
	for _, ta := range expected.Assistants {
		_, err := srv.connPool.Exec(`INSERT INTO teachers(id, name, surname, email, degree, about)
								VALUES ($1, $2, $3, $4, $5, $6)`, ta.ID, ta.Name, ta.Surname,
			ta.Email, ta.Degree, ta.About)
		require.NoError(t, err)
	}

	_, err := srv.connPool.Exec(`INSERT INTO teachers(id, name, surname, email, degree, about)
								VALUES ($1, $2, $3, $4, $5, $6)`, expected.PrimaryLector.ID, expected.PrimaryLector.Name, expected.PrimaryLector.Surname,
		expected.PrimaryLector.Email, expected.PrimaryLector.Degree, expected.PrimaryLector.About)
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO teachers(id, name, surname, email, degree, about)
								VALUES ($1, $2, $3, $4, $5, $6)`, expected.AssistantLector.ID, expected.AssistantLector.Name, expected.AssistantLector.Surname,
		expected.AssistantLector.Email, expected.AssistantLector.Degree, expected.AssistantLector.About)
	require.NoError(t, err)

	_, err = srv.connPool.Exec(`INSERT INTO courses(id, name, primary_lector_id, assistant_lector_id, edu_program) 
						VALUES ($1, $2, $3, $4, $5)`,
		expected.ID, expected.Name, expected.PrimaryLector.ID,
		expected.AssistantLector.ID, expected.Program)
	require.NoError(t, err)

	for _, ta := range expected.Assistants {
		_, err := srv.connPool.Exec(`INSERT INTO courses_teacher_assistants(course_id, assistant_id) VALUES ($1, $2)`,
			expected.ID, ta.ID)
		require.NoError(t, err)
	}

	course, err := srv.GetCourseDetails(expected.ID)
	require.NoError(t, err)
	assert.Equal(t, store.Course{
		ID:      expected.ID,
		Name:    expected.Name,
		Program: expected.Program,
		PrimaryLector: store.Teacher{
			TeacherDetails: store.TeacherDetails{ID: expected.PrimaryLector.ID},
		},
		AssistantLector: store.Teacher{
			TeacherDetails: store.TeacherDetails{ID: expected.AssistantLector.ID},
		},
		Assistants: []store.Teacher{
			{TeacherDetails: store.TeacherDetails{ID: "00000000-0000-0000-0000-100000000001"}},
			{TeacherDetails: store.TeacherDetails{ID: "00000000-0000-0000-0000-100000000002"}},
			{TeacherDetails: store.TeacherDetails{ID: "00000000-0000-0000-0000-100000000003"}},
			{TeacherDetails: store.TeacherDetails{ID: "00000000-0000-0000-0000-100000000004"}},
			{TeacherDetails: store.TeacherDetails{ID: "00000000-0000-0000-0000-100000000005"}},
		},
	}, course)

}

func preparePgStore(t *testing.T) *Postgres {
	// initializing connection with postgres
	connStr := os.Getenv("DB_TEST")
	connConf, err := pgx.ParseConnectionString(connStr)
	require.NoError(t, err)

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 5,
		AcquireTimeout: time.Minute,
	})
	require.NoError(t, err)

	log.Printf("[INFO] initialized postgres connection pool to %s:%d", connConf.Host, connConf.Port)

	p, err := NewPostgres(pool, connConf)
	require.NoError(t, err)

	// setting cleanups
	cleanupStorage(t, p.connPool)
	t.Cleanup(func() {
		cleanupStorage(t, p.connPool)
	})

	return p
}

func cleanupStorage(t *testing.T, p *pgx.ConnPool) {
	tx, err := p.Begin()
	require.NoError(t, err)
	defer func() {
		err := tx.Commit()
		require.NoError(t, err)
	}()

	_, err = tx.Exec(`TRUNCATE study_years CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE groups CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE teachers CASCADE`)
	require.NoError(t, err)
}

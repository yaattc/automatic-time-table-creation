package teacher

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Semior001/timetype"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

func TestPostgres_AddTeacher(t *testing.T) {
	srv := preparePgStore(t)

	expected := store.TeacherDetails{
		ID:      "00000000-0000-0000-0000-000000000001",
		Name:    "Yelshat",
		Surname: "Duskaliyev",
		Email:   "e.duskaliev@innopolis.university",
		Degree:  "nope",
		About:   "Not a teacher but a man",
	}
	id, err := srv.AddTeacher(expected)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, id)

	actual := store.TeacherDetails{}
	row := srv.connPool.QueryRow(`SELECT id, name, surname, email, degree, about FROM teachers`)
	err = row.Scan(&actual.ID, &actual.Name, &actual.Surname, &actual.Email, &actual.Degree, &actual.About)
	require.NoError(t, err)
}

func TestPostgres_GetTeacherFull(t *testing.T) {
	srv := preparePgStore(t)
	setupTestTeachers(t, srv)
	tch, err := srv.GetTeacherFull(tchs[0].ID)
	require.NoError(t, err)
	assert.Equal(t, tchs[0], tch)
}

func TestPostgres_ListTeachers(t *testing.T) {
	srv := preparePgStore(t)
	setupTestTeachers(t, srv)
	l, err := srv.ListTeachers()
	require.NoError(t, err)
	assert.ElementsMatch(t, []store.TeacherDetails{
		tchs[0].TeacherDetails,
		tchs[1].TeacherDetails,
		tchs[2].TeacherDetails,
	}, l)
}

func TestPostgres_DeleteTeacherAndSetPrefs(t *testing.T) {
	srv := preparePgStore(t)
	setupTestTeachers(t, srv)

	err := srv.DeleteTeacher(tchs[0].ID)
	require.NoError(t, err)
	var cnt int
	row := srv.connPool.QueryRow(`SELECT COUNT(*) FROM teachers`)
	err = row.Scan(&cnt)
	require.NoError(t, err)
	assert.Equal(t, 2, cnt)

	row = srv.connPool.QueryRow(`SELECT COUNT(*) FROM teacher_preferences_staff`)
	err = row.Scan(&cnt)
	require.NoError(t, err)
	assert.Zero(t, cnt)

	row = srv.connPool.QueryRow(`SELECT COUNT(*) FROM teacher_preferences_time_slots`)
	err = row.Scan(&cnt)
	require.NoError(t, err)
	assert.Zero(t, cnt)

}

var tchs = []store.Teacher{
	{
		TeacherDetails: store.TeacherDetails{
			ID:      "00000000-0000-0000-0000-000000000001",
			Name:    "foo",
			Surname: "bar",
			Email:   "foo@bar.com",
			Degree:  "graduate",
			About:   "some details about teacher",
		},
		Preferences: store.TeacherPreferences{
			TimeSlots: []store.TimeSlot{
				{
					Weekday:  time.Monday,
					Start:    timetype.NewClock(20, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					Location: "room 108",
				},
				{
					Weekday:  time.Tuesday,
					Start:    timetype.NewClock(10, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					Location: "room 109",
				},
				{
					Weekday:  time.Friday,
					Start:    timetype.NewClock(15, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					Location: "room 102",
				},
			},
			Locations: []store.Location{"108", "102", "109"},
			Staff: []store.TeacherDetails{
				{
					ID:      "00000000-0000-0000-0000-000000000002",
					Name:    "Ivan",
					Surname: "Konyukhov",
					Email:   "i.konyukhov@innopolis.ru",
					Degree:  "Dr.",
					About:   "Good man",
				},
				{
					ID:      "00000000-0000-0000-0000-000000000003",
					Name:    "Nikolay",
					Surname: "Shilov",
					Email:   "n.shilov@innopolis.ru",
					Degree:  "Prof.",
					About:   "???",
				},
			},
		},
	},
	{
		TeacherDetails: store.TeacherDetails{
			ID:      "00000000-0000-0000-0000-000000000002",
			Name:    "Ivan",
			Surname: "Konyukhov",
			Email:   "i.konyukhov@innopolis.ru",
			Degree:  "Dr.",
			About:   "Good man",
		},
	},
	{
		TeacherDetails: store.TeacherDetails{
			ID:      "00000000-0000-0000-0000-000000000003",
			Name:    "Nikolay",
			Surname: "Shilov",
			Email:   "n.shilov@innopolis.ru",
			Degree:  "Prof.",
			About:   "???",
		},
	},
}

func setupTestTeachers(t *testing.T, srv *Postgres) {
	id, err := srv.AddTeacher(tchs[0].TeacherDetails)
	require.NoError(t, err)
	assert.Equal(t, tchs[0].ID, id)

	_, err = srv.AddTeacher(tchs[1].TeacherDetails)
	require.NoError(t, err)

	_, err = srv.AddTeacher(tchs[2].TeacherDetails)
	require.NoError(t, err)

	err = srv.SetPreferences(tchs[0].ID, tchs[0].Preferences)
	require.NoError(t, err)
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

	_, err = tx.Exec(`TRUNCATE teachers CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE teacher_preferences CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE teacher_preferences_staff CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE teacher_preferences_time_slots CASCADE`)
	require.NoError(t, err)
}

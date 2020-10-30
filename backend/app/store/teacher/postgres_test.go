package teacher

import (
	"log"
	"os"
	"testing"
	"time"

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
	err := srv.AddTeacher(expected)
	require.NoError(t, err)

	actual := store.TeacherDetails{}
	row := srv.connPool.QueryRow(`SELECT id, name, surname, email, degree, about FROM teachers`)
	err = row.Scan(&actual.ID, &actual.Name, &actual.Surname, &actual.Email, &actual.Degree, &actual.About)
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

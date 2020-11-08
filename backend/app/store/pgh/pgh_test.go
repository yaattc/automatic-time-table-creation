package pgh

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
)

func TestTx(t *testing.T) {
	// initializing connection with postgres
	connStr := os.Getenv("DB_TEST")
	connConf, err := pgx.ParseConnectionString(connStr)
	require.NoError(t, err)

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 1,
		AcquireTimeout: 5 * time.Millisecond,
	})
	require.NoError(t, err)

	log.Printf("[INFO] initialized postgres connection pool to %s:%d", connConf.Host, connConf.Port)

	// checking the error while calling pgx.ConnPool.Begin
	cn, err := pool.Acquire()
	require.NoError(t, err)

	err = Tx(pool, TxerFunc(func(_ *pgx.Tx) error { return nil }))
	merr, ok := err.(*multierror.Error)
	assert.True(t, ok)
	assert.Equal(t, 1, merr.Len())
	assert.Equal(t, pgx.ErrAcquireTimeout, errors.Cause(merr.Errors[0]))

	pool.Release(cn)

	// checking the handling error when calling f.Do
	werr := errors.New("some weird error")

	err = Tx(pool, TxerFunc(func(_ *pgx.Tx) error { return werr }))
	merr, ok = err.(*multierror.Error)
	assert.True(t, ok)
	assert.Equal(t, 1, merr.Len())
	assert.Equal(t, werr, errors.Cause(merr.Errors[0]))

	// checking the error while committing
	err = Tx(pool, TxerFunc(func(tx *pgx.Tx) error {
		_, _ = tx.Exec(`let's try some weird SQL query'`)
		return nil
	}))
	merr, ok = err.(*multierror.Error)
	assert.True(t, ok)
	assert.Equal(t, 1, merr.Len())
	assert.Equal(t, pgx.ErrTxCommitRollback, errors.Cause(merr.Errors[0]))

	// checking the case when everything is OK
	err = Tx(pool, TxerFunc(func(tx *pgx.Tx) error { return nil }))
	assert.Nil(t, err)

	pool.Close()
}

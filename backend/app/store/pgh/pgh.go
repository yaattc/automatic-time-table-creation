// Package pgh provides some wrappers and helpers for github.com/jackc/pgx methods
package pgh

import (
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

// Txer defines an object that can perform some action in the context
// of database transaction
type Txer interface {
	Do(tx *pgx.Tx) error
}

// TxerFunc is a helper to use ordinary functions as Txer
type TxerFunc func(tx *pgx.Tx) error

// Do implements Txer.Do
func (f TxerFunc) Do(tx *pgx.Tx) error { return f(tx) }

// Tx is a wrapper for transactions to simplify their usage. In case if there was something
// wrong during the transaction, rollback will be issued. This function returns
// a multierror.Error, so it is possible to work with each error separately
func Tx(pool *pgx.ConnPool, fun Txer) error {
	merr := &multierror.Error{}

	tx, err := pool.Begin()
	if err != nil {
		merr = multierror.Append(merr, errors.Wrap(err, "failed to start tx"))
		return merr
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != pgx.ErrTxClosed {
			merr = multierror.Append(merr, err)
			log.Printf("[DEBUG] failed to rollback transaction: %v", err)
		}
	}()

	if err = fun.Do(tx); err != nil {
		merr = multierror.Append(merr, err)
		return merr
	}

	if err := tx.Commit(); err != nil {
		merr = multierror.Append(merr, err)
		log.Printf("[WARN] failed to commit transaction: %v", err)
	}

	if merr.Len() < 1 {
		return nil
	}
	return merr
}

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

// Tx is a wrapper for transactions to simplify their usage
func Tx(pool *pgx.ConnPool, f Txer) error {
	ferr := &multierror.Error{}

	tx, err := pool.Begin()
	if err != nil {
		ferr = multierror.Append(ferr, errors.Wrap(err, "failed to start tx"))
		return ferr
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			ferr = multierror.Append(ferr, err)
			log.Printf("[DEBUG] failed to rollback transaction: %v", err)
		}
	}()

	if err = f.Do(tx); err != nil {
		ferr = multierror.Append(ferr, err)
		return ferr
	}

	if err := tx.Commit(); err != nil {
		ferr = multierror.Append(ferr, err)
		log.Printf("[WARN] failed to commit transaction: %v", err)
	}

	if ferr.Len() < 1 {
		return nil
	}
	return ferr
}

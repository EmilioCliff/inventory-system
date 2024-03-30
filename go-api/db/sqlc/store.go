package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	// "google.golang.org/appengine/log"
)

// type SQLStore struct {
// 	*Queries
// 	connPool *pgxpool.Pool
// }

// type Store interface {
// 	Querier
// }

type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rberr := tx.Rollback(ctx); rberr != nil {
			return fmt.Errorf("Fn erro: %v rb Error %c", err, rberr)
		}
		return err
	}
	return tx.Commit(ctx)
}

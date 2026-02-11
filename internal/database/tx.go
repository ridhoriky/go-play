package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TxFunc func(*sqlx.Tx) error

func WithTransactionX(
	ctx context.Context,
	db *sqlx.DB,
	fn func(*sqlx.Tx) error,
) error {

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

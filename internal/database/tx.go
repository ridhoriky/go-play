package database

import (
	"context"
	"database/sql"
)

type TxFunc func(*sql.Tx) error

func WithTransaction(
	ctx context.Context,
	db *sql.DB,
	fn TxFunc,
) error {

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
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

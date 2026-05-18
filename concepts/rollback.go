package main

import (
	"context"
	"database/sql"
	"time"
)

func BankTransaction(db *sql.DB, fromID, toID int, amount int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly: false,
	})

	if err != nil {
		return err
	}

	defer func(){
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, "UPDATE transactions SET amount = amount - $1 where fromID = $2", amount, fromID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE transactions SET amount = amount + $1 where toID = $2", amount, toID)
	if err != nil {
		return err
	}

	err =  tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all function to exceute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err:%v", err, rbErr)
		}
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the results of transfer transaction
type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountID int64    `json:"from_account_id"`
	ToAccountID   int64    `json:"to_account_id"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

// TranserTx performs a transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a db transaction
func (store *Store) TranserTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTX(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromAccountID = arg.FromAccountID
		result.ToAccountID = arg.ToAccountID

		// TODO: Create two entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// TODO: Update accounts' balance

		return nil
	})

	return result, err
}

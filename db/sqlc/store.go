package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// Run the transactional database
// func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
// 	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	q := New(tx)

// 	err = fn(q)
// 	if err != nil {
// 		return err
// 	}

// 	return tx.Commit(ctx)
// }

// Params to start the transfer
// type TransferTxParams struct {
// 	FromAccountID int64 `json:"from_account_id"`
// 	ToAccountID int64 `json:"to_account_id"`
// 	Amount int64 `json:"amount"`
// }

// Params to show the result from the transfer
// type TransferTxResult struct {
//     Transfer    Transfer `json:"transfer"`
//     FromAccount Account  `json:"from_account"`
//     ToAccount   Account  `json:"to_account"`
//     FromEntry   Entry    `json:"from_entry"`
//     ToEntry     Entry    `json:"to_entry"`
// }

// Start the transaction
// func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
// 	var result TransferTxResult

// 	err := store.execTx(ctx, func(q *Queries) error{
// 		var err error

// 		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
// 			FromAccountID: pgtype.Int8{Int64: arg.FromAccountID, Valid: true},
// 			ToAccountID:   pgtype.Int8{Int64: arg.ToAccountID, Valid: true},
// 			Amount:        arg.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
// 			AccountID: pgtype.Int8{Int64: arg.FromAccountID, Valid: true},
// 			Amount:    -arg.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
// 			AccountID: pgtype.Int8{Int64: arg.ToAccountID, Valid: true},
// 			Amount:    arg.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return result, err
// }
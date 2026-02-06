package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpdateBalance(
	ctx context.Context,
	walletID uuid.UUID,
	amount int64,
	op string,
) error {

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var balance int64
	err = tx.QueryRow(
		ctx,
		`SELECT balance FROM wallets WHERE id=$1 FOR UPDATE`,
		walletID,
	).Scan(&balance)

	if err != nil {
		return err
	}

	switch op {
	case "DEPOSIT":
		balance += amount
	case "WITHDRAW":
		if balance < amount {
			return errors.New("insufficient funds")
		}
		balance -= amount
	default:
		return errors.New("unknown operation")
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE wallets SET balance=$1 WHERE id=$2`,
		balance, walletID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetBalance(
	ctx context.Context,
	walletID uuid.UUID,
) (int64, error) {

	var balance int64
	err := r.db.QueryRow(
		ctx,
		`SELECT balance FROM wallets WHERE id=$1`,
		walletID,
	).Scan(&balance)

	return balance, err
}


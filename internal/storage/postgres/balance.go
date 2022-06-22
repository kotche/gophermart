package postgres

import (
	"context"
	"database/sql"
)

type BalancePostgres struct {
	db *sql.DB
}

func NewBalancePostgres(db *sql.DB) *BalancePostgres {
	return &BalancePostgres{
		db: db,
	}
}

func (b *BalancePostgres) GetCurrentBalance(ctx context.Context) (int, error) {
	return 0, nil
}

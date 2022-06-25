package postgres

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
)

type OrderPostgres struct {
	db *sql.DB
}

func NewOrderPostgres(db *sql.DB) *OrderPostgres {
	return &OrderPostgres{
		db: db,
	}
}

func (o *OrderPostgres) LoadOrder(ctx context.Context, order *model.Order) (int, error) {
	return -1, nil
}

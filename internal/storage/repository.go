package storage

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/storage/postgres"
)

type AuthRepoContract interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUser(ctx context.Context, login, password string) (model.User, error)
}

type OrderRepoContract interface {
	LoadOrder(ctx context.Context, order model.Order) (int, error)
}

type BalanceRepoContract interface {
	GetCurrentBalance(ctx context.Context) (int, error)
}

type Repository struct {
	AuthRepoContract
	OrderRepoContract
	BalanceRepoContract
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepoContract:    postgres.NewAuthPostgres(db),
		OrderRepoContract:   postgres.NewOrderPostgres(db),
		BalanceRepoContract: postgres.NewBalancePostgres(db),
	}
}

package storage

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/storage/postgres"
)

type AuthRepoContract interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUserID(ctx context.Context, user *model.User) (string, error)
}

type OrderRepoContract interface {
	SaveOrder(ctx context.Context, order *model.Order) error
	GetUserIDByNumberOrder(ctx context.Context, number string) string
}

type Repository struct {
	AuthRepoContract
	OrderRepoContract
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepoContract:  postgres.NewAuthPostgres(db),
		OrderRepoContract: postgres.NewOrderPostgres(db),
	}
}

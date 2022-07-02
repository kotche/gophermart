package storage

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/storage/postgres"
)

type AuthRepoContract interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserID(ctx context.Context, user *model.User) (int, error)
}

type AccrualOrderRepoContract interface {
	SaveOrder(ctx context.Context, order *model.AccrualOrder) error
	GetUserIDByNumberOrder(ctx context.Context, number string) int
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type WithdrawOrderRepoContract interface {
	GetAccruals(ctx context.Context, UserID int) float32
	GetWithdrawals(ctx context.Context, UserID int) float32
	DeductPoints(ctx context.Context, order *model.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error)
}

type Repository struct {
	AuthRepoContract
	AccrualOrderRepoContract
	WithdrawOrderRepoContract
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepoContract:          postgres.NewAuthPostgres(db),
		AccrualOrderRepoContract:  postgres.NewAccrualOrderPostgres(db),
		WithdrawOrderRepoContract: postgres.NewWithdrawOrderPostgres(db),
	}
}

package service

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/storage"
)

type AuthServiceContract interface {
	CreateUser(ctx context.Context, user *model.User) error
	AuthenticationUser(ctx context.Context, user *model.User) error
	GenerateToken(user *model.User, tokenAuth *jwtauth.JWTAuth) (string, error)
}

type AccrualOrderServiceContract interface {
	LoadOrder(ctx context.Context, numOrder uint64, userID int) error
	CheckLuhn(number uint64) bool
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type WithdrawOrderServiceContract interface {
	DeductionOfPoints(ctx context.Context, order *model.WithdrawOrder) error
	GetBalance(ctx context.Context, userID int) (float32, float32)
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error)
}

type Service struct {
	AuthServiceContract
	AccrualOrderServiceContract
	WithdrawOrderServiceContract
}

func NewService(repo *storage.Repository) *Service {
	return &Service{
		AuthServiceContract:          NewAuthService(repo),
		AccrualOrderServiceContract:  NewAccrualOrderService(repo),
		WithdrawOrderServiceContract: NewWithdrawOrderService(repo),
	}
}

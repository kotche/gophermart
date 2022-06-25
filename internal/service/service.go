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

type OrderServiceContract interface {
	LoadOrder(ctx context.Context, order *model.Order) (int, error)
}

type BalanceServiceContract interface {
	GetCurrentBalance(ctx context.Context) (int, error)
}

type Service struct {
	AuthServiceContract
	OrderServiceContract
	BalanceServiceContract
}

func NewService(repo *storage.Repository) *Service {
	return &Service{
		AuthServiceContract:    NewAuthService(repo),
		OrderServiceContract:   NewOrderService(repo),
		BalanceServiceContract: NewBalanceService(repo),
	}
}

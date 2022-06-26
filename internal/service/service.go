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
	LoadOrder(ctx context.Context, userID, numOrder string) error
	CheckLuhn(number int) bool
}

type Service struct {
	AuthServiceContract
	OrderServiceContract
}

func NewService(repo *storage.Repository) *Service {
	return &Service{
		AuthServiceContract:  NewAuthService(repo),
		OrderServiceContract: NewOrderService(repo),
	}
}

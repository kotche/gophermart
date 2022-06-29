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
	LoadOrder(ctx context.Context, numOrder string, userID int) error
	CheckLuhn(number int) bool
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type Service struct {
	AuthServiceContract
	AccrualOrderServiceContract
}

func NewService(repo *storage.Repository) *Service {
	return &Service{
		AuthServiceContract:         NewAuthService(repo),
		AccrualOrderServiceContract: NewAccrualOrderService(repo),
	}
}

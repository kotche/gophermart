package service

import (
	"context"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/storage"
)

//Из сервиса убрать интерфейсы. Оставить только в репозитории.
//Логику по доступу в черный ящик тоже через интерфейс в папке accrual

type AuthServiceContract interface {
	CreateUser(ctx context.Context, user model.User) error
}

type OrderServiceContract interface {
	LoadOrder(ctx context.Context, order model.Order) (int, error)
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

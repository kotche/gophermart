package service

import (
	"context"

	"github.com/kotche/gophermart/internal/model"
)

type OrderRepoContract interface {
	LoadOrder(ctx context.Context, order *model.Order) (int, error)
}

type OrderService struct {
	repo OrderRepoContract
}

func NewOrderService(repo OrderRepoContract) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (o *OrderService) LoadOrder(ctx context.Context, order *model.Order) (int, error) {
	return o.repo.LoadOrder(ctx, order)
}

package service

import (
	"context"
)

type BalanceRepoContract interface {
	GetCurrentBalance(ctx context.Context) (int, error)
}

type BalanceService struct {
	repo BalanceRepoContract
}

func NewBalanceService(repo BalanceRepoContract) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

func (b *BalanceService) GetCurrentBalance(ctx context.Context) (int, error) {
	return b.repo.GetCurrentBalance(ctx)
}

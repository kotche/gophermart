package service

import (
	"context"
	"log"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

type WithdrawOrderRepoContract interface {
	GetAccruals(ctx context.Context, UserID int) float32
	GetWithdrawals(ctx context.Context, UserID int) float32
	DeductPoints(ctx context.Context, order *model.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error)
}
type WithdrawOrderService struct {
	repo WithdrawOrderRepoContract
}

func NewWithdrawOrderService(repo WithdrawOrderRepoContract) *WithdrawOrderService {
	return &WithdrawOrderService{
		repo: repo,
	}
}

func (w WithdrawOrderService) GetBalance(ctx context.Context, userID int) (float32, float32) {
	accruals := w.repo.GetAccruals(ctx, userID)
	withdrawn := w.repo.GetWithdrawals(ctx, userID)
	return accruals, withdrawn
}

func (w WithdrawOrderService) DeductionOfPoints(ctx context.Context, order *model.WithdrawOrder) error {
	accruals, withdrawn := w.GetBalance(ctx, order.UserID)

	if order.Sum > accruals-withdrawn {
		return errormodel.NotEnoughPoints{}
	}

	err := w.repo.DeductPoints(ctx, order)
	if err != nil {
		log.Printf("DeductPoints db error: %s", err.Error())
		return err
	}

	return nil
}

func (w *WithdrawOrderService) GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error) {
	orders, err := w.repo.GetWithdrawalOfPoints(ctx, userID)
	if err != nil {
		log.Printf("balance GetWithdrawalOfPoints db error: %s", err.Error())
		return nil, err
	}
	return orders, nil
}

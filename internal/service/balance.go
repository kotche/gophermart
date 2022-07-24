package service

import (
	"context"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
	"github.com/rs/zerolog"
)

type WithdrawOrderRepoContract interface {
	GetAccruals(ctx context.Context, UserID int) float32
	GetWithdrawals(ctx context.Context, UserID int) float32
	DeductPoints(ctx context.Context, order *model.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error)
}
type WithdrawOrderService struct {
	repo WithdrawOrderRepoContract
	log  *zerolog.Logger
}

func NewWithdrawOrderService(repo WithdrawOrderRepoContract, log *zerolog.Logger) *WithdrawOrderService {
	return &WithdrawOrderService{
		repo: repo,
		log:  log,
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
		w.log.Error().Err(err).Msg("WithdrawOrderService.DeductionOfPoints: DeductPoints db error")
		return err
	}

	return nil
}

func (w *WithdrawOrderService) GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error) {
	orders, err := w.repo.GetWithdrawalOfPoints(ctx, userID)
	if err != nil {
		w.log.Error().Err(err).Msg("WithdrawOrderService.GetWithdrawalOfPoints: GetWithdrawalOfPoints db error")
		return nil, err
	}
	return orders, nil
}

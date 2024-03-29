package service

import (
	"context"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
	"github.com/rs/zerolog"
)

type AccrualOrderRepoContract interface {
	SaveOrder(ctx context.Context, order *model.AccrualOrder) error
	GetUserIDByNumberOrder(ctx context.Context, number uint64) int
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type AccrualOrderService struct {
	repo AccrualOrderRepoContract
	log  *zerolog.Logger
}

func NewAccrualOrderService(repo AccrualOrderRepoContract, log *zerolog.Logger) *AccrualOrderService {
	return &AccrualOrderService{
		repo: repo,
		log:  log,
	}
}

func (a *AccrualOrderService) LoadOrder(ctx context.Context, numOrder uint64, userID int) error {

	if !a.CheckLuhn(numOrder) {
		return errormodel.CheckLuhnError{}
	}

	order := model.AccrualOrder{
		Number: numOrder,
		UserID: userID,
		Status: model.StatusNEW,
	}

	userIDinDB := a.repo.GetUserIDByNumberOrder(ctx, order.Number)
	if userIDinDB != 0 {
		if userIDinDB == order.UserID {
			return errormodel.OrderAlreadyUploadedCurrentUserError{}
		} else {
			return errormodel.OrderAlreadyUploadedAnotherUserError{}
		}
	}

	err := a.repo.SaveOrder(ctx, &order)
	if err != nil {
		a.log.Error().Err(err).Msg("AccrualOrderService.LoadOrder: SaveOrder db error")
		return err
	}

	return nil
}

func (a *AccrualOrderService) CheckLuhn(number uint64) bool {
	var sum uint64

	for i := 0; number > 0; i++ {
		cur := number % 10
		if i%2 == 0 {
			sum += cur
			number = number / 10
			continue
		}
		cur = cur * 2
		if cur > 9 {
			cur = cur - 9
		}
		sum += cur
		number = number / 10
	}

	return sum%10 == 0
}

func (a *AccrualOrderService) GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error) {
	orders, err := a.repo.GetUploadedOrders(ctx, userID)
	if err != nil {
		a.log.Error().Err(err).Msg("AccrualOrderService.GetUploadedOrders: GetUploadedOrders db error")
		return nil, err
	}
	return orders, nil
}

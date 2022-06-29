package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
	"github.com/kotche/gophermart/internal/model/status"
)

type AccrualOrderRepoContract interface {
	SaveOrder(ctx context.Context, order *model.AccrualOrder) error
	GetUserIDByNumberOrder(ctx context.Context, number string) int
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type AccrualOrderService struct {
	repo AccrualOrderRepoContract
}

func NewAccrualOrderService(repo AccrualOrderRepoContract) *AccrualOrderService {
	return &AccrualOrderService{
		repo: repo,
	}
}

func (a *AccrualOrderService) LoadOrder(ctx context.Context, numOrder string, userID int) error {

	numOrderInt, err := strconv.Atoi(numOrder)
	if err != nil {
		log.Printf("LoadOrder service - conv from int error: %s", err.Error())
		return err
	}
	if !a.CheckLuhn(numOrderInt) {
		return errormodel.CheckLuhnError{}
	}

	order := model.AccrualOrder{
		Number:     numOrder,
		UserID:     userID,
		Status:     status.NEW,
		UploadedAt: time.Now(),
	}

	err = a.repo.SaveOrder(ctx, &order)
	if err != nil {
		log.Printf("SaveOrder db error: %s", err.Error())
		return err
	}

	return nil
}

func (a *AccrualOrderService) CheckLuhn(number int) bool {
	var sum int

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
		log.Printf("getUploadedOrders db error: %s", err.Error())
		return nil, err
	}
	return orders, nil
}

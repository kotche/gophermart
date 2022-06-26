package service

import (
	"context"
	"strconv"
	"time"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
	"github.com/kotche/gophermart/internal/model/status"
)

type OrderRepoContract interface {
	SaveOrder(ctx context.Context, order *model.Order) error
	GetUserIDByNumberOrder(ctx context.Context, number string) string
}

type OrderService struct {
	repo OrderRepoContract
}

func NewOrderService(repo OrderRepoContract) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (o *OrderService) LoadOrder(ctx context.Context, numOrder, userID string) error {

	numOrderInt, err := strconv.Atoi(numOrder)
	if err != nil {
		return err
	}
	if !o.CheckLuhn(numOrderInt) {
		return errormodel.CheckLuhnError{}
	}

	order := model.Order{
		Number:     numOrder,
		UserID:     userID,
		Status:     status.NEW,
		UploadedAt: time.Now(),
	}

	err = o.repo.SaveOrder(ctx, &order)
	if err != nil {
		return err
	}

	//загрузить в канал для чтения воркерами (проверка начисления баллов)

	return nil
}

func (o *OrderService) CheckLuhn(number int) bool {
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

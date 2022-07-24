package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kotche/gophermart/internal/broker/model"
	"github.com/kotche/gophermart/internal/logger"
	mock_storage "github.com/kotche/gophermart/internal/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestBrokerGetOrdersForProcessing(t *testing.T) {
	type want struct {
		orders []model.Order
	}

	tests := []struct {
		name   string
		err    error
		ticker *time.Ticker
		want   want
	}{
		{
			name:   "correct_get_orders_for_processing",
			ticker: time.NewTicker(timeoutGetOrdersDB*time.Second + time.Second*5),
			want: want{orders: []model.Order{
				{Number: 2377225624, Status: model.StatusNEW},
				{Number: 2377225625, Status: model.StatusPROCESSING},
			},
			},
		},
		{
			name:   "incorrect_get_orders_for_processing",
			ticker: time.NewTicker(timeoutGetOrdersDB*time.Second - timeoutGetOrdersDB*time.Second/2),
			err:    errors.New("error bd"),
			want:   want{orders: []model.Order{}},
		},
	}

	t.Parallel()

	log := logger.Init()

	accrualURL := ""

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			repo := mock_storage.NewMockBrokerRepoContract(control)

			repo.EXPECT().GetOrdersForProcessing(ctx, limitQuery).Return(tt.want.orders, tt.err).Times(1)

			b := NewBroker(repo, accrualURL, log)
			go b.GetOrdersForProcessing(ctx)

			b.chSignalGetOrdersForProcessing <- struct{}{}

			resultOrders := []model.Order{}

			for {
				select {
				case order := <-b.chOrdersForProcessing:
					resultOrders = append(resultOrders, order)
					if len(resultOrders) == len(tt.want.orders) {
						goto next
					}
				case <-tt.ticker.C:
					goto next
				}
			}
		next:
			assert.EqualValues(t, tt.want.orders, resultOrders)
		})
	}
}

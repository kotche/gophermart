package storage

import (
	"context"
	"database/sql"

	modelBroker "github.com/kotche/gophermart/internal/broker/model"
	"github.com/kotche/gophermart/internal/storage/postgres"
)

type BrokerRepoContract interface {
	GetOrdersForProcessing(ctx context.Context, limit int) ([]modelBroker.Order, error)
	UpdateOrderAccruals(ctx context.Context, orderAccruals []modelBroker.OrderAccrual) error
}

type BrokerRepository struct {
	BrokerRepoContract
}

func NewBrokerRepository(db *sql.DB) *BrokerRepository {
	return &BrokerRepository{
		BrokerRepoContract: postgres.NewBrokerPostgres(db),
	}
}

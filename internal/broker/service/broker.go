package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kotche/gophermart/internal/broker/model"
	"github.com/kotche/gophermart/internal/broker/model/status"
)

const (
	timeoutClient       = 10
	maxWorkers          = 3
	bufSizeOrdersRecord = 3
	ordersAPI           = "/api/orders/"
	limitQuery          = 3
	timeoutLoadOrdersDB = 3
	timeoutGetOrdersDB  = 5
)

type RepositoryContract interface {
	GetOrdersForProcessing(ctx context.Context, limit int) ([]model.Order, error)
	UpdateOrderAccruals(ctx context.Context, orderAccruals []model.OrderAccrual) error
}

type Broker struct {
	repo                           RepositoryContract
	client                         *http.Client
	accrualURL                     string
	bufOrderForRecord              []model.OrderAccrual
	chOrdersForProcessing          chan model.Order
	chOrdersAccrual                chan model.OrderAccrual
	chSignalGetOrdersForProcessing chan struct{}
	chLimitWorkers                 chan int
}

func NewBroker(repo RepositoryContract, accrualURL string) *Broker {
	return &Broker{
		repo:                           repo,
		client:                         &http.Client{Timeout: time.Second * timeoutClient},
		accrualURL:                     accrualURL,
		bufOrderForRecord:              make([]model.OrderAccrual, 0, bufSizeOrdersRecord),
		chOrdersForProcessing:          make(chan model.Order),
		chOrdersAccrual:                make(chan model.OrderAccrual),
		chSignalGetOrdersForProcessing: make(chan struct{}),
		chLimitWorkers:                 make(chan int, maxWorkers),
	}
}

func (b *Broker) Start(ctx context.Context) {
	go b.GetOrdersForProcessing(ctx)
	go b.GetOrdersAccrual(ctx)
	go b.LoadOrdersAccrual(ctx)
}

//GetOrdersForProcessing Получаем номера заказов из БД со статусом "NEW" и "PROCESSING" -> кидаем в канал
func (b *Broker) GetOrdersForProcessing(ctx context.Context) {

	ticker := time.NewTicker(timeoutGetOrdersDB * time.Second)

	for {
		select {
		case <-b.chSignalGetOrdersForProcessing:
			b.runGetOrdersForProcessing(ctx)
			ticker.Reset(timeoutGetOrdersDB * time.Second)
		case <-ticker.C:
			b.runGetOrdersForProcessing(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (b *Broker) runGetOrdersForProcessing(ctx context.Context) {
	orders, err := b.repo.GetOrdersForProcessing(ctx, limitQuery)
	if err != nil {
		log.Fatalf("broker runGetOrdersForProcessing db error receiving orders for processing: %s", err.Error())
	}

	for _, numOrder := range orders {
		b.chOrdersForProcessing <- numOrder
	}
}

//GetOrdersAccrual Получаем из канала номера заказов -> обращаемся к сервису accrual за информацией по статусу и начисленных баллов. Ограничиваем обращения до maxWorkers
func (b *Broker) GetOrdersAccrual(ctx context.Context) {
	for {
		select {
		case order := <-b.chOrdersForProcessing:
			b.chLimitWorkers <- 1
			go b.getOrdersAccrualWorker(order)
		case <-ctx.Done():
			return
		}
	}
}

func (b *Broker) getOrdersAccrualWorker(order model.Order) {
	var orderAccrual model.OrderAccrual
	url := fmt.Sprintf("%s%s%s", b.accrualURL, ordersAPI, order.Number)
	err := b.getJSONOrderFromAccrual(url, &orderAccrual)
	if err != nil {
		<-b.chLimitWorkers
		return
	}

	if orderAccrual.Status == status.REGISTERED {
		orderAccrual.Status = status.NEW
	}

	if order.Status != "" && order.Status != orderAccrual.Status {
		b.chOrdersAccrual <- orderAccrual
		<-b.chLimitWorkers
	}
}

func (b *Broker) getJSONOrderFromAccrual(url string, orderAccrual *model.OrderAccrual) error {
	resp, err := b.client.Get(url)
	if err != nil {
		log.Printf("broker getJSONOrderFromAccrual service accrual request error: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&orderAccrual)
	if err != nil {
		log.Printf("broker getJSONOrderFromAccrual decode accrual request error: %s", err.Error())
		return err
	}
	return nil
}

//LoadOrdersAccrual Записываем данные по ордерам в БД (по наполению буфера или по таймауту)
func (b *Broker) LoadOrdersAccrual(ctx context.Context) {
	ticker := time.NewTicker(timeoutGetOrdersDB * time.Second)

	for {
		select {
		case order := <-b.chOrdersAccrual:
			b.bufOrderForRecord = append(b.bufOrderForRecord, order)
			if len(b.bufOrderForRecord) >= bufSizeOrdersRecord {
				b.flush(ctx)
			}
			ticker.Reset(timeoutLoadOrdersDB * time.Second)
		case <-ticker.C:
			if len(b.bufOrderForRecord) > 0 {
				b.flush(ctx)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *Broker) flush(ctx context.Context) {
	ordersUpdate := make([]model.OrderAccrual, len(b.bufOrderForRecord))
	copy(ordersUpdate, b.bufOrderForRecord)
	b.bufOrderForRecord = make([]model.OrderAccrual, 0)
	go func() {
		err := b.repo.UpdateOrderAccruals(ctx, ordersUpdate)
		if err != nil {
			log.Printf("broker flush db error update orders: %s", err.Error())
			return
		}
		b.chSignalGetOrdersForProcessing <- struct{}{}
	}()
}

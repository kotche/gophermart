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
	timeoutLoadOrdersDB = 5
	timeoutGetOrdersDB  = 10
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

func (b *Broker) Start() {
	go b.GetOrdersForProcessing()
	go b.GetOrdersAccrual()
	go b.LoadOrdersAccrual()
}

//GetOrdersForProcessing Получаем номера заказов из БД со статусом "NEW" и "PROCESSING" -> кидаем в канал
func (b *Broker) GetOrdersForProcessing() {

	heartbeat := time.Tick(timeoutGetOrdersDB * time.Second)

	ctx := context.Background()
	for {
		select {
		case <-b.chSignalGetOrdersForProcessing:
			heartbeat = time.Tick(timeoutGetOrdersDB * time.Second)
			b.runGetOrdersForProcessing(ctx)
		case <-heartbeat:
			b.runGetOrdersForProcessing(ctx)
		}
	}
}

func (b *Broker) runGetOrdersForProcessing(ctx context.Context) {
	orders, err := b.repo.GetOrdersForProcessing(ctx, limitQuery)
	if err != nil {
		log.Fatalf("error receiving orders for processing: %s", err.Error())
	}

	for _, numOrder := range orders {
		b.chOrdersForProcessing <- numOrder
	}
}

//GetOrdersAccrual Получаем из канала номера заказов -> обращаемся к сервису accrual за информацией по статусу и начисленных баллов. Ограничиваем обращения до maxWorkers
func (b *Broker) GetOrdersAccrual() {
	for order := range b.chOrdersForProcessing {
		b.chLimitWorkers <- 1
		go b.getOrdersAccrualWorker(order)
	}
}

func (b *Broker) getOrdersAccrualWorker(order model.Order) {
	var orderAccrual model.OrderAccrual
	url := fmt.Sprintf("%s%s%s", b.accrualURL, ordersAPI, order.Number)
	err := b.getJSONOrderFromAccrual(url, &orderAccrual)
	if err != nil {
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
		log.Printf("service accrual request error: %s", err.Error())
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&orderAccrual)
	if err != nil {
		log.Printf("decode accrual request error: %s", err.Error())
		return err
	}
	return nil
}

//LoadOrdersAccrual Записываем данные по ордерам в БД (по наполению буфера или по таймауту)
func (b *Broker) LoadOrdersAccrual() {

	heartbeat := time.Tick(timeoutLoadOrdersDB * time.Second)

	ctx := context.Background()
	for {
		select {
		case order := <-b.chOrdersAccrual:
			heartbeat = time.Tick(timeoutLoadOrdersDB * time.Second)
			b.bufOrderForRecord = append(b.bufOrderForRecord, order)
			if len(b.bufOrderForRecord) >= bufSizeOrdersRecord {
				b.flush(ctx)
			}
		case <-heartbeat:
			if len(b.bufOrderForRecord) > 0 {
				b.flush(ctx)
			}
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
			log.Printf("error update orders: %s", err.Error())
			return
		}
		b.chSignalGetOrdersForProcessing <- struct{}{}
	}()
}

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kotche/gophermart/internal/logger"
	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/service"
	"github.com/kotche/gophermart/internal/storage"
	mock_storage "github.com/kotche/gophermart/internal/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetCurrentBalance(t *testing.T) {
	type want struct {
		status    int
		current   float32
		withdrawn float32
	}

	tests := []struct {
		name      string
		token     string
		userID    int
		accruals  float32
		withdrawn float32
		want      want
	}{
		{
			name:      "correct_get_balance",
			token:     "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID:    1,
			accruals:  100.5,
			withdrawn: 50.5,
			want:      want{status: http.StatusOK, current: 50, withdrawn: 50.5},
		},
	}

	log := logger.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)

			withdraw.EXPECT().GetAccruals(ctx, tt.userID).Return(tt.accruals).Times(1)
			withdraw.EXPECT().GetWithdrawals(ctx, tt.userID).Return(tt.withdrawn).Times(1)

			repo := &storage.Repository{
				Withdraw: withdraw,
			}

			serv := service.NewService(repo, log)
			h := NewHandler(serv, log)

			r := httptest.NewRequest(http.MethodGet, getCurrentBalance, nil)
			w := httptest.NewRecorder()
			r.Header.Add("Authorization", tt.token)

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()

			defer response.Body.Close()
			body, _ := io.ReadAll(response.Body)

			balance := &struct {
				Current   float32 `json:"current"`
				Withdrawn float32 `json:"withdrawn"`
			}{}
			_ = json.Unmarshal(body, balance)

			assert.Equal(t, tt.want.status, response.StatusCode)
			assert.Equal(t, tt.want.current, balance.Current)
			assert.Equal(t, tt.want.withdrawn, balance.Withdrawn)
		})
	}
}

func TestHandlerGetCurrentBalanceUserUnauthorized(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)

	repo := &storage.Repository{
		Withdraw: withdraw,
	}

	serv := service.NewService(repo, log)
	h := NewHandler(serv, log)

	r := httptest.NewRequest(http.MethodGet, getCurrentBalance, nil)
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestHandlerDeductionOfPoints(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name      string
		token     string
		userID    int
		accruals  float32
		withdrawn float32
		order     *model.WithdrawOrder
		body      []byte
		err       error
		want      want
	}{
		{
			name:      "successfully_deducted_points",
			token:     "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID:    1,
			order:     &model.WithdrawOrder{Order: 2377225624, UserID: 1, Sum: 100},
			body:      []byte(`{"order":"2377225624","sum": 100}`),
			accruals:  200.5,
			withdrawn: 50.5,
			want:      want{status: http.StatusOK},
		},
	}

	log := logger.Init()

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)

			withdraw.EXPECT().GetAccruals(ctx, tt.userID).Return(tt.accruals).Times(1)
			withdraw.EXPECT().GetWithdrawals(ctx, tt.userID).Return(tt.withdrawn).Times(1)
			withdraw.EXPECT().DeductPoints(ctx, tt.order).Return(tt.err).Times(1)

			repo := &storage.Repository{
				Withdraw: withdraw,
			}

			serv := service.NewService(repo, log)
			h := NewHandler(serv, log)

			r := httptest.NewRequest(http.MethodPost, deductionOfPoints, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()
			r.Header.Add("Authorization", tt.token)

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()
			assert.Equal(t, tt.want.status, response.StatusCode)
		})
	}
}

func TestHandlerDeductionOfPointsUserUnauthorized(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)

	repo := &storage.Repository{
		Withdraw: withdraw,
	}

	serv := service.NewService(repo, log)
	h := NewHandler(serv, log)

	r := httptest.NewRequest(http.MethodPost, deductionOfPoints, nil)
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestHandlerDeductionOfPointsInsufficientFunds(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	ctx := context.Background()

	withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)
	withdraw.EXPECT().GetAccruals(ctx, 1).Return(float32(100.5)).Times(1)
	withdraw.EXPECT().GetWithdrawals(ctx, 1).Return(float32(50.5)).Times(1)

	repo := &storage.Repository{
		Withdraw: withdraw,
	}

	serv := service.NewService(repo, log)
	h := NewHandler(serv, log)

	body := []byte(`{"order":"2377225624","sum": 100}`)
	r := httptest.NewRequest(http.MethodPost, deductionOfPoints, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.Header.Add("Authorization", "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4")

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusPaymentRequired, response.StatusCode)
}

func TestHandlerDeductionOfPointsBadJSON(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)

	repo := &storage.Repository{
		Withdraw: withdraw,
	}

	serv := service.NewService(repo, log)
	h := NewHandler(serv, log)

	body := []byte(`{"order":"2377225624}`)
	r := httptest.NewRequest(http.MethodPost, deductionOfPoints, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.Header.Add("Authorization", "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4")

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestHandlerGetWithdrawalOfPoints(t *testing.T) {
	type want struct {
		status int
		orders []model.WithdrawOrder
	}

	tests := []struct {
		name   string
		token  string
		userID int
		err    error
		want   want
	}{
		{
			name:   "correct_get_withdrawals",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 1,
			want: want{status: http.StatusOK,
				orders: []model.WithdrawOrder{
					{Order: 2377225624, UserID: 0, Sum: 100, ProcessedAt: time.Now().UTC().Add(time.Second * -10)},
					{Order: 2377225625, UserID: 0, Sum: 200, ProcessedAt: time.Now().UTC()},
				},
			},
		},
		{
			name:   "empty_withdraw",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 1,
			want: want{status: http.StatusNoContent,
				orders: []model.WithdrawOrder{},
			},
		},
		{
			name:   "error_bd",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			err:    errors.New("error bd"),
			userID: 1,
			want: want{status: http.StatusInternalServerError,
				orders: []model.WithdrawOrder{},
			},
		},
	}

	log := logger.Init()

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)

			withdraw.EXPECT().GetWithdrawalOfPoints(ctx, tt.userID).Return(tt.want.orders, tt.err).Times(1)

			repo := &storage.Repository{
				Withdraw: withdraw,
			}

			serv := service.NewService(repo, log)
			h := NewHandler(serv, log)

			r := httptest.NewRequest(http.MethodGet, getWithdrawalOfPoints, nil)
			w := httptest.NewRecorder()
			r.Header.Add("Authorization", tt.token)

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()

			defer response.Body.Close()
			body, _ := io.ReadAll(response.Body)

			var ordersResponse []model.WithdrawOrder
			ordersResponse = []model.WithdrawOrder{}
			_ = json.Unmarshal(body, &ordersResponse)

			assert.Equal(t, tt.want.status, response.StatusCode)
			assert.EqualValues(t, tt.want.orders, ordersResponse)
		})
	}
}

func TestHandlerGetWithdrawalOfPointsUserUnauthorized(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	withdraw := mock_storage.NewMockWithdrawOrderRepoContract(control)

	repo := &storage.Repository{
		Withdraw: withdraw,
	}

	serv := service.NewService(repo, log)
	h := NewHandler(serv, log)

	r := httptest.NewRequest(http.MethodGet, getWithdrawalOfPoints, nil)
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

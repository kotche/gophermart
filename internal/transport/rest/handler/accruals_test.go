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

func TestHandlerLoadOrders(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name   string
		token  string
		err    error
		userID int
		order  *model.AccrualOrder
		body   []byte
		want   want
	}{
		{
			name:   "uploading_new_order",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 0,
			order:  &model.AccrualOrder{Number: 2377225624, UserID: 1, Status: model.StatusNEW},
			body:   []byte("2377225624"),
			want:   want{status: http.StatusAccepted},
		},
		{
			name:   "bd_error",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 0,
			err:    errors.New("bd error"),
			order:  &model.AccrualOrder{Number: 2377225624, UserID: 1, Status: model.StatusNEW},
			body:   []byte("2377225624"),
			want:   want{status: http.StatusInternalServerError},
		},
	}

	t.Parallel()
	log := logger.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			accrual := mock_storage.NewMockAccrualOrderRepoContract(control)

			accrual.EXPECT().GetUserIDByNumberOrder(ctx, tt.order.Number).Return(tt.userID)
			accrual.EXPECT().SaveOrder(ctx, tt.order).Return(tt.err)

			repo := &storage.Repository{
				Accrual: accrual,
			}

			serv := service.NewService(repo, log)
			h := NewHandler(serv, log)

			r := httptest.NewRequest(http.MethodPost, loadOrders, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()
			r.Header.Add("Authorization", tt.token)

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.status, response.StatusCode)
		})
	}
}

func TestHandlerLoadOrdersNoCallsBD(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name  string
		token string
		body  []byte
		want  want
	}{
		{
			name:  "incorrect_number",
			token: "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			body:  []byte("ff237gg7225624"),
			want:  want{status: http.StatusBadRequest},
		},
		{
			name:  "empty_number",
			token: "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			body:  []byte(""),
			want:  want{status: http.StatusBadRequest},
		},
		{
			name:  "failed_Moon_test_number",
			token: "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			body:  []byte("111"),
			want:  want{status: http.StatusUnprocessableEntity},
		},
		{
			name:  "user_unauthorized",
			token: "",
			want:  want{status: http.StatusUnauthorized},
		},
	}

	t.Parallel()
	log := logger.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			accrual := mock_storage.NewMockAccrualOrderRepoContract(control)

			repo := &storage.Repository{
				Accrual: accrual,
			}

			serv := service.NewService(repo, log)
			h := NewHandler(serv, log)

			r := httptest.NewRequest(http.MethodPost, loadOrders, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()
			r.Header.Add("Authorization", tt.token)

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.status, response.StatusCode)
		})
	}
}

func TestHandlerLoadOrdersAlreadyUploaded(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name   string
		token  string
		userID int
		order  *model.AccrualOrder
		body   []byte
		want   want
	}{
		{
			name:   "already_been_uploaded_current_user",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 1,
			order:  &model.AccrualOrder{Number: 2377225624, UserID: 1, Status: model.StatusNEW},
			body:   []byte("2377225624"),
			want:   want{status: http.StatusOK},
		},
		{
			name:   "already_been_uploaded_another_user",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 2,
			order:  &model.AccrualOrder{Number: 2377225624, UserID: 1, Status: model.StatusNEW},
			body:   []byte("2377225624"),
			want:   want{status: http.StatusConflict},
		},
	}

	log := logger.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			accrual := mock_storage.NewMockAccrualOrderRepoContract(control)

			accrual.EXPECT().GetUserIDByNumberOrder(ctx, tt.order.Number).Return(tt.userID)

			repo := &storage.Repository{
				Accrual: accrual,
			}

			serv := service.NewService(repo, log)
			h := NewHandler(serv, log)

			r := httptest.NewRequest(http.MethodPost, loadOrders, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()
			r.Header.Add("Authorization", tt.token)

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.status, response.StatusCode)
		})
	}
}

func TestHandlerGetUploadedOrders(t *testing.T) {
	type want struct {
		status int
		orders []model.AccrualOrder
	}

	tests := []struct {
		name   string
		token  string
		userID int
		err    error
		want   want
	}{
		{
			name:   "correct_get_accruals",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 1,
			want: want{status: http.StatusOK,
				orders: []model.AccrualOrder{
					{Number: 2377225624, UserID: 0, Accrual: 100, UploadedAt: time.Now().UTC().Add(time.Second * -10)},
					{Number: 2377225625, UserID: 0, Accrual: 200, UploadedAt: time.Now().UTC()},
				},
			},
		},
		{
			name:   "empty_accruals",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			userID: 1,
			want: want{status: http.StatusNoContent,
				orders: []model.AccrualOrder{},
			},
		},
		{
			name:   "error_bd",
			token:  "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.AcccNNI5uz9-zKhms5yc5-uWWmwwJ8Bb_qKTM2-pcU4",
			err:    errors.New("error bd"),
			userID: 1,
			want: want{status: http.StatusInternalServerError,
				orders: []model.AccrualOrder{},
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
			accrual := mock_storage.NewMockAccrualOrderRepoContract(control)

			accrual.EXPECT().GetUploadedOrders(ctx, tt.userID).Return(tt.want.orders, tt.err).Times(1)

			repo := &storage.Repository{
				Accrual: accrual,
			}

			serv := service.NewService(repo, log)
			h := NewHandler(serv, log)

			r := httptest.NewRequest(http.MethodGet, getUploadedOrders, nil)
			w := httptest.NewRecorder()
			r.Header.Add("Authorization", tt.token)

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()

			defer response.Body.Close()
			body, _ := io.ReadAll(response.Body)

			ordersResponse := []model.AccrualOrder{}
			_ = json.Unmarshal(body, &ordersResponse)

			assert.Equal(t, tt.want.status, response.StatusCode)
			assert.EqualValues(t, tt.want.orders, ordersResponse)
		})
	}
}

func TestHandlerGetUploadedOrdersUserUnauthorized(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	accrual := mock_storage.NewMockAccrualOrderRepoContract(control)

	repo := &storage.Repository{
		Accrual: accrual,
	}

	serv := service.NewService(repo, log)
	h := NewHandler(serv, log)

	r := httptest.NewRequest(http.MethodGet, getUploadedOrders, nil)
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

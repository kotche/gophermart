package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kotche/gophermart/internal/logger"
	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
	"github.com/kotche/gophermart/internal/service"
	mock_service "github.com/kotche/gophermart/internal/service/mock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerRegistration(t *testing.T) {

	type want struct {
		status int
	}

	tests := []struct {
		name string
		user *model.User
		err  error
		body []byte
		want want
	}{
		{
			name: "correct_registration",
			user: &model.User{Login: "test", Password: "111"},
			body: []byte(`{"login": "test","password": "111"}`),
			want: want{status: http.StatusOK},
		},

		{
			name: "empty_login",
			user: &model.User{Login: "", Password: "111"},
			body: []byte(`{"password": "111"}`),
			want: want{status: http.StatusBadRequest},
		},
		{
			name: "empty_password",
			user: &model.User{Login: "test", Password: ""},
			body: []byte(`{"login": "test"}`),
			want: want{status: http.StatusBadRequest},
		},
	}

	log := logger.Init()

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			auth := mock_service.NewMockAuthServiceContract(control)

			ctx := context.Background()
			auth.EXPECT().CreateUser(ctx, tt.user).Return(tt.err).Times(1)
			auth.EXPECT().GenerateToken(tt.user, nil).Return("1", nil).Times(1)

			serv := &service.Service{
				Auth: auth,
			}

			h := &Handler{
				Service: serv,
				log:     log,
			}

			r := httptest.NewRequest(http.MethodPost, Registration, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()
			assert.Equal(t, tt.want.status, response.StatusCode)
		})
	}
}

func TestHandlerRegistrationConflictLogin(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	auth := mock_service.NewMockAuthServiceContract(control)

	user := &model.User{Login: "test", Password: "111"}

	ctx := context.Background()
	auth.EXPECT().CreateUser(ctx, user).Return(errormodel.ConflictLoginError{}).Times(1)

	serv := &service.Service{
		Auth: auth,
	}

	h := &Handler{
		Service: serv,
		log:     log,
	}

	body := []byte(`{"login": "test","password": "111"}`)
	r := httptest.NewRequest(http.MethodPost, Registration, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestHandlerRegistrationBadBodyJson(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	auth := mock_service.NewMockAuthServiceContract(control)

	serv := &service.Service{
		Auth: auth,
	}

	h := &Handler{
		Service: serv,
		log:     log,
	}

	body := []byte(`{"login}`)
	r := httptest.NewRequest(http.MethodPost, Registration, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestHandlerAuthentication(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name string
		user *model.User
		err  error
		body []byte
		want want
	}{
		{
			name: "correct_authentication",
			user: &model.User{Login: "test", Password: "111"},
			body: []byte(`{"login": "test","password": "111"}`),
			want: want{status: http.StatusOK},
		},
		{
			name: "empty_login",
			user: &model.User{Login: "", Password: "111"},
			body: []byte(`{"password": "111"}`),
			want: want{status: http.StatusBadRequest},
		},
		{
			name: "empty_password",
			user: &model.User{Login: "test", Password: ""},
			body: []byte(`{"login": "test"}`),
			want: want{status: http.StatusBadRequest},
		},
	}

	log := logger.Init()

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			auth := mock_service.NewMockAuthServiceContract(control)

			auth.EXPECT().AuthenticationUser(ctx, tt.user).Return(tt.err).Times(1)
			auth.EXPECT().GenerateToken(tt.user, nil).Return("1", nil).Times(1)

			serv := &service.Service{
				Auth: auth,
			}

			h := &Handler{
				Service: serv,
				log:     log,
			}

			r := httptest.NewRequest(http.MethodPost, Authentication, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			h.InitRoutes().ServeHTTP(w, r)

			response := w.Result()
			assert.Equal(t, tt.want.status, response.StatusCode)
		})
	}
}

func TestHandlerAuthenticationInvalidLoginOrPassword(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	auth := mock_service.NewMockAuthServiceContract(control)

	ctx := context.Background()
	user := &model.User{Login: "test", Password: "111"}
	auth.EXPECT().AuthenticationUser(ctx, user).Return(errormodel.AuthenticationError{}).Times(1)

	serv := &service.Service{
		Auth: auth,
	}

	h := &Handler{
		Service: serv,
		log:     log,
	}

	body := []byte(`{"login": "test","password": "111"}`)
	r := httptest.NewRequest(http.MethodPost, Authentication, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestHandlerAuthenticationBadBodyJSON(t *testing.T) {
	log := logger.Init()

	control := gomock.NewController(t)
	defer control.Finish()

	auth := mock_service.NewMockAuthServiceContract(control)

	serv := &service.Service{
		Auth: auth,
	}

	h := &Handler{
		Service: serv,
		log:     log,
	}

	body := []byte(`{"login}`)
	r := httptest.NewRequest(http.MethodPost, Authentication, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.InitRoutes().ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

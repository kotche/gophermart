package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kotche/gophermart/internal/logger"
	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
	"github.com/kotche/gophermart/internal/storage"
	"github.com/kotche/gophermart/internal/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_CreateUser(t *testing.T) {

	tests := []struct {
		name   string
		user   *model.User
		userID int
		err    error
	}{
		{
			name:   "good_create_user",
			user:   &model.User{Login: "test", Password: "111"},
			userID: 1,
		},
		{
			name: "bad_create_user",
			user: &model.User{Login: "test", Password: "111"},
			err:  errors.New("failed to create user"),
		},
		{
			name: "conflict_login",
			user: &model.User{Login: "test", Password: "111"},
			err:  errormodel.ConflictLoginError{Login: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			authRepo := mock_storage.NewMockAuthRepoContract(control)
			authRepo.EXPECT().CreateUser(ctx, tt.user).Return(tt.userID, tt.err)

			repo := &storage.Repository{Auth: authRepo}
			serv := NewService(repo, logger.Init())

			err := serv.Auth.CreateUser(ctx, tt.user)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAuthService_AuthenticationUser(t *testing.T) {
	tests := []struct {
		name   string
		user   *model.User
		userID int
		err    error
	}{
		{
			name:   "good_authentication_user",
			user:   &model.User{Login: "test", Password: "111"},
			userID: 1,
		},
		{
			name: "invalid_login_or_password",
			user: &model.User{Login: "test", Password: "111"},
			err:  errormodel.AuthenticationError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()
			authRepo := mock_storage.NewMockAuthRepoContract(control)
			authRepo.EXPECT().GetUserID(ctx, tt.user).Return(tt.userID, tt.err)

			repo := &storage.Repository{Auth: authRepo}
			serv := NewService(repo, logger.Init())

			err := serv.Auth.AuthenticationUser(ctx, tt.user)
			assert.Equal(t, tt.err, err)
		})
	}
}

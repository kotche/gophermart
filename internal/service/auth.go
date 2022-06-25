package service

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/kotche/gophermart/internal/model"
)

const (
	secretKey = "be55d1079e6c6167118ac91318fe"
)

type AuthRepoContract interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, login, password string) (model.User, error)
}

type AuthService struct {
	repo AuthRepoContract
}

func NewAuthService(repo AuthRepoContract) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (auth *AuthService) CreateUser(ctx context.Context, user model.User) error {
	user.Password = generatePasswordHash(user.Password)
	return auth.repo.CreateUser(ctx, user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(secretKey)))
}

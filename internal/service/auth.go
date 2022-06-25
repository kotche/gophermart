package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kotche/gophermart/internal/model"
)

const (
	secretKey  = "be55d1079e6c6167118ac91318fe"
	signingKey = "KSFjH$53KSFjH6745u#uEQQjF349%835hFpzA"
	tokenTTL   = 12 * time.Hour
)

type AuthRepoContract interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUserID(ctx context.Context, user *model.User) (string, error)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

type AuthService struct {
	repo AuthRepoContract
}

func NewAuthService(repo AuthRepoContract) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (auth *AuthService) CreateUser(ctx context.Context, user *model.User) error {
	user.Password = generatePasswordHash(user.Password)
	userID, err := auth.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (auth *AuthService) AuthenticationUser(ctx context.Context, user *model.User) error {
	user.Password = generatePasswordHash(user.Password)
	userID, err := auth.repo.GetUserID(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (auth *AuthService) GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})
	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(secretKey)))
}

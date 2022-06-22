package postgres

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (a *AuthPostgres) CreateUser(ctx context.Context, user model.User) (int, error) {
	return 0, nil
}

func (a *AuthPostgres) GetUser(ctx context.Context, login, password string) (model.User, error) {
	var user model.User
	return user, nil
}

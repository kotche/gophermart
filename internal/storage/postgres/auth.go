package postgres

import (
	"context"
	"database/sql"
	"errors"

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

func (a *AuthPostgres) CreateUser(ctx context.Context, user model.User) error {
	stmt, err := a.db.PrepareContext(ctx,
		"INSERT INTO public.users(login,password) VALUES ($1,$2) RETURNING login")
	if err != nil {
		return err
	}
	result := stmt.QueryRowContext(ctx, user.Login, user.Password)
	var output string
	result.Scan(&output)
	if output != user.Login {
		return model.ConflictLoginError{
			Err:   errors.New("duplicate login"),
			Login: user.Login,
		}
	}

	return nil
}

func (a *AuthPostgres) GetUser(ctx context.Context, login, password string) (model.User, error) {
	var user model.User
	return user, nil
}

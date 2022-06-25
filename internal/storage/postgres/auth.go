package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (a *AuthPostgres) CreateUser(ctx context.Context, user *model.User) (string, error) {
	stmt, err := a.db.PrepareContext(ctx,
		"INSERT INTO public.users(login,password) VALUES ($1,$2) RETURNING id")
	if err != nil {
		return "", err
	}
	result := stmt.QueryRowContext(ctx, user.Login, user.Password)
	var output sql.NullInt64
	result.Scan(&output)
	if !output.Valid {
		return "", model.ConflictLoginError{
			Err:   errors.New("duplicate login"),
			Login: user.Login,
		}
	}

	userID := fmt.Sprintf("%d", output.Int64)
	return userID, nil
}

func (a *AuthPostgres) GetUserID(ctx context.Context, user *model.User) (string, error) {
	row := a.db.QueryRowContext(ctx, "SELECT id FROM public.users WHERE login=$1 AND password=$2", user.Login, user.Password)
	var output sql.NullInt64
	row.Scan(&output)
	if !output.Valid {
		return "", model.AuthenticationError{
			Err: errors.New("invalid login/password"),
		}
	}
	userID := fmt.Sprintf("%d", output.Int64)
	return userID, nil
}

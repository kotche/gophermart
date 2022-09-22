package postgres

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
	"github.com/rs/zerolog"
)

type AuthPostgres struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewAuthPostgres(db *sql.DB, log *zerolog.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:  db,
		log: log,
	}
}

func (a *AuthPostgres) CreateUser(ctx context.Context, user *model.User) (int, error) {
	stmt, err := a.db.PrepareContext(ctx,
		"INSERT INTO public.users(login,password) VALUES ($1,$2) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result := stmt.QueryRowContext(ctx, user.Login, user.Password)
	var output sql.NullInt32
	_ = result.Scan(&output)
	if !output.Valid {
		return 0, errormodel.ConflictLoginError{
			Login: user.Login,
		}
	}
	userID := int(output.Int32)
	return userID, nil
}

func (a *AuthPostgres) GetUserID(ctx context.Context, user *model.User) (int, error) {
	row := a.db.QueryRowContext(ctx, "SELECT id FROM public.users WHERE login=$1 AND password=$2", user.Login, user.Password)
	var output sql.NullInt32
	_ = row.Scan(&output)
	if !output.Valid {
		return 0, errormodel.AuthenticationError{}
	}
	userID := int(output.Int32)
	return userID, nil
}

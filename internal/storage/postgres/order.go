package postgres

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

type OrderPostgres struct {
	db *sql.DB
}

func NewOrderPostgres(db *sql.DB) *OrderPostgres {
	return &OrderPostgres{
		db: db,
	}
}

func (o *OrderPostgres) SaveOrder(ctx context.Context, order *model.Order) error {
	userIDinDB := o.GetUserIDByNumberOrder(ctx, order.Number)

	if userIDinDB != 0 {
		if userIDinDB == order.UserID {
			return errormodel.OrderAlreadyUploadedCurrentUserError{}
		} else {
			return errormodel.OrderAlreadyUploadedAnotherUserError{}
		}
	}

	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO public.orders(order_num,user_id) VALUES ($1,$2)", order.Number, order.UserID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO public.accruals(order_num,user_id,status,uploaded_at) VALUES ($1,$2,$3,$4)",
		order.Number, order.UserID, order.Status, order.UploadedAt)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (o *OrderPostgres) GetUserIDByNumberOrder(ctx context.Context, number string) int {
	row := o.db.QueryRowContext(ctx, "SELECT user_id FROM public.accruals WHERE order_num=$1", number)
	var userID int
	row.Scan(&userID)
	return userID
}

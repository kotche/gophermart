package postgres

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errorModel"
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

	if userIDinDB != "" {
		if userIDinDB == order.UserID {
			return errorModel.OrderAlreadyUploadedCurrentUserError{}
		} else {
			return errorModel.OrderAlreadyUploadedAnotherUserError{}
		}
	}

	_, err := o.db.ExecContext(ctx,
		"INSERT INTO public.accruals(order_num,user_id,status,uploaded_at) VALUES ($1,$2,$3,$4)",
		order.Number, order.UserID, order.Status, order.UploadedAt)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderPostgres) GetUserIDByNumberOrder(ctx context.Context, number string) string {
	row := o.db.QueryRowContext(ctx, "SELECT user_id FROM public.accruals WHERE order_num=$1", number)
	var userID string
	row.Scan(&userID)
	return userID
}

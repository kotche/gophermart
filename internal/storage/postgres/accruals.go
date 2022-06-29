package postgres

import (
	"context"
	"database/sql"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

type AccrualOrderPostgres struct {
	db *sql.DB
}

func NewAccrualOrderPostgres(db *sql.DB) *AccrualOrderPostgres {
	return &AccrualOrderPostgres{
		db: db,
	}
}

func (a *AccrualOrderPostgres) SaveOrder(ctx context.Context, order *model.AccrualOrder) error {
	userIDinDB := a.GetUserIDByNumberOrder(ctx, order.Number)

	if userIDinDB != 0 {
		if userIDinDB == order.UserID {
			return errormodel.OrderAlreadyUploadedCurrentUserError{}
		} else {
			return errormodel.OrderAlreadyUploadedAnotherUserError{}
		}
	}

	tx, err := a.db.Begin()
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

func (a *AccrualOrderPostgres) GetUserIDByNumberOrder(ctx context.Context, number string) int {
	row := a.db.QueryRowContext(ctx, "SELECT user_id FROM public.accruals WHERE order_num=$1", number)
	var userID int
	row.Scan(&userID)

	return userID
}

func (a *AccrualOrderPostgres) GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error) {
	rows, err := a.db.QueryContext(ctx, "SELECT order_num,status,amount,uploaded_at FROM public.accruals WHERE user_id =$1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []model.AccrualOrder
	for rows.Next() {
		var order model.AccrualOrder
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

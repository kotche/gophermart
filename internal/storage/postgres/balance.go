package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kotche/gophermart/internal/model"
	"github.com/rs/zerolog"
)

type WithdrawOrderPostgres struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewWithdrawOrderPostgres(db *sql.DB, log *zerolog.Logger) *WithdrawOrderPostgres {
	return &WithdrawOrderPostgres{
		db:  db,
		log: log,
	}
}

func (w *WithdrawOrderPostgres) GetAccruals(ctx context.Context, UserID int) float32 {
	row := w.db.QueryRowContext(ctx, "SELECT SUM(amount) FROM public.accruals WHERE user_id=$1", UserID)
	var accruals float32
	_ = row.Scan(&accruals)

	return accruals
}

func (w *WithdrawOrderPostgres) GetWithdrawals(ctx context.Context, UserID int) float32 {
	row := w.db.QueryRowContext(ctx, "SELECT SUM(amount) FROM public.withdrawals WHERE user_id=$1", UserID)
	var withdrawals float32
	_ = row.Scan(&withdrawals)

	return withdrawals
}

func (w *WithdrawOrderPostgres) DeductPoints(ctx context.Context, order *model.WithdrawOrder) (err error) {
	order.ProcessedAt = time.Now()

	tx, err := w.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			txError := tx.Rollback()
			if txError != nil {
				err = fmt.Errorf("balance DeductPoints rollback error %s: %s", txError.Error(), err.Error())
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO public.orders(order_num,user_id) VALUES ($1,$2)", order.Order, order.UserID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO public.withdrawals(order_num,user_id,amount,processed_at) VALUES ($1,$2,$3,$4)",
		order.Order, order.UserID, order.Sum, order.ProcessedAt)

	if err != nil {
		return err
	}
	return tx.Commit()
}

func (w *WithdrawOrderPostgres) GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error) {
	rows, err := w.db.QueryContext(ctx, "SELECT order_num,amount,processed_at FROM public.withdrawals WHERE user_id =$1 ORDER BY processed_at", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []model.WithdrawOrder
	for rows.Next() {
		var order model.WithdrawOrder
		err = rows.Scan(&order.Order, &order.Sum, &order.ProcessedAt)
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

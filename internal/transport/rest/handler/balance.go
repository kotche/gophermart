package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

type balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

//getCurrentBalance GET /api/user/balance - получение текущего баланса пользователя
func (h *Handler) getCurrentBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserIDFromToken(w, r, "handler.getCurrentBalance")
	if err != nil {
		return
	}

	ctx := context.Background()
	accruals, withdraws := h.Service.GetBalance(ctx, userID)

	b := balance{Current: accruals, Withdrawn: withdraws}

	output, err := json.Marshal(b)
	if err != nil {
		log.Printf("handler.getCurrentBalance - json write error: %s", err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

//deductionOfPoints POST /api/user/balance/withdraw - запрос на списание средств
func (h *Handler) deductionOfPoints(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(w, r, "handler.deductionOfPoints")
	if err != nil {
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("handler.deductionOfPoints - body read error: %s", err.Error())
		http.Error(w, "incorrect input data", http.StatusInternalServerError)
		return
	}

	var order *model.WithdrawOrder
	err = json.Unmarshal(body, &order)
	if err != nil {
		log.Printf("handler.deductionOfPoints - json read error: %s", err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	order.UserID = userID
	order.ProcessedAt = time.Now()

	ctx := context.Background()
	err = h.Service.DeductionOfPoints(ctx, order)

	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusOK)
	case errormodel.NotEnoughPoints:
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

//deductionOfPoints GET /api/user/balance/withdrawals - получение информации о выводе средств
func (h *Handler) getWithdrawalOfPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserIDFromToken(w, r, "handler.getCurrentBalance")
	if err != nil {
		return
	}

	ctx := context.Background()
	orders, err := h.Service.GetWithdrawalOfPoints(ctx, userID)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	output, err := json.Marshal(orders)
	if err != nil {
		log.Printf("handler.getWithdrawalOfPoints json marshal error: %s", err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

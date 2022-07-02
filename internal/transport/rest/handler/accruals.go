package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/kotche/gophermart/internal/model/errormodel"
)

// loadOrders POST /api/user/orders - загрузка номера заказа
func (h *Handler) loadOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(w, r, "handler.loadOrders")
	if err != nil {
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("handler.loadOrders - body read error: %s", err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		log.Println("handler.loadOrders - body empty")
		http.Error(w, "incorrect input data", http.StatusBadRequest)
		return
	}

	numOrder := string(body)

	ctx := context.Background()
	err = h.Service.LoadOrder(ctx, numOrder, userID)

	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	case errormodel.OrderAlreadyUploadedCurrentUserError:
		http.Error(w, err.Error(), http.StatusOK)
		return
	case errormodel.OrderAlreadyUploadedAnotherUserError:
		http.Error(w, err.Error(), http.StatusConflict)
		return
	case errormodel.CheckLuhnError:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

// getUploadedOrders GET /api/user/orders - получение списка загруженных номеров заказов
func (h *Handler) getUploadedOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	userID, err := getUserIDFromToken(w, r, "handler.getUploadedOrders")
	if err != nil {
		return
	}

	ctx := context.Background()
	orders, err := h.Service.GetUploadedOrders(ctx, userID)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	if orders == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	output, err := json.Marshal(orders)
	if err != nil {
		log.Printf("getUploadedOrders -json marshal error: %s", err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

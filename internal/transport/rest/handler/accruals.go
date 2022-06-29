package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

// loadOrders POST /api/user/orders - загрузка номера заказа
func (h *Handler) loadOrders(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Printf("loadOrders - jwt claims error: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		log.Printf("loadOrders - conv user_id to int: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("loadOrders - body read error: %s", err.Error())
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// getUploadedOrders GET /api/user/orders - получение списка загруженных номеров заказов
func (h *Handler) getUploadedOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Printf("getUploadedOrders - jwt claims error: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		log.Printf("getUploadedOrders - conv user_id to int: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	orders, err := h.Service.GetUploadedOrders(ctx, userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if orders == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	output, err := json.Marshal(orders)
	if err != nil {
		log.Printf("getUploadedOrders -json marshal error: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

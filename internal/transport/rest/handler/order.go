package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

// loadOrder POST /api/user/orders - загрузка номера заказа
func (h *Handler) loadOrder(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "incorrect input data", http.StatusBadRequest)
		return
	}
	numOrder := string(body)
	userID := fmt.Sprintf("%v", claims["user_id"])

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
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

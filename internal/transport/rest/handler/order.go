package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// loadOrder POST /api/user/orders - загрузка номера заказа
func (h *Handler) loadOrder(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	userID := fmt.Sprintf("%v", claims["user_id"])
	log.Println(userID)
	//var order model.Order

	//ctx := context.Background()
	//h.Service.LoadOrder(ctx, &order)
}

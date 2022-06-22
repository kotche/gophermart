package handler

import (
	"context"
	"net/http"

	"github.com/kotche/gophermart/internal/model"
)

// LoadOrder POST /api/user/orders - загрузка номера заказа
func (h *Handler) LoadOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	ctx := context.Background()
	h.Service.LoadOrder(ctx, order)
}

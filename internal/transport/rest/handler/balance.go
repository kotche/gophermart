package handler

import (
	"context"
	"net/http"
)

// getCurrentBalance GET /api/user/balance - получение текущего баланса пользователя
func (h *Handler) getCurrentBalance(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.Service.GetCurrentBalance(ctx)
}

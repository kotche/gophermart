package handler

import (
	"context"
	"net/http"
)

// GetCurrentBalance GET /api/user/balance - получение текущего баланса пользователя
func (h *Handler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.Service.GetCurrentBalance(ctx)
}

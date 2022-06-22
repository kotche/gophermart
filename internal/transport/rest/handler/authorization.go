package handler

import (
	"context"
	"net/http"

	"github.com/kotche/gophermart/internal/model"
)

// Register POST /api/user/register - регистрация пользователя
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	ctx := context.Background()
	h.Service.CreateUser(ctx, user)
}

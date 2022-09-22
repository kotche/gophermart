package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

// registration POST /api/user/register - регистрация пользователя
func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := h.readingUserData(w, r, &user, "registration")
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.Auth.CreateUser(ctx, &user)

	if errors.As(err, &errormodel.ConflictLoginError{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		h.log.Error().Err(err).Msg("Handler.registration: CreateUser service error")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.writeToken(w, &user, "registration")
}

//authentication POST /api/user/login - аутентификация пользователя
func (h *Handler) authentication(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := h.readingUserData(w, r, &user, "authentication")
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.Auth.AuthenticationUser(ctx, &user)

	if errors.As(err, &errormodel.AuthenticationError{}) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		h.log.Error().Err(err).Msg("Handler.authentication: AuthenticationUser service error")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.writeToken(w, &user, "authentication")
}

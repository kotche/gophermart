package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/kotche/gophermart/internal/model"
	"github.com/kotche/gophermart/internal/model/errormodel"
)

// registration POST /api/user/register - регистрация пользователя
func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := readingUserData(w, r, &user, "handler.registration")
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.CreateUser(ctx, &user)

	if errors.As(err, &errormodel.ConflictLoginError{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.writeToken(w, &user, "handler.registration")
}

//authentication POST /api/user/login - аутентификация пользователя
func (h *Handler) authentication(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := readingUserData(w, r, &user, "handler.authentication")
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.AuthenticationUser(ctx, &user)

	if errors.As(err, &errormodel.AuthenticationError{}) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.writeToken(w, &user, "handler.authentication")
}

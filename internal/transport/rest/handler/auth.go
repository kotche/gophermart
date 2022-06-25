package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/kotche/gophermart/internal/model"
)

// registration POST /api/user/register - регистрация пользователя
func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := readingUserData(w, r, &user)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.CreateUser(ctx, &user)

	if errors.As(err, &model.ConflictLoginError{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.writeToken(w, &user)
}

//authentication POST /api/user/login - аутентификация пользователя
func (h *Handler) authentication(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := readingUserData(w, r, &user)
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.AuthenticationUser(ctx, &user)
	if errors.As(err, &model.AuthorizationError{}) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.writeToken(w, &user)
}

func (h *Handler) writeToken(w http.ResponseWriter, user *model.User) {
	token, err := h.Service.GenerateToken(user, h.TokenAuth)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "BEARER "+token)
}

func readingUserData(w http.ResponseWriter, r *http.Request, user *model.User) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "incorrect input data", http.StatusBadRequest)
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "json read error", http.StatusInternalServerError)
		return err
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "empty login or password", http.StatusBadRequest)
		return err
	}
	return nil
}

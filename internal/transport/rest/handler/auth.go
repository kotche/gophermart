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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "incorrect input data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "json read error", http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "empty login or password", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = h.Service.CreateUser(ctx, user)

	if errors.As(err, &model.ConflictLoginError{}) {
		http.Error(w, err.Error(), http.StatusConflict)
	} else if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

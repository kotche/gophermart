package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/kotche/gophermart/internal/model"
)

const (
	internalServerError = "internal server error"
)

func (h *Handler) writeToken(w http.ResponseWriter, user *model.User, nameFunc string) {
	token, err := h.Service.Auth.GenerateToken(user, h.TokenAuth)
	if err != nil {
		h.log.Error().Err(err).Msgf("Handler.writeToken: %s - token generate error", nameFunc)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "BEARER "+token)
}

func (h *Handler) readingUserData(w http.ResponseWriter, r *http.Request, user *model.User, nameFunc string) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error().Err(err).Msgf("Handler.readingUserData: %s - body read error", nameFunc)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		h.log.Error().Err(err).Msgf("Handler.readingUserData: %s - json read error", nameFunc)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return err
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "empty login or password", http.StatusBadRequest)
		return err
	}
	return nil
}

func (h *Handler) getUserIDFromToken(w http.ResponseWriter, r *http.Request, nameFunc string) (int, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msgf("Handler.getUserIDFromToken: %s - jwt claims error", nameFunc)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return 0, err
	}

	userID, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		h.log.Error().Err(err).Msgf("Handler.getUserIDFromToken: %s - conv string to int", nameFunc)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return 0, err
	}

	return userID, nil
}

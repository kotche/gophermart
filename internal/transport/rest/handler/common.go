package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/kotche/gophermart/internal/model"
)

const (
	internalServerError = "internal server error"
)

func (h *Handler) writeToken(w http.ResponseWriter, user *model.User, nameFunc string) {
	token, err := h.Service.GenerateToken(user, h.TokenAuth)
	if err != nil {
		log.Printf("%s - token generate error: %s", nameFunc, err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "BEARER "+token)
}

func readingUserData(w http.ResponseWriter, r *http.Request, user *model.User, nameFunc string) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s - body read error: %s", nameFunc, err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("%s - json read error: %s", nameFunc, err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return err
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, "empty login or password", http.StatusBadRequest)
		return err
	}
	return nil
}

func getUserIDFromToken(w http.ResponseWriter, r *http.Request, nameFunc string) (int, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Printf("%s - jwt claims error: %s", nameFunc, err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return 0, err
	}

	userID, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		log.Printf("%s - conv string to int: %s", nameFunc, err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return 0, err
	}

	return userID, nil
}

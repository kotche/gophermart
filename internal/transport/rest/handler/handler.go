package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/kotche/gophermart/internal/service"
	"github.com/kotche/gophermart/internal/transport/rest/middlewares"
)

const (
	signingKey = "KSFjH$53KSFjH6745u#uEQQjF349%835hFpzA"
)

type Handler struct {
	Service   *service.Service
	TokenAuth *jwtauth.JWTAuth
}

func NewHandler(service *service.Service) *Handler {
	tokenAuth := jwtauth.New("HS256", []byte(signingKey), nil)

	handler := &Handler{
		Service:   service,
		TokenAuth: tokenAuth,
	}
	return handler
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares.GzipHandle)

	// Public routes
	router.Group(func(router chi.Router) {
		router.Post("/api/user/register", h.registration)
		router.Post("/api/user/login", h.authentication)
	})

	// Protected routes
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(h.TokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Post("/api/user/orders", h.loadOrder)
	})

	return router
}

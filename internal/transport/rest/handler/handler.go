package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kotche/gophermart/internal/service"
	"github.com/kotche/gophermart/internal/transport/rest/middlewares"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		Service: service,
	}
	return handler
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares.GzipHandle)
	router.Use(middlewares.Authorization)

	router.Post("/api/user/register", h.registration)
	router.Post("/api/user/orders", h.loadOrder)
	router.Get("/api/user/balance", h.getCurrentBalance)

	return router
}

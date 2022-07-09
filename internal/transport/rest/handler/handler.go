package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/kotche/gophermart/internal/service"
	"github.com/kotche/gophermart/internal/transport/rest/middlewares"
)

const (
	registration          = "/api/user/register"
	authentication        = "/api/user/login"
	loadOrders            = "/api/user/orders"
	getUploadedOrders     = "/api/user/orders"
	deductionOfPoints     = "/api/user/balance/withdraw"
	getWithdrawalOfPoints = "/api/user/withdrawals"
	getCurrentBalance     = "/api/user/balance"

	signingKey = "KSFjH$53KSFjH6745u#uEQQjF349%835hFpzA"
)

type Handler struct {
	Service   *service.Service
	TokenAuth *jwtauth.JWTAuth
}

func NewHandler(service *service.Service) *Handler {
	tokenAuth := jwtauth.New("HS256", []byte(signingKey), nil)

	return &Handler{
		Service:   service,
		TokenAuth: tokenAuth,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares.GzipHandle)

	// Public routes
	router.Group(func(router chi.Router) {
		router.Post(registration, h.registration)
		router.Post(authentication, h.authentication)
	})

	// Protected routes
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(h.TokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Post(loadOrders, h.loadOrders)
		router.Get(getUploadedOrders, h.getUploadedOrders)
		router.Post(deductionOfPoints, h.deductionOfPoints)
		router.Get(getWithdrawalOfPoints, h.getWithdrawalOfPoints)
		router.Get(getCurrentBalance, h.getCurrentBalance)
	})

	return router
}

package server

import (
	"context"
	"net/http"
	"time"

	"github.com/kotche/gophermart/internal/config"
)

const (
	idleTimeout  = 60 * time.Second
	readTimeout  = 60 * time.Second
	writeTimeout = 60 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.GophermartAddr,
			Handler:      handler,
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	brokerService "github.com/kotche/gophermart/internal/broker/service"
	"github.com/kotche/gophermart/internal/config"
	"github.com/kotche/gophermart/internal/logger"
	"github.com/kotche/gophermart/internal/server"
	"github.com/kotche/gophermart/internal/service"
	"github.com/kotche/gophermart/internal/storage"
	"github.com/kotche/gophermart/internal/storage/postgres"
	"github.com/kotche/gophermart/internal/transport/rest/handler"
)

func main() {
	log := logger.Init()

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("configuration error")
	}

	pgx, err := postgres.NewPGX(conf.DBConnect)
	if err != nil {
		log.Fatal().Err(err).Msg("db connection error")
	}
	err = pgx.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("creating tables error")
	}

	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	repos := storage.NewRepository(pgx.DB, log)
	services := service.NewService(repos, log)
	handlers := handler.NewHandler(services, log)

	brokerRepos := storage.NewBrokerRepository(pgx.DB, log)
	broker := brokerService.NewBroker(brokerRepos, conf.AccrualAddr, log)
	broker.Start(ctx)

	srv := server.NewServer(conf, handlers.InitRoutes())

	//graceful shutdown
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-termChan
		log.Info().Msg("server shutdown")
		cansel()
		if err = srv.Stop(ctx); err != nil {
			log.Fatal().Err(err).Msg("server shutdown error")
		}
	}()

	if err = srv.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server run error")
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	brokerService "github.com/kotche/gophermart/internal/broker/service"
	"github.com/kotche/gophermart/internal/config"
	"github.com/kotche/gophermart/internal/service"
	"github.com/kotche/gophermart/internal/storage"
	"github.com/kotche/gophermart/internal/storage/postgres"
	"github.com/kotche/gophermart/internal/transport/rest/handler"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Сonfiguration error: %s", err.Error())
	}

	pgx, err := postgres.NewPGX(conf.DBConnect)
	if err != nil {
		log.Fatalf("DB connection error: %s", err.Error())
	}
	err = pgx.Init()
	if err != nil {
		log.Fatalf("Error creating tables: %s", err.Error())
	}

	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	repos := storage.NewRepository(pgx.DB)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	brokerRepos := storage.NewBrokerRepository(pgx.DB)
	broker := brokerService.NewBroker(brokerRepos, conf.AccrualAddr)
	broker.Start(ctx)

	srv := &http.Server{
		Addr:         conf.GophermartAddr,
		Handler:      handlers.InitRoutes(),
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	//graceful shutdown
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-termChan
		log.Println("server shutdown")
		srv.Shutdown(ctx)
		cansel()
	}()

	log.Fatal(srv.ListenAndServe())
}

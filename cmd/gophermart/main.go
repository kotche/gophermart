package main

import (
	"log"
	"net/http"

	service2 "github.com/kotche/gophermart/internal/broker/service"
	storage2 "github.com/kotche/gophermart/internal/broker/storage"
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

	repos := storage.NewRepository(pgx.DB)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	repos2 := storage2.NewRepository(pgx.DB)
	broker := service2.NewBroker(repos2, conf.AccrualAddr)
	broker.Start()

	log.Fatal(http.ListenAndServe(conf.GophermartAddr, handlers.InitRoutes()))
}

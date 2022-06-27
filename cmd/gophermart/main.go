package main

import (
	"log"
	"net/http"

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

	log.Fatal(http.ListenAndServe(conf.GophermartAddr, handlers.InitRoutes()))
}

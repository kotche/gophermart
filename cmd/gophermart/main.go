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
		log.Fatal(err.Error())
		return
	}

	pgx, err := postgres.NewPGX(conf.DBConnect)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	pgx.Init()

	repos := storage.NewRepository(pgx.DB)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	log.Fatal(http.ListenAndServe(conf.GophermartAddr, handlers.InitRoutes()))
}

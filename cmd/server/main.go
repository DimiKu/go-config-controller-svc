package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-config-controller-svc/internal/configs"
	"go-config-controller-svc/internal/handlers"
	"go-config-controller-svc/internal/repos"
	"go-config-controller-svc/internal/service/server_service"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log, _ := zap.NewProduction()
	ctx := context.Background()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	r := chi.NewRouter()

	var cfg configs.ServerConfig
	if err := env.Parse(&cfg); err != nil {
		log.Error("Error parse env: ", zap.Error(err))
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Error("Failed to get connect: ", zap.Error(err))
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Error("Failed to get pool: ", zap.Error(err))
	}

	dbRepo := repos.NewServerDBRepo(conn, pool, log)
	service := server_service.NewServerService(dbRepo, log)

	r.Post("/create_config", handlers.CreateConfigHandler(service, log, ctx))

	err = http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}

}

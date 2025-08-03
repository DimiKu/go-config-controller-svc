package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go-config-controller-svc/internal/configs"
	"go-config-controller-svc/internal/handlers"
	"go-config-controller-svc/internal/middlewares"
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
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)

	// TODO signal.NotifyContext()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	r := chi.NewRouter()
	var cfg configs.ServerConfig

	r.Use(middlewares.AuthMiddleware([]byte(cfg.Secret)))

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

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	})

	redisRepo := repos.NewRedisRepo(rdb, ctx, log)

	dbRepo := repos.NewServerDBRepo(conn, pool, log)
	service := server_service.NewServerService(dbRepo, redisRepo, log)

	r.Post("/login", handlers.LoginUserHandler(service, log, []byte(cfg.Secret), ctx))
	r.Post("/create_user", handlers.CreateUserHandler(service, log, ctx))

	r.Post("/execute", handlers.CreateTaskHandler(service, log, ctx))
	r.Post("/create_config", handlers.CreateConfigHandler(service, log, ctx))
	r.Get("/get_configs", handlers.ListConfigHandler(service, log, ctx))
	r.Post("/delete_config", handlers.DeleteConfigHandler(service, log, ctx))

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	go func() {
		log.Info("Server starting on: ", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error: %v", zap.Error(err))
		}
	}()

	<-signalChan
	log.Info("Start gracefull shutdown and closed db conn")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error: %v", zap.Error(err))
	}

	log.Info("Close db connection")
	conn.Close(ctx)
	pool.Close()
	cancel()
}

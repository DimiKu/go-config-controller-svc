package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-config-controller-svc/internal/configs"
	"go-config-controller-svc/internal/custom_errors"
	"go-config-controller-svc/internal/executors"
	"go-config-controller-svc/internal/repos"
	"go-config-controller-svc/internal/service/controller_service"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var cfg configs.ControllerConfig
	if err = env.Parse(&cfg); err != nil {
		log.Error("Error parse env: ", zap.Error(err))
	}

	ctx := context.Background()

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Error("Failed to get connect: ", zap.Error(err))
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Error("Failed to get pool: ", zap.Error(err))
	}

	dbRepo := repos.NewAgentDBRepo(conn, pool, log)
	fileRepo := repos.NewFileRepo("./config_test", log)
	gitRepo := repos.NewGitControllerRepo("./config_test", cfg.GitUser, cfg.GitToken, cfg.GitRepo, log)
	//simpleExecutor := executors.NewPrintExec()
	simpleNginxExecutor := executors.NewNginxExec()
	configController := controller_service.NewConfigControllerService(dbRepo, gitRepo, fileRepo, simpleNginxExecutor, log)

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	var wg sync.WaitGroup

	for i := 1; i <= cfg.Workers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					log.Warn("Controller stopped")

				case <-ticker.C:
					if err := configController.Work(ctx); err != nil {
						if errors.Is(err, custom_errors.ErrNotLockedConfigNotFound) {
							log.Debug("Controller has error: ", zap.Error(err))
						} else {
							log.Error("Controller has error: ", zap.Error(err))
						}
					}
				}
			}
		}(i)
	}

	go func() {
		<-signalChan
		log.Info("Start gracefull shutdown and closed db conn")

		conn.Close(ctx)
		pool.Close()
		//cancel()

		os.Exit(0)
	}()

	wg.Wait()
}

package repos

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-config-controller-svc/internal/entities"
	"go-config-controller-svc/internal/utils"
	"go.uber.org/zap"
)

type ServerDBRepo struct {
	conn *pgx.Conn
	pool *pgxpool.Pool

	log *zap.Logger
}

func NewServerDBRepo(conn *pgx.Conn, pool *pgxpool.Pool, log *zap.Logger) *ServerDBRepo {
	ctx := context.Background()
	_, err := conn.Exec(ctx, CreateConfigTable)
	if err != nil {
		log.Error("Failed to create config table: ", zap.Error(err))
	}

	return &ServerDBRepo{conn: conn, pool: pool, log: log}
}

func (s *ServerDBRepo) SaveConfig(ctx context.Context, config entities.ServerConfig) error {

	_, err := utils.RetryableExec(ctx, s.pool, s.log, InsertNewConfig, config.ConfigName, config.ConfigValue, config.ConfigBranch)
	if err != nil {
		s.log.Error("Failed to create new config in db: ", zap.Error(err))
		return err
	}

	return nil
}

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
func (s *ServerDBRepo) GetConfigsList(ctx context.Context) ([]entities.ServerConfig, error) {
	r, err := utils.RetryableQuery(ctx, s.pool, s.log, GetAllConfigs)
	if err != nil {
		s.log.Error("Failed to get config: ", zap.Error(err))
		return nil, err
	}

	defer r.Close()
	var rConfigs []entities.ServerConfig

	for r.Next() {
		var c entities.ServerConfig
		if err := r.Scan(&c.ConfigName, &c.ConfigValue, &c.ConfigBranch); err != nil {
			return nil, err
		}

		rConfigs = append(rConfigs, c)
	}
	return rConfigs, nil
}

func (s *ServerDBRepo) DeleteConfig(ctx context.Context, config entities.ServerConfig) error {
	_, err := utils.RetryableExec(ctx, s.pool, s.log, DeleteConfig, config.ConfigName, config.ConfigBranch)
	if err != nil {
		s.log.Error("Failed to create new config in db: ", zap.Error(err))
		return err
	}

	return nil
}

func (s *ServerDBRepo) GetUserByUsername(ctx context.Context, username string) (entities.User, error) {
	r, err := utils.RetryableQuery(ctx, s.pool, s.log, GetUserByName, username)
	if err != nil {
		s.log.Error("Failed to get user from db: ", zap.Error(err))
		return entities.User{}, err
	}

	for r.Next() {
		var u entities.User
		if err := r.Scan(&u.Username, &u.Password); err != nil {
			return entities.User{}, err
		}

		return u, nil
	}

	return entities.User{}, nil
}

func (s *ServerDBRepo) CreateUser(ctx context.Context, user entities.User) error {
	_, err := utils.RetryableExec(ctx, s.pool, s.log, CreateUser, user.Username, user.Password)
	if err != nil {
		s.log.Error("Failed to create new user in db: ", zap.Error(err))
		return err
	}

	return nil
}

func (s *ServerDBRepo) CheckIfUserExists(ctx context.Context, username string) (bool, error) {
	r, err := utils.RetryableQuery(ctx, s.pool, s.log, GetUserByName, username)
	if err != nil {
		s.log.Error("Failed to get user from db: ", zap.Error(err))
		return false, err
	}

	if !r.Next() {
		return false, nil
	}

	return true, nil
}

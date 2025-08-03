package repos

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-config-controller-svc/dto/controller_dto"
	"go-config-controller-svc/internal/custom_errors"
	"go-config-controller-svc/internal/entities"
	"go-config-controller-svc/internal/utils"
	"go.uber.org/zap"
)

type AgentDBRepo struct {
	conn *pgx.Conn
	pool *pgxpool.Pool

	log *zap.Logger
}

func NewAgentDBRepo(conn *pgx.Conn, pool *pgxpool.Pool, log *zap.Logger) *AgentDBRepo {
	ctx := context.Background()
	_, err := conn.Exec(ctx, CreateConfigTable)
	if err != nil {
		log.Error("Failed to create config table: ", zap.Error(err))
	}

	_, err = conn.Exec(ctx, UnlockAllConfigs)
	if err != nil {
		log.Error("Failed to create config table: ", zap.Error(err))
	}

	return &AgentDBRepo{conn: conn, pool: pool, log: log}
}

func (a *AgentDBRepo) GetRepoByName(ctx context.Context, repoName string) (controller_dto.ConfigDBDto, error) {
	var repo controller_dto.ConfigDBDto
	r, err := utils.RetryableQuery(ctx, a.pool, a.log, GetRepoByName, repoName)
	if err != nil {
		a.log.Error("Failed to get config: ", zap.Error(err))
		return controller_dto.ConfigDBDto{}, err
	}

	if !r.Next() {
		return controller_dto.ConfigDBDto{}, custom_errors.ErrNotLockedConfigNotFound
	}

	err = r.Scan(&repo.ConfigName, &repo.ConfigValue, &repo.LastUpdate, &repo.ConfigBranch, &repo.Locked)
	if err != nil {
		a.log.Error("Failed to scan repo value: ", zap.Error(err))
		return controller_dto.ConfigDBDto{}, err
	}

	_, err = utils.RetryableExec(ctx, a.pool, a.log, LockRepo, repoName)
	if err != nil {
		a.log.Error("Failed to lock repo in db: ", zap.Error(err))
		return controller_dto.ConfigDBDto{}, err
	}

	return repo, nil
}

func (a *AgentDBRepo) ChangeRepoUpdateTime(ctx context.Context, repoName string) error {
	_, err := utils.RetryableExec(ctx, a.pool, a.log, ChangeRepoUpdateTime, repoName)
	if err != nil {
		a.log.Error("Failed to update last_update field: ", zap.Error(err))
		return err
	}

	_, err = utils.RetryableExec(ctx, a.pool, a.log, UnlockRepo, repoName)
	if err != nil {
		a.log.Error("Failed to lock repo in db: ", zap.Error(err))
		return err
	}

	return nil
}

func (a *AgentDBRepo) GetNotLockedConfig(ctx context.Context) (entities.ControllerConfig, error) {
	var repo controller_dto.ConfigDBDto
	r, err := utils.RetryableQuery(ctx, a.pool, a.log, GetNotLockedConfig)
	if err != nil {
		a.log.Error("Failed to get config: ", zap.Error(err))
		return entities.ControllerConfig{}, err
	}

	if !r.Next() {
		return entities.ControllerConfig{}, custom_errors.ErrNotLockedConfigNotFound
	}

	err = r.Scan(&repo.ConfigName, &repo.ConfigValue, &repo.LastUpdate, &repo.ConfigBranch, &repo.Locked)
	if err != nil {
		a.log.Error("Failed to scan repo value: ", zap.Error(err))
		return entities.ControllerConfig{}, err
	}

	_, err = utils.RetryableExec(ctx, a.pool, a.log, LockRepo, repo.ConfigName)
	if err != nil {
		a.log.Error("Failed to lock repo in db: ", zap.Error(err))
		return entities.ControllerConfig{}, err
	}

	return entities.ControllerConfig{
		ConfigName:   repo.ConfigName,
		ConfigBranch: repo.ConfigBranch,
		LastUpdate:   repo.LastUpdate,
	}, nil
}

func (a *AgentDBRepo) UnlockRepo(ctx context.Context, repoName string) error {
	_, err := utils.RetryableExec(ctx, a.pool, a.log, UnlockRepo, repoName)
	if err != nil {
		a.log.Error("Failed to lock repo in db: ", zap.Error(err))
		return err
	}

	return nil
}

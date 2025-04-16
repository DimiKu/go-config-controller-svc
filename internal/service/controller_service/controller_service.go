package controller_service

import (
	"context"
	"errors"
	"go-config-controller-svc/dto/controller_dto"
	"go-config-controller-svc/internal/custom_errors"
	"go-config-controller-svc/internal/entities"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Executor interface {
	DoChangeWithNewValues(configMap map[string]map[string]interface{}, valuesName string) error
}

type GitRepo interface {
	InitRepo(branch string) error
	GetLastCommitTime() (time.Time, error)
}

type DBRepo interface {
	GetRepoByName(ctx context.Context, repoName string) (controller_dto.ConfigDBDto, error)
	ChangeRepoUpdateTime(ctx context.Context, repoName string) error
	GetNotLockedConfig(ctx context.Context) (entities.ControllerConfig, error)
	UnlockRepo(ctx context.Context, repoName string) error
}

type FileRepo interface {
	GetValuesFromFile(filePath string) (map[string]map[string]interface{}, error)
}

type ConfigControllerService struct {
	DBRepo   DBRepo
	GitRepo  GitRepo
	Executor Executor
	fileRepo FileRepo
	log      *zap.Logger
}

func NewConfigControllerService(
	DBRepo DBRepo,
	gitRepo GitRepo,
	fileRepo FileRepo,
	executor Executor,
	log *zap.Logger,
) *ConfigControllerService {
	return &ConfigControllerService{
		DBRepo:   DBRepo,
		GitRepo:  gitRepo,
		Executor: executor,
		fileRepo: fileRepo,
		log:      log,
	}
}

func (c *ConfigControllerService) Start(ctx context.Context, workInterval int) {
	ticker := time.NewTicker(time.Duration(workInterval) * time.Second)
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					c.log.Warn("Controller stopped")
					return

				case <-ticker.C:
					if err := c.Work(ctx); err != nil {
						c.log.Error("Controller has error", zap.Error(err))
					}
				}
			}
		}(i)

		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	return
}

func (c *ConfigControllerService) Work(ctx context.Context) error {
	conf, err := c.getNotLockedConfig(ctx)
	if err != nil {
		return err
	}

	if err = c.GitRepo.InitRepo(conf.ConfigBranch); err != nil {
		return err
	}

	commitTime, err := c.GitRepo.GetLastCommitTime()
	if err != nil {
		return err
	}

	//if commitTime.After(conf.LastUpdate) {
	if commitTime.After(conf.LastUpdate) {
		c.log.Info("Repo will be update: ", zap.String("repo", conf.ConfigName))

		configMap, err := c.fileRepo.GetValuesFromFile(conf.ConfigName)
		if err != nil {
			return err
		}

		if err = c.Executor.DoChangeWithNewValues(configMap, conf.ConfigName); err != nil {
			return err
		}

		if err := c.DBRepo.ChangeRepoUpdateTime(ctx, conf.ConfigName); err != nil {
			return err
		}
	}

	if err = c.DBRepo.UnlockRepo(ctx, conf.ConfigName); err != nil {
		return err
	}

	return nil
}

func (c *ConfigControllerService) getNotLockedConfig(ctx context.Context) (entities.ControllerConfig, error) {
	conf, err := c.DBRepo.GetNotLockedConfig(ctx)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNotLockedConfigNotFound) {
			c.log.Debug("not found not locked conf")
			return entities.ControllerConfig{}, err
		} else {
			c.log.Error("Failed to get not locked config", zap.Error(err))
			return entities.ControllerConfig{}, err
		}
	}

	return conf, nil
}

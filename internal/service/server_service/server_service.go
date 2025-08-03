package server_service

import (
	"context"
	"go-config-controller-svc/internal/custom_errors"
	"go-config-controller-svc/internal/entities"
	"go.uber.org/zap"
)

type ServerDBRepo interface {
	SaveConfig(ctx context.Context, config entities.ServerConfig) error
	GetConfigsList(ctx context.Context) ([]entities.ServerConfig, error)
	DeleteConfig(ctx context.Context, config entities.ServerConfig) error
	CheckIfUserExists(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, user entities.User) error
	GetUserByUsername(ctx context.Context, username string) (entities.User, error)
}

type ServerRedisRepo interface {
	AddTask(task string, queue string, ctx context.Context) error
}

type ServerService struct {
	DbRepo    ServerDBRepo
	redisRepo ServerRedisRepo
	log       *zap.Logger
}

func NewServerService(dbRepo ServerDBRepo, redisRepo ServerRedisRepo, log *zap.Logger) *ServerService {
	return &ServerService{DbRepo: dbRepo, redisRepo: redisRepo, log: log}
}

func (s *ServerService) CreateConfig(ctx context.Context, config entities.ServerConfig) error {
	if err := s.DbRepo.SaveConfig(ctx, config); err != nil {
		return err
	}

	return nil
}

func (s *ServerService) GetConfigsList(ctx context.Context) ([]entities.ServerConfig, error) {
	configs, err := s.DbRepo.GetConfigsList(ctx)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (s *ServerService) DeleteConfig(ctx context.Context, config entities.ServerConfig) error {
	if err := s.DbRepo.DeleteConfig(ctx, config); err != nil {
		return err
	}

	return nil
}

func (s *ServerService) GetUser(ctx context.Context, username string) (entities.User, error) {
	check, err := s.DbRepo.CheckIfUserExists(ctx, username)
	if err != nil {
		return entities.User{}, err
	}

	if !check {
		return entities.User{}, custom_errors.ErrLoginError
	}

	userFromDB, err := s.DbRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return entities.User{}, err
	}

	return userFromDB, nil
}

func (s *ServerService) CreateUser(ctx context.Context, user entities.User) error {
	if err := s.DbRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
func (s *ServerService) AddTaskForExecutor(ctx context.Context, task string) error {
	if task != entities.DoExecTask {
		return custom_errors.ErrWrongTask
	}

	if err := s.redisRepo.AddTask(task, entities.CommandQueue, ctx); err != nil {
		return err
	}

	return nil
}

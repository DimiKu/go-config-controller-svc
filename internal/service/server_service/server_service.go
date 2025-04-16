package server_service

import (
	"context"
	"go-config-controller-svc/internal/entities"
	"go.uber.org/zap"
)

type ServerDBRepo interface {
	SaveConfig(ctx context.Context, config entities.ServerConfig) error
	GetConfigsList(ctx context.Context) ([]entities.ServerConfig, error)
	DeleteConfig(ctx context.Context, config entities.ServerConfig) error
}

type ServerService struct {
	DbRepo ServerDBRepo
	log    *zap.Logger
}

func NewServerService(dbRepo ServerDBRepo, log *zap.Logger) *ServerService {
	return &ServerService{DbRepo: dbRepo, log: log}
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

package interfaces

import (
	"context"
	"go-config-controller-svc/internal/entities"
)

type ServerService interface {
	CreateConfig(ctx context.Context, config entities.ServerConfig) error
	GetConfigsList(ctx context.Context) ([]entities.ServerConfig, error)
	DeleteConfig(ctx context.Context, config entities.ServerConfig) error
}

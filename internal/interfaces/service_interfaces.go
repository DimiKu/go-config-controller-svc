package interfaces

import (
	"context"
	"go-config-controller-svc/internal/entities"
)

type ServerService interface {
	CreateConfig(ctx context.Context, config entities.ServerConfig) error
	GetConfigsList(ctx context.Context) ([]entities.ServerConfig, error)
	DeleteConfig(ctx context.Context, config entities.ServerConfig) error
	CreateUser(ctx context.Context, user entities.User) error
	GetUser(ctx context.Context, username string) (entities.User, error)
	AddTaskForExecutor(ctx context.Context, task string) error
}

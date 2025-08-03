package repos

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisRepo struct {
	rdb *redis.Client
	ctx context.Context
	log *zap.Logger
}

func NewRedisRepo(rdb *redis.Client, ctx context.Context, log *zap.Logger) *RedisRepo {
	return &RedisRepo{rdb: rdb, ctx: ctx, log: log}
}

func (repo *RedisRepo) Ping() error {
	_, err := repo.rdb.Ping(repo.ctx).Result()
	if err != nil {
		panic(err)
	}

	return nil
}

func (repo *RedisRepo) AddTask(task string, queue string, ctx context.Context) error {
	err := repo.rdb.LPush(ctx, queue, task).Err()
	if err != nil {
		repo.log.Error("redis LPush err: ", zap.Error(err))
		return err
	}

	return nil
}

func (repo *RedisRepo) GetTask(queue string, ctx context.Context) (string, error) {
	result, err := repo.rdb.BRPop(ctx, 0, queue).Result()
	if err != nil {
		repo.log.Error("redis LPush err: ", zap.Error(err))
		return "", err
	}

	repo.log.Info("Processed task:", zap.String("task", result[1]))

	return result[1], nil
}

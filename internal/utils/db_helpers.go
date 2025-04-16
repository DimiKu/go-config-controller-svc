package utils

import (
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

func CommitOrRollback(tx pgx.Tx, err error, ctx context.Context) {
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}

func isRetryable(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case pgerrcode.ConnectionException, pgerrcode.ConnectionDoesNotExist, pgerrcode.ConnectionFailure, pgerrcode.CannotConnectNow:
			return true
		}
	}
	return false
}

func RetryableQuery(ctx context.Context, pool *pgxpool.Pool, log *zap.Logger, query string, args ...interface{}) (pgx.Rows, error) {
	var row pgx.Rows
	var err error

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer CommitOrRollback(tx, err, ctx)

	for i := 0; i < 3; i++ {
		log.Debug("Im do query: %s, time: %d", zap.String("query", query), zap.Int("i", i))
		row, err = tx.Query(ctx, query, args...)
		if err != nil {
			if isRetryable(err) {
				log.Error("Retrying query due to error: %v\n", zap.Error(err))
				time.Sleep(time.Duration(i) * time.Second)
				log.Debug("Im do query: %s", zap.String("query", query))
				continue
			}

		}
		return row, nil
	}
	return nil, err
}

func RetryableExec(ctx context.Context, pool *pgxpool.Pool, log *zap.Logger, command string, args ...interface{}) (pgconn.CommandTag, error) {
	var tag pgconn.CommandTag
	var err error

	tx, err := pool.Begin(ctx)
	if err != nil {
		return tag, err
	}

	defer CommitOrRollback(tx, err, ctx)

	for i := 0; i < 3; i++ {
		tag, err = tx.Exec(ctx, command, args...)
		if err != nil {
			if isRetryable(err) {
				log.Warn("Retrying exec due to error: %v\n", zap.Error(err))
				time.Sleep(time.Duration(i) * time.Second)
				continue
			}
			return tag, err
		}
		return tag, nil
	}
	return tag, err
}

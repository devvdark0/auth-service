package db

import (
	"context"
	"time"

	"github.com/devvdark0/auth-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPOSTGRESQL(ctx context.Context, cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = int32(cfg.MaxConns)
	pool.Config().MinConns = int32(cfg.MinConns)
	pool.Config().MaxConnLifetime = time.Hour
	pool.Config().MaxConnIdleTime = 30 * time.Minute

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

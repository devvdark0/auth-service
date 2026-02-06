package db

import (
	"context"

	"github.com/devvdark0/auth-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg config.DBConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbPool, nil
}

package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func InitDb(ctx context.Context) *pgxpool.Pool {
	dsn := os.Getenv("DSN")
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("unable to connect to database")
	}
	if err = conn.Ping(ctx); err != nil {
		return nil
	}

	return conn
}

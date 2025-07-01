package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgxpool.Conn, error) {
	godotenv.Load()

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

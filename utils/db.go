package utils

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func ConnectDB() (*pgxpool.Conn, error) {
	once.Do(func() {
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			connStr = fmt.Sprintf(
				"postgres://%s:%s@%s:%s/%s?sslmode=disable",
				os.Getenv("PGUSER"),
				os.Getenv("PGPASSWORD"),
				os.Getenv("PGHOST"),
				os.Getenv("PGPORT"),
				os.Getenv("PGDATABASE"),
			)
		}

		config, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			panic(err)
		}

		config.MaxConns = 10
		config.MinConns = 2

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			panic(err)
		}
	})

	return pool.Acquire(context.Background())
}

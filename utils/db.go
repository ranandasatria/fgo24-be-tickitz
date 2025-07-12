package utils

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func ConnectDB() (*pgxpool.Conn, error) {
	once.Do(func() {
		godotenv.Load()
		connStr := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			os.Getenv("PGUSER"),
			os.Getenv("PGPASSWORD"),
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT"),
			os.Getenv("PGDATABASE"),
		)

		config, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			fmt.Println("Error parsing config:", err)
			return
		}

		config.MaxConns = 20
		config.MinConns = 2
		config.MaxConnLifetime = 30 * time.Minute

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			fmt.Println("Error creating pool:", err)
		}
	})

	if pool == nil {
		return nil, fmt.Errorf("database pool not initialized")
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

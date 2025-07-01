package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func RedisClient() redis.Client {
	godotenv.Load()
	addr := os.Getenv("RDADDRESS")
	pass := os.Getenv("RDPASSWORD")
	db, _ := strconv.Atoi(os.Getenv("RDDB"))
	var RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
	return *RedisClient
}

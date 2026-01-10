package utils

import (
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisOnce sync.Once

func RedisClient() redis.Client {
	redisOnce.Do(func() {
		db, _ := strconv.Atoi(os.Getenv("RDDB"))

		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("RDADDRESS"),
			Password: os.Getenv("RDPASSWORD"),
			DB:       db,
		})
	})
	return *redisClient
}

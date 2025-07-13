package utils

import (
	"context"
)

func DeleteKeysByPrefix(ctx context.Context, prefix string) {
	iter := RedisClient().Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		RedisClient().Del(ctx, iter.Val())
	}
}

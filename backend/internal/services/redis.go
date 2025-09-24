package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil
	}
	return redis.NewClient(opt)
}

var Ctx = context.Background()

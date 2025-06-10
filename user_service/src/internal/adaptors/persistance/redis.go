package persistance

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

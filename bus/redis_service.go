package bus

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "autorack.proxy.rlwy.net:51457",
			Username: "default",
			Password: "xkTittsmFXuykDsWHkXZGiShesVRZegL",
			DB:       0,
		})
	})
	return redisClient
}

func IsRedisConnected() (bool, error) {
	_, err := GetRedisClient().Ping(context.Background()).Result()
	return err == nil, err
}

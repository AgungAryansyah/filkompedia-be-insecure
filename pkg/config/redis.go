package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func loadRedisCredentials() (address string, password string) {
	return fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), os.Getenv("REDIS_PASS")
}

func StartRedis() *redis.Client {
	host, password := loadRedisCredentials()
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	return client
}

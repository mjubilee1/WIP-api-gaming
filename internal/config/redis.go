package config

import (
	"fmt"
	"api-gaming/internal/util"
	"github.com/go-redis/redis/v8"
)

var (
	redisURL = util.ViperEnvVariable("REDIS_URL")
	client *redis.Client
)

// InitRedis - Initializing redis
func InitRedis() {
	initURL, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Println("Error parsing redis URL.")
		panic(err)
	} 
	client = redis.NewClient(initURL)
}

func RedisConn() (*redis.Client) {
	return client
}
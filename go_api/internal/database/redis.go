package database

import "github.com/go-redis/redis/v8"

var RedisDB *redis.Client

func ConnectRedis(redisAddr string, pass string) {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: pass,
	})
}

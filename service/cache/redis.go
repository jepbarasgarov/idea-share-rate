package cache

import "github.com/go-redis/redis/v8"

type RedisService struct {
	Service
	client *redis.Client
}

func NewRedisService(redisConn string, redisDb int, maxIdle int, maxActive int) (r *RedisService) {

	client := redis.NewClient(&redis.Options{
		Addr:         redisConn,
		DB:           redisDb,
		PoolSize:     maxActive,
		MinIdleConns: maxIdle,
	})

	return &RedisService{
		client: client,
	}
}

package initdatabase

import (
	"context"
	"github.com/kawasoki/gowork/configs"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient(redisCfg *configs.RedisCfg) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Addr,
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("redis connect ping failed, err:%s", err.Error())
	} else {
		log.Println("redis connect ping response:" + pong)
	}
	return client
}

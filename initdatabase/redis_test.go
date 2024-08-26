package initdatabase

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	GetDefaultClient()
	GetDLocks()
}

var once sync.Once
var defaultClient *redis.Client

func GetDefaultClient() *redis.Client {
	once.Do(func() {
		defaultClient = NewRedisClient(&RedisCfg{})
	})
	return defaultClient
}

var rs *redsync.Redsync
var once2 sync.Once

func GetDLocks() *redsync.Redsync {
	once2.Do(func() {
		client := NewRedisClient(&RedisCfg{})
		pool := goredis.NewPool(client)
		rs = redsync.New(pool)
		log.Println("初始化redis分布式锁")
	})
	return rs
}

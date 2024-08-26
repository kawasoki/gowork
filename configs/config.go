package configs

import "time"

const (
	DefaultHttpTimeout = 3 * time.Second //http默认的超时时间
)

const MaxPageSize = 1000

type Config struct {
	RedisCfg
	EsCfg
}

type RedisCfg struct {
	Addr         string //host:port address.
	Password     string
	DB           int
	PoolSize     int // Maximum number of socket connections.  Default is 10 connections per every CPU as reported by runtime.NumCPU.
	MinIdleConns int // Minimum number of idle connections which is useful when establishing new connection is slow.
}

type EsCfg struct {
	EsAddress  string //host:port address.
	EsUserName string
	EsPassword string
}

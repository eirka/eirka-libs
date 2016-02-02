package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

// RedisStore holds a handle to the Redis pool
type RedisStore struct {
	Pool  *redis.Pool
	Mutex *Mutex
}

var (
	RedisCache   RedisStore
	ErrCacheMiss = errors.New("cache: key not found.")
)

type Redis struct {
	// Redis address and max pool connections
	Protocol       string
	Address        string
	MaxIdle        int
	MaxConnections int
}

// NewRedisCache creates a new pool
func (r *Redis) NewRedisCache() {

	RedisCache.Pool = &redis.Pool{
		MaxIdle:     r.MaxIdle,
		MaxActive:   r.MaxConnections,
		IdleTimeout: 240 * time.Second,
		Dial: func() (c redis.Conn, err error) {
			c, err = redis.Dial(r.Protocol, r.Address)
			if err != nil {
				return
			}
			return
		},
	}

	// create our distributed lock
	RedisCache.Mutex = NewMutex([]Pool{
		RedisCache.Pool,
	})

	return
}

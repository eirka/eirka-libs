package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/hjr265/redsync.go/redsync"
	"time"
)

// RedisStore holds a handle to the Redis pool
type RedisStore struct {
	Pool  *redis.Pool
	Mutex *redsync.Mutex
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
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(r.Protocol, r.Address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	// create our distributed lock
	RedisCache.Mutex = redsync.NewMutexWithGenericPool("post_lock", []redsync.Pool{
		RedisCache.Pool,
	})

	return
}

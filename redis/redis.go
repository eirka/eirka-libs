package redis

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

var cacheInitialized atomic.Uint32

// SetCacheInitialized marks the cache as initialized
func SetCacheInitialized() {
	cacheInitialized.Store(1)
}

// isCacheInitialized checks if the cache has been initialized
func isCacheInitialized() bool {
	return cacheInitialized.Load() == 1
}

// ErrCacheNotInitialized is returned when Cache is not yet initialized
var ErrCacheNotInitialized = errors.New("redis cache not initialized")

// Pool is a generic connection pool
type Pool interface {
	Get() redis.Conn
	Close() error
}

var _ = Pool(&redis.Pool{})

// Store holds a handle to the Redis pool
type Store struct {
	Pool  Pool
	Mutex *Mutex
	Mock  *redigomock.Conn
}

var (
	// Cache holds a store
	Cache Store
	// ErrCacheMiss is an error for cache misses
	ErrCacheMiss = errors.New("cache: key not found")
)

// Redis holds connection options for redis
type Redis struct {
	// Redis address and max pool connections
	Protocol       string
	Address        string
	MaxIdle        int
	MaxConnections int
}

// NewRedisCache creates a new pool
func (r *Redis) NewRedisCache() {

	Cache.Pool = &redis.Pool{
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
	Cache.Mutex = NewMutex([]Pool{
		Cache.Pool,
	})

	SetCacheInitialized()
}

// NewRedisMock returns a fake redis pool for testing
func NewRedisMock() {

	Cache.Mock = redigomock.NewConn()

	Cache.Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return Cache.Mock, nil
		},
	}

	// create our distributed lock
	Cache.Mutex = NewMutex([]Pool{
		Cache.Pool,
	})

	SetCacheInitialized()
}

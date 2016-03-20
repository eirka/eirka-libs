package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"time"
)

// Pool is a generic connection pool
type Pool interface {
	Get() redis.Conn
	Close() error
}

var _ = Pool(&redis.Pool{})
var _ = Pool(&RedisPoolMock{})

// RedisStore holds a handle to the Redis pool
type RedisStore struct {
	Pool  Pool
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

type RedisPoolMock struct {
	Conn *redigomock.Conn
}

func (r *RedisPoolMock) Get() redis.Conn {
	return r.Conn
}

func (r *RedisPoolMock) GetMock() redigomock.Conn {
	return r.Conn
}

func (r *RedisPoolMock) Close() error {
	return nil
}

// NewRedisMock returns a fake redis pool for testing
func NewRedisMock() {

	RedisCache.Pool = &RedisPoolMock{
		Conn: redigomock.NewConn(),
	}

	// create our distributed lock
	RedisCache.Mutex = NewMutex([]Pool{
		RedisCache.Pool,
	})

	return
}

package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

// Storer defines custom methods for redis operations
type Storer interface {
	Lock(key string) error
	Unlock(key string) bool
	Get(key string) (result []byte, err error)
	HGet(key string, value string) (result []byte, err error)
	Set(key string, result []byte) (err error)
	SetEx(key string, timeout uint, result []byte) (err error)
	HMSet(key string, value string, result []byte) (err error)
	Delete(key ...interface{}) (err error)
	Flush() (err error)
	Incr(key string) (result int, err error)
	Expire(key string, timeout uint) (err error)
}

var _ = Storer(&Store{})

// Lock our shared mutex
func (c *Store) Lock(key string) error {
	return c.Mutex.Lock(key)
}

// Unlock our shared mutex
func (c *Store) Unlock(key string) bool {
	return c.Mutex.Unlock(key)
}

// Get will retrieve a key
func (c *Store) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	if !isCacheInitialized() {
		return nil, ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	result, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

// HGet will retrieve a hash
func (c *Store) HGet(key string, value string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	if value == "" {
		return nil, errors.New("value cannot be empty")
	}

	if !isCacheInitialized() {
		return nil, ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	result, err := redis.Bytes(conn.Do("HGET", key, value))
	if err == redis.ErrNil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Set will set a single record
func (c *Store) Set(key string, result []byte) (err error) {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, result)

	return
}

// SetEx will set a single record with an expiration
func (c *Store) SetEx(key string, timeout uint, result []byte) (err error) {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	if timeout == 0 {
		return errors.New("timeout must be greater than 0")
	}

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, timeout, result)

	return
}

// HMSet will set a hash
func (c *Store) HMSet(key string, value string, result []byte) (err error) {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	if value == "" {
		return errors.New("value cannot be empty")
	}

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("HMSET", key, value, result)

	return
}

// Delete will delete a key
func (c *Store) Delete(key ...interface{}) (err error) {
	if len(key) == 0 {
		return errors.New("at least one key must be provided")
	}

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", key...)

	return
}

// Flush will call flushall and delete all keys
func (c *Store) Flush() (err error) {

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("FLUSHALL")

	return
}

// Incr will increment a redis key
func (c *Store) Incr(key string) (result int, err error) {
	if key == "" {
		return 0, errors.New("key cannot be empty")
	}

	if !isCacheInitialized() {
		return 0, ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", key))
}

// Expire will set expire on a redis key
func (c *Store) Expire(key string, timeout uint) (err error) {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	if timeout == 0 {
		return errors.New("timeout must be greater than 0")
	}

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("EXPIRE", key, timeout)

	return
}

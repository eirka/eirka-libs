package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// lock our shared mutex
func (c *RedisStore) Lock(key string) error {
	return c.Mutex.Lock(key)
}

// unlock our shared mutex
func (c *RedisStore) Unlock(key string) bool {
	return c.Mutex.Unlock(key)
}

// Get will retrieve a key
func (c *RedisStore) Get(key string) (result []byte, err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", key))

}

// HGet will retrieve a hash
func (c *RedisStore) HGet(key string, value string) (result []byte, err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("HGET", key, value))

}

// Set will set a single record
func (c *RedisStore) Set(key string, result []byte) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, result)

	return
}

// Set will set a single record
func (c *RedisStore) SetEx(key string, timeout uint, result []byte) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, timeout, result)

	return
}

// HMSet will set a hash
func (c *RedisStore) HMSet(key string, value string, result []byte) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("HMSET", key, value, result)

	return
}

// Delete will delete a key
func (c *RedisStore) Delete(key ...interface{}) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", key...)

	return
}

// Flush will call flushall and delete all keys
func (c *RedisStore) Flush() (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("FLUSHALL")

	return
}

// will increment a redis key
func (c *RedisStore) Incr(key string) (result int, err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", key))
}

// will set expire on a redis key
func (c *RedisStore) Expire(key string, timeout uint) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("EXPIRE", key, timeout)

	return
}

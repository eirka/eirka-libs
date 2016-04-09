package redis

import (
	"github.com/garyburd/redigo/redis"
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
func (c *Store) Get(key string) (result []byte, err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", key))

}

// HGet will retrieve a hash
func (c *Store) HGet(key string, value string) (result []byte, err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("HGET", key, value))

}

// Set will set a single record
func (c *Store) Set(key string, result []byte) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, result)

	return
}

// SetEx will set a single record with an expiration
func (c *Store) SetEx(key string, timeout uint, result []byte) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, timeout, result)

	return
}

// HMSet will set a hash
func (c *Store) HMSet(key string, value string, result []byte) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("HMSET", key, value, result)

	return
}

// Delete will delete a key
func (c *Store) Delete(key ...interface{}) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", key...)

	return
}

// Flush will call flushall and delete all keys
func (c *Store) Flush() (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("FLUSHALL")

	return
}

// Incr will increment a redis key
func (c *Store) Incr(key string) (result int, err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", key))
}

// Expire will set expire on a redis key
func (c *Store) Expire(key string, timeout uint) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("EXPIRE", key, timeout)

	return
}

package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stvp/tempredis"
)

func TestNewRedisCache(t *testing.T) {

	server, err := tempredis.Start(tempredis.Config{})
	if err != nil {
		panic(err)
	}

	config := Redis{
		Protocol:       "unix",
		Address:        server.Socket(),
		MaxIdle:        1,
		MaxConnections: 5,
	}

	config.NewRedisCache()

	conn := Cache.Pool.Get()

	_, err = conn.Do("PING")

	assert.NoError(t, err, "An error was not expected")

	conn.Close()

	server.Term()

}

// Test cache initialization status
func TestCacheInitialization(t *testing.T) {
	// Reset the initialization flag
	cacheInitialized.Store(0)

	// Check if cache is initialized (should be false)
	assert.False(t, isCacheInitialized(), "Cache should not be initialized")

	// Mark cache as initialized
	SetCacheInitialized()

	// Check if cache is initialized (should be true)
	assert.True(t, isCacheInitialized(), "Cache should be initialized")
}

// Test Redis mocking
func TestNewRedisMock(t *testing.T) {
	// Reset the initialization flag
	cacheInitialized.Store(0)

	// Create a new Redis mock
	NewRedisMock()

	// Check if cache is initialized
	assert.True(t, isCacheInitialized(), "Cache should be initialized after mocking")

	// Verify mock is set up
	assert.NotNil(t, Cache.Mock, "Mock should be initialized")
	assert.NotNil(t, Cache.Pool, "Pool should be initialized")
	assert.NotNil(t, Cache.Mutex, "Mutex should be initialized")

	// Test connection from pool
	conn := Cache.Pool.Get()
	assert.NotNil(t, conn, "Should get a valid connection from pool")

	// Cleanup
	conn.Close()
}

// Test error handling when cache is not initialized
func TestUninitializedCache(t *testing.T) {
	// Store the original state
	originalCache := Cache

	// Reset cache to test uninitialized state
	Cache = Store{}
	cacheInitialized.Store(0)

	// Attempt operations on uninitialized cache
	_, err := Cache.Get("test")
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	_, err = Cache.HGet("test", "field")
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	err = Cache.Set("test", []byte("data"))
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	err = Cache.SetEx("test", 10, []byte("data"))
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	err = Cache.HMSet("test", "field", []byte("data"))
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	err = Cache.Delete("test")
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	err = Cache.Flush()
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	_, err = Cache.Incr("test")
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	err = Cache.Expire("test", 10)
	assert.Equal(t, ErrCacheNotInitialized, err, "Should return 'cache not initialized' error")

	// Restore original state
	Cache = originalCache
}

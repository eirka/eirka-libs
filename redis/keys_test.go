package redis

import (
	"errors"
	"testing"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stvp/tempredis"
)

func TestNewKey(t *testing.T) {

	key := NewKey("image")

	assert.NotEmpty(t, key, "Should not be empty")

	empty := NewKey("blah")

	assert.Empty(t, empty, "Should be empty")

}

func TestSetKeyHash(t *testing.T) {

	key := NewKey("image")

	key = key.SetKey("1", "1")

	assert.NotEmpty(t, key, "Should not be empty")

	assert.True(t, key.keyset, "Key should be set")

	assert.Equal(t, key.key, "image:1", "Key should match")

	assert.Equal(t, key.hashid, "1", "Should have a hash id")

	assert.Equal(t, key.String(), "image:1", "Key should match")

}

func TestSetKeyNoHash(t *testing.T) {

	key := NewKey("new")

	key = key.SetKey("1")

	assert.NotEmpty(t, key, "Should not be empty")

	assert.True(t, key.keyset, "Key should be set")

	assert.Equal(t, key.key, "new:1", "Key should match")

	assert.Empty(t, key.hashid, "Should be empty")

	assert.Equal(t, key.String(), "new:1", "Key should match")

}

func TestSetKeyNoKey(t *testing.T) {

	key := NewKey("tagtypes")

	key = key.SetKey()

	assert.NotEmpty(t, key, "Should not be empty")

	assert.True(t, key.keyset, "Key should be set")

	assert.Equal(t, key.key, "tagtypes", "Key should match")

	assert.Empty(t, key.hashid, "Should be empty")

	assert.Equal(t, key.String(), "tagtypes", "Key should match")

}

func TestSetKeyHashTooManyFields(t *testing.T) {

	key := NewKey("image")

	key = key.SetKey("1", "1", "1")

	assert.NotEmpty(t, key, "Should not be empty")

	assert.False(t, key.keyset, "Key should not be set")

	assert.Empty(t, key.key, "There should not be a key")

	assert.Empty(t, key.hashid, "Should not have a hash id")

	assert.Empty(t, key.String(), "There should not be a key")

}

func TestSetKeyTooManyFields(t *testing.T) {

	key := NewKey("new")

	key = key.SetKey("1", "1", "1")

	assert.NotEmpty(t, key, "Should not be empty")

	assert.False(t, key.keyset, "Key should not be set")

	assert.Empty(t, key.key, "There should not be a key")

	assert.Empty(t, key.hashid, "Should not have a hash id")

	assert.Empty(t, key.String(), "There should not be a key")

}

func TestSetKeyNotEnoughFields(t *testing.T) {

	key := NewKey("thread")

	key = key.SetKey("1")

	assert.NotEmpty(t, key, "Should not be empty")

	assert.False(t, key.keyset, "Key should not be set")

	assert.Empty(t, key.key, "There should not be a key")

	assert.Empty(t, key.hashid, "Should not have a hash id")

	assert.Empty(t, key.String(), "There should not be a key")

}

func TestSetKeyHashNotEnoughFields(t *testing.T) {

	key := NewKey("index")

	key = key.SetKey("1")

	assert.NotEmpty(t, key, "Should not be empty")

	assert.False(t, key.keyset, "Key should not be set")

	assert.Empty(t, key.key, "There should not be a key")

	assert.Empty(t, key.hashid, "Should not have a hash id")

	assert.Empty(t, key.String(), "There should not be a key")

}

func TestKeysGet(t *testing.T) {

	key := NewKey("new")

	key = key.SetKey("1")

	NewRedisMock()

	Cache.Mock.Command("GET", "new:1").Expect("worked!")

	res, err := key.Get()

	assert.NotEmpty(t, res, "Should return data")

	assert.NoError(t, err, "An error was not expected")
}

func TestKeysGetHash(t *testing.T) {

	key := NewKey("thread")

	key = key.SetKey("1", "1", "1")

	NewRedisMock()

	Cache.Mock.Command("HGET", "thread:1:1", "1").Expect("worked!")

	res, err := key.Get()

	assert.NotEmpty(t, res, "Should return data")

	assert.NoError(t, err, "An error was not expected")
}

func TestKeysGetKeyNotSet(t *testing.T) {

	key := NewKey("new")

	key = key.SetKey("1", "1")

	res, err := key.Get()

	assert.Empty(t, res, "Should not return data")

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrKeyNotSet, "Error should be the same")
	}
}

// Test cache not initialized for Get
func TestKeysGetCacheNotInitialized(t *testing.T) {
	// Store the original state
	originalCache := Cache
	originalInitialized := cacheInitialized.Load()

	// Reset cache to test uninitialized state
	Cache = Store{}
	cacheInitialized.Store(0)

	key := NewKey("new")
	key = key.SetKey("1")
	key.keyset = true

	res, err := key.Get()

	assert.Empty(t, res, "Should not return data when cache not initialized")
	assert.Equal(t, ErrCacheNotInitialized, err, "Error should be cache not initialized")

	// Restore original state
	Cache = originalCache
	cacheInitialized.Store(originalInitialized)
}

// Test error from Get operation
func TestKeysGetError(t *testing.T) {
	key := NewKey("new")
	key = key.SetKey("1")

	NewRedisMock()

	Cache.Mock.Command("GET", "new:1").ExpectError(errors.New("redis error"))

	res, err := key.Get()

	assert.Empty(t, res, "Should not return data on error")
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "redis error", err.Error(), "Error message should match")
}

// Test cache miss from Get operation
func TestKeysGetCacheMiss(t *testing.T) {
	key := NewKey("new")
	key = key.SetKey("1")

	NewRedisMock()

	Cache.Mock.Command("GET", "new:1").Expect(nil)

	res, err := key.Get()

	assert.Empty(t, res, "Should not return data on cache miss")
	assert.Equal(t, ErrCacheMiss, err, "Error should be cache miss")
}

func TestKeysSet(t *testing.T) {

	key := NewKey("tagtypes")

	key = key.SetKey()

	NewRedisMock()

	Cache.Mock.Command("SET", "tagtypes", []byte("hello"))

	err := key.Set([]byte("hello"))

	assert.NoError(t, err, "An error was not expected")
}

func TestKeysSetExpire(t *testing.T) {

	key := NewKey("new")

	key = key.SetKey("1")

	NewRedisMock()

	Cache.Mock.Command("SET", "new:1", []byte("hello"))
	Cache.Mock.Command("EXPIRE", "new:1", redigomock.NewAnyData())

	err := key.Set([]byte("hello"))

	assert.NoError(t, err, "An error was not expected")
}

func TestKeysSetHash(t *testing.T) {

	key := NewKey("index")

	key = key.SetKey("1", "1")

	NewRedisMock()

	Cache.Mock.Command("HMSET", "index:1", "1", []byte("hello"))

	err := key.Set([]byte("hello"))

	assert.NoError(t, err, "An error was not expected")
}

func TestKeysSetError(t *testing.T) {

	key := NewKey("index")

	key = key.SetKey("1", "1")

	NewRedisMock()

	Cache.Mock.Command("HMSET", "index:1", "1", []byte("hello")).ExpectError(errors.New("an error"))

	err := key.Set([]byte("hello"))

	assert.Error(t, err, "An error was expected")
}

func TestKeysSetKeyNotSet(t *testing.T) {

	key := NewKey("index")

	key = key.SetKey()

	err := key.Set([]byte("hello"))

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrKeyNotSet, "Error should be the same")
	}
}

// Test cache not initialized for Set
func TestKeysSetCacheNotInitialized(t *testing.T) {
	// Store the original state
	originalCache := Cache
	originalInitialized := cacheInitialized.Load()

	// Reset cache to test uninitialized state
	Cache = Store{}
	cacheInitialized.Store(0)

	key := NewKey("new")
	key = key.SetKey("1")
	key.keyset = true

	err := key.Set([]byte("hello"))

	assert.Equal(t, ErrCacheNotInitialized, err, "Error should be cache not initialized")

	// Restore original state
	Cache = originalCache
	cacheInitialized.Store(originalInitialized)
}

// Test Set with expire error
func TestKeysSetExpireError(t *testing.T) {
	key := NewKey("new")
	key = key.SetKey("1")

	NewRedisMock()

	Cache.Mock.Command("SET", "new:1", []byte("hello"))
	Cache.Mock.Command("EXPIRE", "new:1", redigomock.NewAnyData()).ExpectError(errors.New("expire error"))

	err := key.Set([]byte("hello"))

	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "expire error", err.Error(), "Error should be from EXPIRE command")
}

// Test Set with lock key
func TestKeysSetWithLock(t *testing.T) {
	key := NewKey("index")
	key = key.SetKey("1", "1")

	NewRedisMock()

	Cache.Mock.Command("HMSET", "index:1", "1", []byte("hello"))
	// Expect a call to DEL for unlocking the mutex
	Cache.Mock.Command("DEL", "index:1:mutex")

	err := key.Set([]byte("hello"))

	assert.NoError(t, err, "An error was not expected")
}

func TestKeysDelete(t *testing.T) {

	key := NewKey("thread")

	key = key.SetKey("1", "1", "1")

	NewRedisMock()

	Cache.Mock.Command("DEL", "thread:1:1")

	err := key.Delete()

	assert.NoError(t, err, "An error was not expected")
}

func TestKeysDeleteError(t *testing.T) {

	key := NewKey("thread")

	key = key.SetKey("1", "1", "1")

	NewRedisMock()

	Cache.Mock.Command("DEL", "thread:1:1").ExpectError(errors.New("an error"))

	err := key.Delete()

	assert.Error(t, err, "An error was expected")
}

func TestKeysDeleteLock(t *testing.T) {

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

	key := NewKey("index")

	key = key.SetKey("1", "1")

	err = key.Delete()

	assert.NoError(t, err, "An error was not expected")

	server.Term()
}

// Test Delete when Key is not set
func TestKeysDeleteKeyNotSet(t *testing.T) {
	key := NewKey("index")

	err := key.Delete()

	assert.Error(t, err, "An error was expected")
	assert.Equal(t, ErrKeyNotSet, err, "Error should be key not set")
}

// Test cache not initialized for Delete
func TestKeysDeleteCacheNotInitialized(t *testing.T) {
	// Store the original state
	originalCache := Cache
	originalInitialized := cacheInitialized.Load()

	// Reset cache to test uninitialized state
	Cache = Store{}
	cacheInitialized.Store(0)

	key := NewKey("new")
	key = key.SetKey("1")
	key.keyset = true

	err := key.Delete()

	assert.Equal(t, ErrCacheNotInitialized, err, "Error should be cache not initialized")

	// Restore original state
	Cache = originalCache
	cacheInitialized.Store(originalInitialized)
}

// Test Delete with Lock error
func TestKeysDeleteLockError(t *testing.T) {
	key := NewKey("index")
	key = key.SetKey("1", "1")

	NewRedisMock()

	Cache.Mock.Command("DEL", "index:1")
	// The Mutex.Lock function will return ErrFailed

	err := key.Delete()

	// Since the locking mechanism is complex with multiple retries and timeouts,
	// we can simply verify an error occurs rather than the specific error
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, ErrFailed, err, "Error should be ErrFailed")
}

// Test all RedisKeys defined values
func TestRedisKeysDefinitions(t *testing.T) {
	// The test will simply verify that all predefined keys initialize correctly
	// and can be used to generate keys

	// Define expected results keyed by base name
	expectedResults := map[string]struct {
		fieldCount int
		hash       bool
		expire     bool
		lock       bool
	}{
		"index":       {1, true, false, true},
		"thread":      {2, true, false, false},
		"tag":         {2, true, true, false},
		"image":       {1, true, false, false},
		"post":        {2, true, false, false},
		"tags":        {1, true, false, false},
		"directory":   {1, true, false, false},
		"new":         {1, false, true, false},
		"popular":     {1, false, true, false},
		"favorited":   {1, false, true, false},
		"tagtypes":    {0, false, false, false},
		"imageboards": {0, false, true, false},
	}

	// Verify each key definition matches the expected configuration
	for base, expected := range expectedResults {
		key := NewKey(base)

		assert.NotNil(t, key, "Key should be defined: "+base)
		if key != nil {
			assert.Equal(t, base, key.base, "Base name should match for "+base)
			assert.Equal(t, expected.fieldCount, key.fieldcount, "Field count should match for "+base)
			assert.Equal(t, expected.hash, key.hash, "Hash flag should match for "+base)
			assert.Equal(t, expected.expire, key.expire, "Expire flag should match for "+base)
			assert.Equal(t, expected.lock, key.lock, "Lock flag should match for "+base)
		}
	}

	// Verify RedisKeyIndex contains all keys
	assert.Equal(t, len(expectedResults), len(RedisKeyIndex), "RedisKeyIndex should contain all defined keys")
}

// Test for String on empty keys
func TestStringEmptyKey(t *testing.T) {
	key := &Key{base: "test"}
	assert.Empty(t, key.String(), "String() should return empty string when key is not set")
}

// Test for setting multiple fields with delimiters
func TestSetKeyMultipleFields(t *testing.T) {
	// Create a key with multiple fields (thread has fieldcount 2)
	key := NewKey("thread")

	// Use fields with potential delimiters
	key = key.SetKey("123", "456", "789")

	assert.Equal(t, "thread:123:456", key.key, "Key should properly join fields with delimiters")
	assert.Equal(t, "789", key.hashid, "Hash ID should be set correctly")
}

// Test for verifying the initializations in init function
func TestInitFunction(t *testing.T) {
	// This test verifies that the init function properly maps all keys
	// We check a few keys to ensure they were properly mapped

	// Check for some specific keys
	keys := []string{"index", "thread", "tag", "image", "tagtypes"}

	for _, key := range keys {
		mappedKey, exists := RedisKeyIndex[key]
		assert.True(t, exists, "Key should exist in RedisKeyIndex: "+key)
		assert.Equal(t, key, mappedKey.base, "Base name should match for mapped key: "+key)
	}
}

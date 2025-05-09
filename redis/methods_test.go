package redis

import (
	"errors"
	"testing"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func TestMethodGet(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("GET", "index:1").Expect("worked!")

	res, err := Cache.Get("index:1")

	assert.NotEmpty(t, res, "Should return data")

	assert.NoError(t, err, "An error was not expected")

	Cache.Mock.Command("GET", "index:1").Expect(nil)

	empty, err := Cache.Get("index:1")

	assert.Empty(t, empty, "Should not return data")

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrCacheMiss, "Error should be the same")
	}

	Cache.Mock.Command("GET", "index:1").ExpectError(errors.New("oh shit"))

	bad, err := Cache.Get("index:1")

	assert.Empty(t, bad, "Should not return data")

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, errors.New("oh shit"), "Error should be the same")
	}

}

// Test invalid key for Get
func TestMethodGetInvalidKey(t *testing.T) {
	NewRedisMock()

	// Test with empty key
	empty, err := Cache.Get("")
	assert.Empty(t, empty, "Should not return data for empty key")
	assert.Error(t, err, "An error was expected for empty key")
	assert.Equal(t, "key cannot be empty", err.Error(), "Error should be for empty key")
}

func TestMethodHGet(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("HGET", "index:1", "1").Expect("worked!")

	res, err := Cache.HGet("index:1", "1")

	assert.NotEmpty(t, res, "Should return data")

	assert.NoError(t, err, "An error was not expected")

	Cache.Mock.Command("HGET", "index:1", "1").Expect(nil)

	empty, err := Cache.HGet("index:1", "1")

	assert.Empty(t, empty, "Should not return data")

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, ErrCacheMiss, "Error should be the same")
	}

	Cache.Mock.Command("HGET", "index:1", "1").ExpectError(errors.New("oh shit"))

	bad, err := Cache.HGet("index:1", "1")

	assert.Empty(t, bad, "Should not return data")

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, errors.New("oh shit"), "Error should be the same")
	}

}

// Test invalid key or value for HGet
func TestMethodHGetInvalidParams(t *testing.T) {
	NewRedisMock()

	// Test with empty key
	empty, err := Cache.HGet("", "field")
	assert.Empty(t, empty, "Should not return data for empty key")
	assert.Error(t, err, "An error was expected for empty key")
	assert.Equal(t, "key cannot be empty", err.Error(), "Error should be for empty key")

	// Test with empty value
	empty, err = Cache.HGet("key", "")
	assert.Empty(t, empty, "Should not return data for empty value")
	assert.Error(t, err, "An error was expected for empty value")
	assert.Equal(t, "value cannot be empty", err.Error(), "Error should be for empty value")
}

// Test Set method
func TestMethodSet(t *testing.T) {
	NewRedisMock()

	Cache.Mock.Command("SET", "index:1", []byte("hello"))

	err := Cache.Set("index:1", []byte("hello"))

	assert.NoError(t, err, "An error was not expected")

	// Test with invalid key
	err = Cache.Set("", []byte("hello"))
	assert.Error(t, err, "An error was expected for empty key")
	assert.Equal(t, "key cannot be empty", err.Error(), "Error should be for empty key")

	// Test with SET error
	Cache.Mock.Command("SET", "index:1", []byte("hello")).ExpectError(errors.New("connection error"))
	err = Cache.Set("index:1", []byte("hello"))
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "connection error", err.Error(), "Error should match expected error")
}

func TestMethodSetEx(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("SETEX", "index:1", redigomock.NewAnyData(), redigomock.NewAnyData())

	err := Cache.SetEx("index:1", 600, []byte("hello"))

	assert.NoError(t, err, "An error was not expected")

	// Test with invalid key
	err = Cache.SetEx("", 600, []byte("hello"))
	assert.Error(t, err, "An error was expected for empty key")
	assert.Equal(t, "key cannot be empty", err.Error(), "Error should be for empty key")

	// Test with invalid timeout
	err = Cache.SetEx("index:1", 0, []byte("hello"))
	assert.Error(t, err, "An error was expected for zero timeout")
	assert.Equal(t, "timeout must be greater than 0", err.Error(), "Error should be for invalid timeout")

	// Test with SETEX error
	Cache.Mock.Command("SETEX", "index:1", redigomock.NewAnyData(), redigomock.NewAnyData()).ExpectError(errors.New("connection error"))
	err = Cache.SetEx("index:1", 600, []byte("hello"))
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "connection error", err.Error(), "Error should match expected error")
}

func TestMethodHMSet(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("HMSET", "index:1", "1", redigomock.NewAnyData())

	err := Cache.HMSet("index:1", "1", []byte("hello"))

	assert.NoError(t, err, "An error was not expected")

	// Test with invalid key
	err = Cache.HMSet("", "1", []byte("hello"))
	assert.Error(t, err, "An error was expected for empty key")
	assert.Equal(t, "key cannot be empty", err.Error(), "Error should be for empty key")

	// Test with invalid value
	err = Cache.HMSet("index:1", "", []byte("hello"))
	assert.Error(t, err, "An error was expected for empty value")
	assert.Equal(t, "value cannot be empty", err.Error(), "Error should be for empty value")

	// Test with HMSET error
	Cache.Mock.Command("HMSET", "index:1", "1", redigomock.NewAnyData()).ExpectError(errors.New("connection error"))
	err = Cache.HMSet("index:1", "1", []byte("hello"))
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "connection error", err.Error(), "Error should match expected error")
}

func TestMethodDelete(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("DEL", "index:1", "thread:2")

	err := Cache.Delete("index:1", "thread:2")

	assert.NoError(t, err, "An error was not expected")

	// Test with no keys
	err = Cache.Delete()
	assert.Error(t, err, "An error was expected for no keys")
	assert.Equal(t, "at least one key must be provided", err.Error(), "Error should be for no keys")

	// Test with DEL error
	Cache.Mock.Command("DEL", "index:1").ExpectError(errors.New("connection error"))
	err = Cache.Delete("index:1")
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "connection error", err.Error(), "Error should match expected error")
}

func TestMethodFlush(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("FLUSHALL")

	err := Cache.Flush()

	assert.NoError(t, err, "An error was not expected")

	// Test with FLUSHALL error
	Cache.Mock.Command("FLUSHALL").ExpectError(errors.New("connection error"))
	err = Cache.Flush()
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "connection error", err.Error(), "Error should match expected error")
}

func TestMethodIncr(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("INCR", "login:2").Expect([]byte("2"))

	res, err := Cache.Incr("login:2")

	assert.NotEmpty(t, res, "Should return data")

	assert.NoError(t, err, "An error was not expected")

	// Test with empty key
	res, err = Cache.Incr("")
	assert.Equal(t, 0, res, "Result should be 0 for empty key")
	assert.Error(t, err, "An error was expected for empty key")
	assert.Equal(t, "key cannot be empty", err.Error(), "Error should be for empty key")

	// Test with INCR error
	Cache.Mock.Command("INCR", "login:2").ExpectError(errors.New("connection error"))
	res, err = Cache.Incr("login:2")
	assert.Equal(t, 0, res, "Result should be 0 for error")
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "connection error", err.Error(), "Error should match expected error")
}

func TestMethodExpire(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("EXPIRE", "new:1", redigomock.NewAnyData())

	err := Cache.Expire("new:1", 600)

	assert.NoError(t, err, "An error was not expected")

	// Test with empty key
	err = Cache.Expire("", 600)
	assert.Error(t, err, "An error was expected for empty key")
	assert.Equal(t, "key cannot be empty", err.Error(), "Error should be for empty key")

	// Test with invalid timeout
	err = Cache.Expire("new:1", 0)
	assert.Error(t, err, "An error was expected for zero timeout")
	assert.Equal(t, "timeout must be greater than 0", err.Error(), "Error should be for invalid timeout")

	// Test with EXPIRE error
	Cache.Mock.Command("EXPIRE", "new:1", redigomock.NewAnyData()).ExpectError(errors.New("connection error"))
	err = Cache.Expire("new:1", 600)
	assert.Error(t, err, "An error was expected")
	assert.Equal(t, "connection error", err.Error(), "Error should match expected error")
}
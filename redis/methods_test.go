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

func TestMethodSetEx(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("SETEX", "index:1", redigomock.NewAnyData(), redigomock.NewAnyData())

	err := Cache.SetEx("index:1", 600, []byte("hello"))

	assert.NoError(t, err, "An error was not expected")

}

func TestMethodHMSet(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("HMSET", "index:1", "1", redigomock.NewAnyData())

	err := Cache.HMSet("index:1", "1", []byte("hello"))

	assert.NoError(t, err, "An error was not expected")

}

func TestMethodDelete(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("DEL", "index:1", "thread:2")

	err := Cache.Delete("index:1", "thread:2")

	assert.NoError(t, err, "An error was not expected")

}

func TestMethodFlush(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("FLUSHALL")

	err := Cache.Flush()

	assert.NoError(t, err, "An error was not expected")

}

func TestMethodIncr(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("INCR", "login:2").Expect([]byte("2"))

	res, err := Cache.Incr("login:2")

	assert.NotEmpty(t, res, "Should return data")

	assert.NoError(t, err, "An error was not expected")

}

func TestMethodExpire(t *testing.T) {

	NewRedisMock()

	Cache.Mock.Command("EXPIRE", "new:1", redigomock.NewAnyData())

	err := Cache.Expire("new:1", 600)

	assert.NoError(t, err, "An error was not expected")

}

package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stvp/tempredis"
)

func TestMutex(t *testing.T) {

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

	err = Cache.Lock("test:mutex")

	assert.NoError(t, err, "An error was not expected")

	assert.True(t, Cache.Unlock("test:mutex"), "Mutex should be unlocked")

	err = Cache.Mutex.Lock("test:mutex")

	assert.NoError(t, err, "An error was not expected")

	err = Cache.Lock("test:mutex")

	assert.NoError(t, err, "An error was not expected")

	assert.True(t, Cache.Unlock("test:mutex"), "Mutex should be unlocked")

	server.Term()

}

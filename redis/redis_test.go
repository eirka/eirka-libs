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

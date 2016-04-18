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

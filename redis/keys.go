package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
)

var (
	RedisKeys     []RedisKey
	RedisKeyIndex = make(map[string]RedisKey)
	timeout       = 600
)

type RedisKey struct {
	key        string
	base       string
	fieldcount int
	hash       bool
	hashid     uint
	expire     bool
	keyset     bool
	hashidset  bool
}

func init() {

	// the canonical list of redis keys
	RedisKeys = []RedisKey{
		RedisKey{base: "index", fieldcount: 1, hash: true, expire: false},
		RedisKey{base: "image", fieldcount: 1, hash: true, expire: false},
		RedisKey{base: "tags", fieldcount: 1, hash: true, expire: false},
		RedisKey{base: "tag", fieldcount: 2, hash: true, expire: true},
		RedisKey{base: "thread", fieldcount: 2, hash: true, expire: false},
		RedisKey{base: "post", fieldcount: 2, hash: true, expire: false},
		RedisKey{base: "directory", fieldcount: 1, hash: false, expire: false},
		RedisKey{base: "favorited", fieldcount: 1, hash: false, expire: true},
		RedisKey{base: "new", fieldcount: 1, hash: false, expire: true},
		RedisKey{base: "popular", fieldcount: 1, hash: false, expire: true},
		RedisKey{base: "imageboards", fieldcount: 1, hash: false, expire: true},
	}

	// key index map
	for _, key := range RedisKeys {
		RedisKeyIndex[key.base] = key
	}

}

func (r *RedisKey) SetKey(ids ...uint) *RedisKey {

	var keys []string

	for _, id := range ids {
		keys = append(keys, strconv.Itoa(int(id)))
	}

	r.key = strings.Join([]string{r.base, strings.Join(keys, ":")}, ":")

	r.keyset = true

	return r
}

func (r *RedisKey) SetHashId(id uint) {

	r.hashid = id

	r.hashidset = true

	return
}

func (r *RedisKey) String() string {

	if r.keyset {
		return r.key
	}

	return ""
}

func (r *RedisKey) IsValid() bool {

	if !r.keyset {
		return false
	}

	if r.hash && !r.hashidset {
		return false
	}

	if !r.hash && r.hashidset {
		return false
	}

	return true
}

func (r *RedisKey) Get() (result []byte, err error) {

	if !r.IsValid() {
		err = errors.New("key is not valid")
		return
	}

	conn := RedisCache.Pool.Get()
	defer conn.Close()

	if r.hash {

		result, err = redis.Bytes(conn.Do("HGET", r.key, r.hashid))
		if err != nil {
			return
		}
		if result == nil {
			return nil, ErrCacheMiss
		}

	} else {

		result, err = redis.Bytes(conn.Do("GET", r.key))
		if err != nil {
			return
		}
		if result == nil {
			return nil, ErrCacheMiss
		}

	}

	return
}

func (r *RedisKey) Set(data []byte) (err error) {

	if !r.IsValid() {
		err = errors.New("key is not valid")
		return
	}

	conn := RedisCache.Pool.Get()
	defer conn.Close()

	if r.hash {
		_, err = conn.Do("HMSET", r.key, r.hashid, data)
		if err != nil {
			return
		}
	} else {
		_, err = conn.Do("SET", r.key, data)
		if err != nil {
			return
		}
	}

	if r.expire {
		_, err = conn.Do("EXPIRE", r.key, timeout)
		if err != nil {
			return
		}
	}

	return
}

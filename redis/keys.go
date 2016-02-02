package redis

import (
	"errors"
	"fmt"
	"strings"
)

type RedisKeyer interface {
	String() string
	SetKey(ids ...string) *RedisKey
	Get() (result []byte, err error)
	Set(data []byte) (err error)
	Delete() (err error)
}

type RedisKey struct {
	base       string
	fieldcount int
	hash       bool
	expire     bool
	lock       bool
	key        string
	hashid     string
	keyset     bool
}

var _ = RedisKeyer(&RedisKey{})

var (
	RedisKeyIndex = make(map[string]RedisKey)
	RedisKeys     = []RedisKey{
		{base: "index", fieldcount: 1, hash: true, expire: false, lock: true},
		{base: "thread", fieldcount: 2, hash: true, expire: false},
		{base: "tag", fieldcount: 2, hash: true, expire: true},
		{base: "image", fieldcount: 1, hash: true, expire: false},
		{base: "post", fieldcount: 2, hash: true, expire: false},
		{base: "tags", fieldcount: 1, hash: true, expire: false},
		{base: "directory", fieldcount: 1, hash: true, expire: false},
		{base: "new", fieldcount: 1, hash: false, expire: true},
		{base: "popular", fieldcount: 1, hash: false, expire: true},
		{base: "favorited", fieldcount: 1, hash: false, expire: true},
		{base: "tagtypes", fieldcount: 0, hash: false, expire: false},
		{base: "imageboards", fieldcount: 0, hash: false, expire: true},
	}
)

func init() {
	// key index map
	for _, key := range RedisKeys {
		RedisKeyIndex[key.base] = key
	}
}

// return a string version of the key
func (r *RedisKey) String() string {
	return r.key
}

// populates the fields in a key and sets the hash id
func (r *RedisKey) SetKey(ids ...string) *RedisKey {

	// set the key to the base if theres no fields
	if r.fieldcount == 0 {
		r.key = r.base
		return
	}

	// create our key
	r.key = strings.Join([]string{r.base, strings.Join(ids[:r.fieldcount], ":")}, ":")

	// get our hash id
	if r.hash {
		r.hashid = strings.Join(ids[r.fieldcount:], "")
	}

	r.keyset = true

	return r
}

// gets a key, automatically handles hash types
func (r *RedisKey) Get() (result []byte, err error) {

	if !r.keyset {
		return errors.New("Key is not set")
	}

	if r.hash {
		return RedisCache.HGet(r.key, r.hashid)
	} else {
		return RedisCache.Get(r.key)
	}

	return
}

// sets a key, handles hash types and expiry
func (r *RedisKey) Set(data []byte) (err error) {

	if !r.keyset {
		return errors.New("Key is not set")
	}

	if r.hash {
		err = RedisCache.HMSet(r.key, r.hashid, data)
	} else {
		err = RedisCache.Set(r.key, data)
	}
	if err != nil {
		return
	}

	// unlock this key
	if r.lock {
		RedisCache.Unlock(fmt.Sprintf("%s:mutex", r.key))
	}

	// expire the key if set
	if r.expire {
		return RedisCache.Expire(r.key, 600)
	}

	return
}

// deletes a key
func (r *RedisKey) Delete() (err error) {

	if !r.keyset {
		return errors.New("Key is not set")
	}

	err = RedisCache.Delete(r.key)
	if err != nil {
		return
	}

	// lock this key
	if r.lock {
		RedisCache.Lock(fmt.Sprintf("%s:mutex", r.key))
	}

	return
}

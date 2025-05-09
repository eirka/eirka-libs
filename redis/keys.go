package redis

import (
	"errors"
	"fmt"
	"strings"
)

// Keyer describes the explicit key functions
type Keyer interface {
	SetKey(ids ...string) *Key
	Get() (result []byte, err error)
	Set(data []byte) (err error)
	Delete() (err error)
	String() string
}

// Key holds an explicit keys data
type Key struct {
	base       string
	fieldcount int
	hash       bool
	expire     bool
	lock       bool
	key        string
	hashid     string
	keyset     bool
}

var _ = Keyer(&Key{})

var (
	// ErrKeyNotSet returns if the key was not set properly
	ErrKeyNotSet = errors.New("key not set")
	// RedisKeyIndex holds a searchable index of keys
	RedisKeyIndex = make(map[string]Key)
	// RedisKeys is a slice of all the explicit keys
	RedisKeys = []Key{
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
func (r *Key) String() string {
	return r.key
}

// NewKey returns a key from the index or nil if it doesnt exist
func NewKey(name string) *Key {
	key, ok := RedisKeyIndex[name]
	if !ok {
		return nil
	}
	return &key
}

// SetKey populates the fields in a key and sets the hash id
func (r *Key) SetKey(ids ...string) *Key {

	// set the key to the base if theres no fields
	if r.fieldcount == 0 {
		r.key = r.base
		r.keyset = true
		return r
	}

	// the total mount of fields
	total := r.fieldcount

	if r.hash {
		total++
	}

	// make sure total fields given is not more than allowed
	if len(ids) != total {
		return r
	}

	// create our key
	r.key = strings.Join([]string{r.base, strings.Join(ids[:r.fieldcount], ":")}, ":")

	// get our hash key
	if r.hash {
		r.hashid = strings.Join(ids[r.fieldcount:], "")
	}

	r.keyset = true

	return r
}

// Get gets a key, automatically handles hash types
func (r *Key) Get() (result []byte, err error) {

	if !r.keyset {
		return nil, ErrKeyNotSet
	}

	if !isCacheInitialized() {
		return nil, ErrCacheNotInitialized
	}

	if r.hash {
		return Cache.HGet(r.key, r.hashid)
	}

	return Cache.Get(r.key)

}

// Set sets a key, handles hash types and expiry
func (r *Key) Set(data []byte) (err error) {

	if !r.keyset {
		return ErrKeyNotSet
	}

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	if r.hash {
		err = Cache.HMSet(r.key, r.hashid, data)
	} else {
		err = Cache.Set(r.key, data)
	}
	if err != nil {
		return
	}

	// unlock this key
	if r.lock {
		Cache.Unlock(fmt.Sprintf("%s:mutex", r.key))
	}

	// expire the key if set
	if r.expire {
		return Cache.Expire(r.key, 600)
	}

	return
}

// Delete deletes a key
func (r *Key) Delete() (err error) {

	if !r.keyset {
		return ErrKeyNotSet
	}

	if !isCacheInitialized() {
		return ErrCacheNotInitialized
	}

	err = Cache.Delete(r.key)
	if err != nil {
		return
	}

	// lock this key
	if r.lock {
		err = Cache.Lock(fmt.Sprintf("%s:mutex", r.key))
	}

	return
}

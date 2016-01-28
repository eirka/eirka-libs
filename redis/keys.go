package redis

type RedisKey struct {
	base       string
	fieldcount int
	hash       bool
	expire     bool
	key        string
	hashid     string
}

var (
	RedisKeyIndex = make(map[string]RedisKey)
	RedisKeys     = []RedisKey{
		{base: "index", fieldcount: 1, hash: true, expire: false},
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

// populates the fields in a key and sets the hash id
func (r *RedisKey) SetKey(ids ...string) {

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

	return
}

// gets a key, automatically handles hash types
func (r *RedisKey) Get() (result []byte, err error) {

	if r.hash {
		return RedisCache.HGet(r.key, r.hashid)
	} else {
		return RedisCache.Get(r.key)
	}

	return
}

// sets a key, handles hash types and expiry
func (r *RedisKey) Set(data []byte) (err error) {

	if r.hash {
		err = RedisCache.HMSet(r.key, r.hashid, data)
	} else {
		err = RedisCache.Set(r.key, data)
	}
	if err != nil {
		return
	}

	if r.expire {
		return RedisCache.Expire(r.key, 600)
	}

	return
}

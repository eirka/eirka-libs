package redis

// redis mutex based on https://github.com/hjr265/redsync.go

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	// DefaultExpiry is used when Mutex Duration is 0
	DefaultExpiry = 8 * time.Second
	// DefaultTries is used when Mutex Duration is 0
	DefaultTries = 16
	// DefaultDelay is used when Mutex Delay is 0
	DefaultDelay = 512 * time.Millisecond
	// DefaultFactor is used when Mutex Factor is 0
	DefaultFactor = 0.01
)

var (
	// ErrFailed is returned when lock cannot be acquired
	ErrFailed = errors.New("failed to acquire lock")
)

// Locker interface with Lock returning an error when lock cannot be aquired
type Locker interface {
	Lock(string) error
	Unlock(string) bool
}

// Pool is a generic connection pool
type Pool interface {
	Get() redis.Conn
}

var _ = Pool(&redis.Pool{})

// A Mutex is a mutual exclusion lock.
//
// Fields of a Mutex must not be changed after first use.
type Mutex struct {
	Expiry time.Duration // Duration for which the lock is valid, DefaultExpiry if 0

	Tries int           // Number of attempts to acquire lock before admitting failure, DefaultTries if 0
	Delay time.Duration // Delay between two attempts to acquire lock, DefaultDelay if 0

	Factor float64 // Drift factor, DefaultFactor if 0

	Quorum int // Quorum for the lock, set to len(addrs)/2+1 by NewMutex()

	nodes []Pool
	nodem sync.Mutex
}

var _ = Locker(&Mutex{})

// NewMutexWithGenericPool returns a new Mutex on a named resource connected to the Redis instances at given generic Pools.
// different from NewMutexWithPool to maintain backwards compatibility
func NewMutex(genericNodes []Pool) *Mutex {
	if len(genericNodes) == 0 {
		panic("no pools given")
	}

	return &Mutex{
		Quorum: len(genericNodes)/2 + 1,
		nodes:  genericNodes,
	}
}

// Lock locks m.
// In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock(key string) error {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	value := base64.StdEncoding.EncodeToString(b)

	expiry := m.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}

	retries := m.Tries
	if retries == 0 {
		retries = DefaultTries
	}

	for i := 0; i < retries; i++ {
		n := 0
		start := time.Now()
		for _, node := range m.nodes {
			if node == nil {
				continue
			}

			conn := node.Get()
			reply, err := redis.String(conn.Do("SET", key, value, "NX", "PX", int(expiry/time.Millisecond)))
			conn.Close()
			if err != nil {
				continue
			}
			if reply != "OK" {
				continue
			}
			n++
		}

		factor := m.Factor
		if factor == 0 {
			factor = DefaultFactor
		}

		until := time.Now().Add(expiry - time.Now().Sub(start) - time.Duration(int64(float64(expiry)*factor)) + 2*time.Millisecond)
		if n >= m.Quorum && time.Now().Before(until) {
			return nil
		}

		for _, node := range m.nodes {
			if node == nil {
				continue
			}

			conn := node.Get()
			_, err := conn.Do("DEL", key)
			conn.Close()
			if err != nil {
				continue
			}
		}

		delay := m.Delay
		if delay == 0 {
			delay = DefaultDelay
		}
		time.Sleep(delay)
	}

	return ErrFailed
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
// It returns the status of the unlock
func (m *Mutex) Unlock(key string) bool {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	n := 0
	for _, node := range m.nodes {
		if node == nil {
			continue
		}

		conn := node.Get()
		status, err := conn.Do("DEL", key)
		conn.Close()
		if err != nil {
			continue
		}
		if status == 0 {
			continue
		}
		n++
	}
	if n >= m.Quorum {
		return true
	}
	return false
}

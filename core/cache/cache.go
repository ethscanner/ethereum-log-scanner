package cache

import (
	"errors"
	"sync"

	"github.com/ethscanner/ethereum-log-scanner/core/cache/lru"
)

var (
	mu     sync.RWMutex
	caches = make(map[string]*Cache)
)

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

type Cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	getter     Getter
	cacheBytes int64
}

func (c *Cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.add(key, value)
}

func (c *Cache) add(key string, value ByteView) {
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return c.getLocally(key)
}

func (c *Cache) getLocally(key string) (ByteView, bool) {
	bytes, err := c.getter.Get(key)
	if err != nil {
		return ByteView{}, false
	}
	value := ByteView{b: cloneBytes(bytes)}
	c.add(key, value)
	return value, true
}

func NewCacheByName(name string, cacheBytes int64, getter Getter) *Cache {
	mu.Lock()
	defer mu.Unlock()
	if getter == nil {
		var none GetterFunc = func(key string) ([]byte, error) {
			return nil, errors.New("Not found")
		}
		getter = none
	}
	g := &Cache{cacheBytes: cacheBytes, getter: getter}
	caches[name] = g
	return g
}

func GetCacheByName(name string) *Cache {
	mu.RLock()
	defer mu.RUnlock()
	return caches[name]
}

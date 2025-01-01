package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	ttl   time.Duration
	store *map[string]cacheEntry
	mu    *sync.RWMutex
}

type cacheEntry struct {
	expiration time.Time
	value      []byte
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := (*c.store)[key]
	if !ok {
		return nil, false
	}
	return entry.value, true
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	(*c.store)[key] = cacheEntry{
		expiration: time.Now().Add(c.ttl),
		value:      value,
	}
}

func (c *Cache) startReapLoop() {
	go func() {
		for {
			time.Sleep(c.ttl)
			c.reap()
		}
	}()
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range *c.store {
		if entry.expiration.Compare(time.Now()) < 0 {
			delete(*c.store, key)
		}
	}
}

func NewCache(ttl time.Duration) Cache {
	cache := Cache{
		ttl:   ttl,
		store: &map[string]cacheEntry{},
		mu:    &sync.RWMutex{},
	}

	cache.startReapLoop()
	return cache
}

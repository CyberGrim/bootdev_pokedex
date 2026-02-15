package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mutex    sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, exists := c.cache[key]
	if exists {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mutex.Lock()
		for key, _ := range c.cache {
			cutoff := time.Now().Add(-c.interval)
			if c.cache[key].createdAt.Before(cutoff) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

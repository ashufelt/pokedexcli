package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mem      map[string]cacheEntry
	mu       *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{mem: make(map[string]cacheEntry), mu: &sync.Mutex{}, interval: interval}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mem[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if data, ok := c.mem[key]; ok {
		return data.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	myTicker := time.NewTicker(c.interval)
	for range myTicker.C {
		c.mu.Lock()
		for key, val := range c.mem {
			if time.Since(val.createdAt) > c.interval {
				delete(c.mem, key)
			}
		}
		c.mu.Unlock()
	}
}

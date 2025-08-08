package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val 	  []byte
}

type Cache struct {
	mu 		 sync.Mutex
	items 	 map[string]cacheEntry
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()

	return c

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.items[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	
	for {
		<-ticker.C

		c.mu.Lock()
		for key, item := range c.items {
			if time.Since(item.createdAt) > c.interval {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

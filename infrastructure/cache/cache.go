package cache

import (
	"sync"
	"time"
)

type Cache struct {
	store      sync.Map
	ttlEnabled bool
	ttl        int64
	mu         sync.Mutex
}

func NewCacheWithTTL(ttl int64) *Cache {
	c := Cache{
		store:      sync.Map{},
		ttlEnabled: true,
		ttl:        ttl,
	}
	go c.startTTL()
	return &c
}

func NewCache() *Cache {
	return &Cache{
		store: sync.Map{},
	}
}

func (c *Cache) startTTL() {
	for now := range time.Tick(time.Second) {
		c.mu.Lock()

		c.store.Range(func(key, value interface{}) bool {
			if now.Unix() > c.ttl {
				c.store.Delete(key)
			}
			return true
		})

		c.mu.Unlock()
	}
}

func (c *Cache) Load(k string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.store.Load(k)
}

func (c *Cache) Store(k string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store.Store(k, v)
}

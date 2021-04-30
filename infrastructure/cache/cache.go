package cache

import (
	"sync"
	"time"
)

type Cache struct {
	store      sync.Map
	ttlEnabled bool
	ttl        time.Time
}

func NewCacheWithTTL(ttl time.Time) *Cache {
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
		c.store.Range(func(key, value interface{}) bool {
			if now.Unix() > c.ttl.Unix() {
				c.store.Delete(key)

			}
			return true
		})
	}
}

func (c *Cache) Load(k string) (interface{}, bool) {
	return c.store.Load(k)
}

func (c *Cache) Store(k string, v interface{}) {
	c.store.Store(k, v)
}

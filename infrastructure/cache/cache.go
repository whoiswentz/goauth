package cache

import (
	"errors"
	"sync"
	"time"
)

type element struct {
	key   string
	value interface{}
	ttl   int64
}

type Cache struct {
	store    map[string]*element
	l        sync.Mutex
	override bool
}

var (
	ErrZeroTTL          = errors.New("ttl must be greater than zero")
	ErrCacheMiss        = errors.New("cache miss error")
	ErrKeyAlreadyExists = errors.New("key already exists")
)

func NewCacheWithTTL() *Cache {
	c := Cache{
		store: make(map[string]*element),
	}
	go c.startTTL()
	return &c
}

func NewCacheWithTTLAndOverride(o bool) *Cache {
	c := Cache{
		store:    make(map[string]*element),
		override: o,
	}
	go c.startTTL()
	return &c
}

func (c *Cache) startTTL() {
	for now := range time.Tick(time.Second) {
		c.l.Lock()

		for k, v := range c.store {
			if v.ttl > 0 && now.Unix() > v.ttl {
				delete(c.store, k)
			}
		}

		c.l.Unlock()
	}
}

func (c *Cache) Len() int {
	return len(c.store)
}

func (c *Cache) Put(k string, ttl int, v interface{}) error {
	c.l.Lock()
	defer c.l.Unlock()

	_, ok := c.store[k]
	if !ok {
		if !c.override {
			return ErrKeyAlreadyExists
		}

		c.store[k] = &element{key: k, value: v, ttl: int64(ttl)}
	}

	return nil
}

func (c *Cache) Get(k string) (interface{}, error) {
	c.l.Lock()
	defer c.l.Unlock()

	if it, ok := c.store[k]; ok {
		return it.value, nil
	}

	return nil, ErrCacheMiss
}

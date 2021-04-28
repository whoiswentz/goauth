package cache

import (
	"errors"
	"sync"
	"time"
)

type item struct {
	key   string
	value interface{}
	ttl   int64
}

type Cache struct {
	store    map[string]*item
	l        sync.Mutex
	override bool
}

var (
	ErrZeroTTL          = errors.New("ttl must be greater than zero")
	ErrCacheMiss        = errors.New("cache miss error")
	ErrKeyAlreadyExists = errors.New("key already exists")
)

func NewCacheWithTTL() *Cache {
	c := Cache{store: make(map[string]*item)}
	go func() {
		for now := range time.Tick(time.Second) {
			c.l.Lock()

			for k, v := range c.store {
				if v.ttl > 0 && now.Unix() > v.ttl {
					delete(c.store, k)
				}
			}
			c.l.Unlock()
		}
	}()
	return &c
}

func (c *Cache) Len() int {
	return len(c.store)
}

func (c *Cache) Put(k string, ttl int, v interface{}) error {
	c.l.Lock()

	_, ok := c.store[k]
	if !ok {
		c.store[k] = &item{
			key:   k,
			value: v,
			ttl:   int64(ttl),
		}

		if !c.override {
			return ErrKeyAlreadyExists
		}
	}

	c.l.Unlock()
	return nil
}

func (c *Cache) Get(k string) (interface{}, error) {
	c.l.Lock()

	if it, ok := c.store[k]; ok {
		return it.value, nil
	}

	c.l.Unlock()
	return nil, ErrCacheMiss
}

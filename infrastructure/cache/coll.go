package cache

import (
	"errors"
	"sync"
)

var (
	ErrCollNotExists = errors.New("collection not exists")
)

type CacheCollection struct {
	coll map[string]*Cache
	l    sync.Mutex
}

func NewCacheCollection() *CacheCollection {
	c := &CacheCollection{
		coll: make(map[string]*Cache),
	}

	return c
}

func (c *CacheCollection) CollectionLen() int {
	return len(c.coll)
}

func (c *CacheCollection) Len(s string) int {
	it, ok := c.coll[s]
	if !ok {
		return 0
	}

	return len(it.store)
}

func (c *CacheCollection) Put(cache string, k string, ttl int, v interface{}) error {
	c.l.Lock()
	err := c.coll[cache].Put(k, ttl, v)
	c.l.Unlock()
	return err
}

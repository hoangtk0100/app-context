package cache

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type localCache struct {
	store  map[string]interface{}
	locker *sync.RWMutex
}

func NewLocalCache() *localCache {
	return &localCache{
		store:  make(map[string]interface{}),
		locker: new(sync.RWMutex),
	}
}

func (c *localCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.store[key] = value

	return nil
}

func (c *localCache) Get(ctx context.Context, key string, value interface{}) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	val, ok := c.store[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(val).Elem())

	return nil
}

func (c *localCache) Delete(ctx context.Context, key string) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	delete(c.store, key)

	return nil
}

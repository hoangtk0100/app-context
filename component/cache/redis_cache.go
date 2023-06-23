package cache

import (
	"context"
	"time"

	rdcache "github.com/go-redis/cache/v9"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/common"
)

type redisCache struct {
	store *rdcache.Cache
}

func NewRedisCache(redisDBComponentName string, appCtx appctx.AppContext) *redisCache {
	redisDB := appCtx.MustGet(redisDBComponentName).(common.RedisDBComponent)

	c := rdcache.New(&rdcache.Options{
		Redis:      redisDB.GetDB(),
		LocalCache: rdcache.NewTinyLFU(1000, time.Minute),
	})

	return &redisCache{
		store: c,
	}
}

func (rdc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rdc.store.Set(&rdcache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (rdc *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	return rdc.store.Get(ctx, key, value)
}

func (rdc *redisCache) Delete(ctx context.Context, key string) error {
	return rdc.store.Delete(ctx, key)
}

package caching

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewCacheManager(options *RedisOptions) (*cache.Cache[string], error) {
	if options.Enabled {
		return prepareRedis(options), nil
	} else {
		return prepareInMemory()
	}
}

func prepareInMemory() (*cache.Cache[string], error) {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000,
		MaxCost:     10000,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache, store.WithSynchronousSet())

	cacheManager := cache.New[string](ristrettoStore)

	return cacheManager, nil
}

func prepareRedis(options *RedisOptions) *cache.Cache[string] {
	redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", options.Host, options.Port),
		Password: options.Password,
		DB:       options.Database,
		PoolSize: options.PoolSize,
	}))

	cacheManager := cache.New[string](redisStore)

	return cacheManager
}

func RunCaching(lc fx.Lifecycle, logger *zap.Logger, cache *cache.Cache[string], ctx context.Context) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("starting cache manager...")

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("close cache manager...")
			return nil
		},
	})

	return nil
}

func MarshalCacheItem(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func CheckIsCacheValueNotFound(err error) bool {
	return err == redis.Nil || store.NotFound{}.Is(err)
}

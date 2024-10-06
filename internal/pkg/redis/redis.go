package redis

import (
	"context"
	"fmt"

	redisCli "github.com/redis/go-redis/v9"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"go.uber.org/fx"
)

func NewRedisClient(options *RedisOptions) *redisCli.Client {
	client := redisCli.NewClient(&redisCli.Options{
		Addr:     fmt.Sprintf("%s:%d", options.Host, options.Port),
		Password: options.Password,
		DB:       options.Database,
		PoolSize: options.PoolSize,
	})

	return client
}

func RegisterRedisServer(lc fx.Lifecycle, client *redisCli.Client, log logger.ILogger, ctx context.Context) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			if err := client.Close(); err != nil {
				log.Errorf("error in closing redis: %v", err)
			} else {
				log.Info("redis closed gracefully")
			}

			return nil
		},
	})

	return nil
}

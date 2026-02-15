package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/andreparelho/order-api/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisApi interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

type redisApi struct {
	api RedisApi
}

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) error
}

type redisAdapter struct {
	client *redisApi
}

func NewRedisClient(cfg config.Configuration, ctx context.Context) (RedisClient, error) {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%v:%v@%v/%v", cfg.Redis.User, cfg.Redis.Password, cfg.Redis.Addr, cfg.Redis.DBName))
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	api := &redisApi{
		api: client,
	}

	return &redisAdapter{
		client: api,
	}, nil
}

func (r *redisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.api.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisAdapter) Get(ctx context.Context, key string) error {
	_, err := r.client.api.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

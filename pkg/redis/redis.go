package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/andreparelho/order-api/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisApi interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Close() error
}

type redisApi struct {
	api RedisApi
}

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) error
	Close()
}

type client struct {
	client *redisApi
}

func NewRedisClient(cfg config.Configuration, ctx context.Context) (RedisClient, error) {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%v:%v@%v/%v", cfg.Redis.User, cfg.Redis.Password, cfg.Redis.Addr, 0))
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opt)

	api := &redisApi{
		api: redisClient,
	}

	return &client{
		client: api,
	}, nil
}

func (r *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.api.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *client) Get(ctx context.Context, key string) error {
	_, err := r.client.api.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *client) Close() {
	err := r.client.api.Close()
	if err != nil {
		log.Fatal(err)
	}
}

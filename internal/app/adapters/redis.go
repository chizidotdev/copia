package adapters

import (
	"context"
	"encoding/json"
	"github.com/chizidotdev/copia/config"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

var _ core.RedisRepository = (*RedisClient)(nil)

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: config.EnvVars.RedisUrl,
		DB:   0,
	})

	log.Println("redis client", client)
	return &RedisClient{
		Client: client,
	}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, jsonValue, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal([]byte(result), &value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

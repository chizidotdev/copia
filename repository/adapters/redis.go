package adapters

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(connString string) *RedisStore {
	opt, err := redis.ParseURL(connString)
	if err != nil {
		log.Fatal("Cannot parse redis url:", err)
	}

	client := redis.NewClient(opt)
	return &RedisStore{
		Client: client,
	}
}

func (r *RedisStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, jsonValue, expiration).Err()
}

func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
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

func (r *RedisStore) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

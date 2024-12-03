package redis_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	once   sync.Once
	client *redis.Client
	ctx    = context.Background()
)

func InitializeRedis(addr string, password string, db int) error {
	var initError error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		_, initError = client.Ping(ctx).Result()
		if initError != nil {
			client = nil
		}
	})

	if client == nil {
		return fmt.Errorf("could not initialize Redis client: %w", initError)
	}

	return nil
}

func GetRedisClient() *redis.Client {
	if client == nil {
		panic("Redis client is not initialized. Call InitializeRedis first.")
	}
	return client
}

func Create(key string, value string, expiration time.Duration) error {
	return GetRedisClient().Set(ctx, key, value, expiration).Err()
}

func Read(key string) (string, error) {
	val, err := GetRedisClient().Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", fmt.Errorf("could not get key: %w", err)
	}
	return val, nil
}

func Update(key string, value string, expiration time.Duration) error {
	return Create(key, value, expiration)
}

func Delete(key string) error {
	_, err := GetRedisClient().Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("could not delete key: %w", err)
	}
	return nil
}

func JSONSet(prefix string, id string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("%s:%s", prefix, id)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal data to JSON: %w", err)
	}

	err = GetRedisClient().Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("could not set JSON data: %w", err)
	}
	return nil
}

func JSONGet(prefix string, id string, target interface{}) error {
	key := fmt.Sprintf("%s:%s", prefix, id)

	val, err := GetRedisClient().Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key does not exist")
	} else if err != nil {
		return fmt.Errorf("could not get JSON data: %w", err)
	}

	err = json.Unmarshal([]byte(val), target)
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON data: %w", err)
	}
	return nil
}

func JSONDelete(prefix string, id string) error {
	key := fmt.Sprintf("%s:%s", prefix, id)
	_, err := GetRedisClient().Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("could not delete JSON data: %w", err)
	}
	return nil
}

func JSONExists(prefix string, id string) (bool, error) {
	key := fmt.Sprintf("%s:%s", prefix, id)
	result, err := GetRedisClient().Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("could not check existence of key: %w", err)
	}
	return result > 0, nil
}

package cache

import (
	"time"

	"github.com/go-redis/redis"
)

// NewRedisCache creates a new RedisCache with the given connection string
func NewRedisCache(connection string) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     connection,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisCache{client, 120 * time.Second}
}

// Exists checks to see if a key exists in the cache
func (r *RedisCache) Exists(key string) (bool, error) {
	c, err := r.client.Exists(key).Result()
	return (c > 0), err
}

// Get an image from the Redis store
func (r *RedisCache) Get(key string) ([]byte, error) {
	return r.client.Get(key).Bytes()
}

// Put an image to the Redis store
func (r *RedisCache) Put(key string, data []byte) error {
	return r.client.Set(key, data, r.expiration).Err()
}

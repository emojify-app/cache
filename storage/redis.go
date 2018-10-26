package storage

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisStore implements the Cache API and uses Redis as a backend
type RedisStore struct {
	client     *redis.Client
	expiration time.Duration
}

// NewRedisStore creates a new RedisCache with the given connection string
func NewRedisStore(connection string) Store {
	client := redis.NewClient(&redis.Options{
		Addr:     connection,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisStore{client, 120 * time.Second}
}

// Exists checks to see if a key exists in the cache
func (r *RedisStore) Exists(key string) (bool, error) {
	c, err := r.client.Exists(key).Result()
	return (c > 0), err
}

// Get an image from the Redis store
func (r *RedisStore) Get(key string) ([]byte, error) {
	return r.client.Get(key).Bytes()
}

// Put an image to the Redis store
func (r *RedisStore) Put(key string, data []byte) error {
	return r.client.Set(key, data, r.expiration).Err()
}

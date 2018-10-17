package cache

import (
	"fmt"
	"io"
	"time"

	"crypto/md5"

	"github.com/go-redis/redis"
)

// Cache defines an interface for an image cache
type Cache interface {
	// Exists checks if an item exists in the cache
	Exists(string) (bool, error)
	// Get an image from the cache, returns true, image if found in cache or false, nil if image not found
	Get(string) ([]byte, error)
	// Put an image into the cache, returns an error if unsuccessful
	Put(string, []byte) error
}

// RedisCache implements the Cache API and uses Redis as a backend
type RedisCache struct {
	client     *redis.Client
	expiration time.Duration
}

// HashFilename creates a md5 hash of the given filename
func HashFilename(f string) string {
	h := md5.New()
	io.WriteString(h, f)

	return fmt.Sprintf("%x", h.Sum(nil))
}

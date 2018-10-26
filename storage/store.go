package storage

import (
	"crypto/md5"
	"fmt"
	"io"
)

// Store defines an interface for an image cache
type Store interface {
	// Exists checks if an item exists in the cache
	Exists(string) (bool, error)
	// Get an image from the cache, returns true, image if found in cache or false, nil if image not found
	Get(string) ([]byte, error)
	// Put an image into the cache, returns an error if unsuccessful
	Put(string, []byte) error
}

// HashFilename creates a md5 hash of the given filename
func HashFilename(f string) string {
	h := md5.New()
	io.WriteString(h, f)

	return fmt.Sprintf("%x", h.Sum(nil))
}

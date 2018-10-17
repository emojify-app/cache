package cache

import (
	"io/ioutil"
	"os"
)

// FileCache implements the Cache interface using the local filesystem
type FileCache struct {
	path string
}

// NewFileCache creates a file based cache
func NewFileCache(path string) Cache {
	_, err := os.Open(path)
	if err != nil {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	return &FileCache{path}
}

// Exists checks to see if a file
func (r *FileCache) Exists(key string) (bool, error) {
	_, err := os.Open(r.path + key)
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// Get an image from the File store
func (r *FileCache) Get(key string) ([]byte, error) {
	f, err := os.Open(r.path + key)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// Put an image to the File store
func (r *FileCache) Put(key string, data []byte) error {
	f, err := os.Create(r.path + key)
	if err != nil {
		return err
	}

	_, err = f.Write(data)

	return err
}

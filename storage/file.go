package storage

import (
	"io/ioutil"
	"os"
)

// FileStore implements the Cache interface using the local filesystem
type FileStore struct {
	path string
}

// NewFileStore creates a file based cache
func NewFileStore(path string) Store {
	_, err := os.Open(path)
	if err != nil {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	return &FileStore{path}
}

// Exists checks to see if a file
func (r *FileStore) Exists(key string) (bool, error) {
	_, err := os.Open(r.path + key)
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// Get an image from the File store
func (r *FileStore) Get(key string) ([]byte, error) {
	f, err := os.Open(r.path + key)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// Put an image to the File store
func (r *FileStore) Put(key string, data []byte) error {
	f, err := os.Create(r.path + key)
	if err != nil {
		return err
	}

	_, err = f.Write(data)

	return err
}

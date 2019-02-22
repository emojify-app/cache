package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/emojify-app/cache/logging"
)

// FileStore implements the Cache interface using the local filesystem
type FileStore struct {
	path             string
	cacheDuration    time.Duration
	lastInvalidation time.Time
	logger           logging.Logger
}

// NewFileStore creates a file based cache
func NewFileStore(path string, ci time.Duration, l logging.Logger) Store {
	_, err := os.Open(path)
	if err != nil {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	// if the path does not end in / append
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	f := &FileStore{}
	f.path = path
	f.cacheDuration = ci
	f.logger = l
	go f.invalidateCache()

	return f
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
// when a file is not found data is nil
func (r *FileStore) Get(key string) ([]byte, error) {
	f, err := os.Open(r.path + key)
	if err != nil {
		// if the file does not exist return a nil slice
		if os.IsNotExist(err) {
			return nil, nil
		}

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

func (r *FileStore) startInvalidateCache() {
	r.lastInvalidation = time.Now()
	t := time.NewTicker(r.cacheDuration)

	for range t.C {
		r.invalidateCache()
	}
}

func (r *FileStore) invalidateCache() {
	f := r.logger.CacheInvalidate()

	// find files which have expired
	toDelete := make([]os.FileInfo, 0)
	files, err := ioutil.ReadDir(r.path)
	if err != nil {
		f(http.StatusInternalServerError, fmt.Errorf("Error reading cache directory %s", err))
		return
	}

	for _, f := range files {
		if f.ModTime().Sub(r.lastInvalidation) > r.cacheDuration {
			toDelete = append(toDelete, f)
		}
	}

	// clean up expired files
	for _, f := range toDelete {
		err := os.Remove(r.path + f.Name())
		r.logger.CacheInvalidateItem(f.Name(), err)
	}

	f(http.StatusOK, nil)
}

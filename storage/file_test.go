package storage

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/emojify-app/cache/logging"
	"github.com/stretchr/testify/assert"
)

var tmpDirectory = "/tmp/cache_test"

func setupFileStore(invalidation time.Duration) Store {
	os.Mkdir(tmpDirectory, 0755)
	l, _ := logging.New("test", "test", "localhost:8125", "DEBUG", "text")
	return NewFileStore(tmpDirectory, invalidation, l)
}

func TestPutSavesFile(t *testing.T) {
	c := setupFileStore(1 * time.Second)

	c.Put("abc", []byte("abc1223"))
	fileKey := "abc"

	file, err := os.Open(tmpDirectory + "/" + fileKey)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	defer func() {
		os.Remove(tmpDirectory + "/" + fileKey)
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "abc1223", string(data))
}

func TestExistsWithNoFileReturnsFalse(t *testing.T) {
	c := setupFileStore(1 * time.Second)

	ok, err := c.Exists("abcdefg")
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, ok)
}

func TestRemovesCachedFile(t *testing.T) {
	c := setupFileStore(20 * time.Millisecond)

	c.Put("abc", []byte("abc1223"))
	fileKey := "abc"

	st := time.Now()
	for {
		if time.Now().Sub(st) > 1000*time.Millisecond {
			t.Fatal("Timeout before cache invalidation")
		}

		_, err := os.Open(tmpDirectory + "/" + fileKey)
		ex := os.IsNotExist(err)
		if ex {
			break
		}
	}
}

func TestDoesNotRemoveCachedFile(t *testing.T) {
	c := setupFileStore(1000 * time.Millisecond)

	c.Put("abc", []byte("abc1223"))
	fileKey := "abc"
	filePath := tmpDirectory + "/" + fileKey

	time.Sleep(100 * time.Millisecond)

	_, err := os.Open(filePath)
	ex := os.IsNotExist(err)
	if ex {
		t.Fatal("shoud not have removed file")
	}

	os.Remove(filePath)
}

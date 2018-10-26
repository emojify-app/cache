package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupFileStore() Store {
	return NewFileStore("/tmp/")
}

func TestPutSavesFile(t *testing.T) {
	c := setupFileStore()

	c.Put("abc", []byte("abc1223"))
	fileKey := "abc"

	file, err := os.Open("/tmp/" + fileKey)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	defer func() {
		os.Remove("/tmp/" + fileKey)
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "abc1223", string(data))
}

func TestExistsWithNoFileReturnsFalse(t *testing.T) {
	c := setupFileStore()

	ok, err := c.Exists("abcdefg")
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, ok)
}

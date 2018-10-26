package storage

import (
	"github.com/stretchr/testify/mock"
)

// MockCache is a mock implementation of the Cache interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Exists(key string) (bool, error) {
	args := m.Called(key)

	return args.Bool(0), args.Error(1)
}

// Get calls the mock get function
func (m *MockStore) Get(key string) ([]byte, error) {
	args := m.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

// Put calls the mock put function
func (m *MockStore) Put(key string, data []byte) error {
	args := m.Called(key)
	return args.Error(0)
}

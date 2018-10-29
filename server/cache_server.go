package server

import (
	"context"

	"github.com/emojify-app/cache/protos/cache"
	"github.com/emojify-app/cache/storage"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// CacheServer implements the methofds defined in the gRPC interface
type CacheServer struct {
	store storage.Store
}

// Get an item from the cache
func (c *CacheServer) Get(ctx context.Context, key *wrappers.StringValue) (*cache.CacheItem, error) {
	encodedID := storage.HashFilename(key.Value)

	data, err := c.store.Get(encodedID)
	if err != nil {
		return &cache.CacheItem{}, err
	}

	return &cache.CacheItem{Id: key.Value, Data: data}, nil
}

// Put an item in the cache
func (c *CacheServer) Put(ctx context.Context, item *cache.CacheItem) (*wrappers.StringValue, error) {
	encodedID := storage.HashFilename(item.GetId())

	err := c.store.Put(encodedID, item.GetData())
	return &wrappers.StringValue{Value: encodedID}, err
}

// Exists checks to see if an item already exists in the cache
func (c *CacheServer) Exists(ctx context.Context, key *wrappers.StringValue) (*wrappers.BoolValue, error) {
	encodedID := storage.HashFilename(key.GetValue())

	exists, err := c.store.Exists(encodedID)
	if err != nil {
		return &wrappers.BoolValue{Value: false}, err
	}

	return &wrappers.BoolValue{Value: exists}, nil
}

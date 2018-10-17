package server

import (
	"context"

	cacheHandler "github.com/emojify-app/cache/cache"
	"github.com/emojify-app/cache/protos/cache"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// CacheServer implements the methofds defined in the gRPC interface
type CacheServer struct {
	cache cacheHandler.Cache
}

// Get an item from the cache
func (c *CacheServer) Get(ctx context.Context, key *wrappers.StringValue) (*cache.CacheItem, error) {
	data, err := c.cache.Get(key.Value)
	if err != nil {
		return &cache.CacheItem{}, err
	}

	return &cache.CacheItem{Id: key.Value, Data: data}, nil
}

// Put an item in the cache
func (c *CacheServer) Put(ctx context.Context, item *cache.CacheItem) (*wrappers.StringValue, error) {
	encodedID := cacheHandler.HashFilename(item.GetId())

	err := c.cache.Put(encodedID, item.GetData())
	return &wrappers.StringValue{Value: encodedID}, err
}

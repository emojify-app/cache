package server

import (
	"context"
	"net/http"

	"github.com/emojify-app/cache/logging"
	"github.com/emojify-app/cache/protos/cache"
	"github.com/emojify-app/cache/storage"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CacheServer implements the methofds defined in the gRPC interface
type CacheServer struct {
	store  storage.Store
	logger logging.Logger
}

// Check returns the health status of the api
func (c *CacheServer) Check(ctx context.Context, in *cache.HealthCheckRequest) (*cache.HealthCheckResponse, error) {

	return &cache.HealthCheckResponse{Status: cache.HealthCheckResponse_SERVING}, nil
}

// Get an item from the cache
func (c *CacheServer) Get(ctx context.Context, key *wrappers.StringValue) (*cache.CacheItem, error) {
	f := c.logger.CacheGet(key.Value)

	encodedID := storage.HashFilename(key.Value)

	data, err := c.store.Get(encodedID)
	if err != nil {
		f(http.StatusInternalServerError, err)

		// create a grpc error message
		gerr := status.New(codes.Internal, err.Error())
		return nil, gerr.Err()
	}

	// does the file exist?
	if data == nil || len(data) == 0 {
		f(http.StatusNotFound, nil)

		gerr := status.New(codes.NotFound, "File not found: "+key.Value)
		return nil, gerr.Err()
	}

	f(http.StatusOK, nil)
	return &cache.CacheItem{Id: key.Value, Data: data}, nil
}

// Put an item in the cache
func (c *CacheServer) Put(ctx context.Context, item *cache.CacheItem) (*wrappers.StringValue, error) {
	f := c.logger.CachePut(item.GetId())

	encodedID := storage.HashFilename(item.GetId())

	err := c.store.Put(encodedID, item.GetData())
	if err != nil {
		f(http.StatusInternalServerError, err)

		gerr := status.New(codes.Internal, err.Error())
		return nil, gerr.Err()
	}

	f(http.StatusOK, nil)
	return &wrappers.StringValue{Value: encodedID}, nil
}

// Exists checks to see if an item already exists in the cache
func (c *CacheServer) Exists(ctx context.Context, key *wrappers.StringValue) (*wrappers.BoolValue, error) {
	encodedID := storage.HashFilename(key.GetValue())
	f := c.logger.CachePut(encodedID)

	exists, err := c.store.Exists(encodedID)
	if err != nil {
		f(http.StatusInternalServerError, err)

		gerr := status.New(codes.Internal, err.Error())
		return nil, gerr.Err()
	}

	switch exists {
	case true:
		f(http.StatusFound, nil)
	case false:
		f(http.StatusNotFound, nil)
	}

	return &wrappers.BoolValue{Value: exists}, nil
}

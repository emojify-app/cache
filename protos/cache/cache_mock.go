package cache

import (
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/mock"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// ClientMock is a mock implementation of the gRCP cache client
type ClientMock struct {
	mock.Mock
}

// Put is a mock implementation of the cache put interface
func (c *ClientMock) Put(ctx context.Context, in *CacheItem, opts ...grpc.CallOption) (*wrappers.StringValue, error) {
	args := c.Mock.Called(ctx, in, opts)

	if sv := args.Get(0); sv != nil {
		return sv.(*wrappers.StringValue), args.Error(1)
	}

	return nil, args.Error(0)
}

// Get is a mock implementation of the cache get interface
func (c *ClientMock) Get(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*CacheItem, error) {
	args := c.Mock.Called(ctx, in, opts)

	if sv := args.Get(0); sv != nil {
		return sv.(*CacheItem), args.Error(1)
	}

	return nil, args.Error(0)
}

// Exists is a mock implementation of the cache exists interface
func (c *ClientMock) Exists(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	args := c.Mock.Called(ctx, in, opts)

	if sv := args.Get(0); sv != nil {
		return sv.(*wrappers.BoolValue), args.Error(1)
	}

	return nil, args.Error(0)
}

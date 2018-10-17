package server

import (
	"fmt"
	"net"

	cacheHandler "github.com/emojify-app/cache/cache"
	"github.com/emojify-app/cache/protos/cache"
	"google.golang.org/grpc"
)

var lis net.Listener
var grpcServer = grpc.NewServer()

// Start a new instance of the server
func Start(address, port string, handler cacheHandler.Cache) error {
	cache.RegisterCacheServer(grpcServer, &CacheServer{handler})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return err
	}

	return grpcServer.Serve(lis)
}

// Stop the server
func Stop() error {
	grpcServer.Stop()
	//	lis.Close()
	return nil
}

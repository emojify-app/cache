package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/emojify-app/cache/server"
	"github.com/emojify-app/cache/storage"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

var envBindAddress = env.String("BIND_ADDRESS", false, "localhost", "Bind address for gRPC server, e.g. 127.0.0.1")
var envBindPort = env.Integer("BIND_PORT", false, 9090, "Bind port for gRPC server e.g. 9090")

var envHealthBindAddress = env.String("HEALTH_BIND_ADDRESS", false, "localhost", "Bind address for health endpoint, e.g. 127.0.0.1")
var envHealthBindPort = env.Integer("HEALTH_BIND_PORT", false, 9091, "Bind port for health endpoint e.g. 9091")

var envCacheType = env.String("CACHE_TYPE", false, "file", "Cache type for server e.g. file, cloud_storage")
var envCacheFileLocation = env.String("CACHE_FILE_LOCATION", false, "/files", "Directory to store files for cache type file")
var envCacheInvalidation = env.Duration("CACHE_INVALIDATION", false, "5m", "Cache invalidation period")

var logger hclog.Logger

var version = "0.1"

var help = flag.Bool("help", false, "--help to show help")

func main() {
	flag.Parse()

	// if the help flag is passed show configuration options
	if *help == true {
		fmt.Println("Emojify service version:", version)
		fmt.Println("Configuration values are set using environment variables, for info please see the following list")
		fmt.Println("")
		fmt.Println(env.Help())
	}

	err := env.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger = hclog.Default()
	logger.Info("Started Cache Server", "version", version)

	var c storage.Store
	if *envCacheType == "file" {
		c = storage.NewFileStore(*envCacheFileLocation, *envCacheInvalidation)
	}

	logger.Info("Binding health checks to", "address", *envBindAddress, "port", *envBindPort)
	logger.Info("Starting health server")

	http.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(fmt.Sprintf("%s:%d", *envHealthBindAddress, *envHealthBindPort), nil)

	logger.Info("Binding gRPC to", "address", *envBindAddress, "port", *envBindPort)
	logger.Info("Starting gRPC server")

	err = server.Start(*envBindAddress, *envBindPort, c)
	if err != nil {
		logger.Error("Unable to start server", "error", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	cacheHandler "github.com/emojify-app/cache/cache"
	"github.com/emojify-app/cache/server"
	hclog "github.com/hashicorp/go-hclog"
)

var envBindAddress string
var envBindPort string

var envCacheType string
var envCacheFileLocation string

var logger hclog.Logger

var version = "0.1"

func main() {
	logger = hclog.Default()
	logger.Info("Started Cache Server", "version", version)

	err := processEnvVars()
	if err != nil {
		logger.Error("Error processing environment vars", "error", err)
		os.Exit(1)
	}

	var c cacheHandler.Cache
	if envCacheType == "file" {
		c = cacheHandler.NewFileCache(envCacheFileLocation)
	}

	logger.Info("Binding to", "address", envBindAddress, "port", envBindPort)
	logger.Info("Starting gRPC server")

	err = server.Start(envBindAddress, envBindPort, c)
	if err != nil {
		logger.Error("Unable to start server", "error", err)
		os.Exit(1)
	}
}

func processEnvVars() error {
	envBindAddress = os.Getenv("BIND_ADDRESS")
	if envBindAddress == "" {
		envBindAddress = "localhost"
	}

	envBindPort = os.Getenv("BIND_PORT")
	if envBindPort == "" {
		envBindPort = "8000"
	}

	envCacheType = os.Getenv("CACHE_TYPE")
	if envCacheType != "file" {
		return fmt.Errorf("Invalid cache type, only file is currently supported")
	}

	envCacheFileLocation = os.Getenv("CACHE_FILE_LOCATION")
	if envCacheFileLocation != "" {
		return fmt.Errorf("If using cache type file, you must specifiy a directory to store the cache")
	}

	return nil
}

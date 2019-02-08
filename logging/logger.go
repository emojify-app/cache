package logging

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/hashicorp/go-hclog"
)

var statsPrefix = "cache.service."

// Logger defines an interface for common logging operations
type Logger interface {
	Log() hclog.Logger

	ServiceStart(address, port, version string)

	CacheFileNotFound(string)
	CacheGetFile(string) Finished
	CachePutFile(string) Finished
}

// Finished defines a function to be returned by logging methods which contain timers
type Finished func(status int, err error)

// New creates a new logger with the given name and points it at a statsd server
func New(name, version, statsDServer, logLevel string, logFormat string) (Logger, error) {
	o := hclog.DefaultOptions
	o.Name = name

	// set the log format
	if logFormat == "json" {
		o.JSONFormat = true
	}

	o.Level = hclog.LevelFromString(logLevel)
	l := hclog.New(o)

	c, err := statsd.New(statsDServer)
	c.Tags = []string{fmt.Sprintf("version:%s", version)}

	if err != nil {
		return nil, err
	}

	return &LoggerImpl{l, c}, nil
}

// LoggerImpl is a concrete implementation for the logger function
type LoggerImpl struct {
	l hclog.Logger
	s *statsd.Client
}

// Log returns the underlying logger
func (l *LoggerImpl) Log() hclog.Logger {
	return l.l
}

// ServiceStart logs information about the service start
func (l *LoggerImpl) ServiceStart(address, port, version string) {
	l.s.Incr(statsPrefix+"started", nil, 1)
	l.l.Info("Service started", "address", address, "port", port, "version", version)
}

// CacheFileNotFound logs information when the a file is missing from the cache
func (l *LoggerImpl) CacheFileNotFound(f string) {
	l.s.Incr(statsPrefix+"cache.file_not_found", nil, 1)
	l.l.Info("File not found in cache", "file", f)
}

// CacheGetFile logs information when data is fetched from the cache
func (l *LoggerImpl) CacheGetFile(f string) Finished {
	st := time.Now()
	l.l.Info("Fetching file from cache", "file", f)

	return func(status int, err error) {
		if err != nil {
			l.s.Incr(statsPrefix+"cache.error", nil, 1)
			l.l.Error("Error fetching file from cache", "file", f, "error", err)
		}

		l.s.Timing(statsPrefix+"cache.get", time.Now().Sub(st), getStatusTags(status), 1)
	}
}

// CachePutFile logs information when data is fetched from the cache
func (l *LoggerImpl) CachePutFile(f string) Finished {
	st := time.Now()
	l.l.Info("Putting file to cache", "file", f)

	return func(status int, err error) {
		if err != nil {
			l.s.Incr(statsPrefix+"cache.error", nil, 1)
			l.l.Error("Error fetching file from cache", "file", f, "error", err)
		}

		l.s.Timing(statsPrefix+"cache.get", time.Now().Sub(st), getStatusTags(status), 1)
	}
}

func getStatusTags(status int) []string {
	return []string{
		fmt.Sprintf("status:%d", status),
	}
}

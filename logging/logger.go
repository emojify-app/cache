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

	CacheCheck() Finished
	CacheExists(string) Finished
	CacheGet(string) Finished
	CachePut(string) Finished

	CacheInvalidate() Finished
	CacheInvalidateItem(string, error)
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

// CacheGet logs information when data is fetched from the cache
func (l *LoggerImpl) CacheGet(f string) Finished {
	st := time.Now()
	l.l.Info("Fetching file from cache", "file", f)

	return func(status int, err error) {
		if err != nil {
			l.l.Error("Error fetching file from cache", "file", f, "error", err)
		}

		l.s.Timing(statsPrefix+"cache.get", time.Now().Sub(st), getStatusTags(status), 1)
	}
}

// CachePut logs information when data is fetched from the cache
func (l *LoggerImpl) CachePut(f string) Finished {
	st := time.Now()
	l.l.Info("Putting file to cache", "file", f)

	return func(status int, err error) {
		if err != nil {
			l.l.Error("Error putting file to cache", "file", f, "error", err)
		}

		l.s.Timing(statsPrefix+"cache.put", time.Now().Sub(st), getStatusTags(status), 1)
	}
}

// CacheExists logs information when the cache exists method is called
func (l *LoggerImpl) CacheExists(f string) Finished {
	st := time.Now()
	l.l.Info("Checking file is in cache", "file", f)

	return func(status int, err error) {
		if err != nil {
			l.l.Error("Error checking if file in cache", "file", f, "error", err)
		}

		l.s.Timing(statsPrefix+"cache.exists", time.Now().Sub(st), getStatusTags(status), 1)
	}
}

// CacheCheck logs information when the cache health check is called
func (l *LoggerImpl) CacheCheck() Finished {
	st := time.Now()
	l.l.Info("Health check")

	return func(status int, err error) {
		if err != nil {
			l.l.Error("Error with health check", "error", err)
		}

		l.s.Timing(statsPrefix+"cache.check", time.Now().Sub(st), getStatusTags(status), 1)
	}
}

// CacheInvalidate logs information when the cache invalidation is called
func (l *LoggerImpl) CacheInvalidate() Finished {
	st := time.Now()
	l.l.Info("Invalidate cache items")

	return func(status int, err error) {
		if err != nil {
			l.l.Error("Error invalidating cache items", "error", err)
		}

		l.s.Timing(statsPrefix+"cache.invalidate", time.Now().Sub(st), getStatusTags(status), 1)
	}
}

// CacheInvalidateItem logs information when the cache item is invalidated
func (l *LoggerImpl) CacheInvalidateItem(file string, err error) {
	if err != nil {
		l.l.Error("Unable to invalidate cache item", "file", file, "error", err)
		l.s.Incr(statsPrefix+"cache.invalidation.error", nil, 1.0)
		return
	}

	l.l.Info("Remove expired file", "file", file)
	l.s.Incr(statsPrefix+"cache.invalidated", nil, 1.0)
}

// CacheInvalidated logs the number of items invalidated from the cache
func (l *LoggerImpl) CacheInvalidated(count float64) {
	l.l.Info("Items invalidated from cache", "count", count)
	l.s.Gauge(statsPrefix+"cache.invalidated", count, nil, 1.0)
}

func getStatusTags(status int) []string {
	return []string{
		fmt.Sprintf("status:%d", status),
	}
}

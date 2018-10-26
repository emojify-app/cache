package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/emojify-app/cache/protos/cache"
	"github.com/emojify-app/cache/server"
	"github.com/emojify-app/cache/storage"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

var opt = godog.Options{Output: colors.Colored(os.Stdout)}
var bindAddress = "127.0.0.1"
var bindPort = "9000"

var cacheClient cache.CacheClient
var putReturn string
var getBody string

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func theServerIsRunning() error {
	c := storage.NewFileStore("/tmp/cache/")

	var err error
	go func() {
		err = server.Start(bindAddress, bindPort, c)
	}()
	time.Sleep(1000 * time.Millisecond)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", bindAddress, bindPort), grpc.WithInsecure())
	if err != nil {
		return err
	}
	cacheClient = cache.NewCacheClient(conn)

	return err
}

func iPutAFile() error {
	resp, err := cacheClient.Put(context.Background(), &cache.CacheItem{Id: "123", Data: []byte("abc")})
	if err != nil {
		return err
	}

	putReturn = resp.Value
	if putReturn == "" {
		return fmt.Errorf("Expected id to be returned")
	}

	return nil
}

func theFileShouldExistInTheCache() error {
	f, err := os.Open("/tmp/cache/" + putReturn)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func aFileExistsInTheCache() error {
	f, err := os.OpenFile("/tmp/cache/abcdefg", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("abc"))
	return err
}

func iGetThatFile() error {
	ci, err := cacheClient.Get(context.Background(), &wrappers.StringValue{Value: "abcdefg"})
	if err != nil {
		return err
	}

	getBody = string(ci.GetData())

	return nil
}

func theFileContentsShouldBeReturned() error {
	if getBody != "abc" {
		return fmt.Errorf("expected file contents: abc, got: %s", getBody)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	s.BeforeScenario(func(interface{}) {
		os.MkdirAll("/tmp/cache/", 0755)
	})

	s.AfterScenario(func(interface{}, error) {
		os.Remove("/tmp/cache/")
		server.Stop()
	})

	s.Step(`^the server is running$`, theServerIsRunning)
	s.Step(`^I put a file$`, iPutAFile)
	s.Step(`^the file should exist in the cache$`, theFileShouldExistInTheCache)
	s.Step(`^a file exists in the cache$`, aFileExistsInTheCache)
	s.Step(`^I Get that file`, iGetThatFile)
	s.Step(`^the file contents should be returned`, theFileContentsShouldBeReturned)
}

build_protos:
	protoc -I protos/ protos/cache.proto --go_out=plugins=grpc:protos/cache

build_snapshot: build_protos
	goreleaser --snapshot --rm-dist

test: test_unit test_functional

test_unit:
	go test -v -race `go list ./... | grep -v functional_tests`

test_functional:
	cd functional_tests && go test -v --godog.format=pretty --godog.random -race
	

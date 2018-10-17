build_protos:
	protoc -I protos/ protos/cache.proto --go_out=plugins=grpc:protos/cache

test_unit:
	go test -v -race `go list ./... | grep -v functional_tests`

test_functional:
	cd functional_tests && go test -v --godog.format=pretty --godog.random -race
	

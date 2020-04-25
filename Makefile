.PHONY: test test-short build build-arwen clean

clean:
	go clean -cache -testcache

build:
	go build ./...

build-arwen:
	go build -o ./cmd/arwen/arwen ./cmd/arwen
	cp ./cmd/arwen/arwen ./ipc/tests

test: clean build-arwen
	go test -count=1 ./...

test-short: build-arwen
	go test -short -count=1 ./...

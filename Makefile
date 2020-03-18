.PHONY: test build build-arwen clean

clean:
	go clean -cache -testcache

build:
	go build ./...

build-arwen:
	go build -o ./cmd/arwen/arwen ./cmd/arwen
	cp ./cmd/arwen/arwen ./ipc/tests

test: build-arwen run-tests

run-tests:
	go test -count=1 ./...

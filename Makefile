.PHONY: test build build-arwen clean

clean:
	go clean -cache -testcache

build:
	go build ./...

build-arwen:
	go build -o ./cmd/arwen/arwen ./cmd/arwen
	cp ./cmd/arwen/arwen ./ipc/tests

test: clean build-arwen
	go test -count=1 ./...

build-arwendebug:
	go build -o ./cmd/arwendebug/arwendebug ./cmd/arwendebug

test-arwendebug: build-arwendebug
	ARWENDEBUG=./cmd/arwendebug/arwendebug TESTDATA=./cmd/arwendebug/testdata ./cmd/arwendebug/testdata/simple.sh
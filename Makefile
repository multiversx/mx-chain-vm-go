.PHONY: test build clean

clean:
	go clean -cache -testcache

build:
	go build ./...

test:
	go build -o ./cmd/arwen/arwen ./cmd/arwen
	ARWEN_PATH=${CURDIR}/cmd/arwen/arwen go test -count=1 ./...

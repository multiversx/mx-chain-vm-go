.PHONY: test build clean

clean:
	go clean -cache -testcache

build:
	go build ./...

test:
	go test -count=1 ./...

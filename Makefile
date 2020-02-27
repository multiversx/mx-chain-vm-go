.PHONY: test build

build:
	go build ./...

test:
	go test -count=1 ./...

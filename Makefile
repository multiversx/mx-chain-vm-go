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

build-c-contracts:
	# TODO: rm *.wasm

	erdpy build ./test/contracts/erc20
	erdpy build ./test/contracts/counter

	erdpy build ./test/contracts/init-correct
	erdpy build ./test/contracts/init-simple
	erdpy build ./test/contracts/init-wrong
	erdpy build ./test/contracts/misc
	erdpy build ./test/contracts/num-with-fp
	erdpy build ./test/contracts/signatures

	erdpy build ./test/contracts/exec-same-ctx-child
	erdpy build ./test/contracts/exec-same-ctx-parent
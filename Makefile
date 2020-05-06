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
	# TODO: How to build this one?
	# erdpy build ./test/contracts/num-with-fp
	erdpy build ./test/contracts/signatures
	erdpy build ./test/contracts/elrondei
	erdpy build ./test/contracts/breakpoint

	erdpy build ./test/contracts/exec-same-ctx-simple-parent
	erdpy build ./test/contracts/exec-same-ctx-simple-child
	erdpy build ./test/contracts/exec-same-ctx-child
	erdpy build ./test/contracts/exec-same-ctx-parent
	erdpy build ./test/contracts/exec-dest-ctx-parent
	erdpy build ./test/contracts/exec-dest-ctx-child
	erdpy build ./test/contracts/exec-same-ctx-recursive
	erdpy build ./test/contracts/exec-same-ctx-recursive-parent
	erdpy build ./test/contracts/exec-same-ctx-recursive-child
	erdpy build ./test/contracts/exec-dest-ctx-recursive
	erdpy build ./test/contracts/exec-dest-ctx-recursive-parent
	erdpy build ./test/contracts/exec-dest-ctx-recursive-child
	erdpy build ./test/contracts/async-call-parent
	erdpy build ./test/contracts/async-call-child
	erdpy build ./test/contracts/exec-same-ctx-builtin



build-delegation:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-delegation-rs
	git clone --depth=1 --branch=master https://github.com/ElrondNetwork/sc-delegation-rs.git ${SANDBOX}/sc-delegation-rs
	rm -rf ${SANDBOX}/sc-delegation-rs/.git
	erdpy build ${SANDBOX}/sc-delegation-rs
	erdpy test --directory="tests" ${SANDBOX}/sc-delegation-rs
	cp ${SANDBOX}/sc-delegation-rs/output/delegation.wasm ./test/delegation/delegation.wasm

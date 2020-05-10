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

build-c-contracts:
	erdpy build ./test/contracts/erc20
	erdpy build ./test/contracts/counter

	erdpy build ./test/contracts/init-correct
	erdpy build ./test/contracts/init-simple
	erdpy build ./test/contracts/init-wrong
	erdpy build ./test/contracts/misc
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


build-dns:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-dns-rs
	git clone --depth=1 --branch=master https://github.com/ElrondNetwork/sc-dns-rs.git ${SANDBOX}/sc-dns-rs
	rm -rf ${SANDBOX}/sc-dns-rs/.git
	erdpy build ${SANDBOX}/sc-dns-rs
	erdpy test --directory="tests" ${SANDBOX}/sc-dns-rs
	cp ${SANDBOX}/sc-dns-rs/output/dns.wasm ./test/dns/dns.wasm


build-sc-examples:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples

	erdpy new --template=erc20-c --directory ${SANDBOX}/sc-examples erc20-c
	erdpy build ${SANDBOX}/sc-examples/erc20-c
	cp ${SANDBOX}/sc-examples/erc20-c/output/wrc20_arwen.wasm ./test/erc20/contracts/erc20-c.wasm


build-sc-examples-rs:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples-rs
	
	erdpy new --template=simple-coin --directory ${SANDBOX}/sc-examples-rs simple-coin
	erdpy new --template=adder --directory ${SANDBOX}/sc-examples-rs adder
	erdpy build ${SANDBOX}/sc-examples-rs/adder
	erdpy build ${SANDBOX}/sc-examples-rs/simple-coin
	erdpy test ${SANDBOX}/sc-examples-rs/adder
	erdpy test ${SANDBOX}/sc-examples-rs/simple-coin
	cp ${SANDBOX}/sc-examples-rs/adder/output/adder.wasm ./test/adder/adder.wasm
	cp ${SANDBOX}/sc-examples-rs/simple-coin/output/simple-coin.wasm ./test/erc20/contracts/simple-coin.wasm

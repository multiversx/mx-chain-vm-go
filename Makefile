.PHONY: test test-short build arwendebug clean

ARWEN_VERSION := $(shell git describe --tags --long --dirty --always)

clean:
	go clean -cache -testcache

build:
	go build ./...

arwendebug:
ifndef ARWENDEBUG_PATH
	$(error ARWENDEBUG_PATH is undefined)
endif
	go build -o ./cmd/arwendebug/arwendebug ./cmd/arwendebug
	cp ./cmd/arwendebug/arwendebug ${ARWENDEBUG_PATH}


test: clean
	go test ./...

test-v: clean
	go test ./... -v

test-serial: clean
	go test ./... -failfast -p 1

test-short: clean
	go test ./... -short

test-short-v: clean
	go test ./... -short -v

test-short-serial:
	go test ./... -short -failfast -p 1

print-api-costs:
	@echo "bigIntOps.go:"
	@grep "func v1_5\|GasSchedule" arwen/elrondapi/bigIntOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "elrondei.go:"
	@grep "func v1_5\|GasSchedule" arwen/elrondapi/elrondei.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "managedei.go:"
	@grep "func v1_5\|GasSchedule" arwen/elrondapi/managedei.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "manBufOps.go:"
	@grep "func v1_5\|GasSchedule" arwen/elrondapi/manBufOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"
	@echo "----------------"
	@echo "smallIntOps.go:"
	@grep "func v1_5\|GasSchedule" arwen/elrondapi/smallIntOps.go | sed -e "/func/ s:func v1_5_\(.*\)(.*:\1:" -e "/GasSchedule/ s:metering.GasSchedule()::"


build-test-contracts: build-test-contracts-erdpy build-test-contracts-wat

build-test-contracts-wat:
	cd test/contracts/init-simple-popcnt && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/data-drop/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-copy/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-fill/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/memory-init/output && wat2wasm *.wat
	cd test/contracts/forbidden-opcodes/simd/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/imported-global/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-exceeded-max-pages/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-exceeded-pages/output && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-grow/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-min-pages-greater-than-max-pages/output && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/mem-multiple-max-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-multiple-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-no-max-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-no-pages/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/mem-single-page/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/memoryless/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/multiple-memories/output/ && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/multiple-mutable/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/noglobals/output/ && wat2wasm *.wat
	cd test/contracts/wasmbacking/single-immutable/output/ && wat2wasm --no-check *.wat
	cd test/contracts/wasmbacking/single-mutable/output/ && wat2wasm *.wat

build-test-contracts-erdpy:
	erdpy contract build --no-optimization ./test/contracts/answer
	erdpy contract build ./test/contracts/async-call-builtin
	erdpy contract build ./test/contracts/async-call-child
	erdpy contract build ./test/contracts/async-call-parent
	erdpy contract build ./test/contracts/breakpoint
	erdpy contract build ./test/contracts/big-floats
	erdpy contract build ./test/contracts/counter
	erdpy contract build ./test/contracts/deployer
	erdpy contract build ./test/contracts/deployer-child
	erdpy contract build ./test/contracts/deployer-fromanother-contract
	erdpy contract build ./test/contracts/deployer-parent
	erdpy contract build ./test/contracts/elrondei
	erdpy contract build ./test/contracts/erc20
	erdpy contract build ./test/contracts/exchange
	
	erdpy contract build ./test/contracts/exec-dest-ctx-builtin
	erdpy contract build ./test/contracts/exec-dest-ctx-by-caller/child
	erdpy contract build ./test/contracts/exec-dest-ctx-by-caller/parent
	erdpy contract build ./test/contracts/exec-dest-ctx-child
	erdpy contract build ./test/contracts/exec-dest-ctx-esdt/basic
	erdpy contract build ./test/contracts/exec-dest-ctx-parent
	erdpy contract build ./test/contracts/exec-dest-ctx-recursive
	erdpy contract build ./test/contracts/exec-dest-ctx-recursive-child
	erdpy contract build ./test/contracts/exec-dest-ctx-recursive-parent
	
	erdpy contract build ./test/contracts/exec-same-ctx-child
	erdpy contract build ./test/contracts/exec-same-ctx-parent
	erdpy contract build ./test/contracts/exec-same-ctx-recursive
	erdpy contract build ./test/contracts/exec-same-ctx-recursive-child
	erdpy contract build ./test/contracts/exec-same-ctx-recursive-parent
	erdpy contract build ./test/contracts/exec-same-ctx-simple-child
	erdpy contract build ./test/contracts/exec-same-ctx-simple-parent
	
	erdpy contract build ./test/contracts/exec-sync-ctx-multiple/alpha
	erdpy contract build ./test/contracts/exec-sync-ctx-multiple/beta
	erdpy contract build ./test/contracts/exec-sync-ctx-multiple/delta
	erdpy contract build ./test/contracts/exec-sync-ctx-multiple/gamma
	
	erdpy contract build ./test/contracts/init-correct
	erdpy contract build ./test/contracts/init-simple
	erdpy contract build ./test/contracts/init-wrong
	erdpy contract build ./test/contracts/managed-buffers
	erdpy contract build ./test/contracts/misc
	erdpy contract build --no-optimization ./test/contracts/num-with-fp
	erdpy contract build ./test/contracts/promises
	erdpy contract build ./test/contracts/promises-train
	erdpy contract build ./test/contracts/promises-tracking
	erdpy contract build ./test/contracts/signatures
	erdpy contract build ./test/contracts/timelocks
	erdpy contract build ./test/contracts/upgrader-fromanother-contract

build-delegation:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-delegation-rs
	git clone --depth=1 --branch=master https://github.com/ElrondNetwork/sc-delegation-rs.git ${SANDBOX}/sc-delegation-rs
	rm -rf ${SANDBOX}/sc-delegation-rs/.git
	erdpy contract build ${SANDBOX}/sc-delegation-rs
	erdpy contract test --directory="tests" ${SANDBOX}/sc-delegation-rs
	cp ${SANDBOX}/sc-delegation-rs/output/delegation.wasm ./test/delegation/delegation.wasm


build-dns:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-dns-rs
	git clone --depth=1 --branch=master https://github.com/ElrondNetwork/sc-dns-rs.git ${SANDBOX}/sc-dns-rs
	rm -rf ${SANDBOX}/sc-dns-rs/.git
	erdpy contract build ${SANDBOX}/sc-dns-rs
	erdpy contract test --directory="tests" ${SANDBOX}/sc-dns-rs
	cp ${SANDBOX}/sc-dns-rs/output/dns.wasm ./test/dns/dns.wasm


build-sc-examples:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples

	erdpy contract new --template=erc20-c --directory ${SANDBOX}/sc-examples erc20-c
	erdpy contract build ${SANDBOX}/sc-examples/erc20-c
	cp ${SANDBOX}/sc-examples/erc20-c/output/wrc20_arwen.wasm ./test/erc20/contracts/erc20-c.wasm


build-sc-examples-rs:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples-rs
	
	erdpy contract new --template=simple-coin --directory ${SANDBOX}/sc-examples-rs simple-coin
	erdpy contract new --template=adder --directory ${SANDBOX}/sc-examples-rs adder
	erdpy contract build ${SANDBOX}/sc-examples-rs/adder
	erdpy contract build ${SANDBOX}/sc-examples-rs/simple-coin
	erdpy contract test ${SANDBOX}/sc-examples-rs/adder
	erdpy contract test ${SANDBOX}/sc-examples-rs/simple-coin
	cp ${SANDBOX}/sc-examples-rs/adder/output/adder.wasm ./test/adder/adder.wasm
	cp ${SANDBOX}/sc-examples-rs/simple-coin/output/simple-coin.wasm ./test/erc20/contracts/simple-coin.wasm

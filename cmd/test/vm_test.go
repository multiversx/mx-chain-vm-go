package main

import (
	"path/filepath"
	"testing"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

var excludedTests = []string{}

func TestErc20FromC(t *testing.T) {
	testExec := newArwenTestExecutor().replaceCode(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/erc20-c.wasm"))

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		excludedTests,
		testExec)

	if err != nil {
		t.Error(err)
	}
}

func TestSimpleCoinFromRust(t *testing.T) {
	testExec := newArwenTestExecutor().replaceCode(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/simple-coin-rs.wasm"))

	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		excludedTests,
		testExec)

	if err != nil {
		t.Error(err)
	}
}

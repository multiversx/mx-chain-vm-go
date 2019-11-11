package main

import (
	"path/filepath"
	"testing"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

var excludedTests = []string{}

func TestErc20FromC(t *testing.T) {
	erc20Path := filepath.Join(getTestRoot(), "contracts/erc20-c.wasm")
	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		excludedTests,
		newArwenTestExecutor().setErc20Path(erc20Path))

	if err != nil {
		t.Error(err)
	}
}

func TestErc20FromRustDebug(t *testing.T) {
	erc20Path := filepath.Join(getTestRoot(), "contracts/erc20-rust-debug.wasm")
	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		excludedTests,
		newArwenTestExecutor().setErc20Path(erc20Path))

	if err != nil {
		t.Error(err)
	}
}

func TestErc20FromRustRelease(t *testing.T) {
	erc20Path := filepath.Join(getTestRoot(), "contracts/erc20-rust-release.wasm")
	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		excludedTests,
		newArwenTestExecutor().setErc20Path(erc20Path))

	if err != nil {
		t.Error(err)
	}
}

package main

import (
	"path/filepath"
	"testing"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func TestErc20FromC(t *testing.T) {
	fileResolver := ij.NewDefaultFileResolver().ReplacePath(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/erc20-c.wasm"))
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		fileResolver,
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestErc20FromRust(t *testing.T) {
	fileResolver := ij.NewDefaultFileResolver().ReplacePath(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/simple-coin.wasm"))
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		fileResolver,
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestAdderFromRust(t *testing.T) {
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"adder",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestCryptoBubbles(t *testing.T) {
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"crypto_bubbles_min_v1",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestFeaturesFromRust(t *testing.T) {
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"features",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestAsyncCalls(t *testing.T) {
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"async",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDelegationContract(t *testing.T) {
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"delegation",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDnsContract(t *testing.T) {
	runner := controller.NewRunner(
		newArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"dns",
		".json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

package main

import (
	"path/filepath"
	"testing"

	ajt "github.com/ElrondNetwork/arwen-wasm-vm/arwenjsontest"
	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func TestErc20FromC(t *testing.T) {
	fileResolver := ij.NewDefaultFileResolver().ReplacePath(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/erc20-c.wasm"))
	runner := controller.NewTestRunner(
		ajt.NewArwenTestExecutor(),
		fileResolver,
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".test.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestErc20FromRust(t *testing.T) {
	fileResolver := ij.NewDefaultFileResolver().ReplacePath(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/simple-coin.wasm"))
	runner := controller.NewTestRunner(
		ajt.NewArwenTestExecutor(),
		fileResolver,
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".test.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestAdderFromRust(t *testing.T) {
	runner := controller.NewScenarioRunner(
		ajt.NewArwenScenarioExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"adder",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestCryptoBubbles(t *testing.T) {
	runner := controller.NewTestRunner(
		ajt.NewArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"crypto_bubbles_min_v1",
		".test.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestFeaturesFromRust(t *testing.T) {
	runner := controller.NewTestRunner(
		ajt.NewArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"features",
		".test.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestAsyncCalls(t *testing.T) {
	runner := controller.NewTestRunner(
		ajt.NewArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"async",
		".test.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDelegationContract(t *testing.T) {
	runner := controller.NewTestRunner(
		ajt.NewArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"delegation",
		".test.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDnsContract(t *testing.T) {
	runner := controller.NewTestRunner(
		ajt.NewArwenTestExecutor(),
		ij.NewDefaultFileResolver(),
	)
	err := runner.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"dns",
		".test.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

package vmjsonintegrationtest

import (
	"os"
	"path/filepath"
	"testing"

	ajt "github.com/ElrondNetwork/arwen-wasm-vm/arwenjsontest"
	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
	"github.com/stretchr/testify/require"
)

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func TestErc20FromRust(t *testing.T) {
	fileResolver := ij.NewDefaultFileResolver()
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		fileResolver,
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"erc20",
		".scen.json",
		[]string{})
	if err != nil {
		t.Error(err)
	}
}

func TestErc20FromC(t *testing.T) {
	fileResolver := ij.NewDefaultFileResolver().ReplacePath(
		"contracts/simple-coin.wasm",
		filepath.Join(getTestRoot(), "erc20/contracts/erc20-c.wasm"))
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		fileResolver,
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"erc20",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestAdderFromRust(t *testing.T) {
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		ij.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"adder",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestCryptoBubbles(t *testing.T) {
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		ij.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"crypto_bubbles_min_v1",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestFeaturesFromRust(t *testing.T) {
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		ij.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"features",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestAsyncCalls(t *testing.T) {
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		ij.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"async",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDelegationContract(t *testing.T) {
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		ij.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"delegation",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDnsContract(t *testing.T) {
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		ij.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"dns",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestTimelocks(t *testing.T) {
	executor, err := ajt.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := controller.NewScenarioRunner(
		executor,
		ij.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"timelocks",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

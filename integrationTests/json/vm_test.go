package vmjsonintegrationtest

import (
	"os"
	"path/filepath"
	"testing"

	am "github.com/ElrondNetwork/arwen-wasm-vm/arwenmandos"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/test/test-util/mandos/controller"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/stretchr/testify/require"
)

func init() {
	_ = logger.SetLogLevel("*:DEBUG")
}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func TestRustErc20(t *testing.T) {
	fileResolver := mc.NewDefaultFileResolver()
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		fileResolver,
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"erc20-rust/mandos",
		".scen.json",
		[]string{})
	if err != nil {
		t.Error(err)
	}
}

func TestCErc20(t *testing.T) {
	fileResolver := mc.NewDefaultFileResolver()
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		fileResolver,
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"erc20-c",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestRustAdder(t *testing.T) {
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"adder/mandos",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestCryptoBubbles(t *testing.T) {
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"crypto_bubbles_min_v1/mandos",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestRustFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"features/mandos",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestRustFeaturesNoSmallIntApi(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"features-no-small-int-api/mandos",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

// Backwards compatibility.
func TestRustFeaturesLegacy(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"features-legacy/mandos",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestRustAsyncCalls(t *testing.T) {
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"async/mandos",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDelegation_v0_2(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"delegation/v0_2",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDelegation_v0_3(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"delegation/v0_3",
		".scen.json",
		[]string{
			"delegation/v0_3/test/integration/genesis/genesis.scen.json",
		})

	if err != nil {
		t.Error(err)
	}
}

func TestDelegation_v0_4_genesis(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"delegation/v0_4_genesis",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDelegation_v0_5(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"delegation/v0_5",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestDnsContract(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
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
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
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

func TestRustAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"attestation-rust",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

func TestCAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		"attestation-c",
		".scen.json",
		[]string{})

	if err != nil {
		t.Error(err)
	}
}

// func TestPromises(t *testing.T) {
// 	executor, err := am.NewArwenTestExecutor()
// 	require.Nil(t, err)
// 	runner := mc.NewScenarioRunner(
// 		executor,
// 		mc.NewDefaultFileResolver(),
// 	)
// 	err = runner.RunAllJSONScenariosInDirectory(
// 		getTestRoot(),
// 		"promises",
// 		".scen.json",
// 		[]string{})

// 	if err != nil {
// 		t.Error(err)
// 	}
// }

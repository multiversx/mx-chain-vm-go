package vmjsonintegrationtest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	am "github.com/ElrondNetwork/wasm-vm/arwenmandos"
	executorwrapper "github.com/ElrondNetwork/wasm-vm/executor/wrapper"
	mc "github.com/ElrondNetwork/wasm-vm/mandos-go/controller"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
	"github.com/stretchr/testify/require"
)

func init() {
	_ = logger.SetLogLevel("*:NONE")
}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func runAllTestsInFolder(t *testing.T, folder string) {
	runTestsInFolder(t, folder, []string{})
}

func runTestsInFolder(t *testing.T, folder string, exclusions []string) {
	executor, err := am.NewArwenTestExecutor()
	require.Nil(t, err)
	defer executor.Close()

	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)

	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		folder,
		".scen.json",
		exclusions,
		mc.DefaultRunScenarioOptions())

	if err != nil {
		t.Error(err)
	}
}

func runSingleTestReturnError(folder string, filename string) error {
	executor, err := am.NewArwenTestExecutor()
	if err != nil {
		return err
	}
	defer executor.Close()

	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)

	fullPath := path.Join(getTestRoot(), folder)
	fullPath = path.Join(fullPath, filename)

	return runner.RunSingleJSONScenario(
		fullPath,
		mc.DefaultRunScenarioOptions())
}

func runTestsInFolderWithLog(t *testing.T, folder string, exclusions []string) string {
	executor, err := am.NewArwenTestExecutor()
	logger := executorwrapper.NewStringLogger()
	executor.OverrideVMExecutor = executorwrapper.NewWrappedExecutorFactory(logger, wasmer.ExecutorFactory())
	require.Nil(t, err)
	defer executor.Close()

	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)

	err = runner.RunAllJSONScenariosInDirectory(
		getTestRoot(),
		folder,
		".scen.json",
		exclusions,
		mc.DefaultRunScenarioOptions())

	if err != nil {
		t.Error(err)
	}

	return logger.String()
}

func runTestsInFolderCheckLog(t *testing.T, folder string, exclusions []string, expectedLog string) {
	actualLog := runTestsInFolderWithLog(t, "adder/mandos", []string{})
	if actualLog != expectedLog {
		timestampStr := time.Now().Format("2006_01_02_15_04_05")
		fileExpected, err := os.Create(fmt.Sprintf("executorLog_%s_expected.txt", timestampStr))
		require.Nil(t, err)
		fileExpected.WriteString(expectedLog)
		err = fileExpected.Close()
		require.Nil(t, err)
		fileActual, err := os.Create(fmt.Sprintf("executorLog_%s_actual.txt", timestampStr))
		require.Nil(t, err)
		fileActual.WriteString(actualLog)
		err = fileActual.Close()
		require.Nil(t, err)
		t.Error("log mismatch, see saved logs")
	}

}

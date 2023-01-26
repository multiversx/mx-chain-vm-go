package vmjsonintegrationtest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	executorwrapper "github.com/multiversx/mx-chain-vm-go/executor/wrapper"
	am "github.com/multiversx/mx-chain-vm-go/scenarioexec"
	mc "github.com/multiversx/mx-chain-vm-go/scenarios/controller"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
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
	vmTestRoot := filepath.Join(exePath, "../../test")
	return vmTestRoot
}

// ScenariosTestBuilder defines the Scenarios builder component
type ScenariosTestBuilder struct {
	t              *testing.T
	folder         string
	singleFile     string
	exclusions     []string
	executorLogger *executorwrapper.StringLogger
	currentError   error
}

// ScenariosTest will create a new ScenariosTestBuilder instance
func ScenariosTest(t *testing.T) *ScenariosTestBuilder {
	return &ScenariosTestBuilder{
		t:              t,
		folder:         "",
		singleFile:     "",
		executorLogger: nil,
	}
}

// Folder sets the folder
func (mtb *ScenariosTestBuilder) Folder(folder string) *ScenariosTestBuilder {
	mtb.folder = folder
	return mtb
}

// File sets the file
func (mtb *ScenariosTestBuilder) File(fileName string) *ScenariosTestBuilder {
	mtb.singleFile = fileName
	return mtb
}

// Exclude sets the exclusion path
func (mtb *ScenariosTestBuilder) Exclude(path string) *ScenariosTestBuilder {
	mtb.exclusions = append(mtb.exclusions, path)
	return mtb
}

// WithExecutorLogs sets a StringLogger
func (mtb *ScenariosTestBuilder) WithExecutorLogs() *ScenariosTestBuilder {
	mtb.executorLogger = executorwrapper.NewStringLogger()
	return mtb
}

// Run will start the testing process
func (mtb *ScenariosTestBuilder) Run() *ScenariosTestBuilder {
	executor, err := am.NewVMTestExecutor()
	require.Nil(mtb.t, err)
	defer executor.Close()

	if mtb.executorLogger != nil {
		executor.OverrideVMExecutor = executorwrapper.NewWrappedExecutorFactory(
			mtb.executorLogger,
			wasmer.ExecutorFactory())
	}

	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)

	if len(mtb.singleFile) > 0 {
		fullPath := path.Join(getTestRoot(), mtb.folder)
		fullPath = path.Join(fullPath, mtb.singleFile)

		mtb.currentError = runner.RunSingleJSONScenario(
			fullPath,
			mc.DefaultRunScenarioOptions())
	} else {
		mtb.currentError = runner.RunAllJSONScenariosInDirectory(
			getTestRoot(),
			mtb.folder,
			".scen.json",
			mtb.exclusions,
			mc.DefaultRunScenarioOptions())
	}

	return mtb
}

// CheckNoError does an assert for the containing error
func (mtb *ScenariosTestBuilder) CheckNoError() *ScenariosTestBuilder {
	if mtb.currentError != nil {
		mtb.t.Error(mtb.currentError)
	}
	return mtb
}

// RequireError does an assert for the containing error
func (mtb *ScenariosTestBuilder) RequireError(expectedErrorMsg string) *ScenariosTestBuilder {
	require.EqualError(mtb.t, mtb.currentError, expectedErrorMsg)
	return mtb
}

// CheckLog will check the containing error
func (mtb *ScenariosTestBuilder) CheckLog(expectedLogs string) *ScenariosTestBuilder {
	require.NotNil(mtb.t, mtb.executorLogger)
	actualLog := mtb.executorLogger.String()
	if actualLog != expectedLogs {
		timestampStr := time.Now().Format("2006_01_02_15_04_05")
		fileExpected, err := os.Create(fmt.Sprintf("executorLog_%s_expected.txt", timestampStr))
		require.Nil(mtb.t, err)
		_, _ = fileExpected.WriteString(expectedLogs)
		err = fileExpected.Close()
		require.Nil(mtb.t, err)
		fileActual, err := os.Create(fmt.Sprintf("executorLog_%s_actual.txt", timestampStr))
		require.Nil(mtb.t, err)
		_, _ = fileActual.WriteString(actualLog)
		err = fileActual.Close()
		require.Nil(mtb.t, err)
		mtb.t.Error("log mismatch, see saved logs")
	}
	return mtb
}

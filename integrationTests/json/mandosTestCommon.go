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

type MandosTestBuilder struct {
	t              *testing.T
	folder         string
	singleFile     string
	exclusions     []string
	executorLogger *executorwrapper.StringLogger
	currentError   error
}

func MandosTest(t *testing.T) *MandosTestBuilder {
	return &MandosTestBuilder{
		t:              t,
		folder:         "",
		singleFile:     "",
		executorLogger: nil,
	}
}

func (mtb *MandosTestBuilder) Folder(folder string) *MandosTestBuilder {
	mtb.folder = folder
	return mtb
}

func (mtb *MandosTestBuilder) File(fileName string) *MandosTestBuilder {
	mtb.singleFile = fileName
	return mtb
}

func (mtb *MandosTestBuilder) Exclude(path string) *MandosTestBuilder {
	mtb.exclusions = append(mtb.exclusions, path)
	return mtb
}

func (mtb *MandosTestBuilder) WithExecutorLogs() *MandosTestBuilder {
	mtb.executorLogger = executorwrapper.NewStringLogger()
	return mtb
}

func (mtb *MandosTestBuilder) Run() *MandosTestBuilder {
	executor, err := am.NewArwenTestExecutor()
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

func (mtb *MandosTestBuilder) CheckNoError() *MandosTestBuilder {
	if mtb.currentError != nil {
		mtb.t.Error(mtb.currentError)
	}
	return mtb
}

func (mtb *MandosTestBuilder) RequireError(expectedErrorMsg string) *MandosTestBuilder {
	require.EqualError(mtb.t, mtb.currentError, expectedErrorMsg)
	return mtb
}

func (mtb *MandosTestBuilder) CheckLog(expectedLogs string) *MandosTestBuilder {
	require.NotNil(mtb.t, mtb.executorLogger)
	actualLog := mtb.executorLogger.String()
	if actualLog != expectedLogs {
		timestampStr := time.Now().Format("2006_01_02_15_04_05")
		fileExpected, err := os.Create(fmt.Sprintf("executorLog_%s_expected.txt", timestampStr))
		require.Nil(mtb.t, err)
		fileExpected.WriteString(expectedLogs)
		err = fileExpected.Close()
		require.Nil(mtb.t, err)
		fileActual, err := os.Create(fmt.Sprintf("executorLog_%s_actual.txt", timestampStr))
		require.Nil(mtb.t, err)
		fileActual.WriteString(actualLog)
		err = fileActual.Close()
		require.Nil(mtb.t, err)
		mtb.t.Error("log mismatch, see saved logs")
	}
	return mtb
}

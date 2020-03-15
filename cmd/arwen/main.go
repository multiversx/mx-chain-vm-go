package main

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

func main() {
	arguments, err := common.ParseArguments()
	if err != nil {
		exitWithError(fmt.Sprintf("Bad arguments to Arwen: %v", err), common.ErrCodeBadArguments)
	}

	nodeToArwenFile := os.NewFile(3, "/proc/self/fd/3")
	if nodeToArwenFile == nil {
		exitWithError("Cannot create [nodeToArwenFile] file", common.ErrCodeCannotCreateFile)
	}

	arwenToNodeFile := os.NewFile(4, "/proc/self/fd/4")
	if arwenToNodeFile == nil {
		exitWithError("Cannot create [arwenToNodeFile] file", common.ErrCodeCannotCreateFile)
	}

	logToNodeFile := os.NewFile(5, "/proc/self/fd/5")
	if arwenToNodeFile == nil {
		exitWithError("Cannot create [logToNodeFile] file", common.ErrCodeCannotCreateFile)
	}

	arwenLogger := logger.NewPipeLogger(arguments.LogLevel, logToNodeFile)
	part, err := arwenpart.NewArwenPart(
		arwenLogger,
		nodeToArwenFile,
		arwenToNodeFile,
		arguments.VMType,
		arguments.BlockGasLimit,
		arguments.GasSchedule,
	)
	if err != nil {
		exitWithError(fmt.Sprintf("Cannot create ArwenPart: %v", err), common.ErrCodeInit)
	}

	arwenLogger.Info("Arwen.main() start loop")
	err = part.StartLoop()
	if err != nil {
		exitWithError(fmt.Sprintf("Ended Arwen loop: %v", err), common.ErrCodeTerminated)
	}

	arwenLogger.Info("Arwen.main() ended")
}

func exitWithError(errorMessage string, errorCode int) {
	fmt.Fprintln(os.Stderr, errorCode)
	os.Exit(errorCode)
}

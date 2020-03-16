package main

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

const (
	fileDescriptorNodeToArwen = 3
	fileDescriptorArwenToNode = 4
	fileDescriptorLogToNode   = 5
)

func main() {
	errCode, errMessage := doMain()
	if errCode != common.ErrCodeSuccess {
		fmt.Fprintln(os.Stderr, errMessage)
		os.Exit(errCode)
	}
}

// doMain returns (error code, error message)
func doMain() (int, string) {
	arguments, err := common.ParseArguments()
	if err != nil {
		return common.ErrCodeBadArguments, fmt.Sprintf("Bad arguments to Arwen: %v", err)
	}

	nodeToArwenFile := getPipeFile(fileDescriptorNodeToArwen)
	if nodeToArwenFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [nodeToArwenFile]"
	}

	arwenToNodeFile := getPipeFile(fileDescriptorArwenToNode)
	if arwenToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [arwenToNodeFile]"
	}

	logToNodeFile := getPipeFile(fileDescriptorLogToNode)
	if arwenToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [logToNodeFile]"
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
		return common.ErrCodeInit, fmt.Sprintf("Cannot create ArwenPart: %v", err)
	}

	err = part.StartLoop()
	if err != nil {
		return common.ErrCodeTerminated, fmt.Sprintf("Ended Arwen loop: %v", err)
	}

	return common.ErrCodeSuccess, ""
}

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}

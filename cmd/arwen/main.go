package main

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

const (
	fileDescriptorArwenInit   = 3
	fileDescriptorNodeToArwen = 4
	fileDescriptorArwenToNode = 5
	fileDescriptorLogToNode   = 6
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
	arwenInitFile := getPipeFile(fileDescriptorArwenInit)
	if arwenInitFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [arwenInitFile]"
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
	if logToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [logToNodeFile]"
	}

	pipeArguments, err := common.GetArwenArguments(arwenInitFile)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot receive gasSchedule: %v", err)
	}

	arwenLogger := logger.NewPipeLogger(pipeArguments.LogLevel, logToNodeFile)

	part, err := arwenpart.NewArwenPart(
		arwenLogger,
		nodeToArwenFile,
		arwenToNodeFile,
		pipeArguments.VMType,
		pipeArguments.BlockGasLimit,
		pipeArguments.GasSchedule,
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

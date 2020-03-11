package main

import (
	"log"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

func main() {
	vmType, blockGasLimit, gasSchedule, logLevel, err := common.ParseArguments()
	if err != nil {
		log.Fatalf("Bad arguments to Arwen: %v", err)
	}

	nodeToArwenFile := os.NewFile(3, "/proc/self/fd/3")
	if nodeToArwenFile == nil {
		log.Fatal("Cannot create [nodeToArwenFile] file")
	}

	arwenToNodeFile := os.NewFile(4, "/proc/self/fd/4")
	if arwenToNodeFile == nil {
		log.Fatal("Cannot create [arwenToNodeFile] file")
	}

	logToNodeFile := os.NewFile(5, "/proc/self/fd/5")
	if arwenToNodeFile == nil {
		log.Fatal("Cannot create [logToNodeFile] file")
	}

	arwenLogger := logger.NewPipeLogger(logLevel, logToNodeFile)
	part, err := arwenpart.NewArwenPart(arwenLogger, nodeToArwenFile, arwenToNodeFile, vmType, blockGasLimit, gasSchedule)
	if err != nil {
		log.Fatalf("Cannot create ArwenPart: %v", err)
	}

	arwenLogger.Info("Arwen.main() start loop")
	err = part.StartLoop()
	if err != nil {
		log.Fatalf("Ended Arwen loop: %v", err)
	}

	arwenLogger.Info("Arwen.main() ended")
}

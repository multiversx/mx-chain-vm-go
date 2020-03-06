package main

import (
	"log"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

func main() {
	common.LogInfo("Arwen.main()")

	vmType, blockGasLimit, gasSchedule, err := common.ParseArguments()
	if err != nil {
		log.Fatalf("Bad arguments to Arwen: %v", err)
	}

	nodeToArwenFile := os.NewFile(3, "/proc/self/fd/3")
	if nodeToArwenFile == nil {
		log.Fatal("Cannot create file")
	}

	arwenToNodeFile := os.NewFile(4, "/proc/self/fd/4")
	if arwenToNodeFile == nil {
		log.Fatal("Cannot create file")
	}

	part, err := arwenpart.NewArwenPart(nodeToArwenFile, arwenToNodeFile, vmType, blockGasLimit, gasSchedule)
	if err != nil {
		log.Fatalf("Cannot create ArwenPart: %v", err)
	}

	err = part.StartLoop()
	if err != nil {
		log.Fatalf("Ended Arwen loop: %v", err)
	}

	common.LogInfo("Arwen.main() ended")
}

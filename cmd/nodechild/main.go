package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	beginMessageLoop(reader, writer)
}

func beginMessageLoop(reader *bufio.Reader, writer *bufio.Writer) {
	messenger := NewMessenger(reader, writer)
	blockchain := NewBlockchainHookGateway(messenger)
	arwenVirtualMachineType := []byte{5, 0}
	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)

	host, err := host.NewArwenVM(blockchain, nil, arwenVirtualMachineType, blockGasLimit, gasSchedule)
	if err != nil {
		log.Fatal(err)
	}

	for {
		command := messenger.WaitContractCommand()
		fmt.Println("Command", command)
		fmt.Println(host)
	}
}

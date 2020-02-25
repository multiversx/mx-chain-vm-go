package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	err := beginMessageLoop(reader, writer)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func beginMessageLoop(reader *bufio.Reader, writer *bufio.Writer) error {
	messenger := NewMessenger(reader, writer)
	blockchain := NewBlockchainHookGateway(messenger)
	arwenVirtualMachineType := []byte{5, 0}
	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)

	_, err := host.NewArwenVM(blockchain, nil, arwenVirtualMachineType, blockGasLimit, gasSchedule)
	if err != nil {
		log.Fatal(err)
	}

	for {
		command := messenger.WaitContractCommand()
		err := executeCommand(command)
		if errors.Is(err, ErrCriticalError) {
			return err
		}
	}
}

func executeCommand(command *ContractCommand) error {
	fmt.Println("executeCommand()", command)

	switch command.Tag {
	case "Deploy":
		fmt.Println("Deploy smart contract")
	case "Call":
		fmt.Println("Call smart contract")
	default:
		return ErrBadCommandFromNode
	}

	return nil
}

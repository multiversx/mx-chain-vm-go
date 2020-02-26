package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

func main() {
	doMain(os.Stdin, os.Stdout)
}

func doMain(input *os.File, output *os.File) {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	err := beginMessageLoop(reader, writer)
	if err != nil {
		fmt.Println(err)
	}
}

func beginMessageLoop(reader *bufio.Reader, writer *bufio.Writer) error {
	fmt.Println("Arwen: Begin message loop.")
	defer fmt.Println("Arwen: End message loop.")

	messenger := NewChildMessenger(reader, writer)
	blockchain := NewBlockchainHookGateway(messenger)
	arwenVirtualMachineType := []byte{5, 0}
	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)

	_, err := host.NewArwenVM(blockchain, nil, arwenVirtualMachineType, blockGasLimit, gasSchedule)
	if err != nil {
		return err
	}

	var endingError error
	for {
		request, err := messenger.ReceiveContractRequest()
		if err != nil {
			endingError = err
			break
		}

		err = executeRequest(request)
		if errors.Is(err, ErrCriticalError) {
			endingError = err
			break
		}

		fmt.Println("Non critical error:", err)
	}

	messenger.SendResponseIHaveCriticalError(endingError)
	return endingError
}

func executeRequest(request *ContractRequest) error {
	fmt.Println("Arwen: executeRequest()", request)

	switch request.Tag {
	case "Deploy":
		fmt.Println("Deploy smart contract")
	case "Call":
		fmt.Println("Call smart contract")
	default:
		return ErrBadRequestFromNode
	}

	return nil
}

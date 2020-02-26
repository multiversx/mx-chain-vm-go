package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

// Server is
type Server struct {
	Messenger *ChildMessenger
	Host      arwen.VMHost
}

// NewServer creates
func NewServer(input *os.File, output *os.File) (*Server, error) {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	messenger := NewChildMessenger(reader, writer)
	blockchain := NewBlockchainHookGateway(messenger)
	arwenVirtualMachineType := []byte{5, 0}
	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)

	host, err := host.NewArwenVM(blockchain, nil, arwenVirtualMachineType, blockGasLimit, gasSchedule)
	if err != nil {
		return nil, err
	}

	return &Server{
		Messenger: messenger,
		Host:      host,
	}, nil
}

// Start runs the main loop
func (server *Server) Start() error {
	var endingError error
	for {
		request, err := server.Messenger.ReceiveContractRequest()
		if err != nil {
			endingError = err
			break
		}

		err = executeRequest(request)
		if err != nil {
			if errors.Is(err, ErrCriticalError) {
				endingError = err
				break
			} else {
				fmt.Println("Non critical error:", err)
			}
		}
	}

	server.Messenger.SendResponseIHaveCriticalError(endingError)
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

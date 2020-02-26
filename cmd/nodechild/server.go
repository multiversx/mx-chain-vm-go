package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

// Server is
type Server struct {
	Messenger *ChildMessenger
	VMHost    VMHost
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
		VMHost:    host,
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

		response, err := server.executeRequest(request)
		if err != nil {
			if errors.Is(err, ErrCriticalError) {
				endingError = err
				break
			} else {
				fmt.Println("Non critical error:", err)
			}
		}

		// Successful execution, send response
		server.Messenger.SendContractResponse(response)
	}

	server.Messenger.SendResponseIHaveCriticalError(endingError)
	return endingError
}

func (server *Server) executeRequest(request *ContractRequest) (*ContractResponse, error) {
	fmt.Println("Arwen: executeRequest()", request)

	switch request.Tag {
	case "Deploy":
		return server.doRunSmartContractCreate(request), nil
	case "Call":
		fmt.Println("Call smart contract")
	default:
		return nil, ErrBadRequestFromNode
	}

	return nil, nil
}

func (server *Server) doRunSmartContractCreate(request *ContractRequest) *ContractResponse {
	fmt.Println("doRunSmartContractCreate")
	vmOutput, err := server.VMHost.RunSmartContractCreate(nil)

	return &ContractResponse{
		Tag:      request.Tag,
		VMOutput: vmOutput,
		Response: Response{ErrorMessage: err.Error(), HasCriticalError: false},
	}
}

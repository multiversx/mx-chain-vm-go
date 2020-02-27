package arwenpart

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

// ChildServer is
type ChildServer struct {
	Messenger *ChildMessenger
	VMHost    VMHost
}

// NewChildServer creates
func NewChildServer(input *os.File, output *os.File) (*ChildServer, error) {
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

	return &ChildServer{
		Messenger: messenger,
		VMHost:    host,
	}, nil
}

// Start runs the main loop
func (server *ChildServer) Start() error {
	var endingError error
	for {
		request, err := server.Messenger.ReceiveContractRequest()
		if err != nil {
			endingError = err
			break
		}

		response, err := server.executeRequest(request)
		if err != nil {
			if errors.Is(err, common.ErrCriticalError) {
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

func (server *ChildServer) executeRequest(request *common.ContractRequest) (*common.ContractResponse, error) {
	fmt.Println("Arwen: executeRequest()", request)

	switch request.Tag {
	case "Deploy":
		return server.doRunSmartContractCreate(request), nil
	case "Call":
		fmt.Println("Call smart contract")
	case "Stop":
		return nil, common.ErrStopPerNodeRequest
	default:
		return nil, common.ErrBadRequestFromNode
	}

	return nil, nil
}

func (server *ChildServer) doRunSmartContractCreate(request *common.ContractRequest) *common.ContractResponse {
	vmOutput, err := server.VMHost.RunSmartContractCreate(request.CreateInput)

	return &common.ContractResponse{
		Tag:      request.Tag,
		VMOutput: vmOutput,
		Response: common.Response{ErrorMessage: err.Error(), HasCriticalError: false},
	}
}

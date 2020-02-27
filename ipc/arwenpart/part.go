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

// ArwenPart is
type ArwenPart struct {
	Messenger *ChildMessenger
	VMHost    VMHost
}

// NewArwenPart creates
func NewArwenPart(input *os.File, output *os.File) (*ArwenPart, error) {
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

	return &ArwenPart{
		Messenger: messenger,
		VMHost:    host,
	}, nil
}

// StartLoop runs the main loop
func (part *ArwenPart) StartLoop() error {
	var endingError error
	for {
		request, err := part.Messenger.ReceiveContractRequest()
		if err != nil {
			endingError = err
			break
		}

		response, err := part.handleContractRequest(request)
		if err != nil {
			if errors.Is(err, common.ErrCriticalError) {
				endingError = err
				break
			} else {
				fmt.Println("Non critical error:", err)
			}
		}

		// Successful execution, send response
		part.Messenger.SendContractResponse(response)
	}

	part.Messenger.SendResponseIHaveCriticalError(endingError)
	return endingError
}

func (part *ArwenPart) handleContractRequest(request *common.ContractRequest) (*common.HookCallRequestOrContractResponse, error) {
	fmt.Println("Arwen: handleContractRequest()", request)

	switch request.Tag {
	case "Deploy":
		return part.doRunSmartContractCreate(request), nil
	case "Call":
		fmt.Println("Call smart contract")
	case "Stop":
		return nil, common.ErrStopPerNodeRequest
	default:
		return nil, common.ErrBadRequestFromNode
	}

	return nil, nil
}

func (part *ArwenPart) doRunSmartContractCreate(request *common.ContractRequest) *common.HookCallRequestOrContractResponse {
	vmOutput, err := part.VMHost.RunSmartContractCreate(request.CreateInput)
	return common.NewContractResponse(vmOutput, err.Error())
}

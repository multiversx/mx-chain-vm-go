package arwenpart

import (
	"errors"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ArwenPart is
type ArwenPart struct {
	Messenger *ChildMessenger
	VMHost    vmcommon.VMExecutionHandler
}

// NewArwenPart creates
func NewArwenPart(input *os.File, output *os.File) (*ArwenPart, error) {
	messenger := NewChildMessenger(input, output)
	blockchain := NewBlockchainHookGateway(messenger)
	crypto := NewCryptoHookGateway()
	arwenVirtualMachineType := []byte{5, 0}
	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)

	host, err := host.NewArwenVM(blockchain, crypto, arwenVirtualMachineType, blockGasLimit, gasSchedule)
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
				common.LogDebug("Non critical error:", err)
			}
		}

		// Successful execution, send response
		part.Messenger.SendContractResponse(response)
		part.Messenger.Nonce = 0
	}

	part.Messenger.SendResponseIHaveCriticalError(endingError)
	return endingError
}

func (part *ArwenPart) handleContractRequest(request *common.ContractRequest) (*common.HookCallRequestOrContractResponse, error) {
	common.LogDebug("Arwen: handleContractRequest() %v", request)

	switch request.Action {
	case "Deploy":
		return part.doRunSmartContractCreate(request), nil
	case "Call":
		return part.doRunSmartContractCall(request), nil
	case "Stop":
		return nil, common.ErrStopPerNodeRequest
	default:
		return nil, common.ErrBadRequestFromNode
	}

	// TODO: for Deploy and Call, return the actual errors.
}

func (part *ArwenPart) doRunSmartContractCreate(request *common.ContractRequest) *common.HookCallRequestOrContractResponse {
	vmOutput, err := part.VMHost.RunSmartContractCreate(request.CreateInput)
	return common.NewContractResponse(vmOutput, err)
}

func (part *ArwenPart) doRunSmartContractCall(request *common.ContractRequest) *common.HookCallRequestOrContractResponse {
	vmOutput, err := part.VMHost.RunSmartContractCall(request.CallInput)
	return common.NewContractResponse(vmOutput, err)
}

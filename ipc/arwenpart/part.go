package arwenpart

import (
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
	Handlers  []common.MessageCallback
}

// NewArwenPart creates
func NewArwenPart(input *os.File, output *os.File) (*ArwenPart, error) {
	messenger := NewChildMessenger(input, output)
	blockchain := NewBlockchainHookGateway(messenger)
	crypto := NewCryptoHookGateway()
	arwenVirtualMachineType := []byte{5, 0} // TODO
	blockGasLimit := uint64(10000000)       // TODO
	gasSchedule := config.MakeGasMap(1)     // TODO

	host, err := host.NewArwenVM(blockchain, crypto, arwenVirtualMachineType, blockGasLimit, gasSchedule)
	if err != nil {
		return nil, err
	}

	part := &ArwenPart{
		Messenger: messenger,
		VMHost:    host,
	}

	part.Handlers = common.CreateHandlerSlots()
	part.Handlers[common.Stop] = part.handleStop
	part.Handlers[common.ContractDeployRequest] = part.handleRunSmartContractCreate
	part.Handlers[common.ContractCallRequest] = part.handleRunSmartContractCall

	return part, nil
}

// StartLoop runs the main loop
func (part *ArwenPart) StartLoop() error {
	err := part.doLoop()
	part.Messenger.Shutdown()
	common.LogError("Arwen: end of loop, err=%v", err)
	return err
}

// doLoop ends only when a critical failure takes place
func (part *ArwenPart) doLoop() error {
	for {
		request, err := part.Messenger.ReceiveContractRequest()
		if err != nil {
			return err
		}

		response, err := part.handleContractRequest(request)
		if err != nil {
			return err
		}

		// Successful execution, send response
		part.Messenger.SendContractResponse(response)
		part.Messenger.EndDialogue()
	}
}

func (part *ArwenPart) handleContractRequest(request common.MessageHandler) (common.MessageHandler, error) {
	common.LogDebug("Arwen: handleContractRequest() %v", request)
	handler := part.Handlers[request.GetKind()]
	return handler(request)
}

func (part *ArwenPart) handleRunSmartContractCreate(request common.MessageHandler) (common.MessageHandler, error) {
	typedRequest := request.(*common.MessageContractDeployRequest)
	vmOutput, err := part.VMHost.RunSmartContractCreate(typedRequest.CreateInput)
	common.LogDebug("doRunSmartContractCreate, err=%v", err)
	common.LogDebugJSON("VMOutput", vmOutput)
	return common.NewMessageContractResponse(vmOutput, err), nil
}

func (part *ArwenPart) handleRunSmartContractCall(request common.MessageHandler) (common.MessageHandler, error) {
	typedRequest := request.(*common.MessageContractCallRequest)
	vmOutput, err := part.VMHost.RunSmartContractCall(typedRequest.CallInput)
	common.LogDebug("doRunSmartContractCall, err=%v", err)
	return common.NewMessageContractResponse(vmOutput, err), nil
}

func (part *ArwenPart) handleStop(request common.MessageHandler) (common.MessageHandler, error) {
	return nil, common.ErrStopPerNodeRequest
}

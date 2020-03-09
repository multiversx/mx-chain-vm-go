package arwenpart

import (
	"os"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ArwenPart is
type ArwenPart struct {
	Messenger *ChildMessenger
	VMHost    vmcommon.VMExecutionHandler
	Repliers  []common.MessageReplier
}

// NewArwenPart creates
func NewArwenPart(input *os.File, output *os.File, vmType []byte, blockGasLimit uint64, gasSchedule map[string]map[string]uint64) (*ArwenPart, error) {
	messenger := NewChildMessenger(input, output)
	blockchain := NewBlockchainHookGateway(messenger)
	crypto := NewCryptoHookGateway()

	host, err := host.NewArwenVM(blockchain, crypto, vmType, blockGasLimit, gasSchedule)
	if err != nil {
		return nil, err
	}

	part := &ArwenPart{
		Messenger: messenger,
		VMHost:    host,
	}

	part.Repliers = common.CreateReplySlots()
	part.Repliers[common.ContractDeployRequest] = part.replyToRunSmartContractCreate
	part.Repliers[common.ContractCallRequest] = part.replyToRunSmartContractCall
	part.Repliers[common.DiagnoseWaitRequest] = part.replyToDiagnoseWait

	return part, nil
}

// StartLoop runs the main loop
func (part *ArwenPart) StartLoop() error {
	err := part.doLoop()
	part.Messenger.Shutdown()
	common.LogError("[ARWEN]: end of loop, err=%v", err)
	return err
}

// doLoop ends only when a critical failure takes place
func (part *ArwenPart) doLoop() error {
	for {
		request, err := part.Messenger.ReceiveNodeRequest()
		if err != nil {
			return err
		}
		if common.IsStopRequest(request) {
			return common.ErrStopPerNodeRequest
		}

		response := part.replyToNodeRequest(request)

		// Successful execution, send response
		part.Messenger.SendContractResponse(response)
		part.Messenger.EndDialogue()
	}
}

func (part *ArwenPart) replyToNodeRequest(request common.MessageHandler) common.MessageHandler {
	common.LogInfo("[ARWEN]: replyToNodeRequest() %v", request)
	replier := part.Repliers[request.GetKind()]
	return replier(request)
}

func (part *ArwenPart) replyToRunSmartContractCreate(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractDeployRequest)
	vmOutput, err := part.VMHost.RunSmartContractCreate(typedRequest.CreateInput)
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *ArwenPart) replyToRunSmartContractCall(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractCallRequest)
	vmOutput, err := part.VMHost.RunSmartContractCall(typedRequest.CallInput)
	common.LogInfo("[ARWEN]: replyToRunSmartContractCall() done")
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *ArwenPart) replyToDiagnoseWait(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageDiagnoseWaitRequest)
	time.Sleep(time.Duration(typedRequest.Milliseconds) * time.Millisecond)
	return common.NewMessageDiagnoseWaitResponse()
}

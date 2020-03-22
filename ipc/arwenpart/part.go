package arwenpart

import (
	"os"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ArwenPart is the endpoint that implements the message loop on Arwen's side
type ArwenPart struct {
	Messenger *ArwenMessenger
	Logger    logger.Logger
	VMHost    vmcommon.VMExecutionHandler
	Repliers  []common.MessageReplier
}

// NewArwenPart creates the Arwen part
func NewArwenPart(
	mainLogger logger.Logger,
	dialogueLogger logger.Logger,
	input *os.File,
	output *os.File,
	vmHostArguments *common.VMHostArguments,
	marshalizer marshaling.Marshalizer,
) (*ArwenPart, error) {
	messenger := NewArwenMessenger(dialogueLogger, input, output, marshalizer)
	blockchain := NewBlockchainHookGateway(messenger)
	crypto := NewCryptoHookGateway()

	host, err := host.NewArwenVM(blockchain, crypto, vmHostArguments.VMType, vmHostArguments.BlockGasLimit, vmHostArguments.GasSchedule)
	if err != nil {
		return nil, err
	}

	part := &ArwenPart{
		Messenger: messenger,
		Logger:    mainLogger,
		VMHost:    host,
	}

	part.Repliers = common.CreateReplySlots(part.noopReplier)
	part.Repliers[common.ContractDeployRequest] = part.replyToRunSmartContractCreate
	part.Repliers[common.ContractCallRequest] = part.replyToRunSmartContractCall
	part.Repliers[common.DiagnoseWaitRequest] = part.replyToDiagnoseWait

	return part, nil
}

func (part *ArwenPart) noopReplier(message common.MessageHandler) common.MessageHandler {
	part.Logger.Error("noopReplier called")
	return common.CreateMessage(common.UndefinedRequestOrResponse)
}

// StartLoop runs the main loop
func (part *ArwenPart) StartLoop() error {
	part.Messenger.Reset()
	err := part.doLoop()
	part.Messenger.Shutdown()
	part.Logger.Error("[ARWEN]: end of loop", "err", err)
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
		part.Messenger.ResetDialogue()
	}
}

func (part *ArwenPart) replyToNodeRequest(request common.MessageHandler) common.MessageHandler {
	part.Logger.Debug("[ARWEN]: replyToNodeRequest()", "req", request)
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
	part.Logger.Debug("[ARWEN]: replyToRunSmartContractCall() done")
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *ArwenPart) replyToDiagnoseWait(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageDiagnoseWaitRequest)
	duration := time.Duration(int64(typedRequest.Milliseconds) * int64(time.Millisecond))
	time.Sleep(duration)
	return common.NewMessageDiagnoseWaitResponse()
}

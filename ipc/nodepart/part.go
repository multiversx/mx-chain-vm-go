package nodepart

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// NodePart is
type NodePart struct {
	Messenger  *NodeMessenger
	blockchain vmcommon.BlockchainHook
}

// NewNodePart creates
func NewNodePart(input *os.File, output *os.File, blockchain vmcommon.BlockchainHook) (*NodePart, error) {
	messenger := NewNodeMessenger(input, output)

	return &NodePart{
		Messenger:  messenger,
		blockchain: blockchain,
	}, nil
}

// StartLoop runs the main loop
func (part *NodePart) StartLoop(request *common.ContractRequest) (*common.HookCallRequestOrContractResponse, error) {
	part.Messenger.SendContractRequest(request)

	var endingError error
	var message *common.HookCallRequestOrContractResponse

	for {
		// TODO: start with initial timeout, decrement with "time.Since".
		// TODO: Allow a total max of 1 second (accumulated wait).
		message, endingError = part.Messenger.ReceiveHookCallRequestOrContractResponse(1000)
		if endingError != nil {
			message = nil
			break
		} else if message.IsCriticalError() {
			endingError = message.GetError()
			message = nil
			break
		} else if message.IsHookCallRequest() {
			err := part.handleHookCallRequest(message)
			if err != nil {
				endingError = err
				break
			}
		} else if message.IsContractResponse() {
			break
		} else {
			endingError = common.ErrBadMessageFromArwen
			message = nil
			break
		}
	}

	common.LogDebug("Node: {{{End loop}}}.", endingError)
	part.Messenger.Nonce = 0
	return message, endingError
}

func (part *NodePart) handleHookCallRequest(request *common.HookCallRequestOrContractResponse) error {
	hook := request.Hook
	function := request.Function

	common.LogDebug("Node: handleHookCallRequest, %s.%s()", hook, function)

	response := &common.HookCallResponse{}
	var hookError error

	// TODO: Map (Slice, actually) of functions. Enum of functions, with iota.
	if hook == "blockchain" {
		switch function {
		case "NewAddress":
			response.Bytes1, hookError = part.blockchain.NewAddress(request.Bytes1, request.Uint64_1, request.Bytes2)
		case "GetCode":
			response.Bytes1, hookError = part.blockchain.GetCode(request.Bytes1)
		case "GetNonce":
			response.Uint64_1, hookError = part.blockchain.GetNonce(request.Bytes1)
		case "GetStorageData":
			response.Bytes1, hookError = part.blockchain.GetStorageData(request.Bytes1, request.Bytes2)
		default:
			common.LogError("unknown function hook: %s", function)
		}
	} else {
		common.LogError("unknown hook: %s", hook)
	}

	if hookError != nil {
		response.ErrorMessage = hookError.Error()
	}

	err := part.Messenger.SendHookCallResponse(response)
	return err
}

// SendStopSignal sends a stop signal to Arwen
// Should only be used for tests!
func (part *NodePart) SendStopSignal() error {
	request := &common.ContractRequest{
		Action: "Stop",
	}

	err := part.Messenger.SendContractRequest(request)
	if err != nil {
		return err
	}

	common.LogDebug("Node: sent stop signal to Arwen.")
	return nil
}

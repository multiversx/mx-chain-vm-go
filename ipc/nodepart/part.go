package nodepart

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// NodePart is
type NodePart struct {
	Messenger  *NodeMessenger
	blockchain vmcommon.BlockchainHook
	cryptoHook vmcommon.CryptoHook
}

// NewNodePart creates
func NewNodePart(input *os.File, output *os.File, blockchain vmcommon.BlockchainHook, cryptoHook vmcommon.CryptoHook) (*NodePart, error) {
	reader := bufio.NewReaderSize(input, 8096*16)
	writer := bufio.NewWriter(output)

	messenger := NewNodeMessenger(reader, writer)

	return &NodePart{
		Messenger:  messenger,
		blockchain: blockchain,
		cryptoHook: cryptoHook,
	}, nil
}

// StartLoop runs the main loop
func (part *NodePart) StartLoop(request *common.ContractRequest) (*common.HookCallRequestOrContractResponse, error) {
	part.Messenger.SendContractRequest(request)

	var endingError error
	var isCriticalError bool
	var message *common.HookCallRequestOrContractResponse

	for {
		message, endingError = part.Messenger.ReceiveHookCallRequestOrContractResponse()
		if endingError != nil {
			isCriticalError = true
			message = nil
			break
		} else if message.IsCriticalError() {
			endingError = message.GetError()
			isCriticalError = true
			message = nil
			break
		} else if message.IsHookCallRequest() {
			err := part.handleHookCallRequest(message)
			if err != nil {
				endingError = err
				isCriticalError = true
				break
			}
		} else if message.IsContractResponse() {
			break
		} else {
			endingError = common.ErrBadMessageFromArwen
			isCriticalError = true
			message = nil
			break
		}
	}

	// If critical error, node should know that Arwen should be reset / restarted.
	fmt.Println("Node: End loop. IsCriticalError?", isCriticalError)

	return message, endingError
}

func (part *NodePart) handleHookCallRequest(request *common.HookCallRequestOrContractResponse) error {
	hook := request.Hook
	function := request.Function

	fmt.Printf("Node: handleHookCallRequest, %s.%s()\n", hook, function)

	response := &common.HookCallResponse{}
	var hookError error

	if hook == "blockchain" {
		switch function {
		case "NewAddress":
			response.Bytes1, hookError = part.blockchain.NewAddress(request.Bytes1, request.Uint64_1, request.Bytes2)
		case "GetCode":
			response.Bytes1, hookError = part.blockchain.GetCode(request.Bytes1)
		default:
			panic(fmt.Sprintf("unknown function hook: %s", function))
		}
	} else {
		panic(fmt.Sprintf("unknown hook: %s", hook))
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

	fmt.Println("Node: sent stop signal to Arwen.")
	return nil
}

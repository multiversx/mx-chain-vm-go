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
	Blockchain vmcommon.BlockchainHook
}

// NewNodePart creates
func NewNodePart(input *os.File, output *os.File, blockchain vmcommon.BlockchainHook) (*NodePart, error) {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	messenger := NewNodeMessenger(reader, writer)

	return &NodePart{
		Messenger:  messenger,
		Blockchain: blockchain,
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

	if hook == "blockchain" {
		if function == "NewAddress" {
			address, err := part.Blockchain.NewAddress(request.Bytes1, request.Uint64_1, request.Bytes2)
			if err != nil {
				response.ErrorMessage = err.Error()
			}

			response.Bytes1 = address
		}
	} else {
		panic("unknown hook")
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

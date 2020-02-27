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
func NewNodePart(input *os.File, output *os.File) (*NodePart, error) {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	messenger := NewNodeMessenger(reader, writer)

	return &NodePart{
		Messenger:  messenger,
		Blockchain: nil,
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
			part.handleHookCallRequest(message)
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

func (part *NodePart) handleHookCallRequest(request *common.HookCallRequestOrContractResponse) {
	fmt.Println("Node: handleHookCallRequest()", request)
	panic("TODO")
	// execute, send response.
}

package nodepart

import (
	"bufio"
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

// NodeMessenger is
type NodeMessenger struct {
	common.Messenger
}

// NewNodeMessenger creates
func NewNodeMessenger(reader *bufio.Reader, writer *bufio.Writer) *NodeMessenger {
	return &NodeMessenger{
		Messenger: *common.NewMessenger("Node", reader, writer),
	}
}

// SendContractRequest sends
func (messenger *NodeMessenger) SendContractRequest(request *common.ContractRequest) error {
	err := messenger.Send(request)
	if err != nil {
		return common.ErrCannotSendContractRequest
	}

	fmt.Printf("Node: %v sent\n", request)
	return nil
}

// ReceiveHookCallRequestOrContractResponse waits
func (messenger *NodeMessenger) ReceiveHookCallRequestOrContractResponse() (*common.HookCallRequestOrContractResponse, error) {
	message := &common.HookCallRequestOrContractResponse{}

	err := messenger.Receive(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

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
func (messenger *NodeMessenger) SendContractRequest(request *common.ContractRequest) (*common.ContractResponse, error) {
	fmt.Println("Node: Sending contract request...")

	err := messenger.Send(request)
	if err != nil {
		return nil, common.ErrCannotSendContractRequest
	}

	fmt.Println("Node: Request sent, waiting for response...")

	response := &common.ContractResponse{}
	err = messenger.Receive(response)
	if err != nil {
		return nil, err
	}

	if response.HasError() {
		return nil, response.GetError()
	}

	return response, nil
}

package main

import (
	"bufio"
	"fmt"
)

// NodeMessenger is
type NodeMessenger struct {
	Messenger
}

// NewNodeMessenger creates
func NewNodeMessenger(reader *bufio.Reader, writer *bufio.Writer) *NodeMessenger {
	return &NodeMessenger{
		Messenger: *NewMessenger("Node", reader, writer),
	}
}

// SendContractRequest sends
func (messenger *NodeMessenger) SendContractRequest(request *ContractRequest) (*ContractResponse, error) {
	fmt.Println("Node: Sending contract request...")

	err := messenger.send(request)
	if err != nil {
		return nil, ErrCannotSendContractRequest
	}

	fmt.Println("Node: Request sent, waiting for response...")

	response := &ContractResponse{}
	err = messenger.receive(response)
	if err != nil {
		return nil, err
	}
	if response.HasError() {
		return nil, response.GetError()
	}

	return response, nil
}

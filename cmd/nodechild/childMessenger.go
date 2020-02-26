package main

import (
	"bufio"
	"fmt"
	"log"
)

// ChildMessenger is
type ChildMessenger struct {
	Messenger
}

// NewChildMessenger creates
func NewChildMessenger(reader *bufio.Reader, writer *bufio.Writer) *ChildMessenger {
	return &ChildMessenger{
		Messenger: *NewMessenger("Arwen", reader, writer),
	}
}

// ReceiveContractRequest waits
func (messenger *ChildMessenger) ReceiveContractRequest() (*ContractRequest, error) {
	request := &ContractRequest{}

	err := messenger.receive(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// SendHookCallRequest calls
func (messenger *ChildMessenger) SendHookCallRequest(request *HookCallRequest) *HookCallResponse {
	response := &HookCallResponse{}

	err := messenger.send(request)
	if err != nil {
		log.Fatal("SendHookCallRequest: send receive")
	}

	err = messenger.receive(response)
	if err != nil {
		log.Fatal("SendHookCallRequest: cannot receive")
	}

	if response.Tag != request.Tag {
		log.Fatal("SendHookCallRequest: bad tag")
	}

	return response
}

// SendResponseIHaveCriticalError calls
func (messenger *ChildMessenger) SendResponseIHaveCriticalError(endingError error) error {
	fmt.Println("Arwen: Sending end message...")
	err := messenger.send(&Response{ErrorMessage: endingError.Error(), HasCriticalError: true})
	return err
}

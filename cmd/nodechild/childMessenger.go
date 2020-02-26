package main

import (
	"bufio"
	"fmt"
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

// SendContractResponse sends
func (messenger *ChildMessenger) SendContractResponse(response *ContractResponse) error {
	err := messenger.send(response)
	if err != nil {
		return err
	}

	return nil
}

// SendHookCallRequest calls
func (messenger *ChildMessenger) SendHookCallRequest(request *HookCallRequest) (*HookCallResponse, error) {
	response := &HookCallResponse{}

	err := messenger.send(request)
	if err != nil {
		return nil, ErrCannotSendHookCallRequest
	}

	err = messenger.receive(response)
	if err != nil {
		return nil, ErrCannotReceiveHookCallResponse
	}

	if response.Tag != request.Tag {
		return nil, ErrBadResponseTag
	}

	return response, nil
}

// SendResponseIHaveCriticalError calls
func (messenger *ChildMessenger) SendResponseIHaveCriticalError(endingError error) error {
	fmt.Println("Arwen: Sending end message...")
	err := messenger.send(&Response{ErrorMessage: endingError.Error(), HasCriticalError: true})
	return err
}

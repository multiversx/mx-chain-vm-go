package arwenpart

import (
	"bufio"
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

// ChildMessenger is
type ChildMessenger struct {
	common.Messenger
}

// NewChildMessenger creates
func NewChildMessenger(reader *bufio.Reader, writer *bufio.Writer) *ChildMessenger {
	return &ChildMessenger{
		Messenger: *common.NewMessenger("Arwen", reader, writer),
	}
}

// ReceiveContractRequest waits
func (messenger *ChildMessenger) ReceiveContractRequest() (*common.ContractRequest, error) {
	request := &common.ContractRequest{}

	err := messenger.Receive(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// SendContractResponse sends
func (messenger *ChildMessenger) SendContractResponse(response *common.HookCallRequestOrContractResponse) error {
	err := messenger.Send(response)
	if err != nil {
		return err
	}

	return nil
}

// SendHookCallRequest calls
func (messenger *ChildMessenger) SendHookCallRequest(request *common.HookCallRequestOrContractResponse) (*common.HookCallResponse, error) {
	fmt.Printf("%s: CallHook [%s.%s()]\n", messenger.Name, request.Hook, request.Function)

	response := &common.HookCallResponse{}

	err := messenger.Send(request)
	if err != nil {
		return nil, common.ErrCannotSendHookCallRequest
	}

	err = messenger.Receive(response)
	if err != nil {
		return nil, common.ErrCannotReceiveHookCallResponse
	}

	if response.HasError() {
		return nil, response.GetError()
	}

	return response, nil
}

// SendResponseIHaveCriticalError calls
func (messenger *ChildMessenger) SendResponseIHaveCriticalError(endingError error) error {
	fmt.Println("Arwen: Sending end message...")
	err := messenger.Send(common.NewCriticalError(endingError.Error()))
	return err
}

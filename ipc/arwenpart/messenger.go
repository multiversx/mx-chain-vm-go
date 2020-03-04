package arwenpart

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

// ChildMessenger is
type ChildMessenger struct {
	common.Messenger
}

// NewChildMessenger creates
func NewChildMessenger(reader *os.File, writer *os.File) *ChildMessenger {
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
	common.LogDebug("%s: CallHook [%s.%s()]", messenger.Name, request.Hook, request.Function)

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

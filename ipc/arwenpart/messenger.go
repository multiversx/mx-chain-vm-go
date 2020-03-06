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
		Messenger: *common.NewMessenger("ARWEN", reader, writer),
	}
}

// ReceiveNodeRequest waits
func (messenger *ChildMessenger) ReceiveNodeRequest() (common.MessageHandler, error) {
	message, err := messenger.Receive(0)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// SendContractResponse sends
func (messenger *ChildMessenger) SendContractResponse(response common.MessageHandler) error {
	err := messenger.Send(response)
	if err != nil {
		return err
	}

	return nil
}

// SendHookCallRequest calls
func (messenger *ChildMessenger) SendHookCallRequest(request common.MessageHandler) (common.MessageHandler, error) {
	common.LogDebug("[ARWEN]: CallHook %s", request)

	err := messenger.Send(request)
	if err != nil {
		return nil, common.ErrCannotSendHookCallRequest
	}

	response, err := messenger.Receive(0)
	if err != nil {
		return nil, common.ErrCannotReceiveHookCallResponse
	}

	return response, nil
}

package arwenpart

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

// ArwenMessenger is the messenger on Arwen's part of the pipe
type ArwenMessenger struct {
	common.Messenger
}

// NewArwenMessenger creates a new messenger
func NewArwenMessenger(logger logger.Logger, reader *os.File, writer *os.File) *ArwenMessenger {
	return &ArwenMessenger{
		Messenger: *common.NewMessengerPipes("ARWEN", logger, reader, writer),
	}
}

// ReceiveNodeRequest waits for a request from Node
func (messenger *ArwenMessenger) ReceiveNodeRequest() (common.MessageHandler, error) {
	message, err := messenger.Receive(0)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// SendContractResponse sends a contract response to the Node
func (messenger *ArwenMessenger) SendContractResponse(response common.MessageHandler) error {
	err := messenger.Send(response)
	if err != nil {
		return err
	}

	return nil
}

// SendHookCallRequest makes a hook call (over the pipe) and waits for the response
func (messenger *ArwenMessenger) SendHookCallRequest(request common.MessageHandler) (common.MessageHandler, error) {
	messenger.Logger.Trace("[ARWEN]: CallHook", request)

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

package common

import (
	"fmt"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
)

// Messenger intermediates communication (message exchange) via pipes
type Messenger struct {
	Name     string
	Nonce    uint32
	Logger   logger.Logger
	receiver *Receiver
	sender   *Sender
}

// NewMessengerPipes creates a new messenger from pipes
func NewMessengerPipes(name string, logger logger.Logger, reader *os.File, writer *os.File, marshalizer marshaling.Marshalizer) *Messenger {
	return &Messenger{
		Name:     name,
		Logger:   logger,
		receiver: NewReceiver(reader, marshalizer),
		sender:   NewSender(writer, marshalizer),
	}
}

// NewMessenger creates a new messenger
func NewMessenger(name string, logger logger.Logger, receiver *Receiver, sender *Sender) *Messenger {
	return &Messenger{
		Name:     name,
		Logger:   logger,
		receiver: receiver,
		sender:   sender,
	}
}

// Send sends a message over the pipe
func (messenger *Messenger) Send(message MessageHandler) error {
	messenger.Nonce++
	message.SetNonce(messenger.Nonce)
	length, err := messenger.sender.Send(message)
	messenger.Logger.Trace(fmt.Sprintf("[%s][#%d]: SENT message", messenger.Name, message.GetNonce()), "size", length, "msg", message)
	return err
}

// Receive receives a message, reads it from the pipe
func (messenger *Messenger) Receive(timeout int) (MessageHandler, error) {
	messenger.Logger.Trace(fmt.Sprintf("[%s]: Receive message...", messenger.Name))
	message, length, err := messenger.receiver.Receive(timeout)
	if err != nil {
		return nil, err
	}

	messenger.Logger.Trace(fmt.Sprintf("[%s][#%d]: RECEIVED message", messenger.Name, message.GetNonce()), "size", length, "msg", message)
	messageNonce := message.GetNonce()
	if messageNonce != messenger.Nonce+1 {
		return nil, ErrInvalidMessageNonce
	}

	messenger.Nonce = messageNonce
	return message, nil
}

// Reset resets the messenger
func (messenger *Messenger) Reset() {
	messenger.ResetDialogue()
}

// ResetDialogue resets the dialogue nonce
func (messenger *Messenger) ResetDialogue() {
	messenger.Nonce = 0
}

// Shutdown closes the pipes
func (messenger *Messenger) Shutdown() {
	messenger.Logger.Debug("%s:  Messenger::Shutdown", messenger.Name)

	err := messenger.receiver.Shutdown()
	if err != nil {
		messenger.Logger.Error("Cannot close receiver: %v", err)
	}

	err = messenger.sender.Shutdown()
	if err != nil {
		messenger.Logger.Error("Cannot close sender: %v", err)
	}
}

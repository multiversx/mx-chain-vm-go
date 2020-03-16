package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

// Messenger intermediates communication (message exchange) via pipes
type Messenger struct {
	Name   string
	Nonce  uint32
	Logger logger.Logger
	reader *os.File
	writer *os.File
}

// NewMessenger creates a new messenger
func NewMessenger(name string, logger logger.Logger, reader *os.File, writer *os.File) *Messenger {
	return &Messenger{
		Name:   name,
		Logger: logger,
		reader: reader,
		writer: writer,
	}
}

// Send sends a message over the pipe
func (messenger *Messenger) Send(message MessageHandler) error {
	messenger.Nonce++
	message.SetNonce(messenger.Nonce)

	dataBytes, err := marshalMessage(message)
	if err != nil {
		return err
	}

	err = messenger.sendMessageLengthAndKind(len(dataBytes), message.GetKind())
	if err != nil {
		return err
	}

	_, err = messenger.writer.Write(dataBytes)
	if err != nil {
		return err
	}

	messenger.Logger.Trace(fmt.Sprintf("[%s][#%d]: SENT message", messenger.Name, message.GetNonce()), "size", len(dataBytes), "msg", message)
	return err
}

func (messenger *Messenger) sendMessageLengthAndKind(length int, kind MessageKind) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint32(buffer[0:4], uint32(length))
	binary.LittleEndian.PutUint32(buffer[4:8], uint32(kind))
	_, err := messenger.writer.Write(buffer)
	return err
}

// Receive receives a message, reads it from the pipe
func (messenger *Messenger) Receive(timeout int) (MessageHandler, error) {
	messenger.Logger.Trace(fmt.Sprintf("[%s]: Receive message...", messenger.Name))

	if timeout != 0 {
		messenger.setReceiveDeadline(timeout)
		defer messenger.resetReceiveDeadline()
	}

	length, kind, err := messenger.receiveMessageLengthAndKind()
	if err != nil {
		return nil, err
	}

	message, err := messenger.readMessage(kind, length)
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

func (messenger *Messenger) setReceiveDeadline(timeout int) {
	duration := time.Duration(timeout) * time.Millisecond
	future := time.Now().Add(duration)
	messenger.reader.SetDeadline(future)
}

func (messenger *Messenger) resetReceiveDeadline() {
	messenger.reader.SetDeadline(time.Time{})
}

func (messenger *Messenger) receiveMessageLengthAndKind() (int, MessageKind, error) {
	buffer := make([]byte, 8)
	_, err := io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return 0, FirstKind, err
	}

	length := binary.LittleEndian.Uint32(buffer[0:4])
	kind := MessageKind(binary.LittleEndian.Uint32(buffer[4:8]))
	return int(length), kind, nil
}

func (messenger *Messenger) readMessage(kind MessageKind, length int) (MessageHandler, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return nil, err
	}

	message := CreateMessage(kind)
	err = unmarshalMessage(buffer, message)
	if err != nil {
		return nil, err
	}

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

	err := messenger.writer.Close()
	if err != nil {
		messenger.Logger.Error("Cannot close writer: %v", err)
	}

	err = messenger.reader.Close()
	if err != nil {
		messenger.Logger.Error("Cannot close reader: %v", err)
	}
}

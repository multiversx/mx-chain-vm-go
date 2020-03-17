package common

import (
	"encoding/binary"
	"os"
)

// Sender intermediates communication (message sending) via pipes
type Sender struct {
	Name   string
	Nonce  uint32
	writer *os.File
}

// NewSender creates a new sender
func NewSender(writer *os.File) *Sender {
	return &Sender{
		writer: writer,
	}
}

// Send sends a message over the pipe
func (sender *Sender) Send(message MessageHandler) (int, error) {
	dataBytes, err := marshalMessage(message)
	if err != nil {
		return 0, err
	}

	length := len(dataBytes)
	err = sender.sendMessageLengthAndKind(length, message.GetKind())
	if err != nil {
		return 0, err
	}

	_, err = sender.writer.Write(dataBytes)
	if err != nil {
		return 0, err
	}

	return length, err
}

func (sender *Sender) sendMessageLengthAndKind(length int, kind MessageKind) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint32(buffer[0:4], uint32(length))
	binary.LittleEndian.PutUint32(buffer[4:8], uint32(kind))
	_, err := sender.writer.Write(buffer)
	return err
}

// Shutdown closes the pipe
func (sender *Sender) Shutdown() error {
	err := sender.writer.Close()
	return err
}

package common

import (
	"encoding/binary"
	"io"
	"os"
	"time"
)

// Receiver intermediates communication (message receiving) via pipes
type Receiver struct {
	reader *os.File
}

// NewReceiver creates a new receiver
func NewReceiver(reader *os.File) *Receiver {
	return &Receiver{
		reader: reader,
	}
}

// Receive receives a message, reads it from the pipe
func (receiver *Receiver) Receive(timeout int) (MessageHandler, int, error) {
	if timeout != 0 {
		receiver.setReceiveDeadline(timeout)
		defer receiver.resetReceiveDeadline()
	}

	length, kind, err := receiver.receiveMessageLengthAndKind()
	if err != nil {
		return nil, 0, err
	}

	message, err := receiver.readMessage(kind, length)
	if err != nil {
		return nil, 0, err
	}

	return message, length, nil
}

func (receiver *Receiver) setReceiveDeadline(timeout int) {
	duration := time.Duration(timeout) * time.Millisecond
	future := time.Now().Add(duration)
	receiver.reader.SetDeadline(future)
}

func (receiver *Receiver) resetReceiveDeadline() {
	receiver.reader.SetDeadline(time.Time{})
}

func (receiver *Receiver) receiveMessageLengthAndKind() (int, MessageKind, error) {
	buffer := make([]byte, 8)
	_, err := io.ReadFull(receiver.reader, buffer)
	if err != nil {
		return 0, FirstKind, err
	}

	length := binary.LittleEndian.Uint32(buffer[0:4])
	kind := MessageKind(binary.LittleEndian.Uint32(buffer[4:8]))
	return int(length), kind, nil
}

func (receiver *Receiver) readMessage(kind MessageKind, length int) (MessageHandler, error) {
	buffer := make([]byte, length)
	_, err := io.ReadFull(receiver.reader, buffer)
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

// Shutdown closes the pipe
func (receiver *Receiver) Shutdown() error {
	err := receiver.reader.Close()
	return err
}

package common

import (
	"encoding/binary"
	"io"
	"os"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
)

// Receiver intermediates communication (message receiving) via pipes
type Receiver struct {
	reader      *os.File
	marshalizer marshaling.Marshalizer
}

// NewReceiver creates a new receiver
func NewReceiver(reader *os.File, marshalizer marshaling.Marshalizer) *Receiver {
	return &Receiver{
		reader:      reader,
		marshalizer: marshalizer,
	}
}

// Receive receives a message, reads it from the pipe
func (receiver *Receiver) Receive(timeout int) (MessageHandler, int, error) {
	if timeout > 0 {
		err := receiver.setReceiveDeadline(timeout)
		if err != nil {
			return nil, 0, err
		}

		defer receiver.resetReceiveDeadlineQuietly()
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

func (receiver *Receiver) setReceiveDeadline(timeout int) error {
	duration := time.Duration(timeout) * time.Millisecond
	future := time.Now().Add(duration)
	return receiver.reader.SetDeadline(future)
}

func (receiver *Receiver) resetReceiveDeadlineQuietly() {
	_ = receiver.reader.SetDeadline(time.Time{})
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
	err = receiver.marshalizer.Unmarshal(message, buffer)
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

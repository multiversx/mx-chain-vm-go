package common

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
	"time"
)

// Messenger intermediates communication (message exchange) via pipes
type Messenger struct {
	Name   string
	Nonce  uint32
	reader *os.File
	writer *os.File
}

// NewMessenger creates a new messenger
func NewMessenger(name string, reader *os.File, writer *os.File) *Messenger {
	return &Messenger{
		Name:   name,
		reader: reader,
		writer: writer,
	}
}

// Send sends a message over the pipe
func (messenger *Messenger) Send(message MessageHandler) error {
	messenger.Nonce++
	message.SetNonce(messenger.Nonce)

	dataBytes, err := messenger.marshal(message)
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

	LogDebug("[%s][#%d]: SENT message of size %d %s", messenger.Name, message.GetNonce(), len(dataBytes), message)
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
	LogDebug("[%s]: Receive message...", messenger.Name)

	if timeout != 0 {
		messenger.setReceiveDeadline(timeout)
		defer messenger.resetReceiveDeadline()
	}

	length, kind, err := messenger.receiveMessageLengthAndKind()
	if err != nil {
		return nil, err
	}

	message := CreateMessage(kind)

	// Now read the body
	buffer := make([]byte, length)
	_, err = io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return nil, err
	}

	err = messenger.unmarshal(buffer, message)
	if err != nil {
		return nil, err
	}

	LogDebug("[%s][#%d]: RECEIVED message of size %d %s", messenger.Name, message.GetNonce(), length, message)
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

// Shutdown closes the pipes
func (messenger *Messenger) Shutdown() {
	LogDebug("%s:  Messenger::Shutdown", messenger.Name)

	err := messenger.writer.Close()
	if err != nil {
		LogError("Cannot close writer: %v", err)
	}

	err = messenger.reader.Close()
	if err != nil {
		LogError("Cannot close reader: %v", err)
	}
}

// EndDialogue resets the dialogue nonce
func (messenger *Messenger) EndDialogue() {
	messenger.Nonce = 0
}

func (messenger *Messenger) marshal(data interface{}) ([]byte, error) {
	return marshalJSON(data)
}

func (messenger *Messenger) unmarshal(dataBytes []byte, data interface{}) error {
	return unmarshalJSON(dataBytes, data)
}

func marshalGob(data interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func unmarshalGob(dataBytes []byte, data interface{}) error {
	buffer := bytes.NewBuffer(dataBytes)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func marshalJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func unmarshalJSON(dataBytes []byte, data interface{}) error {
	return json.Unmarshal(dataBytes, data)
}

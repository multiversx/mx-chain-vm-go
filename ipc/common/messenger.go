package common

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
)

// Messenger is
type Messenger struct {
	Name   string
	Nonce  uint32
	reader *os.File
	writer *os.File
}

// NewMessenger creates
func NewMessenger(name string, reader *os.File, writer *os.File) *Messenger {
	return &Messenger{
		Name:   name,
		reader: reader,
		writer: writer,
	}
}

// Send sends
func (messenger *Messenger) Send(message Message) error {
	messenger.Nonce++
	message.SetNonce(messenger.Nonce)

	dataBytes, err := messenger.marshal(message)
	if err != nil {
		return err
	}

	err = messenger.sendMessageLength(dataBytes)
	if err != nil {
		return err
	}

	_, err = messenger.writer.Write(dataBytes)
	if err != nil {
		return err
	}

	LogDebug("[MSG %d] %s: SENT message of size %d", message.GetNonce(), messenger.Name, len(dataBytes))
	return err
}

func (messenger *Messenger) sendMessageLength(marshalizedMessage []byte) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(len(marshalizedMessage)))
	_, err := messenger.writer.Write(buffer)
	return err
}

// Receive receives
func (messenger *Messenger) Receive(message Message) error {
	LogDebug("%s: Receive message...", messenger.Name)

	length, err := messenger.receiveMessageLength()
	if err != nil {
		return err
	}

	// Now read the body of [length]
	buffer := make([]byte, length)
	_, err = io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return err
	}

	err = messenger.unmarshal(buffer, message)
	if err != nil {
		return err
	}

	LogDebug("[MSG %d] %s: RECEIVED message of size %d\n", message.GetNonce(), messenger.Name, length)
	messageNonce := message.GetNonce()
	if messageNonce != messenger.Nonce+1 {
		return ErrInvalidMessageNonce
	}

	messenger.Nonce = messageNonce
	return nil
}

func (messenger *Messenger) receiveMessageLength() (int, error) {
	buffer := make([]byte, 4)
	_, err := io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return 0, err
	}

	length := binary.LittleEndian.Uint32(buffer)
	return int(length), nil
}

// Shutdown does
func (messenger *Messenger) Shutdown() {
	LogDebug("%s:  Messenger:Shutdown", messenger.Name)
	err := messenger.writer.Close()
	if err != nil {
		LogError("Cannot close writer: %v", err)
	}

	err = messenger.reader.Close()
	if err != nil {
		LogError("Cannot close reader: %v", err)
	}
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

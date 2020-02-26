package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// Messenger is
type Messenger struct {
	name   string
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewMessenger creates
func NewMessenger(name string, reader *bufio.Reader, writer *bufio.Writer) *Messenger {
	return &Messenger{
		name:   name,
		reader: reader,
		writer: writer,
	}
}

func (messenger *Messenger) send(message interface{}) error {
	jsonData, err := messenger.marshal(message)
	if err != nil {
		return err
	}

	err = messenger.sendMessageLength(jsonData)
	if err != nil {
		return err
	}

	fmt.Printf("%s: Send: %s\n", messenger.name, jsonData)
	_, err = messenger.writer.Write(jsonData)
	if err != nil {
		return err
	}

	err = messenger.writer.Flush()
	return err
}

func (messenger *Messenger) sendMessageLength(marshalizedMessage []byte) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(len(marshalizedMessage)))
	_, err := messenger.writer.Write(buffer)
	return err
}

func (messenger *Messenger) receive(message interface{}) error {
	// Wait for the start of a message
	messenger.blockingPeek(4)

	length, err := messenger.receiveMessageLength()
	if err != nil {
		return err
	}

	// Now read the body of [length]
	messenger.blockingPeek(length)
	buffer := make([]byte, length)
	_, err = io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return err
	}

	fmt.Printf("%s: Received: %s\n", messenger.name, string(buffer))
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

func (messenger *Messenger) blockingPeek(n int) {
	fmt.Printf("%s: blockingPeek %d bytes\n", messenger.name, n)
	for {
		_, err := messenger.reader.Peek(n)
		if err == nil {
			break
		}
	}
	fmt.Printf("%s: peeked %d bytes\n", messenger.name, n)
}

func (messenger *Messenger) marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (messenger *Messenger) unmarshal(jsonData []byte, data interface{}) error {
	return json.Unmarshal(jsonData, data)
}

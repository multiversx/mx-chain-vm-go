package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// ChildMessenger is
type ChildMessenger struct {
	Messenger
}

// NewChildMessenger creates
func NewChildMessenger(reader *bufio.Reader, writer *bufio.Writer) *ChildMessenger {
	return &ChildMessenger{
		Messenger: *NewMessenger("Arwen", reader, writer),
	}
}

// ReceiveContractRequest waits
func (messenger *ChildMessenger) ReceiveContractRequest() (*ContractRequest, error) {
	request := &ContractRequest{}

	err := messenger.receive(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// SendHookCallRequest calls
func (messenger *ChildMessenger) SendHookCallRequest(request *HookCallRequest) *HookCallResponse {
	response := &HookCallResponse{}

	err := messenger.send(request)
	if err != nil {
		log.Fatal("SendHookCallRequest: send receive")
	}

	err = messenger.receive(response)
	if err != nil {
		log.Fatal("SendHookCallRequest: cannot receive")
	}

	if response.Tag != request.Tag {
		log.Fatal("SendHookCallRequest: bad tag")
	}

	return response
}

// SendResponseIHaveCriticalError calls
func (messenger *ChildMessenger) SendResponseIHaveCriticalError(endingError error) error {
	fmt.Println("Arwen: Sending end message...")
	err := messenger.send(&Response{ErrorMessage: endingError.Error(), HasCriticalError: true})
	return err
}

// NodeMessenger is
type NodeMessenger struct {
	Messenger
}

// NewNodeMessenger creates
func NewNodeMessenger(reader *bufio.Reader, writer *bufio.Writer) *NodeMessenger {
	return &NodeMessenger{
		Messenger: *NewMessenger("Node", reader, writer),
	}
}

// SendContractRequest sends
func (messenger *NodeMessenger) SendContractRequest(request *ContractRequest) (*ContractResponse, error) {
	fmt.Println("Node: Sending contract request...")

	err := messenger.send(request)
	if err != nil {
		return nil, ErrCannotSendContractRequest
	}

	fmt.Println("Node: Request sent, waiting for response...")

	response := &ContractResponse{}
	err = messenger.receive(response)
	if err != nil {
		return nil, err
	}
	if response.HasError() {
		return nil, response.GetError()
	}

	return response, nil
}

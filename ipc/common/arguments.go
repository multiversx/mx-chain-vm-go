package common

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
)

// ArwenArguments represents the initialization arguments required by Arwen, passed through the initialization pipe
type ArwenArguments struct {
	VMType              []byte
	BlockGasLimit       uint64
	LogLevel            logger.LogLevel
	GasSchedule         GasScheduleMap
	LogsMarshalizer     marshaling.MarshalizerKind
	MessagesMarshalizer marshaling.MarshalizerKind
}

// GasScheduleMap is an alias
type GasScheduleMap = map[string]map[string]uint64

// SendArwenArguments sends initialization arguments through a pipe
// For the arguments, the marshalizer is hardcoded, JSON
func SendArwenArguments(pipe *os.File, pipeArguments ArwenArguments) error {
	sender := NewSender(pipe, marshaling.CreateMarshalizer(marshaling.JSON))
	message := NewMessageInitialize(pipeArguments)
	_, err := sender.Send(message)
	return err
}

// GetArwenArguments reads initialization arguments from the pipe
// For the arguments, the marshalizer is hardcoded, JSON
func GetArwenArguments(pipe *os.File) (*ArwenArguments, error) {
	receiver := NewReceiver(pipe, marshaling.CreateMarshalizer(marshaling.JSON))
	message, _, err := receiver.Receive(0)
	if err != nil {
		return nil, err
	}

	typedMessage := message.(*MessageInitialize)
	return &typedMessage.Arguments, nil
}

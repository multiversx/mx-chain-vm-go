package common

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
)

// ArwenArguments represents the initialization arguments required by Arwen, passed through the initialization pipe
type ArwenArguments struct {
	VMHostArguments
	LogLevel            logger.LogLevel
	LogsMarshalizer     marshaling.MarshalizerKind
	MessagesMarshalizer marshaling.MarshalizerKind
}

// VMHostArguments represents the arguments to be passed to VMHost
type VMHostArguments struct {
	VMType              []byte
	BlockGasLimit       uint64
	GasSchedule         GasScheduleMap
}

// GasScheduleMap is an alias
type GasScheduleMap = map[string]map[string]uint64

// SendArwenArguments sends initialization arguments through a pipe
func SendArwenArguments(pipe *os.File, pipeArguments ArwenArguments) error {
	sender := NewSender(pipe, createArgumentsMarshalizer())
	message := NewMessageInitialize(pipeArguments)
	_, err := sender.Send(message)
	return err
}

// GetArwenArguments reads initialization arguments from the pipe
func GetArwenArguments(pipe *os.File) (*ArwenArguments, error) {
	receiver := NewReceiver(pipe, createArgumentsMarshalizer())
	message, _, err := receiver.Receive(0)
	if err != nil {
		return nil, err
	}

	typedMessage := message.(*MessageInitialize)
	return &typedMessage.Arguments, nil
}

// For the arguments, the marshalizer is fixed to JSON
func createArgumentsMarshalizer() marshaling.Marshalizer {
	return marshaling.CreateMarshalizer(marshaling.JSON)
}

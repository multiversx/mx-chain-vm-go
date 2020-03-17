package common

import (
	"encoding/hex"
	"os"
	"strconv"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

// NumCommandLineArguments is the number of arguments required by the CLI
const NumCommandLineArguments = 4

// CommandLineArguments represents the command-line arguments required by Arwen
type CommandLineArguments struct {
	VMType        []byte
	BlockGasLimit uint64
	LogLevel      logger.LogLevel
}

// PipeArguments represents the initialization arguments required by Arwen, passed through the initialization pipe
type PipeArguments struct {
	GasSchedule GasScheduleMap
}

// GasScheduleMap is an alias
type GasScheduleMap = map[string]map[string]uint64

// PrepareCommandLineArguments prepares the list of arguments (command line) to be sent by the Node to Arwen when Arwen should be started
func PrepareCommandLineArguments(arguments CommandLineArguments) ([]string, error) {
	stringArguments := []string{
		hex.EncodeToString(arguments.VMType),
		strconv.FormatUint(arguments.BlockGasLimit, 10),
		strconv.FormatUint(uint64(arguments.LogLevel), 10),
	}

	return stringArguments, nil
}

// ParseCommandLineArguments parses the arguments (command line) received by Arwen from the Node
func ParseCommandLineArguments() (*CommandLineArguments, error) {
	arguments := os.Args
	if len(arguments) != NumCommandLineArguments {
		return nil, ErrBadArwenArguments
	}

	vmType, err := hex.DecodeString(arguments[1])
	if err != nil {
		return nil, err
	}

	blockGasLimit, err := strconv.ParseUint(arguments[2], 10, 64)
	if err != nil {
		return nil, err
	}

	logLevelUint, err := strconv.ParseUint(arguments[3], 10, 8)
	if err != nil {
		return nil, err
	}

	logLevel := logger.LogLevel(logLevelUint)

	return &CommandLineArguments{
		VMType:        vmType,
		BlockGasLimit: blockGasLimit,
		LogLevel:      logLevel,
	}, nil
}

// SendPipeArguments sends initialization arguments through a pipe
func SendPipeArguments(pipe *os.File, pipeArguments PipeArguments) error {
	sender := NewSender(pipe)
	message := NewMessageInitialize(pipeArguments)
	_, err := sender.Send(message)
	return err
}

// GetPipeArguments reads initialization arguments from the pipe
func GetPipeArguments(pipe *os.File) (*PipeArguments, error) {
	receiver := NewReceiver(pipe)
	message, _, err := receiver.Receive(0)
	if err != nil {
		return nil, err
	}

	typedMessage := message.(*MessageInitialize)
	return &typedMessage.Arguments, nil
}

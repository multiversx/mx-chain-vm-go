package common

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

const NumArguments = 5

// Arguments represents the command-line arguments required by Arwen
type Arguments struct {
	VMType        []byte
	BlockGasLimit uint64
	GasSchedule   map[string]map[string]uint64
	LogLevel      logger.LogLevel
}

// PrepareArguments prepares the list of arguments (command line) to be sent by the Node to Arwen when the latter should be started
func PrepareArguments(arguments Arguments) ([]string, error) {
	file, err := ioutil.TempFile("", "gasScheduleToArwen")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := json.Marshal(arguments.GasSchedule)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	stringArguments := []string{
		hex.EncodeToString(arguments.VMType),
		strconv.FormatUint(arguments.BlockGasLimit, 10),
		file.Name(),
		strconv.FormatUint(uint64(arguments.LogLevel), 10),
	}

	return stringArguments, nil
}

// ParseArguments parses the arguments (command line) received by Arwen from the Node
func ParseArguments() (*Arguments, error) {
	arguments := os.Args
	if len(arguments) != NumArguments {
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

	gasSchedule := make(map[string]map[string]uint64)
	gasSchedulePath := arguments[3]
	gasScheduleBytes, err := ioutil.ReadFile(gasSchedulePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(gasScheduleBytes, &gasSchedule)
	if err != nil {
		return nil, err
	}

	err = os.Remove(gasSchedulePath)
	if err != nil {
		return nil, err
	}

	logLevelUint, err := strconv.ParseUint(arguments[4], 10, 8)
	if err != nil {
		return nil, err
	}

	logLevel := logger.LogLevel(logLevelUint)

	return &Arguments{
		VMType:        vmType,
		BlockGasLimit: blockGasLimit,
		GasSchedule:   gasSchedule,
		LogLevel:      logLevel,
	}, nil
}

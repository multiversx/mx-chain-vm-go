package common

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
)

// PrepareArguments prepares the list of arguments (command line) to be sent by the Node to Arwen when the latter should be started
func PrepareArguments(vmType []byte, blockGasLimit uint64, gasSchedule map[string]map[string]uint64, logLevel logger.LogLevel) ([]string, error) {
	file, err := ioutil.TempFile("", "gasScheduleToArwen")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := json.Marshal(gasSchedule)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	arguments := []string{
		hex.EncodeToString(vmType),
		strconv.FormatUint(blockGasLimit, 10),
		file.Name(),
		strconv.FormatUint(uint64(logLevel), 10),
	}

	return arguments, nil
}

// ParseArguments parses the arguments (command line) received by Arwen from the Node
func ParseArguments() (vmType []byte, blockGasLimit uint64, gasSchedule map[string]map[string]uint64, logLevel logger.LogLevel, err error) {
	arguments := os.Args
	if len(arguments) != 5 {
		return nil, 0, nil, 0, ErrBadArwenArguments
	}

	vmType, err = hex.DecodeString(arguments[1])
	if err != nil {
		return
	}

	blockGasLimit, err = strconv.ParseUint(arguments[2], 10, 64)
	if err != nil {
		return
	}

	gasSchedule = make(map[string]map[string]uint64)
	gasSchedulePath := arguments[3]
	gasScheduleBytes, err := ioutil.ReadFile(gasSchedulePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(gasScheduleBytes, &gasSchedule)
	if err != nil {
		return
	}

	err = os.Remove(gasSchedulePath)
	if err != nil {
		return
	}

	logLevelUint, err := strconv.ParseUint(arguments[4], 10, 8)
	if err != nil {
		return
	}

	logLevel = logger.LogLevel(logLevelUint)

	return
}

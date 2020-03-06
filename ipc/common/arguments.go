package common

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

// PrepareArguments prepares
func PrepareArguments(vmType []byte, blockGasLimit uint64, gasSchedule map[string]map[string]uint64) ([]string, error) {
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
	}

	return arguments, nil
}

// ParseArguments parses
func ParseArguments() (vmType []byte, blockGasLimit uint64, gasSchedule map[string]map[string]uint64, err error) {
	arguments := os.Args
	if len(arguments) != 4 {
		return nil, 0, nil, ErrBadArwenArguments
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

	errRemoveTemp := os.Remove(gasSchedulePath)
	if errRemoveTemp != nil {
		LogError("Could not remoce temporary file: %v", errRemoveTemp)
	}

	return
}

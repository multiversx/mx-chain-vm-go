package common

import (
	"encoding/hex"
	"os"
	"strconv"
)

// PrepareArguments prepares
func PrepareArguments(vmType []byte, blockGasLimit uint64) []string {
	arguments := []string{
		hex.EncodeToString(vmType),
		strconv.FormatUint(blockGasLimit, 10),
	}

	return arguments
}

// ParseArguments parses
func ParseArguments() (vmType []byte, blockGasLimit uint64, err error) {
	arguments := os.Args
	if len(arguments) != 3 {
		return nil, 0, ErrBadArwenArguments
	}

	vmType, err = hex.DecodeString(arguments[1])
	if err != nil {
		return
	}

	blockGasLimit, err = strconv.ParseUint(arguments[2], 10, 64)
	if err != nil {
		return
	}

	return
}

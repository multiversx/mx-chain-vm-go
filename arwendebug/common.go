package arwendebug

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
)

func fixTestAddress(address string) []byte {
	if len(address) > arwen.AddressLen {
		address = address[0:arwen.AddressLen]
	}

	address = fmt.Sprintf("%032s", address)
	return []byte(address)
}

func decodeArguments(arguments []string) ([][]byte, error) {
	result := make([][]byte, len(arguments))

	for i := 0; i < len(arguments); i++ {
		decoded, err := hex.DecodeString(arguments[i])
		if err != nil {
			return nil, ErrInvalidArgumentEncoding
		}

		result[i] = decoded
	}

	return result, nil
}

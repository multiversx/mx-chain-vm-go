package arwendebug

import (
	"encoding/hex"
)

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

package vmserver

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
)

func decodeArguments(arguments []string) ([][]byte, error) {
	result := make([][]byte, len(arguments))

	for i := 0; i < len(arguments); i++ {
		decoded, err := fromHex(arguments[i])
		if err != nil {
			return nil, ErrInvalidArgumentEncoding
		}

		result[i] = decoded
	}

	return result, nil
}

func parseValue(value string) (*big.Int, error) {
	valueAsBigInt := big.NewInt(0)

	if len(value) > 0 {
		_, ok := valueAsBigInt.SetString(value, 10)
		if !ok {
			return nil, NewRequestError("invalid value")
		}
	}

	return valueAsBigInt, nil
}

func prettyJson(request interface{}) string {
	data, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		log.Error("prettyJson", "err", err)
	}

	return string(data)
}

func toHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func fromHex(encoded string) ([]byte, error) {
	return hex.DecodeString(encoded)
}

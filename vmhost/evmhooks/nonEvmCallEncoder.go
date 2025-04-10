package evmhooks

import (
	"encoding/binary"
)

const lengthPrefixSize = 4

func lengthPrefixEncode(arguments [][]byte) []byte {
	totalLength := 0
	for _, argument := range arguments {
		totalLength += lengthPrefixSize + len(argument)
	}

	currentOffset := 0
	encoded := make([]byte, totalLength)
	for _, argument := range arguments {
		argumentStart := currentOffset + lengthPrefixSize
		binary.BigEndian.PutUint32(encoded[currentOffset:argumentStart], uint32(len(argument)))
		copy(encoded[argumentStart:], argument)
		currentOffset = argumentStart + len(argument)
	}
	return encoded
}

func lengthPrefixDecode(encoded []byte) ([][]byte, error) {
	currentOffset := 0
	var arguments [][]byte
	for currentOffset < len(encoded) {
		argumentStart := currentOffset + lengthPrefixSize
		if argumentStart > len(encoded) {
			return nil, ErrInvalidEncodedData
		}

		length := binary.BigEndian.Uint32(encoded[currentOffset:argumentStart])
		argumentEnd := argumentStart + int(length)
		if argumentEnd > len(encoded) {
			return nil, ErrInvalidEncodedData
		}

		currentArgument := encoded[argumentStart:argumentEnd]
		arguments = append(arguments, currentArgument)

		currentOffset = argumentEnd
	}
	return arguments, nil
}

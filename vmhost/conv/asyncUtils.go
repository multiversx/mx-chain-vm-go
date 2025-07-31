package conv

import (
	"encoding/binary"
)

func EncodeGasLimits(gasLimits []uint64) []byte {
	if len(gasLimits) == 0 {
		return nil
	}

	// 1 byte for length, then 8 bytes for each uint64
	buf := make([]byte, 1+8*len(gasLimits))
	buf[0] = byte(len(gasLimits))
	for i, limit := range gasLimits {
		binary.LittleEndian.PutUint64(buf[1+8*i:], limit)
	}
	return buf
}

func DecodeGasLimits(data []byte) ([]uint64, []byte) {
	if len(data) == 0 {
		return nil, data
	}

	numLimits := int(data[0])
	if len(data) < 1+8*numLimits {
		// Not enough data, assume no gas limits encoded
		return nil, data
	}

	gasLimits := make([]uint64, numLimits)
	for i := 0; i < numLimits; i++ {
		gasLimits[i] = binary.LittleEndian.Uint64(data[1+8*i:])
	}

	return gasLimits, data[1+8*numLimits:]
}

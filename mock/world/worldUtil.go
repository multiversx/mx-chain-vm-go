package worldmock

import "github.com/ElrondNetwork/elrond-go-core/hashing/blake2b"

// DefaultHasher is an exposed value to use in tests
var DefaultHasher = blake2b.NewBlake2b()

// GenerateMockAddress simulates creation of a new address by the protocol.
func GenerateMockAddress(creatorAddress []byte, creatorNonce uint64) []byte {
	result := make([]byte, 32)
	result[10] = 0x11
	result[11] = 0x11
	result[12] = 0x11
	result[13] = 0x11
	copy(result[14:29], creatorAddress)

	result[29] = byte(creatorNonce)

	copy(result[30:], creatorAddress[30:])
	return result
}

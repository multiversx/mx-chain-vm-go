package worldmock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// DefaultHasher is an exposed value to use in tests
var DefaultHasher = blake2b.NewBlake2b()

// DefaultVMType is an exposed value to use in tests
var DefaultVMType = []byte{0xF, 0xF}

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
	copy(result[vmcommon.NumInitCharactersForScAddress-core.VMTypeLen:], DefaultVMType)
	return result
}

package mandosexpressioninterpreter

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ElrondNetwork/elrond-go-core/core"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"golang.org/x/crypto/sha3"
)

// SCAddressNumLeadingZeros is the number of zero bytes every smart contract address begins with.
// Its value is 10.
// 10 = 8 zeros for all SC addresses + 2 zeros as placeholder for the VM type.
const SCAddressNumLeadingZeros = vmcommon.NumInitCharactersForScAddress

// Keccak256 cryptographic function
// TODO: externalize the same way as the file resolver
func Keccak256(data []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	result := hash.Sum(nil)
	return result, nil
}

func decodeShardId(shardIdRaw string) (byte, error) {
	shardId, err := hex.DecodeString(shardIdRaw)
	if err != nil {
		return 0, fmt.Errorf("could not parse address shard id: %w", err)
	}
	if len(shardId) != 1 {
		return 0, fmt.Errorf("bad address shard id length: %s", shardIdRaw)
	}
	return shardId[0], nil
}

func createAddressFromPrefix(prefix []byte, startIndex, endIndex int) *[32]byte {
	var result [32]byte
	for i := 0; i < len(prefix) && i < endIndex-startIndex; i++ {
		result[i+startIndex] = prefix[i]
	}
	for i := len(prefix) + startIndex; i < endIndex; i++ {
		result[i] = byte('_')
	}
	return &result
}

func createAddressOptionalShardId(input string, numLeadingZeros int) ([]byte, error) {
	tokens := strings.Split(input, "#")
	switch len(tokens) {
	case 1:
		address := createAddressFromPrefix([]byte(tokens[0]), numLeadingZeros, 32)
		return address[:], nil
	case 2:
		shardId, err := decodeShardId(tokens[1])
		if err != nil {
			return []byte{}, err
		}
		address := createAddressFromPrefix([]byte(tokens[0]), numLeadingZeros, 32)
		address[31] = shardId
		return address[:], nil
	default:
		return []byte{}, fmt.Errorf("only one shard id separator allowed in address expression. Got: `%s`", input)
	}
}

// Generates a 32-byte EOA address based on the input.
func addressExpression(input string) ([]byte, error) {
	return createAddressOptionalShardId(input, 0)
}

// Generates a 32-byte smart contract address based on the input.
func (ei *ExprInterpreter) scExpression(input string) ([]byte, error) {
	address, err := createAddressOptionalShardId(input, SCAddressNumLeadingZeros)
	copy(address[vmcommon.NumInitCharactersForScAddress-core.VMTypeLen:], ei.GetVMType()[:])
	return address, err
}

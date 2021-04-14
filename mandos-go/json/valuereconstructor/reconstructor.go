package mandosvaluereconstructor

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	vi "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/valueinterpreter"
)

type ValueReconstructorHint uint64

const (
	// NoHint indicates that the type if not known
	NoHint ValueReconstructorHint = iota

	// NumberHint hints that value should be a number
	NumberHint

	// AddressHint hints that value should be an address
	AddressHint

	// StrHint hints that value should be a string expression, e.g. a username, "str:..."
	StrHint
)

// ValueReconstructor is a component that attempts to convert raw bytes to a human-readable format.
type ValueReconstructor struct{}

func (vr *ValueReconstructor) Reconstruct(value []byte, hint ValueReconstructorHint) string {
	switch hint {
	case NumberHint:
		return fmt.Sprintf("%d", big.NewInt(0).SetBytes(value))
	case StrHint:
		return fmt.Sprintf("str:%s", string(value))
	case AddressHint:
		return addressPretty((value))
	default:
		return unknownByteArrayPretty(value)
	}
}

func (vr *ValueReconstructor) ReconstructFromBigInt(value *big.Int) string {
	return vr.Reconstruct(value.Bytes(), NumberHint)
}

func (vr *ValueReconstructor) ReconstructFromUint64(value uint64) string {
	return vr.Reconstruct(big.NewInt(0).SetUint64(value).Bytes(), NumberHint)
}

func unknownByteArrayPretty(bytes []byte) string {
	if len(bytes) == 0 {
		return "[]"
	}

	if canInterpretAsString(bytes) {
		return fmt.Sprintf("0x%s (``%s)", hex.EncodeToString(bytes), string(bytes))
	}

	asInt := big.NewInt(0).SetBytes(bytes)
	return fmt.Sprintf("0x%s (%d)", hex.EncodeToString(bytes), asInt)
}

func addressPretty(value []byte) string {
	if len(value) != 32 {
		return unknownByteArrayPretty(value)
	}

	// smart contract addresses
	leadingZeros := make([]byte, vi.SCAddressNumLeadingZeros)
	if bytes.Equal(value[:vi.SCAddressNumLeadingZeros], leadingZeros) {
		if value[31] == byte('_') {
			addrStr := string(value[vi.SCAddressNumLeadingZeros:])
			addrStr = strings.TrimRight(addrStr, "_")
			return fmt.Sprintf("sc:%s", addrStr)
		} else {
			// last byte is the shard id and is explicit
			addrStr := string(value[vi.SCAddressNumLeadingZeros:31])
			addrStr = strings.TrimRight(addrStr, "_")
			shard_id := value[31]
			return fmt.Sprintf("sc:%s#%x", addrStr, shard_id)
		}
	}

	// regular addresses
	if value[31] == byte('_') {
		addrStr := string(value)
		addrStr = strings.TrimRight(addrStr, "_")
		return fmt.Sprintf("address:%s", addrStr)
	} else {
		// last byte is the shard id and is explicit
		addrStr := string(value[:31])
		addrStr = strings.TrimRight(addrStr, "_")
		shard_id := value[31]
		return fmt.Sprintf("address:%s#%02x", addrStr, shard_id)
	}
}

func canInterpretAsString(bytes []byte) bool {
	if len(bytes) == 0 {
		return false
	}
	for _, b := range bytes {
		if b < 32 || b > 126 {
			return false
		}
	}
	return true
}

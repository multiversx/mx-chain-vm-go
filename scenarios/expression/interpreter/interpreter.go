package mandosexpressioninterpreter

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
	fr "github.com/multiversx/wasm-vm/scenarios/fileresolver"
	oj "github.com/multiversx/wasm-vm/scenarios/orderedjson"
)

var strPrefixes = []string{"str:", "``", "''"}

const addrPrefix = "address:"
const scAddrPrefix = "sc:"

const filePrefix = "file:"
const keccak256Prefix = "keccak256:"

const u64Prefix = "u64:"
const u32Prefix = "u32:"
const u16Prefix = "u16:"
const u8Prefix = "u8:"
const i64Prefix = "i64:"
const i32Prefix = "i32:"
const i16Prefix = "i16:"
const i8Prefix = "i8:"

const bigFloatPrefix = "bigfloat:"
const biguintPrefix = "biguint:"
const nestedPrefix = "nested:"

// ExprInterpreter provides context for computing Mandos values.
type ExprInterpreter struct {
	FileResolver fr.FileResolver
}

// InterpretSubTree attempts to produce a value based on a JSON subtree.
// Subtrees are composed of strings, lists and maps.
// The idea is to intuitively represent serialized objects.
// Lists are evaluated by concatenating their items' representations.
// Maps are evaluated by concatenating their values' representations (keys are ignored).
// See InterpretString on how strings are being interpreted.
func (ei *ExprInterpreter) InterpretSubTree(obj oj.OJsonObject) ([]byte, error) {
	if str, isStr := obj.(*oj.OJsonString); isStr {
		return ei.InterpretString(str.Value)
	}

	if list, isList := obj.(*oj.OJsonList); isList {
		var concat []byte
		for _, item := range list.AsList() {
			value, err := ei.InterpretSubTree(item)
			if err != nil {
				return []byte{}, err
			}
			concat = append(concat, value...)
		}
		return concat, nil
	}

	if mp, isMap := obj.(*oj.OJsonMap); isMap {
		var concat []byte

		// keys are ignored, they do not form the value but act like documentation
		// but we sort by keys, because other JSON implementations cannot retain key order
		// and we need consistency
		sortedKVP := mp.KeyValuePairsSortedByKey()
		for _, kvp := range sortedKVP {
			value, err := ei.InterpretSubTree(kvp.Value)
			if err != nil {
				return []byte{}, err
			}
			concat = append(concat, value...)
		}
		return concat, nil
	}

	return []byte{}, errors.New("cannot interpret given JSON subtree as value")
}

// InterpretString resolves a string to a byte slice according to the Mandos value format.
// Supported rules are:
// - numbers: decimal, hex, binary, signed/unsigned
// - fixed length numbers: "u32:5", "i8:-3", etc.
// - ascii strings as "str:...", "``...", "''..."
// - "true"/"false"
// - "address:..."
// - "sc:..." (also an address)
// - "file:..."
// - "keccak256:..."
// - concatenation using |
//
func (ei *ExprInterpreter) InterpretString(strRaw string) ([]byte, error) {
	if len(strRaw) == 0 {
		return []byte{}, nil
	}

	// file contents
	// TODO: make this part of a proper parser
	if strings.HasPrefix(strRaw, filePrefix) {
		if ei.FileResolver == nil {
			return []byte{}, errors.New("parser FileResolver not provided")
		}
		fileContents, err := ei.FileResolver.ResolveFileValue(strRaw[len(filePrefix):])
		if err != nil {
			return []byte{}, err
		}
		return fileContents, nil
	}

	// keccak256
	// TODO: make this part of a proper parser
	if strings.HasPrefix(strRaw, keccak256Prefix) {
		arg, err := ei.InterpretString(strRaw[len(keccak256Prefix):])
		if err != nil {
			return []byte{}, fmt.Errorf("cannot parse keccak256 argument: %w", err)
		}
		hash, err := Keccak256(arg)
		if err != nil {
			return []byte{}, fmt.Errorf("error computing keccak256: %w", err)
		}
		return hash, nil
	}

	// concatenate values of different formats
	// TODO: make this part of a proper parser
	parts := strings.Split(strRaw, "|")
	if len(parts) > 1 {
		concat := make([]byte, 0)
		for _, part := range parts {
			eval, err := ei.InterpretString(part)
			if err != nil {
				return []byte{}, err
			}
			concat = append(concat, eval...)
		}
		return concat, nil
	}

	if strRaw == "false" {
		return []byte{}, nil
	}

	if strRaw == "true" {
		return []byte{0x01}, nil
	}

	// allow ascii strings, for readability
	for _, strPrefix := range strPrefixes {
		if strings.HasPrefix(strRaw, strPrefix) {
			str := strRaw[len(strPrefix):]
			return []byte(str), nil
		}
	}

	// address
	if strings.HasPrefix(strRaw, addrPrefix) {
		addrArgument := strRaw[len(addrPrefix):]
		return addressExpression(addrArgument)
	}

	// smart contract address (different format)
	if strings.HasPrefix(strRaw, scAddrPrefix) {
		addrArgument := strRaw[len(scAddrPrefix):]
		return scExpression(addrArgument)
	}

	// fixed width numbers
	parsed, result, err := ei.tryInterpretFixedWidth(strRaw)
	if err != nil {
		return nil, err
	}
	if parsed {
		return result, nil
	}

	// general numbers, arbitrary length
	return ei.interpretNumber(strRaw, 0)
}

func (ei *ExprInterpreter) interpretFloatingPointNumber(strRaw string) ([]byte, error) {
	str := strings.ReplaceAll(strRaw, "_", "")
	strRaw = strings.ReplaceAll(str, ",", "")

	// hex, the usual representation
	if strings.HasPrefix(strRaw, "0x") || strings.HasPrefix(strRaw, "0X") {
		str := strRaw[2:]
		if len(str)%2 == 1 {
			str = "0" + str
		}
		return hex.DecodeString(str)
	} else {
		returnFloatValue := big.NewFloat(0)
		mandosBigFloatVal, ok := big.NewFloat(0).SetString(strRaw)
		if !ok {
			return make([]byte, 0), fmt.Errorf("could not parse %s to big float", strRaw)
		}
		_ = returnFloatValue.Add(returnFloatValue, mandosBigFloatVal)
		encodedFloatValue, err := returnFloatValue.GobEncode()
		if err != nil {
			return make([]byte, 0), fmt.Errorf("could not parse %s to big float", strRaw)
		}
		return encodedFloatValue, nil
	}
}

// targetWidth = 0 means minimum length that can contain the result
func (ei *ExprInterpreter) interpretNumber(strRaw string, targetWidth int) ([]byte, error) {
	if strings.Contains(strRaw, ".") {
		bfBytes, err := ei.interpretFloatingPointNumber(strRaw[:])
		return bfBytes, err
	}

	// signed numbers
	if strRaw[0] == '-' || strRaw[0] == '+' {
		return ei.interpretNumberWithSign(strRaw, targetWidth)
	}

	// unsigned numbers
	if targetWidth == 0 {
		return ei.interpretUnsignedNumber(strRaw)
	}

	return ei.interpretUnsignedNumberFixedWidth(strRaw, targetWidth)
}

func (ei *ExprInterpreter) interpretNumberWithSign(strRaw string, targetWidth int) ([]byte, error) {
	numberBytes, err := ei.interpretUnsignedNumber(strRaw[1:])
	if err != nil {
		return []byte{}, err
	}
	number := big.NewInt(0).SetBytes(numberBytes)
	if strRaw[0] == '-' {
		number = number.Neg(number)
	}
	if targetWidth == 0 {
		return twos.ToBytes(number), nil
	}

	return twos.ToBytesOfLength(number, targetWidth)
}

func (ei *ExprInterpreter) interpretUnsignedNumber(strRaw string) ([]byte, error) {
	str := strings.ReplaceAll(strRaw, "_", "") // allow underscores, to group digits
	strRaw = strings.ReplaceAll(str, ",", "")  // also allow commas to group digits

	// hex, the usual representation
	if strings.HasPrefix(strRaw, "0x") || strings.HasPrefix(strRaw, "0X") {
		str := strRaw[2:]
		if len(str)%2 == 1 {
			str = "0" + str
		}
		return hex.DecodeString(str)
	}

	// binary representation
	if strings.HasPrefix(strRaw, "0b") || strings.HasPrefix(strRaw, "0B") {
		result := new(big.Int)
		var parseOk bool
		result, parseOk = result.SetString(str[2:], 2)
		if !parseOk {
			return []byte{}, fmt.Errorf("could not parse binary value: %s", strRaw)
		}

		return result.Bytes(), nil
	}

	// default: parse as BigInt, base 10
	result := new(big.Int)
	var parseOk bool
	result, parseOk = result.SetString(strRaw, 10)
	if !parseOk {
		return []byte{}, fmt.Errorf("could not parse base 10 value: %s", strRaw)
	}

	if result.Sign() < 0 {
		return []byte{}, fmt.Errorf("negative numbers not allowed in this context: %s", strRaw)
	}

	return result.Bytes(), nil
}

func (ei *ExprInterpreter) interpretUnsignedNumberFixedWidth(strRaw string, targetWidth int) ([]byte, error) {
	numberBytes, err := ei.interpretUnsignedNumber(strRaw)
	if err != nil {
		return []byte{}, err
	}
	if targetWidth == 0 {
		return numberBytes, nil
	}

	if len(numberBytes) > targetWidth {
		return []byte{}, fmt.Errorf("representation of %s does not fit in %d bytes", strRaw, targetWidth)
	}
	return twos.CopyAlignRight(numberBytes, targetWidth), nil
}

func (ei *ExprInterpreter) tryInterpretFixedWidth(strRaw string) (bool, []byte, error) {
	if strings.HasPrefix(strRaw, u64Prefix) {
		r, err := ei.interpretUnsignedNumberFixedWidth(strRaw[len(u64Prefix):], 8)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, u32Prefix) {
		r, err := ei.interpretUnsignedNumberFixedWidth(strRaw[len(u32Prefix):], 4)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, u16Prefix) {
		r, err := ei.interpretUnsignedNumberFixedWidth(strRaw[len(u16Prefix):], 2)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, u8Prefix) {
		r, err := ei.interpretUnsignedNumberFixedWidth(strRaw[len(u8Prefix):], 1)
		return true, r, err
	}

	if strings.HasPrefix(strRaw, i64Prefix) {
		r, err := ei.interpretNumber(strRaw[len(i64Prefix):], 8)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, i32Prefix) {
		r, err := ei.interpretNumber(strRaw[len(i32Prefix):], 4)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, i16Prefix) {
		r, err := ei.interpretNumber(strRaw[len(i16Prefix):], 2)
		return true, r, err
	}
	if strings.HasPrefix(strRaw, i8Prefix) {
		r, err := ei.interpretNumber(strRaw[len(i8Prefix):], 1)
		return true, r, err
	}

	if strings.HasPrefix(strRaw, biguintPrefix) {
		return ei.interpretExplicitBigUintNumber(strRaw)
	}

	if strings.HasPrefix(strRaw, bigFloatPrefix) {
		return ei.interpretExplicitFloatingPointNumber(strRaw)
	}

	if strings.HasPrefix(strRaw, nestedPrefix) {
		return ei.interpretNestedBytes(strRaw)
	}

	return false, []byte{}, nil
}

func (ei *ExprInterpreter) interpretExplicitFloatingPointNumber(strRaw string) (bool, []byte, error) {
	bfBytes, err := ei.interpretFloatingPointNumber(strRaw[len(bigFloatPrefix):])
	lengthBytes := big.NewInt(int64(len(bfBytes))).Bytes()
	encodedLength := twos.CopyAlignRight(lengthBytes, 4)
	return true, append(encodedLength, bfBytes...), err
}

func (ei *ExprInterpreter) interpretExplicitBigUintNumber(strRaw string) (bool, []byte, error) {
	biBytes, err := ei.interpretUnsignedNumber(strRaw[len(biguintPrefix):])
	lengthBytes := big.NewInt(int64(len(biBytes))).Bytes()
	encodedLength := twos.CopyAlignRight(lengthBytes, 4)
	return true, append(encodedLength, biBytes...), err
}

func (ei *ExprInterpreter) interpretNestedBytes(strRaw string) (bool, []byte, error) {
	nestedBytes, err := ei.InterpretString(strRaw[len(nestedPrefix):])
	lengthBytes := big.NewInt(int64(len(nestedBytes))).Bytes()
	encodedLength := twos.CopyAlignRight(lengthBytes, 4)
	return true, append(encodedLength, nestedBytes...), err
}

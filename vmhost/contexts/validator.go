package contexts

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const allowedCharsInFunctionName = "abcdefghijklmnopqrstuvwxyz0123456789_"

var reservedFunctionsActivationFlag = map[string]core.EnableEpochFlag{
	"mbufferToSmallIntUnsigned":                    vmhost.BarnardOpcodesFlag,
	"mbufferToSmallIntSigned":                      vmhost.BarnardOpcodesFlag,
	"mbufferFromSmallIntUnsigned":                  vmhost.BarnardOpcodesFlag,
	"mbufferFromSmallIntSigned":                    vmhost.BarnardOpcodesFlag,
	"getBlockRoundTimeMs":                          vmhost.BarnardOpcodesFlag,
	"epochStartBlockTimeStamp":                     vmhost.BarnardOpcodesFlag,
	"epochStartBlockNonce":                         vmhost.BarnardOpcodesFlag,
	"epochStartBlockRound":                         vmhost.BarnardOpcodesFlag,
	"managedGetAllTransfersCallValue":              vmhost.BarnardOpcodesFlag,
	"managedExecuteOnDestContextWithErrorReturn":   vmhost.BarnardOpcodesFlag,
	"managedMultiTransferESDTNFTExecuteWithReturn": vmhost.BarnardOpcodesFlag,
	"managedGetCodeHash":                           vmhost.BarnardOpcodesFlag,
	"managedGetESDTTokenType":                      vmhost.BarnardOpcodesFlag,
}

// wasmValidator is a validator for WASM SmartContracts
type wasmValidator struct {
	reserved *reservedFunctions
}

// newWASMValidator creates a new WASMValidator
func newWASMValidator(
	scAPINames vmcommon.FunctionNames,
	builtInFuncContainer vmcommon.BuiltInFunctionContainer,
	enableEpochsHandler vmcommon.EnableEpochsHandler,
) *wasmValidator {
	return &wasmValidator{
		reserved: NewReservedFunctions(scAPINames, builtInFuncContainer, reservedFunctionsActivationFlag, enableEpochsHandler),
	}
}

func (validator *wasmValidator) verifyMemoryDeclaration(instance executor.Instance) error {
	if !instance.HasMemory() {
		return vmhost.ErrMemoryDeclarationMissing
	}

	return nil
}

func (validator *wasmValidator) verifyFunctions(instance executor.Instance) error {
	functionNames := instance.GetFunctionNames()
	if functionNames == nil {
		return vmhost.ErrContractInvalid
	}

	for _, functionName := range functionNames {
		err := validator.verifyValidFunctionName(functionName)
		if err != nil {
			return err
		}
	}

	return instance.ValidateFunctionArities()
}

var protectedFunctions = map[string]bool{
	"internalVMErrors":  true,
	"transferValueOnly": true,
	"writeLog":          true,
	"signalError":       true,
	"completedTxEvent":  true}

func (validator *wasmValidator) verifyProtectedFunctions(instance executor.Instance) error {
	functionNames := instance.GetFunctionNames()
	if functionNames == nil {
		return vmhost.ErrContractInvalid
	}

	for _, functionName := range functionNames {
		_, found := protectedFunctions[functionName]
		if found {
			return vmhost.ErrContractInvalid
		}

	}

	return nil
}

func (validator *wasmValidator) verifyValidFunctionName(functionName string) error {
	err := verifyCallFunction(functionName)
	if err != nil {
		return err
	}

	errInvalidName := fmt.Errorf("%w: %x", vmhost.ErrInvalidFunctionName, functionName)
	if validator.reserved.IsReserved(functionName) {
		return errInvalidName
	}

	return nil
}

func verifyCallFunction(functionName string) error {
	const maxLengthOfFunctionName = 256

	errInvalidName := fmt.Errorf("%w: %s", vmhost.ErrInvalidFunctionName, functionName)

	if len(functionName) == 0 {
		return errInvalidName
	}
	if len(functionName) >= maxLengthOfFunctionName {
		return errInvalidName
	}
	if isFirstCharacterNumeric(functionName) {
		return errInvalidName
	}
	if !validCharactersOnly(functionName) {
		return errInvalidName
	}

	return nil
}

func validCharactersOnly(input string) bool {
	for i := 0; i < len(input); i++ {
		b := input[i]
		isAlphaNum := (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
		if !isAlphaNum {
			return false
		}
	}

	return true
}

func isFirstCharacterNumeric(name string) bool {
	return name[0] >= '0' && name[0] <= '9'
}

const maxAllowedTableElements = 10000

var (
	wasmMagicNumber            = []byte{0x00, 0x61, 0x73, 0x6d}
	wasmVersion1               = []byte{0x01, 0x00, 0x00, 0x00}
	errUnexpectedEndOfBytecode = errors.New("unexpected end of bytecode")
	errMalformedLeb128         = errors.New("malformed LEB128 integer")
)

func (validator *wasmValidator) verifyContractBytecode(bytecode []byte) error {
	if len(bytecode) < 8 {
		return fmt.Errorf("%w: bytecode too short for WASM header", vmhost.ErrContractInvalid)
	}
	if !bytes.Equal(bytecode[:4], wasmMagicNumber) || !bytes.Equal(bytecode[4:8], wasmVersion1) {
		return fmt.Errorf("%w: invalid WASM magic or version", vmhost.ErrContractInvalid)
	}

	offset := 8
	for offset < len(bytecode) {
		sectionID := bytecode[offset]
		offset++

		sectionSize, newOffset, err := decodeVarUint32(bytecode, offset)
		if err != nil {
			return fmt.Errorf("%w: failed to decode section size: %v", vmhost.ErrContractInvalid, err)
		}
		offset = newOffset

		if int(sectionSize) < 0 || offset+int(sectionSize) > len(bytecode) {
			return fmt.Errorf("%w: section size exceeds bytecode length", vmhost.ErrContractInvalid)
		}

		sectionEnd := offset + int(sectionSize)
		sectionPayload := bytecode[offset:sectionEnd]

		switch sectionID {
		case 2: // Import section
			err = validator.verifyImportSection(sectionPayload)
			if err != nil {
				return err
			}
		case 4: // Table section
			err = validator.verifyTableSection(sectionPayload)
			if err != nil {
				return err
			}
		case 10: // Code section
			err = validator.verifyCodeSection(sectionPayload)
			if err != nil {
				return err
			}
		}

		offset = sectionEnd
	}

	return nil
}

func (validator *wasmValidator) verifyImportSection(payload []byte) error {
	count, offset, err := decodeVarUint32(payload, 0)
	if err != nil {
		return fmt.Errorf("%w: malformed import count: %v", vmhost.ErrContractInvalid, err)
	}

	for i := uint32(0); i < count; i++ {
		modLen, newOffset, errMod := decodeVarUint32(payload, offset)
		if errMod != nil {
			return fmt.Errorf("%w: malformed import module length: %v", vmhost.ErrContractInvalid, errMod)
		}
		offset = newOffset + int(modLen)
		if offset > len(payload) {
			return fmt.Errorf("%w: import module name exceeds payload", vmhost.ErrContractInvalid)
		}

		fieldLen, newOffset, errField := decodeVarUint32(payload, offset)
		if errField != nil {
			return fmt.Errorf("%w: malformed import field length: %v", vmhost.ErrContractInvalid, errField)
		}
		offset = newOffset + int(fieldLen)
		if offset > len(payload) {
			return fmt.Errorf("%w: import field name exceeds payload", vmhost.ErrContractInvalid)
		}

		if offset >= len(payload) {
			return fmt.Errorf("%w: unexpected end of import payload", vmhost.ErrContractInvalid)
		}
		kind := payload[offset]
		offset++

		switch kind {
		case 0x00: // Function
			_, newOffset, errFunc := decodeVarUint32(payload, offset)
			if errFunc != nil {
				return fmt.Errorf("%w: malformed import func type: %v", vmhost.ErrContractInvalid, errFunc)
			}
			offset = newOffset
		case 0x01: // Table
			return fmt.Errorf("%w: imported tables are not supported", vmhost.ErrContractInvalid)
		case 0x02: // Memory
			if offset >= len(payload) {
				return fmt.Errorf("%w: unexpected end of memory import", vmhost.ErrContractInvalid)
			}
			flags := payload[offset]
			offset++
			_, newOffset, errMem := decodeVarUint32(payload, offset)
			if errMem != nil {
				return fmt.Errorf("%w: malformed memory min: %v", vmhost.ErrContractInvalid, errMem)
			}
			offset = newOffset
			if (flags & 0x01) != 0 {
				_, newOffset, errMemMax := decodeVarUint32(payload, offset)
				if errMemMax != nil {
					return fmt.Errorf("%w: malformed memory max: %v", vmhost.ErrContractInvalid, errMemMax)
				}
				offset = newOffset
			}
		case 0x03: // Global
			offset += 2 // valtype + mut
			if offset > len(payload) {
				return fmt.Errorf("%w: unexpected end of global import", vmhost.ErrContractInvalid)
			}
		default:
			return fmt.Errorf("%w: unknown import kind 0x%02x", vmhost.ErrContractInvalid, kind)
		}
	}

	return nil
}

func (validator *wasmValidator) verifyTableSection(payload []byte) error {
	count, offset, err := decodeVarUint32(payload, 0)
	if err != nil {
		return fmt.Errorf("%w: malformed table count: %v", vmhost.ErrContractInvalid, err)
	}

	for i := uint32(0); i < count; i++ {
		if offset >= len(payload) {
			return fmt.Errorf("%w: unexpected end of table section", vmhost.ErrContractInvalid)
		}
		elemType := payload[offset]
		offset++
		if elemType != 0x70 { // 0x70 = funcref
			return fmt.Errorf("%w: invalid table element type 0x%02x", vmhost.ErrContractInvalid, elemType)
		}

		if offset >= len(payload) {
			return fmt.Errorf("%w: unexpected end of table section", vmhost.ErrContractInvalid)
		}
		flags := payload[offset]
		offset++

		if (flags & 0x01) == 0 {
			return fmt.Errorf("%w: table must specify a maximum size", vmhost.ErrContractInvalid)
		}

		min, newOffset, errMin := decodeVarUint32(payload, offset)
		if errMin != nil {
			return fmt.Errorf("%w: malformed table min: %v", vmhost.ErrContractInvalid, errMin)
		}
		offset = newOffset

		max, newOffset, errMax := decodeVarUint32(payload, offset)
		if errMax != nil {
			return fmt.Errorf("%w: malformed table max: %v", vmhost.ErrContractInvalid, errMax)
		}
		offset = newOffset

		if min > max {
			return fmt.Errorf("%w: table minimum %d exceeds maximum %d", vmhost.ErrContractInvalid, min, max)
		}

		if max > maxAllowedTableElements {
			return fmt.Errorf("%w: table maximum %d exceeds limit %d", vmhost.ErrContractInvalid, max, maxAllowedTableElements)
		}
	}

	return nil
}

func (validator *wasmValidator) verifyCodeSection(payload []byte) error {
	count, offset, err := decodeVarUint32(payload, 0)
	if err != nil {
		return fmt.Errorf("%w: malformed code section function count: %v", vmhost.ErrContractInvalid, err)
	}

	for i := uint32(0); i < count; i++ {
		bodySize, newOffset, errSize := decodeVarUint32(payload, offset)
		if errSize != nil {
			return fmt.Errorf("%w: malformed function body size: %v", vmhost.ErrContractInvalid, errSize)
		}
		offset = newOffset

		bodyEnd := offset + int(bodySize)
		if bodyEnd > len(payload) || bodyEnd < offset {
			return fmt.Errorf("%w: function body exceeds code section", vmhost.ErrContractInvalid)
		}

		numLocals, newOffset, errLocals := decodeVarUint32(payload, offset)
		if errLocals != nil {
			return fmt.Errorf("%w: malformed locals count: %v", vmhost.ErrContractInvalid, errLocals)
		}
		offset = newOffset

		for j := uint32(0); j < numLocals; j++ {
			_, newOffset, errLocal := decodeVarUint32(payload, offset)
			if errLocal != nil {
				return fmt.Errorf("%w: malformed local count: %v", vmhost.ErrContractInvalid, errLocal)
			}
			offset = newOffset + 1 // skip valtype byte
			if offset > bodyEnd {
				return fmt.Errorf("%w: locals exceed function body", vmhost.ErrContractInvalid)
			}
		}

		err = validator.scanInstructions(payload[offset:bodyEnd])
		if err != nil {
			return err
		}

		offset = bodyEnd
	}

	return nil
}

func (validator *wasmValidator) scanInstructions(code []byte) error {
	p := 0
	for p < len(code) {
		op := code[p]
		p++

		switch {
		case op == 0x00 || op == 0x01 || op == 0x05 || op == 0x0B || op == 0x0F || op == 0x1A || op == 0x1B:
			continue

		case op == 0x02 || op == 0x03 || op == 0x04:
			if p >= len(code) {
				return fmt.Errorf("%w: unexpected end of blocktype", vmhost.ErrContractInvalid)
			}
			b := code[p]
			if b == 0x40 || b == 0x7F || b == 0x7E || b == 0x7D || b == 0x7C {
				p++
			} else {
				_, newP, err := decodeVarInt32(code, p)
				if err != nil {
					return fmt.Errorf("%w: malformed blocktype: %v", vmhost.ErrContractInvalid, err)
				}
				p = newP
			}

		case op == 0x0C || op == 0x0D || op == 0x10 || op == 0x20 || op == 0x21 || op == 0x22 || op == 0x23 || op == 0x24:
			_, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed index operand: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0x0E:
			cnt, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed br_table count: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP
			for k := uint32(0); k < cnt; k++ {
				_, newP, err = decodeVarUint32(code, p)
				if err != nil {
					return fmt.Errorf("%w: malformed br_table target: %v", vmhost.ErrContractInvalid, err)
				}
				p = newP
			}
			_, newP, err = decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed br_table default: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0x11:
			_, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed call_indirect typeidx: %v", vmhost.ErrContractInvalid, err)
			}
			_, newP, err = decodeVarUint32(code, newP)
			if err != nil {
				return fmt.Errorf("%w: malformed call_indirect tableidx: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0x1C:
			cnt, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed typed select count: %v", vmhost.ErrContractInvalid, err)
			}
			if uint64(newP)+uint64(cnt) > uint64(len(code)) {
				return fmt.Errorf("%w: unexpected end of typed select types", vmhost.ErrContractInvalid)
			}
			p = newP + int(cnt)

		case op == 0x25:
			_, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed table.get operand: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0x26:
			return fmt.Errorf("%w: opcode table.set is forbidden", vmhost.ErrContractInvalid)

		case op >= 0x28 && op <= 0x3E:
			_, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed mem align: %v", vmhost.ErrContractInvalid, err)
			}
			_, newP, err = decodeVarUint32(code, newP)
			if err != nil {
				return fmt.Errorf("%w: malformed mem offset: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0x3F || op == 0x40:
			if p >= len(code) {
				return fmt.Errorf("%w: unexpected end of memory instruction", vmhost.ErrContractInvalid)
			}
			p++

		case op == 0x41:
			_, newP, err := decodeVarInt32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed i32.const: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0x42:
			_, newP, err := decodeVarInt64(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed i64.const: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0x43:
			p += 4
			if p > len(code) {
				return fmt.Errorf("%w: unexpected end of f32.const", vmhost.ErrContractInvalid)
			}

		case op == 0x44:
			p += 8
			if p > len(code) {
				return fmt.Errorf("%w: unexpected end of f64.const", vmhost.ErrContractInvalid)
			}

		case op >= 0x45 && op <= 0xC4:
			continue

		case op == 0xD0:
			p++
			if p > len(code) {
				return fmt.Errorf("%w: unexpected end of ref.null", vmhost.ErrContractInvalid)
			}

		case op == 0xD1:
			continue

		case op == 0xD2:
			_, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed ref.func: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

		case op == 0xFC:
			subOp, newP, err := decodeVarUint32(code, p)
			if err != nil {
				return fmt.Errorf("%w: malformed 0xFC sub-opcode: %v", vmhost.ErrContractInvalid, err)
			}
			p = newP

			if subOp == 0x0F {
				return fmt.Errorf("%w: opcode table.grow is forbidden", vmhost.ErrContractInvalid)
			}

			switch {
			case subOp <= 0x07:
				continue
			case subOp == 0x08:
				_, newP, err = decodeVarUint32(code, p)
				if err != nil {
					return fmt.Errorf("%w: malformed memory.init: %v", vmhost.ErrContractInvalid, err)
				}
				p = newP + 1
			case subOp == 0x09:
				_, newP, err = decodeVarUint32(code, p)
				if err != nil {
					return fmt.Errorf("%w: malformed data.drop: %v", vmhost.ErrContractInvalid, err)
				}
				p = newP
			case subOp == 0x0A:
				p += 2
			case subOp == 0x0B:
				p++
			case subOp == 0x0C || subOp == 0x0E:
				_, newP, err = decodeVarUint32(code, p)
				if err != nil {
					return fmt.Errorf("%w: malformed table op: %v", vmhost.ErrContractInvalid, err)
				}
				_, newP, err = decodeVarUint32(code, newP)
				if err != nil {
					return fmt.Errorf("%w: malformed table op: %v", vmhost.ErrContractInvalid, err)
				}
				p = newP
			case subOp == 0x0D || subOp == 0x10 || subOp == 0x11:
				_, newP, err = decodeVarUint32(code, p)
				if err != nil {
					return fmt.Errorf("%w: malformed table op: %v", vmhost.ErrContractInvalid, err)
				}
				p = newP
			default:
				return fmt.Errorf("%w: unknown 0xFC sub-opcode 0x%02x", vmhost.ErrContractInvalid, subOp)
			}

			if p > len(code) {
				return fmt.Errorf("%w: unexpected end of 0xFC instruction", vmhost.ErrContractInvalid)
			}

		case op == 0xFD:
			return fmt.Errorf("%w: SIMD opcode 0xFD is not supported", vmhost.ErrContractInvalid)

		default:
			return fmt.Errorf("%w: unknown opcode 0x%02x", vmhost.ErrContractInvalid, op)
		}
	}

	return nil
}

func decodeVarUint32(data []byte, offset int) (uint32, int, error) {
	var val uint32
	var shift uint
	for {
		if offset >= len(data) {
			return 0, offset, errUnexpectedEndOfBytecode
		}
		b := data[offset]
		offset++
		val |= uint32(b&0x7f) << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
		if shift >= 35 {
			return 0, offset, errMalformedLeb128
		}
	}
	return val, offset, nil
}

func decodeVarInt32(data []byte, offset int) (int32, int, error) {
	var val int32
	var shift uint
	for {
		if offset >= len(data) {
			return 0, offset, errUnexpectedEndOfBytecode
		}
		b := data[offset]
		offset++
		val |= int32(b&0x7f) << shift
		shift += 7
		if (b & 0x80) == 0 {
			if shift < 32 && (b&0x40) != 0 {
				val |= -(int32(1) << shift)
			}
			break
		}
		if shift >= 35 {
			return 0, offset, errMalformedLeb128
		}
	}
	return val, offset, nil
}

func decodeVarInt64(data []byte, offset int) (int64, int, error) {
	var val int64
	var shift uint
	for {
		if offset >= len(data) {
			return 0, offset, errUnexpectedEndOfBytecode
		}
		b := data[offset]
		offset++
		val |= int64(b&0x7f) << shift
		shift += 7
		if (b & 0x80) == 0 {
			if shift < 64 && (b&0x40) != 0 {
				val |= -(int64(1) << shift)
			}
			break
		}
		if shift >= 70 {
			return 0, offset, errMalformedLeb128
		}
	}
	return val, offset, nil
}

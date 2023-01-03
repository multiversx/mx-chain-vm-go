package config

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

type operations struct {
	OperationA uint64
	OperationB uint64
	OperationC uint64
	OperationD uint64
	OperationE uint64
}

func TestDecode(t *testing.T) {
	gasMap := make(map[string]uint64)
	gasMap["OperationB"] = 4
	gasMap["OperationA"] = 3
	gasMap["OperationC"] = 100
	gasMap["OperationD"] = 1000

	op := &operations{}
	err := mapstructure.Decode(gasMap, op)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", op)
}

func TestDecode_ArwenGas(t *testing.T) {
	gasMap := make(map[string]uint64)
	gasMap["StorePerByte"] = 4
	gasMap["GetSCAddress"] = 4
	gasMap["GetExternalBalance"] = 4
	gasMap["BigIntByteLength"] = 4

	bigIntOp := &BigIntAPICost{}
	err := mapstructure.Decode(gasMap, bigIntOp)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", bigIntOp)

	bigFloatOp := &BigFloatAPICost{}
	err = mapstructure.Decode(gasMap, bigFloatOp)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", bigFloatOp)

	erdOp := &ElrondAPICost{}
	err = mapstructure.Decode(gasMap, erdOp)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", erdOp)

	ethOp := &EthAPICost{}
	err = mapstructure.Decode(gasMap, ethOp)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", ethOp)
}

func TestDecode_ZeroGasCostError(t *testing.T) {
	gasMap := FillGasMapWASMOpcodeValues(1)

	wasmCosts := &WASMOpcodeCost{}
	err := mapstructure.Decode(gasMap, wasmCosts)
	assert.Nil(t, err)

	err = checkForZeroUint64Fields(*wasmCosts)
	assert.Nil(t, err)

	gasMap["BrIf"] = 0
	wasmCosts = &WASMOpcodeCost{}
	err = mapstructure.Decode(gasMap, wasmCosts)
	assert.Nil(t, err)

	err = checkForZeroUint64Fields(*wasmCosts)
	assert.Error(t, err)
}

func Test_getSignedCoefficient(t *testing.T) {
	gasScheduleMap := MakeGasMap(1, 1)

	A := uint64(687)
	B := uint64(30483)
	C := uint64(15883)

	gasMap := make(map[string]uint64)
	gasMap["A"] = A
	gasMap["SignOfA"] = 0
	gasMap["B"] = B
	gasMap["SignOfB"] = 0
	gasMap["C"] = C
	gasMap["SignOfC"] = 1
	gasScheduleMap["DynamicStorageLoad"] = gasMap

	gasCost, err := CreateGasConfig(gasScheduleMap)
	assert.Nil(t, err)
	assert.Equal(t, int64(A), gasCost.DynamicStorageLoad.Quadratic)
	assert.Equal(t, int64(B), gasCost.DynamicStorageLoad.Linear)
	assert.Equal(t, int64(C)*-1, gasCost.DynamicStorageLoad.Constant)
}

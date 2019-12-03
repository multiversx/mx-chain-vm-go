package ethapi

import (
	"encoding/hex"
	"fmt"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

// The mapping between system contracts and their addresses is defined here:
// https://ewasm.readthedocs.io/en/mkdocs/system_contracts/
var contractsMap = map[string]func(unsafe.Pointer, []byte) ([]byte, error){
	"0000000000000000000000000000000000000001": ecrecover,
	"0000000000000000000000000000000000000002": sha2,
	"0000000000000000000000000000000000000003": ripemd160,
	"0000000000000000000000000000000000000004": identity,
	"0000000000000000000000000000000000000009": keccak256,
}

// IsAddressForPredefinedContract returns whether the address is recognized as a eth "precompiled contract" (predefined contract)
func IsAddressForPredefinedContract(address []byte) bool {
	contractKey := hex.EncodeToString(address)
	_, ok := contractsMap[contractKey]
	return ok
}

// CallPredefinedContract executes a predefined contract specified by address
func CallPredefinedContract(ctx unsafe.Pointer, address []byte, data []byte) error {
	instCtx := wasmer.IntoInstanceContext(ctx)
	ethCtx := arwen.GetEthContext(instCtx.Data())

	contractKey := hex.EncodeToString(address)
	contract, ok := contractsMap[contractKey]
	if !ok {
		return fmt.Errorf("invalid EEI system contract call - missing: %s", contractKey)
	}

	returnData, err := contract(ctx, data)
	if err != nil {
		return fmt.Errorf("erroneous EEI system contract call: %s", err.Error())
	}

	ethCtx.ClearReturnData()
	ethCtx.Finish(returnData)
	return nil
}

func ecrecover(context unsafe.Pointer, data []byte) ([]byte, error) {
	return nil, fmt.Errorf("EEI system contract not implemented: ecrecover")
}

func sha2(context unsafe.Pointer, data []byte) ([]byte, error) {
	instCtx := wasmer.IntoInstanceContext(context)
	cryptoCtx := arwen.GetCryptoContext(instCtx.Data())
	return cryptoCtx.CryptoHooks().Sha256(data)
}

func ripemd160(context unsafe.Pointer, data []byte) ([]byte, error) {
	instCtx := wasmer.IntoInstanceContext(context)
	cryptoCtx := arwen.GetCryptoContext(instCtx.Data())
	return cryptoCtx.CryptoHooks().Ripemd160(data)
}

func identity(context unsafe.Pointer, data []byte) ([]byte, error) {
	return data, nil
}

func keccak256(context unsafe.Pointer, data []byte) ([]byte, error) {
	instCtx := wasmer.IntoInstanceContext(context)
	cryptoCtx := arwen.GetCryptoContext(instCtx.Data())
	return cryptoCtx.CryptoHooks().Keccak256(data)
}

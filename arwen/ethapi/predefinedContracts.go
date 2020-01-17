package ethapi

import (
	"encoding/hex"
	"fmt"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
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
func CallPredefinedContract(context unsafe.Pointer, address []byte, data []byte) error {
	output := arwen.GetOutputContext(context)

	contractKey := hex.EncodeToString(address)
	contract, ok := contractsMap[contractKey]
	if !ok {
		return fmt.Errorf("invalid EEI system contract call - missing: %s", contractKey)
	}

	returnData, err := contract(context, data)
	if err != nil {
		return fmt.Errorf("erroneous EEI system contract call: %s", err.Error())
	}

	output.ClearReturnData()
	output.Finish(returnData)
	return nil
}

func ecrecover(context unsafe.Pointer, data []byte) ([]byte, error) {
	return nil, fmt.Errorf("EEI system contract not implemented: ecrecover")
}

func sha2(context unsafe.Pointer, data []byte) ([]byte, error) {
	crypto := arwen.GetCryptoContext(context)
	return crypto.Sha256(data)
}

func ripemd160(context unsafe.Pointer, data []byte) ([]byte, error) {
	crypto := arwen.GetCryptoContext(context)
	return crypto.Ripemd160(data)
}

func identity(context unsafe.Pointer, data []byte) ([]byte, error) {
	return data, nil
}

func keccak256(context unsafe.Pointer, data []byte) ([]byte, error) {
	crypto := arwen.GetCryptoContext(context)
	result, err := crypto.Keccak256(data)
	if err != nil {
		fmt.Printf("Error Keccak256: %s\n", err.Error())
	}
	return result, err
}

package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// reservedFunctions holds the reserved function names
type reservedFunctions struct {
	functionNames vmcommon.FunctionNames
}

// NewReservedFunctions creates a new reservedFunctions
func NewReservedFunctions(scAPINames vmcommon.FunctionNames, protocolBuiltinFunctions vmcommon.FunctionNames) *reservedFunctions {
	result := &reservedFunctions{
		functionNames: make(vmcommon.FunctionNames),
	}

	for name, value := range protocolBuiltinFunctions {
		result.functionNames[name] = value
	}

	for name, value := range scAPINames {
		result.functionNames[name] = value
	}

	var empty struct{}
	result.functionNames[arwen.UpgradeFunctionName] = empty

	return result
}

// IsReserved returns whether a function is reserved
func (reservedFunctions *reservedFunctions) IsReserved(functionName string) bool {
	if _, ok := reservedFunctions.functionNames[functionName]; ok {
		return true
	}

	return false
}

// GetReserved gets the reserved functions as a slice of strings
func (reservedFunctions *reservedFunctions) GetReserved() []string {
	keys := make([]string, len(reservedFunctions.functionNames))

	i := 0
	for key := range reservedFunctions.functionNames {
		keys[i] = key
		i++
	}

	return keys
}

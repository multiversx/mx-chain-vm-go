package contexts

import (
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/arwen"
)

// reservedFunctions holds the reserved function names
type reservedFunctions struct {
	functionNames vmcommon.FunctionNames
}

// NewReservedFunctions creates a new reservedFunctions
func NewReservedFunctions(scAPINames vmcommon.FunctionNames, builtInFuncContainer vmcommon.BuiltInFunctionContainer) *reservedFunctions {
	result := &reservedFunctions{
		functionNames: make(vmcommon.FunctionNames),
	}

	protocolFuncNames := builtInFuncContainer.Keys()
	for name := range protocolFuncNames {
		function, err := builtInFuncContainer.Get(name)
		if err != nil || !function.IsActive() {
			continue
		}

		result.functionNames[name] = struct{}{}
	}

	for name, value := range scAPINames {
		result.functionNames[name] = value
	}

	var empty struct{}
	result.functionNames[arwen.UpgradeFunctionName] = empty
	result.functionNames[arwen.DeleteFunctionName] = empty

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

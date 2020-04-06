package contexts

import "github.com/ElrondNetwork/arwen-wasm-vm/arwen"

// ReservedFunctions holds the reserved function names
type ReservedFunctions struct {
	functionNames map[string]struct{}
}

// NewReservedFunctions creates a new ReservedFunctions
func NewReservedFunctions(scAPINames []string, protocolReservedFunctions []string) *ReservedFunctions {
	result := &ReservedFunctions{
		functionNames: make(map[string]struct{}),
	}

	var empty struct{}

	for _, name := range protocolReservedFunctions {
		result.functionNames[name] = empty
	}

	for _, name := range scAPINames {
		result.functionNames[name] = empty
	}

	result.functionNames[arwen.UpgradeFunctionName] = empty

	return result
}

// IsReserved returns whether a function is reserved
func (reservedFunctions *ReservedFunctions) IsReserved(functionName string) bool {
	if _, ok := reservedFunctions.functionNames[functionName]; ok {
		return true
	}

	return false
}

// GetReserved gets the reserved functions as a slice of strings
func (reservedFunctions *ReservedFunctions) GetReserved() []string {
	keys := make([]string, len(reservedFunctions.functionNames))

	i := 0
	for key := range reservedFunctions.functionNames {
		keys[i] = key
		i++
	}

	return keys
}

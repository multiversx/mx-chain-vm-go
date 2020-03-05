package contexts

// ReservedFunctions holds the reserved function names
type ReservedFunctions struct {
	functionNames map[string]interface{}
}

// NewReservedFunctions creates a new ReservedFunctions
func NewReservedFunctions() *ReservedFunctions {
	result := &ReservedFunctions{
		functionNames: make(map[string]interface{}, 0),
	}

	result.functionNames["claimDeveloperRewards"] = nil
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

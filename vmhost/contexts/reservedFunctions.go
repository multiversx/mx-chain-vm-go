package contexts

import (
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// reservedFunctions holds the reserved function names.
//
// The VM hook names are included here only to reject contracts that export a function
// with the same name as a VM hook (see wasmValidator.verifyValidFunctionName).
// It plays no part in import resolution: the executor validates the hooks a contract
// imports separately, against the hooks it actually provides.
type reservedFunctions struct {
	epochsHandler           vmcommon.EnableEpochsHandler
	functionNames           vmcommon.FunctionNames
	functionActivationFlags map[string]core.EnableEpochFlag
	builtInFuncContainer    vmcommon.BuiltInFunctionContainer
}

// NewReservedFunctions creates a new reservedFunctions
func NewReservedFunctions(
	scAPINames vmcommon.FunctionNames,
	builtInFuncContainer vmcommon.BuiltInFunctionContainer,
	functionActivationFlags map[string]core.EnableEpochFlag,
	epochsHandler vmcommon.EnableEpochsHandler,
) *reservedFunctions {
	result := &reservedFunctions{
		functionNames:           make(vmcommon.FunctionNames),
		functionActivationFlags: functionActivationFlags,
		epochsHandler:           epochsHandler,
		builtInFuncContainer:    builtInFuncContainer,
	}

	for name, value := range scAPINames {
		result.functionNames[name] = value
	}

	var empty struct{}
	result.functionNames[vmhost.UpgradeFunctionName] = empty
	result.functionNames[vmhost.DeleteFunctionName] = empty

	return result
}

// IsReserved returns whether a function is reserved
func (reservedFunctions *reservedFunctions) IsReserved(functionName string) bool {
	function, err := reservedFunctions.builtInFuncContainer.Get(functionName)
	if err == nil && function.IsActive() {
		return true
	}

	if _, ok := reservedFunctions.functionNames[functionName]; ok {
		barnardOpcodesEnabled := reservedFunctions.epochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag)

		if !barnardOpcodesEnabled {
			if flag, ok := reservedFunctions.functionActivationFlags[functionName]; ok {
				return reservedFunctions.epochsHandler.IsFlagEnabled(flag)
			}
		}

		return true
	}

	return false
}

// GetReserved gets the reserved functions as a slice of strings
func (reservedFunctions *reservedFunctions) GetReserved() []string {
	keys := make([]string, 0, len(reservedFunctions.functionNames))

	barnardOpcodesEnabled := reservedFunctions.epochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag)

	for key := range reservedFunctions.functionNames {
		if !barnardOpcodesEnabled {
			if flag, ok := reservedFunctions.functionActivationFlags[key]; ok {
				if !reservedFunctions.epochsHandler.IsFlagEnabled(flag) {
					continue
				}
			}
		}

		keys = append(keys, key)
	}

	protocolFuncNames := reservedFunctions.builtInFuncContainer.Keys()
	for name := range protocolFuncNames {
		function, err := reservedFunctions.builtInFuncContainer.Get(name)
		if err != nil || !function.IsActive() {
			continue
		}

		keys = append(keys, name)
	}

	return keys
}

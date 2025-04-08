package contexts

import (
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// reservedFunctions holds the reserved function names
type reservedFunctions struct {
	epochsHandler           vmcommon.EnableEpochsHandler
	functionNames           vmcommon.FunctionNames
	functionActivationFlags map[string]core.EnableEpochFlag
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
	result.functionNames[vmhost.UpgradeFunctionName] = empty
	result.functionNames[vmhost.DeleteFunctionName] = empty

	return result
}

// IsReserved returns whether a function is reserved
func (reservedFunctions *reservedFunctions) IsReserved(functionName string) bool {
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
	keys := make([]string, len(reservedFunctions.functionNames))

	barnardOpcodesEnabled := reservedFunctions.epochsHandler.IsFlagEnabled(vmhost.BarnardOpcodesFlag)

	i := 0
	for key := range reservedFunctions.functionNames {
		if !barnardOpcodesEnabled {
			if flag, ok := reservedFunctions.functionActivationFlags[key]; ok {
				if !reservedFunctions.epochsHandler.IsFlagEnabled(flag) {
					continue
				}
			}
		}

		keys[i] = key
		i++
	}

	return keys
}

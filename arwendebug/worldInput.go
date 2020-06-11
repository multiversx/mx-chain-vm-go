package arwendebug

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (world *world) prepareDeployInput(request DeployRequest) *vmcommon.ContractCreateInput {
	createInput := &vmcommon.ContractCreateInput{}
	createInput.CallerAddr = request.ImpersonatedAsBytes
	createInput.CallValue = request.ValueAsBigInt
	createInput.ContractCode = request.CodeAsBytes
	createInput.ContractCodeMetadata = request.CodeMetadataAsBytes
	createInput.Arguments = request.ArgumentsAsBytes
	createInput.GasProvided = request.GasLimit
	createInput.GasPrice = request.GasPrice

	return createInput
}

func (world *world) prepareUpgradeInput(request UpgradeRequest) *vmcommon.ContractCallInput {
	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = []byte(request.ContractAddress)
	callInput.CallerAddr = request.ImpersonatedAsBytes
	callInput.CallValue = request.ValueAsBigInt
	callInput.Function = arwen.UpgradeFunctionName
	allArguments := make([][]byte, 0)
	allArguments = append(allArguments, request.CodeAsBytes, request.CodeMetadataAsBytes)
	allArguments = append(allArguments, request.ArgumentsAsBytes...)

	callInput.Arguments = allArguments
	callInput.GasProvided = request.GasLimit
	callInput.GasPrice = request.GasPrice

	return callInput
}

func (world *world) prepareCallInput(request RunRequest) *vmcommon.ContractCallInput {
	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = []byte(request.ContractAddressAsHex)
	callInput.CallerAddr = request.ImpersonatedAsBytes
	callInput.CallValue = request.ValueAsBigInt
	callInput.Function = request.Function
	callInput.Arguments = request.ArgumentsAsBytes
	callInput.GasProvided = request.GasLimit
	callInput.GasPrice = request.GasPrice

	return callInput
}

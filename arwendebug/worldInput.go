package arwendebug

import (
	"github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

func (w *world) prepareDeployInput(request DeployRequest) *vmcommon.ContractCreateInput {
	createInput := &vmcommon.ContractCreateInput{}
	createInput.CallerAddr = request.Impersonated
	createInput.CallValue = request.ValueAsBigInt
	createInput.ContractCode = request.Code
	createInput.ContractCodeMetadata = request.CodeMetadataBytes
	createInput.Arguments = request.Arguments
	createInput.GasProvided = request.GasLimit
	createInput.GasPrice = request.GasPrice

	return createInput
}

func (w *world) prepareUpgradeInput(request UpgradeRequest) *vmcommon.ContractCallInput {
	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = request.ContractAddress
	callInput.CallerAddr = request.Impersonated
	callInput.CallValue = request.ValueAsBigInt
	callInput.Function = arwen.UpgradeFunctionName
	allArguments := make([][]byte, 0)
	allArguments = append(allArguments, request.Code, request.CodeMetadataBytes)
	allArguments = append(allArguments, request.Arguments...)

	callInput.Arguments = allArguments
	callInput.GasProvided = request.GasLimit
	callInput.GasPrice = request.GasPrice

	return callInput
}

func (w *world) prepareCallInput(request RunRequest) *vmcommon.ContractCallInput {
	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = request.ContractAddress
	callInput.CallerAddr = request.Impersonated
	callInput.CallValue = request.ValueAsBigInt
	callInput.Function = request.Function
	callInput.Arguments = request.Arguments
	callInput.GasProvided = request.GasLimit
	callInput.GasPrice = request.GasPrice

	return callInput
}

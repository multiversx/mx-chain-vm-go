package arwendebug

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (world *world) prepareDeployInput(request DeployRequest) (*vmcommon.ContractCreateInput, error) {
	var err error

	createInput := &vmcommon.ContractCreateInput{}
	createInput.CallerAddr = request.getImpersonated()
	createInput.CallValue = request.getValue()
	createInput.ContractCode, err = request.getCode()
	if err != nil {
		return nil, err
	}

	createInput.ContractCodeMetadata, err = request.getCodeMetadata()
	if err != nil {
		return nil, err
	}

	createInput.Arguments, err = request.getArguments()
	if err != nil {
		return nil, err
	}

	createInput.GasProvided = request.getGasLimit()
	createInput.GasPrice = request.getGasPrice()

	return createInput, nil
}

func (world *world) prepareCallInput(request RunRequest) (*vmcommon.ContractCallInput, error) {
	var err error

	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = []byte(request.ContractAddress)
	callInput.CallerAddr = request.getImpersonated()
	callInput.CallValue = request.getValue()
	callInput.Function = request.Function
	callInput.Arguments, err = request.getArguments()
	if err != nil {
		return nil, err
	}

	callInput.GasProvided = request.getGasLimit()
	callInput.GasPrice = request.getGasPrice()

	return callInput, nil
}

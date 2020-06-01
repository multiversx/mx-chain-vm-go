package arwendebug

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (world *world) prepareDeployInput(request DeployRequest) (*vmcommon.ContractCreateInput, error) {
	var err error

	createInput := &vmcommon.ContractCreateInput{}
	createInput.CallerAddr, err = request.getImpersonated()
	if err != nil {
		return nil, err
	}

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

func (world *world) prepareUpgradeInput(request UpgradeRequest) (*vmcommon.ContractCallInput, error) {
	var err error

	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = []byte(request.ContractAddress)
	callInput.CallerAddr, err = request.getImpersonated()
	if err != nil {
		return nil, err
	}

	callInput.CallValue = request.getValue()
	callInput.Function = arwen.UpgradeFunctionName

	contractCode, err := request.getCode()
	if err != nil {
		return nil, err
	}

	contractCodeMetadata, err := request.getCodeMetadata()
	if err != nil {
		return nil, err
	}

	arguments, err := request.getArguments()
	if err != nil {
		return nil, err
	}

	allArguments := make([][]byte, 0)
	allArguments = append(allArguments, contractCode, contractCodeMetadata)
	allArguments = append(allArguments, arguments...)

	callInput.Arguments = allArguments
	callInput.GasProvided = request.getGasLimit()
	callInput.GasPrice = request.getGasPrice()

	return callInput, nil
}

func (world *world) prepareCallInput(request RunRequest) (*vmcommon.ContractCallInput, error) {
	var err error

	callInput := &vmcommon.ContractCallInput{}
	callInput.RecipientAddr = []byte(request.ContractAddress)
	callInput.CallerAddr, err = request.getImpersonated()
	if err != nil {
		return nil, err
	}

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

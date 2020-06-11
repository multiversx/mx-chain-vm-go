package arwendebug

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type worldDataModel struct {
	ID       string
	Accounts AccountsMap
}

type world struct {
	id             string
	blockchainHook *BlockchainHookMock
	vm             vmcommon.VMExecutionHandler
}

func newWorldDataModel(worldID string) *worldDataModel {
	return &worldDataModel{
		ID:       worldID,
		Accounts: make(AccountsMap),
	}
}

// newWorld creates a new debugging world
func newWorld(dataModel *worldDataModel) (*world, error) {
	blockchainHook := NewBlockchainHookMock()
	blockchainHook.Accounts = dataModel.Accounts

	vm, err := host.NewArwenVM(
		blockchainHook,
		arwenpart.NewCryptoHookGateway(),
		getHostParameters(),
	)
	if err != nil {
		return nil, err
	}

	return &world{
		id:             dataModel.ID,
		blockchainHook: blockchainHook,
		vm:             vm,
	}, nil
}

func getHostParameters() *arwen.VMHostParameters {
	return &arwen.VMHostParameters{
		VMType:                   []byte{5, 0},
		BlockGasLimit:            uint64(10000000),
		GasSchedule:              config.MakeGasMap(1, 1),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
	}
}

func (world *world) deploySmartContract(request DeployRequest) (*DeployResponse, error) {
	input := world.prepareDeployInput(request)
	log.Trace("world.deploySmartContract()", "input", prettyJson(input))

	vmOutput, err := world.vm.RunSmartContractCreate(input)
	if err == nil {
		world.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &DeployResponse{}
	response.Input = &input.VMInput
	response.Output = vmOutput
	response.Error = err
	response.ContractAddress = string(world.blockchainHook.LastCreatedContractAddress)
	return response, nil
}

func (world *world) upgradeSmartContract(request UpgradeRequest) (*UpgradeResponse, error) {
	input := world.prepareUpgradeInput(request)
	log.Trace("world.upgradeSmartContract()", "input", prettyJson(input))

	vmOutput, err := world.vm.RunSmartContractCall(input)
	if err == nil {
		world.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &UpgradeResponse{}
	response.Input = &input.VMInput
	response.Output = vmOutput
	response.Error = err

	return response, nil
}

func (world *world) runSmartContract(request RunRequest) (*RunResponse, error) {
	input := world.prepareCallInput(request)
	log.Trace("world.runSmartContract()", "input", prettyJson(input))

	vmOutput, err := world.vm.RunSmartContractCall(input)
	if err == nil {
		world.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &RunResponse{}
	response.Input = &input.VMInput
	response.Output = vmOutput
	response.Error = err

	return response, nil
}

func (world *world) querySmartContract(request QueryRequest) (*QueryResponse, error) {
	input := world.prepareCallInput(request.RunRequest)
	log.Trace("world.querySmartContract()", "input", prettyJson(input))

	vmOutput, err := world.vm.RunSmartContractCall(input)

	response := &QueryResponse{}
	response.Input = &input.VMInput
	response.Output = vmOutput
	response.Error = err

	return response, nil
}

func (world *world) createAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	log.Trace("world.createAccount()", "request", prettyJson(request))

	account := NewAccount(request.Address, request.Nonce, request.BalanceAsBigInt)
	world.blockchainHook.AddAccount(account)
	return &CreateAccountResponse{Account: account}, nil
}

func (world *world) toDataModel() *worldDataModel {
	return &worldDataModel{
		ID:       world.id,
		Accounts: world.blockchainHook.Accounts,
	}
}

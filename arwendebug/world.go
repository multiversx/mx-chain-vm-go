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

// NewWorld -
func NewWorld(dataModel *worldDataModel) (*world, error) {
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

// DeploySmartContract -
func (world *world) DeploySmartContract(request DeployRequest) (*DeployResponse, error) {
	log.Debug("world.DeploySmartContract()")

	createInput, err := world.prepareDeployInput(request)
	if err != nil {
		return nil, err
	}

	vmOutput, err := world.vm.RunSmartContractCreate(createInput)
	if err == nil {
		world.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &DeployResponse{}
	response.Input = &createInput.VMInput
	response.Output = vmOutput
	response.Error = err
	response.ContractAddress = string(world.blockchainHook.LastCreatedContractAddress)
	return response, nil
}

// UpgradeSmartContract -
func (world *world) UpgradeSmartContract() (*UpgradeResponse, error) {
	return &UpgradeResponse{}, nil
}

// RunSmartContract -
func (world *world) RunSmartContract(request RunRequest) (*RunResponse, error) {
	callInput, err := world.prepareCallInput(request)
	if err != nil {
		return nil, err
	}

	vmOutput, err := world.vm.RunSmartContractCall(callInput)
	if err == nil {
		world.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts)
	}

	response := &RunResponse{}
	response.Input = &callInput.VMInput
	response.Output = vmOutput
	response.Error = err

	return response, nil
}

// QuerySmartContract -
func (world *world) QuerySmartContract(request QueryRequest) (*QueryResponse, error) {
	callInput, err := world.prepareCallInput(request.RunRequest)
	if err != nil {
		return nil, err
	}

	vmOutput, err := world.vm.RunSmartContractCall(callInput)

	response := &QueryResponse{}
	response.Input = &callInput.VMInput
	response.Output = vmOutput
	response.Error = err

	return response, nil
}

func (world *world) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	balance, err := request.getBalance()
	if err != nil {
		return nil, err
	}

	account := &Account{
		Address: request.getAddress(),
		Nonce:   request.Nonce,
		Balance: balance,
	}

	world.blockchainHook.AddAccount(account)
	return &CreateAccountResponse{Account: account}, nil
}

func (world *world) toDataModel() *worldDataModel {
	return &worldDataModel{
		ID:       world.id,
		Accounts: world.blockchainHook.Accounts,
	}
}

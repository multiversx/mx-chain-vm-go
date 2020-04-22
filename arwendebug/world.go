package arwendebug

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type worldDataModel struct {
	id        string
	createdOn string
	accounts  mock.AccountsMap
}

type world struct {
	blockchainHook mock.BlockchainHookMock
	vm             vmcommon.VMExecutionHandler
}

func newWorldDataModel(worldID string) *worldDataModel {
	return &worldDataModel{
		id:        worldID,
		createdOn: "now",
		accounts:  make(mock.AccountsMap),
	}
}

// NewWorld -
func NewWorld(dataModel *worldDataModel) (*world, error) {
	blockchainHook := mock.NewBlockchainHookMock()
	blockchainHook.Accounts = dataModel.accounts

	vm, err := host.NewArwenVM(
		blockchainHook,
		arwenpart.NewCryptoHookGateway(),
		getHostParameters(),
	)
	if err != nil {
		return nil, err
	}

	return &world{
		blockchainHook: *blockchainHook,
		vm:             vm,
	}, nil
}

func getHostParameters() *arwen.VMHostParameters {
	return &arwen.VMHostParameters{
		VMType:        []byte{5, 0},
		BlockGasLimit: uint64(10000000),
		GasSchedule:   config.MakeGasMap(1),
	}
}

// DeploySmartContract -
func (world *world) DeploySmartContract(request DeployRequest) (*DeployResponse, error) {
	log.Debug("world.DeploySmartContract()")

	createInput := &vmcommon.ContractCreateInput{}
	vmOutput, err := world.vm.RunSmartContractCreate(createInput)

	response := &DeployResponse{}
	response.Output = *vmOutput
	response.Error = err

	return response, nil
}

// UpgradeSmartContract -
func (world *world) UpgradeSmartContract() (*UpgradeResponse, error) {
	return &UpgradeResponse{}, nil
}

// RunSmartContract -
func (world *world) RunSmartContract() error {
	return nil
}

// QuerySmartContract -
func (world *world) QuerySmartContract() error {
	return nil
}

func (world *world) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	balance, err := request.parseBalance()
	if err != nil {
		return nil, err
	}

	account := &mock.Account{
		Address: []byte(request.Address),
		Nonce:   request.Nonce,
		Balance: balance,
	}

	world.blockchainHook.AddAccount(account)
	return &CreateAccountResponse{}, nil
}

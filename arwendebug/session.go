package arwendebug

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type session struct {
	blockchainHook mock.BlockchainHookMock
	vm             vmcommon.VMExecutionHandler
}

// NewSession -
func NewSession(blockchainHook *mock.BlockchainHookMock) (*session, error) {
	vm, err := host.NewArwenVM(
		blockchainHook,
		arwenpart.NewCryptoHookGateway(),
		getHostParameters(),
	)
	if err != nil {
		return nil, err
	}

	return &session{
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
func (session *session) DeploySmartContract(request DeployRequest) (*DeployResponse, error) {
	log.Debug("Session.DeploySmartContract()")

	createInput := &vmcommon.ContractCreateInput{}
	vmOutput, err := session.vm.RunSmartContractCreate(createInput)

	response := &DeployResponse{}
	response.Output = *vmOutput
	response.Error = err

	return response, nil
}

// UpgradeSmartContract -
func (session *session) UpgradeSmartContract() (*UpgradeResponse, error) {
	return &UpgradeResponse{}, nil
}

// RunSmartContract -
func (session *session) RunSmartContract() error {
	return nil
}

// QuerySmartContract -
func (session *session) QuerySmartContract() error {
	return nil
}

func (session *session) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	balance, err := request.parseBalance()
	if err != nil {
		return nil, err
	}

	account := &mock.Account{
		Address: []byte(request.Address),
		Nonce:   request.Nonce,
		Balance: balance,
	}

	session.blockchainHook.AddAccount(account)
	return &CreateAccountResponse{}, nil
}

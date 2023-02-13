package vmserver

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/config"
	worldmock "github.com/multiversx/mx-chain-vm-go/mock/world"
	"github.com/multiversx/mx-chain-vm-go/testcommon/testexecutor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
)

var defaultHasher = blake2b.NewBlake2b()

type worldDataModel struct {
	ID       string
	Accounts worldmock.AccountMap
}

type world struct {
	id             string
	blockchainHook *worldmock.MockWorld
	vm             vmcommon.VMExecutionHandler
}

func newWorldDataModel(worldID string) *worldDataModel {
	return &worldDataModel{
		ID:       worldID,
		Accounts: worldmock.NewAccountMap(),
	}
}

// newWorld creates a new debugging world
func newWorld(tb testing.TB, dataModel *worldDataModel) (*world, error) {
	blockchainHook := worldmock.NewMockWorld()
	blockchainHook.AcctMap = dataModel.Accounts

	vm, err := hostCore.NewVMHost(
		blockchainHook,
		getHostParameters(tb),
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

func getHostParameters(tb testing.TB) *vmhost.VMHostParameters {
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	return &vmhost.VMHostParameters{
		VMType:                   []byte{5, 0},
		OverrideVMExecutor:       testexecutor.NewDefaultTestExecutorFactory(tb),
		BlockGasLimit:            uint64(10000000),
		GasSchedule:              config.MakeGasMap(1, 1),
		ProtectedKeyPrefix:       []byte("E" + "L" + "R" + "O" + "N" + "D"),
		BuiltInFuncContainer:     builtInFunctions.NewBuiltInFunctionContainer(),
		ESDTTransferParser:       esdtTransferParser,
		EpochNotifier:            &mock.EpochNotifierStub{},
		EnableEpochsHandler:      worldmock.EnableEpochsHandlerStubNoFlags(),
		WasmerSIGSEGVPassthrough: false,
		Hasher:                   defaultHasher,
	}
}

func (w *world) deploySmartContract(request DeployRequest) *DeployResponse {
	input := w.prepareDeployInput(request)
	log.Trace("w.deploySmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCreate(input)
	if err == nil {
		_ = w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts, nil)
	}

	response := &DeployResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err
	response.ContractAddress = w.blockchainHook.LastCreatedContractAddress
	response.ContractAddressHex = toHex(response.ContractAddress)
	return response
}

func (w *world) upgradeSmartContract(request UpgradeRequest) *UpgradeResponse {
	input := w.prepareUpgradeInput(request)
	log.Trace("w.upgradeSmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)
	if err == nil {
		_ = w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts, nil)
	}

	response := &UpgradeResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) runSmartContract(request RunRequest) *RunResponse {
	input := w.prepareCallInput(request)
	log.Trace("w.runSmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)
	if err == nil {
		_ = w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts, nil)
	}

	response := &RunResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) querySmartContract(request QueryRequest) *QueryResponse {
	input := w.prepareCallInput(request.RunRequest)
	log.Trace("w.querySmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)

	response := &QueryResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) createAccount(request CreateAccountRequest) *CreateAccountResponse {
	log.Trace("w.createAccount()", "request", prettyJson(request))

	account := worldmock.Account{
		Address:         request.Address,
		Nonce:           request.Nonce,
		Balance:         request.BalanceAsBigInt,
		BalanceDelta:    big.NewInt(0),
		DeveloperReward: big.NewInt(0),
	}
	w.blockchainHook.AcctMap.PutAccount(&account)
	return &CreateAccountResponse{Account: &account}
}

func (w *world) toDataModel() *worldDataModel {
	accounts := w.blockchainHook.AcctMap.Clone()
	for _, account := range accounts {
		account.MockWorld = nil
	}

	return &worldDataModel{
		ID:       w.id,
		Accounts: accounts,
	}
}

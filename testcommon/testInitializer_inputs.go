package testcommon

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen"
	arwenHost "github.com/multiversx/mx-chain-vm-v1_4-go/arwen/host"
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen/mock"
	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	"github.com/multiversx/mx-chain-vm-v1_4-go/crypto/hashing"
	contextmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	"github.com/stretchr/testify/require"
)

// DefaultVMType is an exposed value to use in tests
var DefaultVMType = []byte{0xF, 0xF}

// ErrAccountNotFound is an exposed value to use in tests
var ErrAccountNotFound = errors.New("account not found")

// UserAddress is an exposed value to use in tests
var UserAddress = []byte("userAccount.....................")

// UserAddress2 is an exposed value to use in tests
var UserAddress2 = []byte("userAccount2....................")

// AddressSize is the size of an account address, in bytes.
const AddressSize = 32

// SCAddressPrefix is the prefix of any smart contract address used for testing.
var SCAddressPrefix = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0f")

// ParentAddress is an exposed value to use in tests
var ParentAddress = MakeTestSCAddress("parentSC")

// ChildAddress is an exposed value to use in tests
var ChildAddress = MakeTestSCAddress("childSC")

// NephewAddress is an exposed value to use in tests
var NephewAddress = MakeTestSCAddress("NephewAddress")

// ESDTTransferGasCost is an exposed value to use in tests
var ESDTTransferGasCost = uint64(1)

// ESDTTestTokenName is an exposed value to use in tests
var ESDTTestTokenName = []byte("TTT-010101")

// DefaultCodeMetadata is an exposed value to use in tests
var DefaultCodeMetadata = []byte{3, 0}

// MakeTestSCAddress generates a new smart contract address to be used for
// testing based on the given identifier.
func MakeTestSCAddress(identifier string) []byte {
	numberOfTrailingDots := AddressSize - len(SCAddressPrefix) - len(identifier)
	leftBytes := SCAddressPrefix
	rightBytes := []byte(identifier + strings.Repeat(".", numberOfTrailingDots))
	return append(leftBytes, rightBytes...)
}

// GetSCCode retrieves the bytecode of a WASM module from a file
func GetSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("GetSCCode(): %s", fileName))
	}

	return code
}

// GetTestSCCode retrieves the bytecode of a WASM testing contract
func GetTestSCCode(scName string, prefixToTestSCs ...string) []byte {
	var searchedPaths []string
	for _, prefixToTestSC := range prefixToTestSCs {
		pathToSC := prefixToTestSC + "test/contracts/" + scName + "/output/" + scName + ".wasm"
		searchedPaths = append(searchedPaths, pathToSC)
		code, err := ioutil.ReadFile(filepath.Clean(pathToSC))
		if err == nil {
			return code
		}
	}
	panic(fmt.Sprintf("GetSCCode(): %s", searchedPaths))
}

// GetTestSCCodeModule retrieves the bytecode of a WASM testing contract, given
// a specific name of the WASM module
func GetTestSCCodeModule(scName string, moduleName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/output/" + moduleName + ".wasm"
	return GetSCCode(pathToSC)
}

// DefaultTestArwenForDeployment creates an Arwen vmHost configured for testing deployments
func DefaultTestArwenForDeployment(t *testing.T, _ uint64, newAddress []byte) (arwen.VMHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host := DefaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

// DefaultTestArwenForCall creates a BlockchainHookStub
func DefaultTestArwenForCall(tb testing.TB, code []byte, balance *big.Int) (arwen.VMHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, ParentAddress) {
			return &contextmock.StubAccount{
				Balance: balance,
			}, nil
		}
		return nil, ErrAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		return code
	}

	host := DefaultTestArwen(tb, stubBlockchainHook)
	return host, stubBlockchainHook
}

// DefaultTestArwenForCallSigSegv creates a BlockchainHookStub and a host
func DefaultTestArwenForCallSigSegv(tb testing.TB, code []byte, balance *big.Int) (arwen.VMHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, ParentAddress) {
			return &contextmock.StubAccount{
				Balance: balance,
			}, nil
		}
		return nil, ErrAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		return code
	}

	customGasSchedule := config.GasScheduleMap(nil)
	host := DefaultTestArwenWithGasSchedule(tb, stubBlockchainHook, customGasSchedule, true)
	return host, stubBlockchainHook
}

// DefaultTestArwenForCallWithInstanceRecorderMock creates an InstanceBuilderRecorderMock
func DefaultTestArwenForCallWithInstanceRecorderMock(tb testing.TB, code []byte, balance *big.Int) (arwen.VMHost, *contextmock.InstanceBuilderRecorderMock) {
	// this uses a Blockchain Hook Stub that does not cache the compiled code
	host, _ := DefaultTestArwenForCall(tb, code, balance)

	instanceBuilderRecorderMock := contextmock.NewInstanceBuilderRecorderMock()
	host.Runtime().ReplaceInstanceBuilder(instanceBuilderRecorderMock)

	return host, instanceBuilderRecorderMock
}

// DefaultTestArwenForCallWithInstanceMocks creates an InstanceBuilderMock
func DefaultTestArwenForCallWithInstanceMocks(tb testing.TB) (arwen.VMHost, *worldmock.MockWorld, *contextmock.InstanceBuilderMock) {
	world := worldmock.NewMockWorld()
	host := DefaultTestArwen(tb, world)

	instanceBuilderMock := contextmock.NewInstanceBuilderMock(world)
	host.Runtime().ReplaceInstanceBuilder(instanceBuilderMock)

	return host, world, instanceBuilderMock
}

// DefaultTestArwenForCallWithWorldMock creates a MockWorld
func DefaultTestArwenForCallWithWorldMock(tb testing.TB, code []byte, balance *big.Int) (arwen.VMHost, *worldmock.MockWorld) {
	world := worldmock.NewMockWorld()
	host := DefaultTestArwen(tb, world)

	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	host.SetBuiltInFunctionsContainer(world.BuiltinFuncs.Container)

	parentAccount := world.AcctMap.CreateSmartContractAccount(UserAddress, ParentAddress, code, world)
	parentAccount.Balance = balance

	return host, world
}

// DefaultTestArwenForTwoSCs creates an Arwen vmHost configured for testing calls between 2 SmartContracts
func DefaultTestArwenForTwoSCs(
	t *testing.T,
	parentCode []byte,
	childCode []byte,
	parentSCBalance *big.Int,
	childSCBalance *big.Int,
) (arwen.VMHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	if parentSCBalance == nil {
		parentSCBalance = big.NewInt(1000)
	}

	if childSCBalance == nil {
		childSCBalance = big.NewInt(1000)
	}

	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, ParentAddress) {
			return &contextmock.StubAccount{
				Address: ParentAddress,
				Balance: parentSCBalance,
			}, nil
		}
		if bytes.Equal(scAddress, ChildAddress) {
			return &contextmock.StubAccount{
				Address: ChildAddress,
				Balance: childSCBalance,
			}, nil
		}

		return nil, ErrAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		if bytes.Equal(account.AddressBytes(), ParentAddress) {
			return parentCode
		}
		if bytes.Equal(account.AddressBytes(), ChildAddress) {
			return childCode
		}
		return nil
	}

	host := DefaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwenForContracts(
	tb testing.TB,
	contracts []*InstanceTestSmartContract,
	gasSchedule config.GasScheduleMap,
	wasmerSIGSEGVPassthrough bool,
) (arwen.VMHost, *contextmock.BlockchainHookStub) {

	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	contractsMap := make(map[string]*contextmock.StubAccount)
	codeMap := make(map[string]*[]byte)
	compiledCodeMap := make(map[string][]byte)

	for _, contract := range contracts {
		codeHash, _ := hashing.NewHasher().Sha256(contract.code)
		contractsMap[string(contract.address)] = &contextmock.StubAccount{
			Address:      contract.address,
			Balance:      big.NewInt(contract.balance),
			CodeHash:     codeHash,
			CodeMetadata: DefaultCodeMetadata,
			OwnerAddress: ParentAddress,
		}
		codeMap[string(contract.address)] = &contract.code
	}

	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		contract, found := contractsMap[string(scAddress)]
		if found {
			return contract, nil
		}
		return nil, ErrAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		code, found := codeMap[string(account.AddressBytes())]
		if found {
			return *code
		}
		return nil
	}

	stubBlockchainHook.SaveCompiledCodeCalled = func(codehash []byte, code []byte) {
		compiledCodeMap[string(codehash)] = code
	}
	stubBlockchainHook.GetCompiledCodeCalled = func(codeHash []byte) (bool, []byte) {
		compiledCode, ok := compiledCodeMap[string(codeHash)]
		if ok {
			return ok, compiledCode
		}
		return false, nil
	}

	host := DefaultTestArwenWithGasSchedule(tb, stubBlockchainHook, gasSchedule, wasmerSIGSEGVPassthrough)
	return host, stubBlockchainHook
}

// DefaultTestArwenWithWorldMock creates a host configured with a mock world
func DefaultTestArwenWithWorldMock(tb testing.TB) (arwen.VMHost, *worldmock.MockWorld) {
	customGasSchedule := config.GasScheduleMap(nil)
	return DefaultTestArwenWithWorldMockWithGasSchedule(tb, customGasSchedule)
}

// DefaultTestArwenWithWorldMockWithGasSchedule creates a host configured with a mock world
func DefaultTestArwenWithWorldMockWithGasSchedule(tb testing.TB, customGasSchedule config.GasScheduleMap) (arwen.VMHost, *worldmock.MockWorld) {
	world := worldmock.NewMockWorld()
	gasSchedule := customGasSchedule
	if gasSchedule == nil {
		gasSchedule = config.MakeGasMapForTests()
	}
	err := world.InitBuiltinFunctions(gasSchedule)
	require.Nil(tb, err)

	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	host, err := arwenHost.NewArwenVM(world, &arwen.VMHostParameters{
		VMType:               DefaultVMType,
		BlockGasLimit:        uint64(1000),
		GasSchedule:          gasSchedule,
		BuiltInFuncContainer: world.BuiltinFuncs.Container,
		ProtectedKeyPrefix:   []byte("ELROND"),
		ESDTTransferParser:   esdtTransferParser,
		EpochNotifier:        &mock.EpochNotifierStub{},
		EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
			IsStorageAPICostOptimizationFlagEnabledField:     true,
			IsMultiESDTTransferFixOnCallBackFlagEnabledField: true,
			IsFixOOGReturnCodeFlagEnabledField:               true,
			IsRemoveNonUpdatedStorageFlagEnabledField:        true,
			IsCreateNFTThroughExecByCallerFlagEnabledField:   true,
			IsManagedCryptoAPIsFlagEnabledField:              true,
			IsFailExecutionOnEveryAPIErrorFlagEnabledField:   true,
			IsESDTTransferRoleFlagEnabledField:               true,
			IsSendAlwaysFlagEnabledField:                     true,
			IsGlobalMintBurnFlagEnabledField:                 true,
			IsCheckFunctionArgumentFlagEnabledField:          true,
			IsCheckExecuteOnReadOnlyFlagEnabledField:         true,
		},
		WasmerSIGSEGVPassthrough: false,
		Hasher:                   worldmock.DefaultHasher,
	})
	require.Nil(tb, err)
	require.NotNil(tb, host)

	return host, world
}

// DefaultTestArwen creates a host configured with a blockchain hook
func DefaultTestArwen(tb testing.TB, blockchain vmcommon.BlockchainHook) arwen.VMHost {
	customGasSchedule := config.GasScheduleMap(nil)
	return DefaultTestArwenWithGasSchedule(tb, blockchain, customGasSchedule, false)
}

// DefaultTestArwenWithGasSchedule creates a host with the provided blockchain hook and gas schedule
func DefaultTestArwenWithGasSchedule(
	tb testing.TB,
	blockchain vmcommon.BlockchainHook,
	customGasSchedule config.GasScheduleMap,
	wasmerSIGSEGVPassthrough bool,
) arwen.VMHost {
	gasSchedule := customGasSchedule
	if gasSchedule == nil {
		gasSchedule = config.MakeGasMapForTests()
	}

	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	host, err := arwenHost.NewArwenVM(blockchain, &arwen.VMHostParameters{
		VMType:               DefaultVMType,
		BlockGasLimit:        uint64(1000),
		GasSchedule:          gasSchedule,
		BuiltInFuncContainer: builtInFunctions.NewBuiltInFunctionContainer(),
		ProtectedKeyPrefix:   []byte("ELROND"),
		ESDTTransferParser:   esdtTransferParser,
		EpochNotifier:        &mock.EpochNotifierStub{},
		EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
			IsStorageAPICostOptimizationFlagEnabledField:         true,
			IsMultiESDTTransferFixOnCallBackFlagEnabledField:     true,
			IsFixOOGReturnCodeFlagEnabledField:                   true,
			IsRemoveNonUpdatedStorageFlagEnabledField:            true,
			IsCreateNFTThroughExecByCallerFlagEnabledField:       true,
			IsManagedCryptoAPIsFlagEnabledField:                  true,
			IsFailExecutionOnEveryAPIErrorFlagEnabledField:       true,
			IsRefactorContextFlagEnabledField:                    true,
			IsCheckCorrectTokenIDForTransferRoleFlagEnabledField: true,
			IsDisableExecByCallerFlagEnabledField:                true,
			IsESDTTransferRoleFlagEnabledField:                   true,
			IsSendAlwaysFlagEnabledField:                         true,
			IsGlobalMintBurnFlagEnabledField:                     true,
			IsCheckFunctionArgumentFlagEnabledField:              true,
			IsCheckExecuteOnReadOnlyFlagEnabledField:             true,
		},
		WasmerSIGSEGVPassthrough: wasmerSIGSEGVPassthrough,
		Hasher:                   worldmock.DefaultHasher,
	})
	require.Nil(tb, err)
	require.NotNil(tb, host)

	return host
}

// AddTestSmartContractToWorld directly deploys the provided code into the
// given MockWorld under a SC address built with the given identifier.
func AddTestSmartContractToWorld(world *worldmock.MockWorld, identifier string, code []byte) *worldmock.Account {
	address := MakeTestSCAddress(identifier)
	return world.AcctMap.CreateSmartContractAccount(UserAddress, address, code, world)
}

// DefaultTestContractCreateInput creates a vmcommon.ContractCreateInput struct
// with default values.
func DefaultTestContractCreateInput() *vmcommon.ContractCreateInput {
	return &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments: [][]byte{
				[]byte("argument 1"),
				[]byte("argument 2"),
			},
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		ContractCode: []byte("contract"),
	}
}

// DefaultTestContractCallInput creates a vmcommon.ContractCallInput struct
// with default values.
func DefaultTestContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  UserAddress,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: ParentAddress,
		Function:      "function",
	}
}

// ContractCallInputBuilder extends a ContractCallInput for extra building functionality during testing
type ContractCallInputBuilder struct {
	vmcommon.ContractCallInput
}

// CreateTestContractCallInputBuilder is a builder for ContractCallInputBuilder
func CreateTestContractCallInputBuilder() *ContractCallInputBuilder {
	return &ContractCallInputBuilder{
		ContractCallInput: *DefaultTestContractCallInput(),
	}
}

// WithRecipientAddr provides the recepient address of ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithRecipientAddr(address []byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.RecipientAddr = address
	return contractInput
}

// WithCallerAddr provides the caller address of ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithCallerAddr(address []byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.CallerAddr = address
	return contractInput
}

// WithCallValue provides the value transferred to the called contract
func (contractInput *ContractCallInputBuilder) WithCallValue(value int64) *ContractCallInputBuilder {
	contractInput.ContractCallInput.CallValue = big.NewInt(value)
	return contractInput
}

// WithGasProvided provides the gas of ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithGasProvided(gas uint64) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.GasProvided = gas
	return contractInput
}

// WithFunction provides the function to be called for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithFunction(function string) *ContractCallInputBuilder {
	contractInput.ContractCallInput.Function = function
	return contractInput
}

// WithArguments provides the arguments to be called for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithArguments(arguments ...[]byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.Arguments = arguments
	return contractInput
}

// WithCallType provides the arguments to be called for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithCallType(callType vm.CallType) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.CallType = callType
	return contractInput
}

// WithCurrentTxHash provides the CurrentTxHash for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithCurrentTxHash(txHash []byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.CurrentTxHash = txHash
	return contractInput
}

func (contractInput *ContractCallInputBuilder) initESDTTransferIfNeeded() {
	if len(contractInput.ESDTTransfers) == 0 {
		contractInput.ESDTTransfers = make([]*vmcommon.ESDTTransfer, 1)
		contractInput.ESDTTransfers[0] = &vmcommon.ESDTTransfer{}
	}
}

// WithESDTValue provides the ESDTValue for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithESDTValue(esdtValue *big.Int) *ContractCallInputBuilder {
	contractInput.initESDTTransferIfNeeded()
	contractInput.ContractCallInput.ESDTTransfers[0].ESDTValue = esdtValue
	return contractInput
}

// WithESDTTokenName provides the ESDTTokenName for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithESDTTokenName(esdtTokenName []byte) *ContractCallInputBuilder {
	contractInput.initESDTTransferIfNeeded()
	contractInput.ContractCallInput.ESDTTransfers[0].ESDTTokenName = esdtTokenName
	return contractInput
}

// Build completes the build of a ContractCallInput
func (contractInput *ContractCallInputBuilder) Build() *vmcommon.ContractCallInput {
	return &contractInput.ContractCallInput
}

// ContractCreateInputBuilder extends a ContractCreateInput for extra building functionality during testing
type ContractCreateInputBuilder struct {
	vmcommon.ContractCreateInput
}

// CreateTestContractCreateInputBuilder is a builder for ContractCreateInputBuilder
func CreateTestContractCreateInputBuilder() *ContractCreateInputBuilder {
	return &ContractCreateInputBuilder{
		ContractCreateInput: *DefaultTestContractCreateInput(),
	}
}

// WithGasProvided provides the GasProvided for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithGasProvided(gas uint64) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.GasProvided = gas
	return contractInput
}

// WithContractCode provides the ContractCode for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithContractCode(code []byte) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.ContractCode = code
	return contractInput
}

// WithCallerAddr provides the CallerAddr for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithCallerAddr(address []byte) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.CallerAddr = address
	return contractInput
}

// WithCallValue provides the CallValue for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithCallValue(callValue int64) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.CallValue = big.NewInt(callValue)
	return contractInput
}

// WithArguments provides the Arguments for a ContractCreateInputBuilder
func (contractInput *ContractCreateInputBuilder) WithArguments(arguments ...[]byte) *ContractCreateInputBuilder {
	contractInput.ContractCreateInput.Arguments = arguments
	return contractInput
}

// Build completes the build of a ContractCreateInput
func (contractInput *ContractCreateInputBuilder) Build() *vmcommon.ContractCreateInput {
	return &contractInput.ContractCreateInput
}

// MakeVMOutput creates a vmcommon.VMOutput struct with default values
func MakeVMOutput() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		ReturnData:      make([][]byte, 0),
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
	}
}

// AddNewOutputAccount creates a new vmcommon.OutputAccount from the provided arguments and adds it to OutputAccounts of the provided vmOutput
func AddNewOutputAccount(vmOutput *vmcommon.VMOutput, sender []byte, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	account := &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(balanceDelta),
		Balance:        nil,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		Code:           nil,
	}
	if data != nil {
		account.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				Value:         big.NewInt(balanceDelta),
				SenderAddress: sender,
			},
		}
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

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

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	worldmock "github.com/multiversx/mx-chain-vm-go/mock/world"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

var defaultHasher = blake2b.NewBlake2b()

// DefaultVMType is an exposed value to use in tests
var DefaultVMType = []byte{0xF, 0xF}

// ErrAccountNotFound is an exposed value to use in tests
var ErrAccountNotFound = errors.New("account not found")

// UserAddress is an exposed value to use in tests
var UserAddress = MakeTestSCAddressWithDefaultVM("userAccount")

// UserAddress2 is an exposed value to use in tests
var UserAddress2 = []byte("userAccount2....................")

// AddressSize is the size of an account address, in bytes.
const AddressSize = 32

// SCAddressPrefix is the prefix of any smart contract address used for testing.
var SCAddressPrefix = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0f")

// ParentAddress is an exposed value to use in tests
var ParentAddress = MakeTestSCAddressWithDefaultVM("parentSC")

// ChildAddress is an exposed value to use in tests
var ChildAddress = MakeTestSCAddressWithDefaultVM("childSC")

// NephewAddress is an exposed value to use in tests
var NephewAddress = MakeTestSCAddressWithDefaultVM("NephewAddress")

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

// MakeTestSCAddressWithDefaultVM generates a new smart contract address to be used for
// testing based on the given identifier.
func MakeTestSCAddressWithDefaultVM(identifier string) []byte {
	return MakeTestSCAddressWithVMType(identifier, worldmock.DefaultVMType)
}

// MakeTestSCAddressWithVMType generates a new smart contract address to be used for
// testing based on the given identifier.
func MakeTestSCAddressWithVMType(identifier string, vmType []byte) []byte {
	address := MakeTestSCAddress(identifier)
	copy(address[vmcommon.NumInitCharactersForScAddress-core.VMTypeLen:], vmType)
	return address
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

// TestHostBuilder allows tests to configure and initialize the VM host and blockhain mock on which they operate.
type TestHostBuilder struct {
	tb               testing.TB
	blockchainHook   vmcommon.BlockchainHook
	vmHostParameters *vmhost.VMHostParameters
	host             vmhost.VMHost
}

// NewTestHostBuilder commences a test host builder pattern.
func NewTestHostBuilder(tb testing.TB) *TestHostBuilder {
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	return &TestHostBuilder{
		tb: tb,
		vmHostParameters: &vmhost.VMHostParameters{
			VMType:                   DefaultVMType,
			BlockGasLimit:            uint64(1000),
			GasSchedule:              nil,
			BuiltInFuncContainer:     nil,
			ProtectedKeyPrefix:       []byte("E" + "L" + "R" + "O" + "N" + "D"),
			ESDTTransferParser:       esdtTransferParser,
			EpochNotifier:            &mock.EpochNotifierStub{},
			EnableEpochsHandler:      worldmock.EnableEpochsHandlerStubAllFlags(),
			WasmerSIGSEGVPassthrough: false,
			Hasher:                   defaultHasher,
		},
	}
}

// Ensures gas costs are initialized.
func (thb *TestHostBuilder) initializeGasCosts() {
	if thb.vmHostParameters.GasSchedule == nil {
		thb.vmHostParameters.GasSchedule = config.MakeGasMapForTests()
	}
}

// Ensures the built-in function container is initialized.
func (thb *TestHostBuilder) initializeBuiltInFuncContainer() {
	if thb.vmHostParameters.BuiltInFuncContainer == nil {
		thb.vmHostParameters.BuiltInFuncContainer = builtInFunctions.NewBuiltInFunctionContainer()
	}

}

// WithBlockchainHook sets a pre-built blockchain hook for the VM to work with.
func (thb *TestHostBuilder) WithBlockchainHook(blockchainHook vmcommon.BlockchainHook) *TestHostBuilder {
	thb.blockchainHook = blockchainHook
	return thb
}

// WithBuiltinFunctions sets up builtin functions in the blockchain hook.
// Only works if the blockchain hook is of type worldmock.MockWorld.
func (thb *TestHostBuilder) WithBuiltinFunctions() *TestHostBuilder {
	thb.initializeGasCosts()
	mockWorld, ok := thb.blockchainHook.(*worldmock.MockWorld)
	require.True(thb.tb, ok, "builtin functions can only be injected into blockchain hooks of type MockWorld")
	err := mockWorld.InitBuiltinFunctions(thb.vmHostParameters.GasSchedule)
	require.Nil(thb.tb, err)
	thb.vmHostParameters.BuiltInFuncContainer = mockWorld.BuiltinFuncs.Container
	return thb
}

// WithExecutorFactory allows tests to choose what executor to use. The default is wasmer 1.
func (thb *TestHostBuilder) WithExecutorFactory(executorFactory executor.ExecutorAbstractFactory) *TestHostBuilder {
	thb.vmHostParameters.OverrideVMExecutor = executorFactory
	return thb
}

// WithWasmerSIGSEGVPassthrough allows tests to configure the WasmerSIGSEGVPassthrough flag.
func (thb *TestHostBuilder) WithWasmerSIGSEGVPassthrough(wasmerSIGSEGVPassthrough bool) *TestHostBuilder {
	thb.vmHostParameters.WasmerSIGSEGVPassthrough = wasmerSIGSEGVPassthrough
	return thb
}

// WithGasSchedule allows tests to use the gas costs. The default is config.MakeGasMapForTests().
func (thb *TestHostBuilder) WithGasSchedule(gasSchedule config.GasScheduleMap) *TestHostBuilder {
	thb.vmHostParameters.GasSchedule = gasSchedule
	return thb
}

// Build initializes the VM host with all configured options.
func (thb *TestHostBuilder) Build() vmhost.VMHost {
	thb.initializeHost()
	return thb.host
}

func (thb *TestHostBuilder) initializeHost() {
	thb.initializeGasCosts()
	if thb.host == nil {
		thb.host = thb.newHost()
	}
}

func (thb *TestHostBuilder) newHost() vmhost.VMHost {
	thb.initializeBuiltInFuncContainer()
	host, err := hostCore.NewVMHost(
		thb.blockchainHook,
		thb.vmHostParameters,
	)
	require.Nil(thb.tb, err)
	require.NotNil(thb.tb, host)

	return host
}

// BlockchainHookStubForCallSigSegv -
func BlockchainHookStubForCallSigSegv(code []byte, balance *big.Int) *contextmock.BlockchainHookStub {
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
	return stubBlockchainHook
}

// BlockchainHookStubForCall creates a BlockchainHookStub
func BlockchainHookStubForCall(code []byte, balance *big.Int) *contextmock.BlockchainHookStub {
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

	return stubBlockchainHook
}

// BlockchainHookStubForTwoSCs creates a world stub configured for testing calls between 2 SmartContracts
func BlockchainHookStubForTwoSCs(
	parentCode []byte,
	childCode []byte,
	parentSCBalance *big.Int,
	childSCBalance *big.Int,
) *contextmock.BlockchainHookStub {
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

	return stubBlockchainHook
}

// BlockchainHookStubForContracts -
func BlockchainHookStubForContracts(
	contracts []*InstanceTestSmartContract,
) *contextmock.BlockchainHookStub {

	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	contractsMap := make(map[string]*contextmock.StubAccount)
	codeMap := make(map[string]*[]byte)

	for _, contract := range contracts {
		codeHash := defaultHasher.Compute(string(contract.code))
		contractsMap[string(contract.address)] = &contextmock.StubAccount{
			Address:      contract.address,
			Balance:      big.NewInt(contract.balance),
			CodeHash:     codeHash,
			CodeMetadata: DefaultCodeMetadata,
			OwnerAddress: contract.ownerAddress,
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

	return stubBlockchainHook
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
			OriginalCallerAddr: UserAddress,
			CallerAddr:         UserAddress,
			Arguments:          make([][]byte, 0),
			CallValue:          big.NewInt(0),
			CallType:           vm.DirectCall,
			GasPrice:           0,
			GasProvided:        0,
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

// WithGasLocked provides the locked gas of ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithGasLocked(gas uint64) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.GasLocked = gas
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

// WithAsyncArguments provides the async arguments to be called for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithAsyncArguments(arguments *vmcommon.AsyncArguments) *ContractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.AsyncArguments = arguments
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

// WithPrevTxHash provides the PrevTxHash for ContractCallInputBuilder
func (contractInput *ContractCallInputBuilder) WithPrevTxHash(txHash []byte) *ContractCallInputBuilder {
	contractInput.ContractCallInput.PrevTxHash = txHash
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

package arwen

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/config"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
)

// DefaultVMType is an exposed value to use in tests
var DefaultVMType = []byte{0xF, 0xF}

// ErrAccountNotFound is an exposed value to use in tests
var ErrAccountNotFound = errors.New("account not found")

// UserAddress is an exposed value to use in tests
var UserAddress = []byte("userAccount.....................")

// AddressSize is the size of an account address, in bytes.
const AddressSize = 32

// SCAddressPrefix is the prefix of any smart contract address used for testing.
var SCAddressPrefix = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0f")

// WalletAddressPrefix is the prefix of any smart contract address used for testing.
var WalletAddressPrefix = []byte("..........")

// ParentAddress is an exposed value to use in tests
var ParentAddress = MakeTestSCAddress("parentSC")

// ChildAddress is an exposed value to use in tests
var ChildAddress = MakeTestSCAddress("childSC")

var customGasSchedule = config.GasScheduleMap(nil)

// ESDTTransferGasCost is an exposed value to use in tests
var ESDTTransferGasCost = uint64(1)

// ESDTTestTokenName is an exposed value to use in tests
var ESDTTestTokenName = []byte("TT")

// DefaultCodeMetadata is an exposed value to use in tests
var DefaultCodeMetadata = []byte{3, 0}

// MakeTestSCAddress generates a new smart contract address to be used for
// testing based on the given identifier.
func MakeTestSCAddress(identifier string) []byte {
	return makeTestAddress(SCAddressPrefix, identifier)
}

// MakeTestWalletAddress generates a new wallet address to be used for
// testing based on the given identifier.
func MakeTestWalletAddress(identifier string) []byte {
	return makeTestAddress(WalletAddressPrefix, identifier)
}

func makeTestAddress(prefix []byte, identifier string) []byte {
	numberOfTrailingDots := AddressSize - len(SCAddressPrefix) - len(identifier)
	leftBytes := SCAddressPrefix
	rightBytes := []byte(identifier + strings.Repeat(".", numberOfTrailingDots))
	return append(leftBytes, rightBytes...)
}

// AddTestSmartContractToWorld directly deploys the provided code into the
// given MockWorld under a SC address built with the given identifier.
func AddTestSmartContractToWorld(world *worldmock.MockWorld, identifier string, code []byte) *worldmock.Account {
	address := MakeTestSCAddress(identifier)
	return world.AcctMap.CreateSmartContractAccount(UserAddress, address, code, world)
}

// MakeEmptyContractCallInput instantiates an empty ContractCallInput
func MakeEmptyContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:           nil,
			Arguments:            make([][]byte, 0),
			CallValue:            big.NewInt(0),
			CallType:             vm.DirectCall,
			GasPrice:             1,
			GasProvided:          0,
			ReturnCallAfterError: false,
		},
		RecipientAddr: nil,
		Function:      "",
	}
}

// MakeContractCallInput creates a ContractCallInput and sets the provided arguments
func MakeContractCallInput(
	caller []byte,
	recipient []byte,
	function string,
	value int,
) *vmcommon.ContractCallInput {
	input := MakeEmptyContractCallInput()
	SetCallParties(input, caller, recipient)
	input.Function = function
	input.CallValue = big.NewInt(int64(value))
	return input
}

// SetCallParties sets the caller and recipient of the given ContractCallInput
func SetCallParties(input *vmcommon.ContractCallInput, caller []byte, recipient []byte) {
	input.CallerAddr = caller
	input.RecipientAddr = recipient
}

// AddArgument adds the provided argument to the ContractCallInput
func AddArgument(input *vmcommon.ContractCallInput, argument []byte) {
	if input.Arguments == nil {
		input.Arguments = make([][]byte, 0)
	}
	input.Arguments = append(input.Arguments, argument)
}

// CopyTxHashes copies the tx hashes from a source ContractCallInput into another
func CopyTxHashes(input *vmcommon.ContractCallInput, sourceInput *vmcommon.ContractCallInput) {
	input.CurrentTxHash = sourceInput.CurrentTxHash
	input.PrevTxHash = sourceInput.PrevTxHash
	input.OriginalTxHash = sourceInput.OriginalTxHash
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

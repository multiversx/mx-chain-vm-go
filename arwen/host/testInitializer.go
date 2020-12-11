package host

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
)

var defaultVMType = []byte{0xF, 0xF}
var errAccountNotFound = errors.New("account not found")

var userAddress = []byte("userAccount.....................")
var parentAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fparentSC..............")
var childAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fchildSC...............")

var customGasSchedule = config.GasScheduleMap(nil)

// GetSCCode retrieves the bytecode of a WASM module from a file
func GetSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("GetSCCode(): %s", fileName))
	}

	return code
}

// GetTestSCCode retrieves the bytecode of a WASM testing contract
func GetTestSCCode(scName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/output/" + scName + ".wasm"
	return GetSCCode(pathToSC)
}

// GetTestSCCodeModule retrieves the bytecode of a WASM testing contract, given
// a specific name of the WASM module
func GetTestSCCodeModule(scName string, moduleName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/output/" + moduleName + ".wasm"
	return GetSCCode(pathToSC)
}

// defaultTestArwenForDeployment creates an Arwen vmHost configured for testing deployments
func defaultTestArwenForDeployment(t *testing.T, _ uint64, newAddress []byte) *vmHost {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host, _ := defaultTestArwen(t, stubBlockchainHook)
	return host
}

func defaultTestArwenForCall(tb testing.TB, code []byte, balance *big.Int) (*vmHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &contextmock.StubAccount{
				Code:    code,
				Balance: balance,
			}, nil
		}
		return nil, errAccountNotFound
	}

	host, _ := defaultTestArwen(tb, stubBlockchainHook)
	return host, stubBlockchainHook
}

// defaultTestArwenForTwoSCs creates an Arwen vmHost configured for testing calls between 2 SmartContracts
func defaultTestArwenForTwoSCs(
	t *testing.T,
	parentCode []byte,
	childCode []byte,
	parentSCBalance *big.Int,
	childSCBalance *big.Int,
) (*vmHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	if parentSCBalance == nil {
		parentSCBalance = big.NewInt(1000)
	}

	if childSCBalance == nil {
		childSCBalance = big.NewInt(1000)
	}

	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &contextmock.StubAccount{
				Code:    parentCode,
				Balance: parentSCBalance,
			}, nil
		}
		if bytes.Equal(scAddress, childAddress) {
			return &contextmock.StubAccount{
				Code:    childCode,
				Balance: childSCBalance,
			}, nil
		}

		return nil, errAccountNotFound
	}

	host, _ := defaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwen(tb testing.TB, blockchain vmcommon.BlockchainHook) (*vmHost, error) {
	gasSchedule := customGasSchedule
	if gasSchedule == nil {
		gasSchedule = config.MakeGasMapForTests()
	}

	host, err := NewArwenVM(blockchain, &arwen.VMHostParameters{
		VMType:                   defaultVMType,
		BlockGasLimit:            uint64(1000),
		GasSchedule:              gasSchedule,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
		UseWarmInstance:          false,
		DynGasLockEnableEpoch:    0,
	})
	require.Nil(tb, err)
	require.NotNil(tb, host)
	return host, err
}

// DefaultTestContractCreateInput creates a vmcommon.ContractCreateInput struct with default values
func DefaultTestContractCreateInput() *vmcommon.ContractCreateInput {
	return &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments: [][]byte{
				[]byte("argument 1"),
				[]byte("argument 2"),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		ContractCode: []byte("contract"),
	}
}

// DefaultTestContractCallInput creates a vmcommon.ContractCallInput struct with default values
func DefaultTestContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  userAddress,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: parentAddress,
		Function:      "function",
	}
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

// AddFinishData appends the provided []byte to the ReturnData of the given vmOutput
func AddFinishData(vmOutput *vmcommon.VMOutput, data []byte) {
	vmOutput.ReturnData = append(vmOutput.ReturnData, data)
}

// AddNewOutputAccount creates a new vmcommon.OutputAccount from the provided arguments and adds it to OutputAccounts of the provided vmOutput
func AddNewOutputAccount(vmOutput *vmcommon.VMOutput, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
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
				Data:  data,
				Value: big.NewInt(balanceDelta),
			},
		}
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

// SetStorageUpdate sets a storage update to the provided vmcommon.OutputAccount
func SetStorageUpdate(account *vmcommon.OutputAccount, key []byte, data []byte) {
	keyString := string(key)
	update, exists := account.StorageUpdates[keyString]
	if !exists {
		update = &vmcommon.StorageUpdate{}
		account.StorageUpdates[keyString] = update
	}
	update.Offset = key
	update.Data = data
}

// SetStorageUpdateStrings sets a storage update to the provided vmcommon.OutputAccount, from string arguments
func SetStorageUpdateStrings(account *vmcommon.OutputAccount, key string, data string) {
	SetStorageUpdate(account, []byte(key), []byte(data))
}

// OpenFile method opens the file from given path - does not close the file
func OpenFile(relativePath string) (*os.File, error) {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Printf("cannot create absolute path for the provided file: %s", err.Error())
		return nil, err
	}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return f, nil
}

// LoadTomlFileToMap opens and decodes a toml file as a map[string]interface{}
func LoadTomlFileToMap(relativePath string) (map[string]interface{}, error) {
	f, err := OpenFile(relativePath)
	if err != nil {
		return nil, err
	}

	fileinfo, err := f.Stat()
	if err != nil {
		fmt.Printf("cannot stat file: %s", err.Error())
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = f.Read(buffer)
	if err != nil {
		fmt.Printf("cannot read from file: %s", err.Error())
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Printf("cannot close file: %s", err.Error())
		}
	}()

	loadedTree, err := toml.Load(string(buffer))
	if err != nil {
		fmt.Printf("cannot interpret file contents as toml: %s", err.Error())
		return nil, err
	}

	loadedMap := loadedTree.ToMap()

	return loadedMap, nil
}

func loadGasScheduleConfig(filepath string) (config.GasScheduleMap, error) {
	gasScheduleConfig, err := LoadTomlFileToMap(filepath)
	if err != nil {
		return nil, err
	}

	flattenedGasSchedule := make(config.GasScheduleMap)
	for libType, costs := range gasScheduleConfig {
		flattenedGasSchedule[libType] = make(map[string]uint64)
		costsMap := costs.(map[string]interface{})
		for operationName, cost := range costsMap {
			flattenedGasSchedule[libType][operationName] = uint64(cost.(int64))
		}
	}

	return flattenedGasSchedule, nil
}

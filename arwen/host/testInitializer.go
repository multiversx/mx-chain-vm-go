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

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
)

var defaultVmType = []byte{0xF, 0xF}
var ErrCodeNotFound = errors.New("code not found")
var firstAddress = []byte("firstSC.........................")
var secondAddress = []byte("secondSC........................")

func GetSCCode(fileName string) []byte {
	code, _ := ioutil.ReadFile(filepath.Clean(fileName))

	return code
}

func GetTestSCCode(scName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/" + scName + ".wasm"
	return GetSCCode(pathToSC)
}

func DefaultTestArwenForDeployment(t *testing.T, ownerNonce uint64, newAddress []byte) *vmHost {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetNonceCalled = func(address []byte) (uint64, error) {
		return 24, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	return host
}

func DefaultTestArwenForCall(tb testing.TB, code []byte) (*vmHost, *mock.BlockchainHookStub) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetCodeCalled = func(address []byte) ([]byte, error) {
		return code, nil
	}
	host, _ := DefaultTestArwen(tb, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

func DefaultTestArwenForTwoSCs(t *testing.T, firstCode []byte, secondCode []byte) (*vmHost, *mock.BlockchainHookStub) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetCodeCalled = func(scAddress []byte) ([]byte, error) {
		if bytes.Equal(scAddress, firstAddress) {
			return firstCode, nil
		}
		if bytes.Equal(scAddress, secondAddress) {
			return secondCode, nil
		}
		return nil, ErrCodeNotFound
	}
	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

func DefaultTestArwen(tb testing.TB, blockchain vmcommon.BlockchainHook, crypto vmcommon.CryptoHook) (*vmHost, error) {
	host, err := NewArwenVM(blockchain, crypto, defaultVmType, uint64(1000), config.MakeGasMap(1))
	require.Nil(tb, err)
	require.NotNil(tb, host)
	return host, err
}

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

func DefaultTestContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  []byte("caller"),
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: []byte("smartcontract"),
		Function:      "function",
	}
}

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

func AddFinishData(vmOutput *vmcommon.VMOutput, data []byte) {
	vmOutput.ReturnData = append(vmOutput.ReturnData, data)
}

func AddNewOutputAccount(vmOutput *vmcommon.VMOutput, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	account := &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(balanceDelta),
		Balance:        nil,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		Code:           nil,
		Data:           data,
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

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

// LoadTomlFile method to open and decode toml file
func LoadTomlFile(dest interface{}, relativePath string) error {
	f, err := OpenFile(relativePath)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Printf("cannot close file: %s", err.Error())
		}
	}()

	return toml.NewDecoder(f).Decode(dest)
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

func LoadGasScheduleConfig(filepath string) (map[string]map[string]uint64, error) {
	gasScheduleConfig, err := LoadTomlFileToMap(filepath)
	if err != nil {
		return nil, err
	}

	flattenedGasSchedule := make(map[string]map[string]uint64)
	for libType, costs := range gasScheduleConfig {
		flattenedGasSchedule[libType] = make(map[string]uint64)
		costsMap := costs.(map[string]interface{})
		for operationName, cost := range costsMap {
			flattenedGasSchedule[libType][operationName] = uint64(cost.(int64))
		}
	}

	return flattenedGasSchedule, nil
}

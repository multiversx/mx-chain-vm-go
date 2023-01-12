package testcommon

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"unicode"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/stretchr/testify/require"
)

// VMOutputVerifier holds the output to be verified
type VMOutputVerifier struct {
	VmOutput  *vmcommon.VMOutput
	AllErrors arwen.WrappableError
	T         testing.TB
}

// NewVMOutputVerifier builds a new verifier
func NewVMOutputVerifier(t testing.TB, vmOutput *vmcommon.VMOutput, err error) *VMOutputVerifier {
	return NewVMOutputVerifierWithAllErrors(t, vmOutput, err, nil)
}

// NewVMOutputVerifierWithAllErrors builds a new verifier with all errors included
func NewVMOutputVerifierWithAllErrors(t testing.TB, vmOutput *vmcommon.VMOutput, err error, allErrors error) *VMOutputVerifier {
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	var allErrorsAsWrappable arwen.WrappableError
	if allErrors != nil {
		allErrorsAsWrappable = allErrors.(arwen.WrappableError)
	}

	return &VMOutputVerifier{
		VmOutput:  vmOutput,
		AllErrors: allErrorsAsWrappable,
		T:         t,
	}
}

// Ok verifies if return code is vmcommon.Ok
func (v *VMOutputVerifier) Ok() *VMOutputVerifier {
	return v.ReturnCode(vmcommon.Ok)
}

// ExecutionFailed verifies if return code is vmcommon.ExecutionFailed
func (v *VMOutputVerifier) ExecutionFailed() *VMOutputVerifier {
	return v.ReturnCode(vmcommon.ExecutionFailed)
}

// OutOfGas verifies if return code is vmcommon.OutOfGas
func (v *VMOutputVerifier) OutOfGas() *VMOutputVerifier {
	return v.ReturnCode(vmcommon.OutOfGas)
}

// ContractInvalid verifies if return code is vmcommon.ContractInvalid
func (v *VMOutputVerifier) ContractInvalid() *VMOutputVerifier {
	return v.ReturnCode(vmcommon.ContractInvalid)
}

// ContractNotFound verifies if return code is vmcommon.ContractNotFound
func (v *VMOutputVerifier) ContractNotFound() *VMOutputVerifier {
	return v.ReturnCode(vmcommon.ContractNotFound)
}

// UserError verifies if return code is vmcommon.UserError
func (v *VMOutputVerifier) UserError() *VMOutputVerifier {
	return v.ReturnCode(vmcommon.UserError)
}

// FunctionNotFound verifies if return code is vmcommon.FunctionNotFound
func (v *VMOutputVerifier) FunctionNotFound() *VMOutputVerifier {
	return v.ReturnCode(vmcommon.FunctionNotFound)
}

// ReturnCode verifies if ReturnCode of output is the same as the provided one
func (v *VMOutputVerifier) ReturnCode(code vmcommon.ReturnCode) *VMOutputVerifier {
	require.Equal(v.T, code, v.VmOutput.ReturnCode, "ReturnCode")
	return v
}

// ReturnMessage verifies if ReturnMessage of output is the same as the provided one
func (v *VMOutputVerifier) ReturnMessage(message string) *VMOutputVerifier {
	require.Equal(v.T, message, v.VmOutput.ReturnMessage, "ReturnMessage")
	return v
}

// ReturnMessageContains verifies if ReturnMessage of output contains the provided one
func (v *VMOutputVerifier) ReturnMessageContains(message string) *VMOutputVerifier {
	require.Contains(v.T, v.VmOutput.ReturnMessage, message, "ReturnMessage")
	return v
}

// HasRuntimeErrors verifies if the provided errors are present in the runtime context
func (v *VMOutputVerifier) HasRuntimeErrors(messages ...string) *VMOutputVerifier {
	for _, message := range messages {
		errorFound := false
		require.NotNil(v.T, v.AllErrors)
		for _, err := range v.AllErrors.GetAllErrors() {
			if strings.HasPrefix(err.Error(), message) {
				errorFound = true
			}
		}
		require.True(v.T, errorFound, fmt.Sprintf("No error with message '%s' found", message))
	}
	return v
}

// HasRuntimeErrorAndInfo verifies if the provided errors are present in the runtime context
func (v *VMOutputVerifier) HasRuntimeErrorAndInfo(message string, otherInfo string) *VMOutputVerifier {
	errorFound := false
	require.NotNil(v.T, v.AllErrors)
	errors, otherInfos := v.AllErrors.GetAllErrorsAndOtherInfo()
	for index, err := range errors {
		if strings.HasPrefix(err.Error(), message) && strings.HasPrefix(otherInfos[index], otherInfo) {
			errorFound = true
			break
		}
	}
	require.True(v.T, errorFound, fmt.Sprintf("No error with message '%s' found", message))
	return v
}

// GasUsed verifies if GasUsed of the specified account is the same as the provided one
func (v *VMOutputVerifier) GasUsed(address []byte, gas uint64) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("GasUsed", address)
	require.NotNil(v.T, account, errMsg)
	require.Equal(v.T, int(gas), int(account.GasUsed), errMsg)
	return v
}

// GasRemaining verifies if GasRemaining of the specified account is the same as the provided one
func (v *VMOutputVerifier) GasRemaining(gas uint64) *VMOutputVerifier {
	require.Equal(v.T, int(gas), int(v.VmOutput.GasRemaining), "GasRemaining")
	return v
}

// Balance verifies if Balance of the specified account is the same as the provided one
func (v *VMOutputVerifier) Balance(address []byte, balance int64) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("Balance", address)
	require.NotNil(v.T, account, errMsg)
	require.NotNil(v.T, account.Balance, errMsg)
	require.Equal(v.T, balance, account.Balance.Int64(), errMsg)
	return v
}

// BalanceDelta verifies if BalanceDelta of the specified account is the same as the provided one
func (v *VMOutputVerifier) BalanceDelta(address []byte, balanceDelta int64) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("BalanceDelta", address)
	require.NotNil(v.T, account, errMsg)
	require.NotNil(v.T, account.BalanceDelta, errMsg)
	require.Equal(v.T, balanceDelta, account.BalanceDelta.Int64(), errMsg)
	return v
}

// Nonce verifies if Nonce of the specified account is the same as the provided one
func (v *VMOutputVerifier) Nonce(address []byte, nonce uint64) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("Nonce", address)
	require.NotNil(v.T, account, errMsg)
	require.Equal(v.T, nonce, account.Nonce, errMsg)
	return v
}

// Code verifies if Code of the specified account is the same as the provided one
func (v *VMOutputVerifier) Code(address []byte, code []byte) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("Code", address)
	require.NotNil(v.T, account, errMsg)
	require.Equal(v.T, code, account.Code, errMsg)
	return v
}

// CodeMetadata if CodeMetadata of the specified account is the same as the provided one
func (v *VMOutputVerifier) CodeMetadata(address []byte, codeMetadata []byte) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("CodeMetadata", address)
	require.NotNil(v.T, account, errMsg)
	require.Equal(v.T, codeMetadata, account.CodeMetadata, errMsg)
	return v
}

// CodeDeployerAddress if CodeDeployerAddress of the specified account is the same as the provided one
func (v *VMOutputVerifier) CodeDeployerAddress(address []byte, codeDeployerAddress []byte) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("CodeDeployerAddress", address)
	require.NotNil(v.T, account, errMsg)
	require.Equal(v.T, codeDeployerAddress, account.CodeDeployerAddress, errMsg)
	return v
}

// ReturnData verifies if ReturnData is the same as the provided one
func (v *VMOutputVerifier) ReturnData(returnData ...[]byte) *VMOutputVerifier {
	require.Equal(v.T, len(returnData), len(v.VmOutput.ReturnData), "ReturnData")
	for idx := range v.VmOutput.ReturnData {
		require.Equal(v.T, returnData[idx], v.VmOutput.ReturnData[idx], "ReturnData")
	}
	return v
}

// StoreEntry holds the data for a storage assertion
type StoreEntry struct {
	address []byte
	key     []byte
	value   []byte
	written bool
}

// CreateStoreEntry creates the data for a storage assertion
func CreateStoreEntry(address []byte) *StoreEntry {
	return &StoreEntry{address: address}
}

// WithKey sets the key for a storage assertion
func (storeEntry *StoreEntry) WithKey(key []byte) *StoreEntry {
	storeEntry.key = key
	return storeEntry
}

// WithValue sets the value for a storage assertion
func (storeEntry *StoreEntry) WithValue(value []byte) StoreEntry {
	storeEntry.value = value
	storeEntry.written = true
	return *storeEntry
}

// Storage verifies if StorageUpdate(s) for the speficied accounts are the same as the provided ones
func (v *VMOutputVerifier) Storage(expectedEntries ...StoreEntry) *VMOutputVerifier {

	storage := make(map[string]map[string]vmcommon.StorageUpdate)

	for _, storeEntry := range expectedEntries {
		account := string(storeEntry.address)
		accountStorageMap, exists := storage[account]
		if !exists {
			accountStorageMap = make(map[string]vmcommon.StorageUpdate)
			storage[account] = accountStorageMap
		}
		accountStorageMap[string(storeEntry.key)] = vmcommon.StorageUpdate{Offset: storeEntry.key, Data: storeEntry.value, Written: storeEntry.written}
	}

	for _, outputAccount := range v.VmOutput.OutputAccounts {
		accountStorageMap := storage[string(outputAccount.Address)]
		require.Equal(v.T, len(accountStorageMap), len(outputAccount.StorageUpdates), "Storage")
		for key, value := range accountStorageMap {
			require.Equal(v.T, value, *outputAccount.StorageUpdates[key], "Storage")
		}
		delete(storage, string(outputAccount.Address))
	}
	require.Equal(v.T, 0, len(storage), "Storage")

	return v
}

// TransferEntry holds the data for an output transfer assertion
type TransferEntry struct {
	vmcommon.OutputTransfer
	address []byte
}

// CreateTransferEntry creates the data for an output transfer assertion
func CreateTransferEntry(senderAddress []byte, receiverAddress []byte) *TransferEntry {
	return &TransferEntry{
		OutputTransfer: vmcommon.OutputTransfer{SenderAddress: senderAddress},
		address:        receiverAddress,
	}
}

// WithData create sets the data for an output transfer assertion
func (transferEntry *TransferEntry) WithData(data []byte) *TransferEntry {
	transferEntry.Data = data
	return transferEntry
}

// WithGasLimit create sets the data for an output transfer assertion
func (transferEntry *TransferEntry) WithGasLimit(gas uint64) *TransferEntry {
	transferEntry.GasLimit = gas
	return transferEntry
}

// WithGasLocked create sets the data for an output transfer assertion
func (transferEntry *TransferEntry) WithGasLocked(gas uint64) *TransferEntry {
	transferEntry.GasLocked = gas
	return transferEntry
}

// WithCallType create sets the data for an output transfer assertion
func (transferEntry *TransferEntry) WithCallType(callType vm.CallType) *TransferEntry {
	transferEntry.CallType = callType
	return transferEntry
}

// WithValue create sets the value for an output transfer assertion
func (transferEntry *TransferEntry) WithValue(value *big.Int) TransferEntry {
	transferEntry.Value = value
	return *transferEntry
}

// Transfers verifies if OutputTransfer(s) for the speficied accounts are the same as the provided ones
func (v *VMOutputVerifier) Transfers(transfers ...TransferEntry) *VMOutputVerifier {
	transfersMap := make(map[string][]vmcommon.OutputTransfer)
	for _, transferEntry := range transfers {
		account := string(transferEntry.address)
		accountTransfers, exists := transfersMap[account]
		if !exists {
			accountTransfers = make([]vmcommon.OutputTransfer, 0)
		}
		transfersMap[account] = append(accountTransfers, transferEntry.OutputTransfer)
	}

	for _, outputAccount := range v.VmOutput.OutputAccounts {
		transfersForAccount := transfersMap[string(outputAccount.Address)]
		errMsg := formatErrorForAccount("Transfers to ", outputAccount.Address)
		require.Equal(v.T, len(transfersForAccount), len(outputAccount.OutputTransfers), errMsg)
		for idx := range transfersForAccount {
			errMsg = formatErrorForAccount("Transfers from / to ",
				outputAccount.OutputTransfers[idx].SenderAddress,
				outputAccount.Address)
			require.Equal(v.T, transfersForAccount[idx], outputAccount.OutputTransfers[idx], errMsg)
		}
		delete(transfersMap, string(outputAccount.Address))
	}
	require.Equal(v.T, 0, len(transfersMap), "Transfers asserted, but not present in VMOutput")

	return v
}

// Logs verifies if Logs is the same as the provided one
func (v *VMOutputVerifier) Logs(logs ...vmcommon.LogEntry) *VMOutputVerifier {
	require.Equal(v.T, len(logs), len(v.VmOutput.Logs), "Logs")
	for idx := range v.VmOutput.Logs {
		require.Equal(v.T, logs[idx].Address, v.VmOutput.Logs[idx].Address, "Logs.Address")
		require.Equal(v.T, logs[idx].Topics, v.VmOutput.Logs[idx].Topics, "Logs.Topics")
		require.Equal(v.T, logs[idx].Data, v.VmOutput.Logs[idx].Data, "Logs.Data")
		require.Equal(v.T, logs[idx].Identifier, v.VmOutput.Logs[idx].Identifier, "Logs.Identifier")
	}
	return v
}

// BytesAddedToStorage verifies the number of bytes added to storage
func (v *VMOutputVerifier) BytesAddedToStorage(address []byte, bytesAdded int) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("BytesAddedToStorage", address)
	require.NotNil(v.T, account, errMsg)
	require.Equal(v.T, bytesAdded, int(account.BytesAddedToStorage), errMsg)
	return v
}

// BytesDeletedFromStorage verifies the number of bytes deleted from storage
func (v *VMOutputVerifier) BytesDeletedFromStorage(address []byte, bytesDelted int) *VMOutputVerifier {
	account := v.VmOutput.OutputAccounts[string(address)]
	errMsg := formatErrorForAccount("BytesAddedToStorage", address)
	require.NotNil(v.T, account, errMsg)
	require.Equal(v.T, bytesDelted, int(account.BytesDeletedFromStorage), errMsg)
	return v
}

func formatErrorForAccount(field string, address ...[]byte) string {
	return fmt.Sprintf("%s %s", field, humanReadable(address...))
}

func humanReadable(addresses ...[]byte) string {
	var result []byte
	for _, address := range addresses {
		for _, c := range address {
			if unicode.IsPrint(rune(c)) {
				result = append(result, c)
			}
		}
		result = append(result, '|')
	}
	return string(result)
}

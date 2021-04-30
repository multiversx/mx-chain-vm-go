package host

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

// VMOutputVerifier holds the output to be verified
type VMOutputVerifier struct {
	vmOutput *vmcommon.VMOutput
	T        testing.TB
}

// NewVMOutputVerifier builds a new verifier
func NewVMOutputVerifier(t testing.TB, vmOutput *vmcommon.VMOutput, err error) *VMOutputVerifier {
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	return &VMOutputVerifier{
		vmOutput: vmOutput,
		T:        t,
	}
}

// Ok verifies if return code is vmcommon.Ok
func (v *VMOutputVerifier) Ok() *VMOutputVerifier {
	return v.RetCode(vmcommon.Ok)
}

// RetCode verifies if ReturnCode of output is the same as the provided one
func (v *VMOutputVerifier) RetCode(code vmcommon.ReturnCode) *VMOutputVerifier {
	require.Equal(v.T, code, v.vmOutput.ReturnCode, "ReturnCode")
	return v
}

// ReturnMessage verifies if ReturnMessage of output is the same as the provided one
func (v *VMOutputVerifier) ReturnMessage(message string) *VMOutputVerifier {
	require.Equal(v.T, message, v.vmOutput.ReturnMessage, "ReturnMessage")
	return v
}

// NoMsg verifies that ReturnMessage is empty
func (v *VMOutputVerifier) NoMsg() *VMOutputVerifier {
	require.Equal(v.T, "", v.vmOutput.ReturnMessage, "ReturnMessage")
	return v
}

// Msg verifies if ReturnMessage of output is the same as the provided one
func (v *VMOutputVerifier) Msg(message string) *VMOutputVerifier {
	require.Equal(v.T, message, v.vmOutput.ReturnMessage, "ReturnMessage")
	return v
}

// GasUsed verifies if GasUsed of the specified account is the same as the provided one
func (v *VMOutputVerifier) GasUsed(address []byte, gas uint64) *VMOutputVerifier {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account, "GasUsed")
	require.Equal(v.T, int(gas), int(account.GasUsed), "GasUsed")
	return v
}

// GasRemaining verifies if GasRemaining of the specified account is the same as the provided one
func (v *VMOutputVerifier) GasRemaining(gas uint64) *VMOutputVerifier {
	require.Equal(v.T, int(gas), int(v.vmOutput.GasRemaining), "GasRemaining")
	return v
}

// Balance verifies if Balance of the specified account is the same as the provided one
func (v *VMOutputVerifier) Balance(address []byte, balance int64) *VMOutputVerifier {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account, "Balance")
	require.Equal(v.T, balance, account.Balance.Int64(), "Balance")
	return v
}

// BalanceDelta verifies if BalanceDelta of the specified account is the same as the provided one
func (v *VMOutputVerifier) BalanceDelta(address []byte, balanceDelta int64) *VMOutputVerifier {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account, "BalanceDelta")
	require.Equal(v.T, balanceDelta, account.BalanceDelta.Int64(), "BalanceDelta")
	return v
}

// ReturnData verifies if ReturnData is the same as the provided one
func (v *VMOutputVerifier) ReturnData(returnData ...[]byte) *VMOutputVerifier {
	require.Equal(v.T, len(returnData), len(v.vmOutput.ReturnData))
	for idx := range v.vmOutput.ReturnData {
		require.Equal(v.T, returnData[idx], v.vmOutput.ReturnData[idx], "ReturnData")
	}
	return v
}

type storeEntry struct {
	address []byte
	key     []byte
	value   []byte
}

// Storage verifies if StorageUpdate(s) for the speficied accounts are the same as the provided ones
func (v *VMOutputVerifier) Storage(returnData ...storeEntry) *VMOutputVerifier {

	storage := make(map[string]map[string]vmcommon.StorageUpdate)

	for _, storeEntry := range returnData {
		account := string(storeEntry.address)
		accountStorageMap, exists := storage[account]
		if !exists {
			accountStorageMap = make(map[string]vmcommon.StorageUpdate)
			storage[account] = accountStorageMap
		}
		accountStorageMap[string(storeEntry.key)] = vmcommon.StorageUpdate{Offset: storeEntry.key, Data: storeEntry.value}
	}

	for _, outputAccount := range v.vmOutput.OutputAccounts {
		accountStorageMap := storage[string(outputAccount.Address)]
		require.Equal(v.T, len(accountStorageMap), len(outputAccount.StorageUpdates), "Storage")
		for key, value := range accountStorageMap {
			require.Equal(v.T, value, *outputAccount.StorageUpdates[key], "Storage")
		}
	}

	return v
}

type transferEntry struct {
	address  []byte
	transfer vmcommon.OutputTransfer
}

// Transfers verifies if OutputTransfer(s) for the speficied accounts are the same as the provided ones
func (v *VMOutputVerifier) Transfers(transfers ...transferEntry) *VMOutputVerifier {

	transfersMap := make(map[string][]vmcommon.OutputTransfer)

	for _, transferEntry := range transfers {
		account := string(transferEntry.address)
		accountTransfers, exists := transfersMap[account]
		if !exists {
			accountTransfers = make([]vmcommon.OutputTransfer, 0)
		}
		transfersMap[account] = append(accountTransfers, transferEntry.transfer)
	}

	for _, outputAccount := range v.vmOutput.OutputAccounts {
		transfersForAccount := transfersMap[string(outputAccount.Address)]
		require.Equal(v.T, len(transfersForAccount), len(outputAccount.OutputTransfers), "Transfers")
		for idx := range transfersForAccount {
			require.Equal(v.T, transfersForAccount[idx], outputAccount.OutputTransfers[idx], "Transfers")
		}
	}

	return v
}

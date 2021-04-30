package host

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

type VMOutputVerifier struct {
	vmOutput *vmcommon.VMOutput

	T testing.TB
}

func NewVMOutputVerifier(t testing.TB, vmOutput *vmcommon.VMOutput, err error) *VMOutputVerifier {
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	return &VMOutputVerifier{
		vmOutput: vmOutput,
		T:        t,
	}
}

func (v *VMOutputVerifier) Ok() *VMOutputVerifier {
	return v.RetCode(vmcommon.Ok)
}

func (v *VMOutputVerifier) RetCode(code vmcommon.ReturnCode) *VMOutputVerifier {
	require.Equal(v.T, code, v.vmOutput.ReturnCode)
	return v
}

func (v *VMOutputVerifier) NoMsg() *VMOutputVerifier {
	require.Equal(v.T, "", v.vmOutput.ReturnMessage)
	return v
}
func (v *VMOutputVerifier) Msg(message string) *VMOutputVerifier {
	require.Equal(v.T, message, v.vmOutput.ReturnMessage)
	return v
}

func (v *VMOutputVerifier) GasUsed(address []byte, gas uint64) *VMOutputVerifier {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account)
	require.Equal(v.T, int(gas), int(account.GasUsed))
	return v
}

func (v *VMOutputVerifier) GasRemaining(gas uint64) *VMOutputVerifier {
	require.Equal(v.T, int(gas), int(v.vmOutput.GasRemaining))
	return v
}

func (v *VMOutputVerifier) Balance(address []byte, balance int64) *VMOutputVerifier {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account)
	require.Equal(v.T, balance, account.Balance.Int64())
	return v
}

func (v *VMOutputVerifier) BalanceDelta(address []byte, balanceDelta int64) *VMOutputVerifier {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account)
	require.Equal(v.T, balanceDelta, account.BalanceDelta.Int64())
	return v
}

func (v *VMOutputVerifier) ReturnData(returnData ...[]byte) *VMOutputVerifier {
	require.Equal(v.T, len(returnData), len(v.vmOutput.ReturnData))
	for idx := range v.vmOutput.ReturnData {
		require.Equal(v.T, returnData[idx], v.vmOutput.ReturnData[idx])
	}
	return v
}

type storeEntry struct {
	address []byte
	key     []byte
	value   []byte
}

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
		require.Equal(v.T, len(accountStorageMap), len(outputAccount.StorageUpdates))
		for key, value := range accountStorageMap {
			require.Equal(v.T, value, *outputAccount.StorageUpdates[key])
		}
	}

	return v
}

type transferEntry struct {
	address  []byte
	transfer vmcommon.OutputTransfer
}

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
		require.Equal(v.T, len(transfersForAccount), len(outputAccount.OutputTransfers))
		for idx := range transfersForAccount {
			require.Equal(v.T, transfersForAccount[idx], outputAccount.OutputTransfers[idx])
		}
	}

	return v
}

package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type storageContext struct {
	host           arwen.VMHost
	blockChainHook vmcommon.BlockchainHook
}

func NewStorageContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
) (*storageContext, error) {
	context := &storageContext{
		host:           host,
		blockChainHook: blockChainHook,
	}

	return context, nil
}

func (context *storageContext) GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate {
	account, _ := context.host.Output().GetOutputAccount(address)
	return account.StorageUpdates
}

func (context *storageContext) GetStorage(address []byte, key []byte) []byte {
	storageUpdates := context.GetStorageUpdates(address)
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		return storageUpdate.Data
	}

	value, _ := context.blockChainHook.GetStorageData(address, key)
	return value
}

func (context *storageContext) SetStorage(address []byte, key []byte, value []byte) int32 {
	if context.host.Runtime().ReadOnly() {
		return 0
	}

	metering := context.host.Metering()
	zero := []byte{}
	strKey := string(key)
	length := len(value)

	var oldValue []byte
	storageUpdates := context.GetStorageUpdates(address)
	if update, ok := storageUpdates[strKey]; !ok {
		oldValue = context.GetStorage(address, key)
		storageUpdates[strKey] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   oldValue,
		}
	} else {
		oldValue = update.Data
	}

	if bytes.Equal(oldValue, value) {
		useGas := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
		metering.UseGas(useGas)
		return int32(arwen.StorageUnchanged)
	}

	newUpdate := &vmcommon.StorageUpdate{
		Offset: key,
		Data:   make([]byte, length),
	}
	copy(newUpdate.Data[:length], value[:length])
	storageUpdates[strKey] = newUpdate

	if bytes.Equal(oldValue, zero) {
		useGas := metering.GasSchedule().BaseOperationCost.StorePerByte * uint64(length)
		metering.UseGas(useGas)
		return int32(arwen.StorageAdded)
	}
	if bytes.Equal(value, zero) {
		freeGas := metering.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(oldValue))
		metering.FreeGas(freeGas)
		return int32(arwen.StorageDeleted)
	}

	useGas := metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(useGas)

	return int32(arwen.StorageModified)
}

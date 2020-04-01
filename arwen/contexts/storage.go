package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type storageContext struct {
	host           arwen.VMHost
	blockChainHook vmcommon.BlockchainHook
	address        []byte
	stateStack     [][]byte
}

// NewStorageContext creates a new storageContext
func NewStorageContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
) (*storageContext, error) {
	context := &storageContext{
		host:           host,
		blockChainHook: blockChainHook,
		stateStack:     make([][]byte, 0),
	}

	return context, nil
}

func (context *storageContext) InitState() {
}

func (context *storageContext) PushState() {
	context.stateStack = append(context.stateStack, context.address)
}

func (context *storageContext) PopState() {
	stateStackLen := len(context.stateStack)
	prevAddress := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.address = prevAddress
}

func (context *storageContext) ClearStateStack() {
	context.stateStack = make([][]byte, 0)
}

func (context *storageContext) SetAddress(address []byte) {
	context.address = address
}

func (context *storageContext) GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate {
	account, _ := context.host.Output().GetOutputAccount(address)
	return account.StorageUpdates
}

func (context *storageContext) GetStorage(key []byte) []byte {
	storageUpdates := context.GetStorageUpdates(context.address)
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		return storageUpdate.Data
	}

	value, _ := context.blockChainHook.GetStorageData(context.address, key)
	if value != nil {
		storageUpdates[string(key)] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   value,
		}
	}

	return value
}

func (context *storageContext) SetStorage(key []byte, value []byte) int32 {
	if context.host.Runtime().ReadOnly() {
		return int32(arwen.StorageUnchanged)
	}

	metering := context.host.Metering()
	var zero []byte
	strKey := string(key)
	length := len(value)

	var oldValue []byte
	storageUpdates := context.GetStorageUpdates(context.address)
	if update, ok := storageUpdates[strKey]; !ok {
		oldValue = context.GetStorage(key)
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

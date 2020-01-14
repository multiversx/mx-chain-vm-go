package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Storage struct {
	host           arwen.VMHost
	blockChainHook vmcommon.BlockchainHook
}

func NewStorageContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
) (*Storage, error) {
	storage := &Storage{
		host:           host,
		blockChainHook: blockChainHook,
	}

	return storage, nil
}

func (storage *Storage) InitState() {
}

func (storage *Storage) GetStorage(addr []byte, key []byte) []byte {
	storageUpdate := storage.host.Output().GetStorageUpdates()
	strAdr := string(addr)
	if _, ok := storageUpdate[strAdr]; ok {
		if value, ok := storageUpdate[strAdr][string(key)]; ok {
			return value
		}
	}

	hash, _ := storage.blockChainHook.GetStorageData(addr, key)
	return hash
}

func (storage *Storage) SetStorage(addr []byte, key []byte, value []byte) int32 {
	if storage.host.Runtime().ReadOnly() {
		return 0
	}

	strAdr := string(addr)

	storageUpdate := storage.host.Output().GetStorageUpdates()
	if _, ok := storageUpdate[strAdr]; !ok {
		storageUpdate[strAdr] = make(map[string][]byte, 0)
	}
	if _, ok := storageUpdate[strAdr][string(key)]; !ok {
		oldValue := storage.GetStorage(addr, key)
		storageUpdate[strAdr][string(key)] = oldValue
	}

	oldValue := storageUpdate[strAdr][string(key)]
	lengthOldValue := len(oldValue)
	length := len(value)
	storageUpdate[strAdr][string(key)] = make([]byte, length)
	copy(storageUpdate[strAdr][string(key)][:length], value[:length])

	metering := storage.host.Metering()
	if bytes.Equal(oldValue, value) {
		useGas := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
		metering.UseGas(useGas)
		return int32(arwen.StorageUnchanged)
	}

	zero := []byte{}
	if bytes.Equal(oldValue, zero) {
		useGas := metering.GasSchedule().BaseOperationCost.StorePerByte * uint64(length)
		metering.UseGas(useGas)
		return int32(arwen.StorageAdded)
	}
	if bytes.Equal(value, zero) {
		freeGas := metering.GasSchedule().BaseOperationCost.StorePerByte * uint64(lengthOldValue)
		metering.FreeGas(freeGas)
		return int32(arwen.StorageDeleted)
	}

	useGas := metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(useGas)

	return int32(arwen.StorageModified)
}

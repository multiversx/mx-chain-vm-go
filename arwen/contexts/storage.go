package contexts

import (
	"bytes"
	"errors"
	builtinMath "math"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type storageContext struct {
	host                     arwen.VMHost
	blockChainHook           vmcommon.BlockchainHook
	address                  []byte
	stateStack               [][]byte
	elrondProtectedKeyPrefix []byte
}

// NewStorageContext creates a new storageContext
func NewStorageContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
	elrondProtectedKeyPrefix []byte,
) (*storageContext, error) {
	if len(elrondProtectedKeyPrefix) == 0 {
		return nil, errors.New("elrondProtectedKeyPrefix cannot be empty")
	}
	context := &storageContext{
		host:                     host,
		blockChainHook:           blockChainHook,
		stateStack:               make([][]byte, 0),
		elrondProtectedKeyPrefix: elrondProtectedKeyPrefix,
	}

	return context, nil
}

// InitState does nothing
func (context *storageContext) InitState() {
}

// PushState appends the current address to the state stack.
func (context *storageContext) PushState() {
	context.stateStack = append(context.stateStack, context.address)
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current address
func (context *storageContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevAddress := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.address = prevAddress
}

// PopDiscard removes the latest entry from the state stack
func (context *storageContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack clears the state stack from the current context.
func (context *storageContext) ClearStateStack() {
	context.stateStack = make([][]byte, 0)
}

// SetAddress sets the given address as the address for the current context.
func (context *storageContext) SetAddress(address []byte) {
	context.address = address
}

// GetStorageUpdates returns the storage updates for the account mapped to the given address.
func (context *storageContext) GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate {
	account, _ := context.host.Output().GetOutputAccount(address)
	return account.StorageUpdates
}

// GetStorage returns the storage data mapped to the given key.
func (context *storageContext) GetStorage(key []byte) []byte {
	metering := context.host.Metering()

	extraBytes := len(key) - arwen.AddressLen
	if extraBytes > 0 {
		gasToUse, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		if err != nil {
			log.Error("GetStorage overflow",
				"DataCopyPerByte = ", metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
				"extraBytes = ", uint64(extraBytes),
			)
			gasToUse = builtinMath.MaxUint64
		}
		metering.UseGas(gasToUse)
	}

	value := context.GetStorageUnmetered(key)

	gasToUse, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(value)))
	if err != nil {
		log.Error("GetStorage overflow",
			"DataCopyPerByte = ", metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
			"value = ", uint64(len(value)),
		)
		gasToUse = builtinMath.MaxUint64
	}

	metering.UseGas(gasToUse)

	return value
}

// GetStorageFromAddress returns the data under the given key from the account mapped to the given address.
func (context *storageContext) GetStorageFromAddress(address []byte, key []byte) []byte {
	metering := context.host.Metering()

	extraBytes := len(key) - arwen.AddressLen
	if extraBytes > 0 {
		gasToUse, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		if err != nil {
			log.Error("GetStorageFromAddress overflow",
				"DataCopyPerByte = ", metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
				"extraBytes = ", uint64(extraBytes),
			)
			gasToUse = builtinMath.MaxUint64
		}
		metering.UseGas(gasToUse)
	}

	if !bytes.Equal(address, context.address) {
		userAcc, err := context.blockChainHook.GetUserAccount(address)
		if err != nil || check.IfNil(userAcc) {
			return nil
		}

		metadata := vmcommon.CodeMetadataFromBytes(userAcc.GetCodeMetadata())
		if !metadata.Readable {
			return nil
		}
	}

	value := context.getStorageFromAddressUnmetered(address, key)

	gasToUse, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(len(value)))
	if err != nil {
		log.Error("GetStorageFromAddress overflow",
			"DataCopyPerByte = ", metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
			"value = ", uint64(len(value)),
		)
		gasToUse = builtinMath.MaxUint64
	}
	metering.UseGas(gasToUse)

	return value
}

func (context *storageContext) getStorageFromAddressUnmetered(address []byte, key []byte) []byte {
	var value []byte

	storageUpdates := context.GetStorageUpdates(address)
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		value = storageUpdate.Data
	} else {
		value, _ = context.blockChainHook.GetStorageData(address, key)
		storageUpdates[string(key)] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   value,
		}
	}

	return value
}

// GetStorageUnmetered returns the data under the given key.
func (context *storageContext) GetStorageUnmetered(key []byte) []byte {
	return context.getStorageFromAddressUnmetered(context.address, key)
}

func (context *storageContext) isElrondReservedKey(key []byte) bool {
	return bytes.HasPrefix(key, context.elrondProtectedKeyPrefix)
}

// SetStorage sets the given value at the given key.
func (context *storageContext) SetStorage(key []byte, value []byte) (arwen.StorageStatus, error) {
	if context.isElrondReservedKey(key) {
		return arwen.StorageUnchanged, arwen.ErrStoreElrondReservedKey
	}

	if context.host.Runtime().ReadOnly() {
		return arwen.StorageUnchanged, nil
	}

	metering := context.host.Metering()

	extraBytes := len(key) - arwen.AddressLen
	if extraBytes > 0 {
		gasToUse, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		if err != nil {
			log.Error("SetStorage overflow",
				"DataCopyPerByte = ", metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
				"extraBytes = ", uint64(extraBytes),
			)
			gasToUse = builtinMath.MaxUint64
		}
		metering.UseGas(gasToUse)
	}

	var zero []byte
	strKey := string(key)
	length := len(value)

	var oldValue []byte
	storageUpdates := context.GetStorageUpdates(context.address)
	if update, ok := storageUpdates[strKey]; !ok {
		oldValue = context.GetStorageUnmetered(key)
		storageUpdates[strKey] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   oldValue,
		}
	} else {
		oldValue = update.Data
	}

	lengthOldValue := len(oldValue)
	if bytes.Equal(oldValue, value) {
		useGas, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
		if err != nil {
			log.Error("SetStorage overflow",
				"DataCopyPerByte = ", metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
				"length = ", uint64(length),
			)
			useGas = builtinMath.MaxUint64
		}
		metering.UseGas(useGas)
		return arwen.StorageUnchanged, nil
	}

	newUpdate := &vmcommon.StorageUpdate{
		Offset: key,
		Data:   make([]byte, length),
	}
	copy(newUpdate.Data[:length], value[:length])
	storageUpdates[strKey] = newUpdate

	if bytes.Equal(oldValue, zero) {
		useGas, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.StorePerByte, uint64(length))
		if err != nil {
			log.Error("SetStorage overflow",
				"StorePerByte = ", metering.GasSchedule().BaseOperationCost.StorePerByte,
				"length = ", uint64(length),
			)
			useGas = builtinMath.MaxUint64
		}
		metering.UseGas(useGas)
		return arwen.StorageAdded, nil
	}
	if bytes.Equal(value, zero) {
		freeGas, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.ReleasePerByte, uint64(lengthOldValue))
		if err != nil {
			log.Error("SetStorage overflow",
				"ReleasePerByte = ", metering.GasSchedule().BaseOperationCost.ReleasePerByte,
				"lengthOldValue = ", uint64(lengthOldValue),
			)
			freeGas = builtinMath.MaxUint64
		}
		metering.FreeGas(freeGas)
		return arwen.StorageDeleted, nil
	}

	newValueExtraLength, err := math.SubInt(length, lengthOldValue)
	if err != nil {
		log.Error("SetStorage underflow",
			"length = ", length,
			"lengthOldValue = ", lengthOldValue,
		)
		newValueExtraLength = builtinMath.MinInt32
	}

	if newValueExtraLength > 0 {
		useGas, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(lengthOldValue))
		if err != nil {
			log.Error("SetStorage overflow",
				"PersistPerByte = ", metering.GasSchedule().BaseOperationCost.PersistPerByte,
				"lengthOldValue = ", uint64(lengthOldValue),
			)
			useGas = builtinMath.MaxUint64
		}

		newValStoreUseGas, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.StorePerByte, uint64(newValueExtraLength))
		if err != nil {
			log.Error("SetStorage overflow",
				"StorePerByte = ", metering.GasSchedule().BaseOperationCost.StorePerByte,
				"newValueExtraLength = ", uint64(newValueExtraLength),
			)
			newValStoreUseGas = builtinMath.MaxUint64
		}

		gasUsed, err := math.AddUint64(useGas, newValStoreUseGas)
		if err != nil {
			log.Error("SetStorage overflow",
				"useGas = ", useGas,
				"newValStoreUseGas = ", newValStoreUseGas,
			)
			gasUsed = builtinMath.MaxUint64
		}

		metering.UseGas(gasUsed)
	}

	if newValueExtraLength < 0 {
		newValueExtraLength = -newValueExtraLength

		useGas, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
		if err != nil {
			log.Error("SetStorage overflow",
				"PersistPerByte = ", metering.GasSchedule().BaseOperationCost.PersistPerByte,
				"length = ", uint64(length),
			)
			useGas = builtinMath.MaxUint64
		}
		metering.UseGas(useGas)

		freeGas, err := math.MulUint64(metering.GasSchedule().BaseOperationCost.ReleasePerByte, uint64(newValueExtraLength))
		if err != nil {
			log.Error("SetStorage overflow",
				"ReleasePerByte = ", metering.GasSchedule().BaseOperationCost.ReleasePerByte,
				"newValueExtraLength = ", uint64(newValueExtraLength),
			)
			freeGas = builtinMath.MaxUint64
		}
		metering.FreeGas(freeGas)
	}

	return arwen.StorageModified, nil
}

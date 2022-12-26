package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/math"
)

var logStorage = logger.GetOrCreate("arwen/storage")

type storageContext struct {
	host                          arwen.VMHost
	blockChainHook                vmcommon.BlockchainHook
	address                       []byte
	stateStack                    [][]byte
	protectedKeyPrefix            []byte
	arwenStorageProtectionEnabled bool
}

// NewStorageContext creates a new storageContext
func NewStorageContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
	protectedKeyPrefix []byte,
) (*storageContext, error) {
	if len(protectedKeyPrefix) == 0 {
		return nil, arwen.ErrEmptyProtectedKeyPrefix
	}
	if check.IfNil(host) {
		return nil, arwen.ErrNilVMHost
	}
	if check.IfNil(blockChainHook) {
		return nil, arwen.ErrNilBlockChainHook
	}

	context := &storageContext{
		host:                          host,
		blockChainHook:                blockChainHook,
		stateStack:                    make([][]byte, 0),
		protectedKeyPrefix:            protectedKeyPrefix,
		arwenStorageProtectionEnabled: true,
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
	logStorage.Trace("storage under address set", "address", address)
}

// GetStorageUpdates returns the storage updates for the account mapped to the given address.
func (context *storageContext) GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate {
	account, _ := context.host.Output().GetOutputAccount(address)
	return account.StorageUpdates
}

// GetStorage returns the storage data mapped to the given key.
func (context *storageContext) GetStorage(key []byte) ([]byte, bool, error) {
	value, usedCache, err := context.GetStorageUnmetered(key)
	if err != nil {
		return nil, false, err
	}
	context.useExtraGasForKeyIfNeeded(key, usedCache)
	context.useGasForValueIfNeeded(value, usedCache)
	logStorage.Trace("get", "key", key, "value", value)

	return value, usedCache, nil
}

func (context *storageContext) useGasForValueIfNeeded(value []byte, usedCache bool) {
	metering := context.host.Metering()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	gasFlagSet := enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled()
	if !usedCache || !gasFlagSet {
		costPerByte := metering.GasSchedule().BaseOperationCost.DataCopyPerByte
		gasToUse := math.MulUint64(costPerByte, uint64(len(value)))
		metering.UseGas(gasToUse)
	}
}

func (context *storageContext) useExtraGasForKeyIfNeeded(key []byte, usedCache bool) {
	metering := context.host.Metering()
	extraBytes := len(key) - arwen.AddressLen
	if extraBytes <= 0 {
		return
	}
	enableEpochsHandler := context.host.EnableEpochsHandler()
	gasFlagSet := enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled()
	if !gasFlagSet || !usedCache {
		gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		metering.UseGas(gasToUse)
	}
}

// GetStorageFromAddress returns the data under the given key from the account mapped to the given address.
func (context *storageContext) GetStorageFromAddress(address []byte, key []byte) ([]byte, bool, error) {
	if !bytes.Equal(address, context.address) {
		userAcc, err := context.blockChainHook.GetUserAccount(address)
		if err != nil || check.IfNil(userAcc) {
			context.useExtraGasForKeyIfNeeded(key, false)
			return nil, false, nil
		}

		metadata := vmcommon.CodeMetadataFromBytes(userAcc.GetCodeMetadata())
		if !metadata.Readable {
			context.useExtraGasForKeyIfNeeded(key, false)
			return nil, false, nil
		}
	}

	// If the requested key is protected by the Elrond node, the stored value
	// could have been changed by a built-in function in the meantime, even if
	// contracts themselves cannot change protected values. Values stored under
	// protected keys must always be retrieved from the node, not from the cached
	// StorageUpdates.
	value, usedCache, err := context.getStorageFromAddressUnmetered(address, key)

	context.useExtraGasForKeyIfNeeded(key, usedCache)
	context.useGasForValueIfNeeded(value, usedCache)

	logStorage.Trace("get from address", "address", address, "key", key, "value", value)
	return value, usedCache, err
}

func (context *storageContext) getStorageFromAddressUnmetered(address []byte, key []byte) ([]byte, bool, error) {
	var value []byte
	var err error

	enableEpochsHandler := context.host.EnableEpochsHandler()
	if context.isElrondReservedKey(key) && enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled() {
		value, err = context.readFromBlockchain(address, key)
		return value, false, err
	}

	storageUpdates := context.GetStorageUpdates(address)
	usedCache := true
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		value = storageUpdate.Data
	} else {
		value, err = context.readFromBlockchain(address, key)
		if err != nil {
			return nil, false, err
		}
		storageUpdates[string(key)] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   value,
		}
		usedCache = false
	}

	return value, usedCache, nil
}

func (context *storageContext) readFromBlockchain(address []byte, key []byte) ([]byte, error) {
	value, _, err := context.blockChainHook.GetStorageData(address, key)

	return value, err
}

// GetStorageUnmetered returns the data under the given key.
func (context *storageContext) GetStorageUnmetered(key []byte) ([]byte, bool, error) {
	return context.getStorageFromAddressUnmetered(context.address, key)
}

// enableStorageProtection will prevent writing to protected keys
func (context *storageContext) enableStorageProtection() {
	context.arwenStorageProtectionEnabled = true
}

// disableStorageProtection will prevent writing to protected keys
func (context *storageContext) disableStorageProtection() {
	context.arwenStorageProtectionEnabled = false
}

func (context *storageContext) isArwenProtectedKey(key []byte) bool {
	return bytes.HasPrefix(key, []byte(arwen.ProtectedStoragePrefix))
}

func (context *storageContext) isElrondReservedKey(key []byte) bool {
	return bytes.HasPrefix(key, context.protectedKeyPrefix)
}

// SetProtectedStorage sets storage for timelocks and promises
func (context *storageContext) SetProtectedStorage(key []byte, value []byte) (arwen.StorageStatus, error) {
	context.disableStorageProtection()
	defer context.enableStorageProtection()

	return context.SetStorage(key, value)
}

// SetStorage sets the given value at the given key.
func (context *storageContext) SetStorage(key []byte, value []byte) (arwen.StorageStatus, error) {
	if context.host.Runtime().ReadOnly() {
		logStorage.Trace("storage set", "error", "cannot set storage in readonly mode")
		if context.host.CheckExecuteReadOnly() {
			return arwen.StorageUnchanged, arwen.ErrCannotWriteOnReadOnly
		}

		return arwen.StorageUnchanged, nil
	}
	if context.isElrondReservedKey(key) {
		logStorage.Trace("storage set", "error", arwen.ErrStoreElrondReservedKey, "key", key)
		return arwen.StorageUnchanged, arwen.ErrStoreElrondReservedKey
	}
	if context.isArwenProtectedKey(key) && context.arwenStorageProtectionEnabled {
		logStorage.Trace("storage set", "error", arwen.ErrCannotWriteProtectedKey, "key", key)
		return arwen.StorageUnchanged, arwen.ErrCannotWriteProtectedKey
	}

	metering := context.host.Metering()

	length := len(value)

	storageUpdates := context.GetStorageUpdates(context.address)
	oldValue, usedCache, err := context.getOldValue(storageUpdates, key)
	if err != nil {
		return arwen.StorageUnchanged, err
	}

	gasForKey := context.computeGasForKey(key, usedCache)
	metering.UseGas(gasForKey)

	if bytes.Equal(oldValue, value) {
		return context.storageUnchanged(length, usedCache)
	}

	deltaBytes := len(value) - len(oldValue)
	context.addDeltaBytes(deltaBytes)

	context.changeStorageUpdate(key, value, storageUpdates)

	if len(oldValue) == 0 {
		return context.storageAdded(length, key, value)
	}

	lengthOldValue := len(oldValue)
	if len(value) == 0 {
		return context.storageDeleted(lengthOldValue, key)
	}

	newValueExtraLength := math.SubInt(length, lengthOldValue)

	var gasToUseForValue, gasToFreeForValue uint64
	switch {
	case newValueExtraLength > 0:
		gasToUseForValue, gasToFreeForValue = context.computeGasForBiggerValues(lengthOldValue, newValueExtraLength)
	case newValueExtraLength < 0:
		gasToUseForValue, gasToFreeForValue = context.computeGasForSmallerValues(newValueExtraLength, length)
	case newValueExtraLength == 0:
		gasToUseForValue, gasToFreeForValue = 0, 0
	}

	metering.UseGas(gasToUseForValue)
	metering.FreeGas(gasToFreeForValue)

	logStorage.Trace("storage modified", "key", key, "value", value, "lengthDelta", newValueExtraLength)
	return arwen.StorageModified, nil
}

func (context *storageContext) addDeltaBytes(deltaBytes int) {
	account, _ := context.host.Output().GetOutputAccount(context.address)
	if deltaBytes > 0 {
		account.BytesAddedToStorage += uint64(deltaBytes)
	} else {
		account.BytesDeletedFromStorage += uint64(-deltaBytes)
	}
}

func (context *storageContext) changeStorageUpdate(key []byte, value []byte, storageUpdates map[string]*vmcommon.StorageUpdate) {
	length := len(value)
	newUpdate := &vmcommon.StorageUpdate{
		Offset:  key,
		Data:    make([]byte, length),
		Written: true,
	}
	copy(newUpdate.Data[:length], value[:length])
	storageUpdates[string(key)] = newUpdate
}

func (context *storageContext) computeGasForSmallerValues(newValueExtraLength int, length int) (uint64, uint64) {
	metering := context.host.Metering()
	newValueExtraLength = -newValueExtraLength
	useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	freeGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.ReleasePerByte, uint64(newValueExtraLength))
	return useGas, freeGas
}

func (context *storageContext) computeGasForBiggerValues(lengthOldValue int, newValueExtraLength int) (uint64, uint64) {
	metering := context.host.Metering()
	useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(lengthOldValue))
	newValStoreUseGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.StorePerByte, uint64(newValueExtraLength))
	useGas = math.AddUint64(useGas, newValStoreUseGas)
	return useGas, 0
}

func (context *storageContext) storageAdded(length int, key []byte, value []byte) (arwen.StorageStatus, error) {
	metering := context.host.Metering()
	useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.StorePerByte, uint64(length))
	metering.UseGas(useGas)
	logStorage.Trace("storage added", "key", key, "value", value)
	return arwen.StorageAdded, nil
}

func (context *storageContext) storageDeleted(lengthOldValue int, key []byte) (arwen.StorageStatus, error) {
	metering := context.host.Metering()
	freeGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.ReleasePerByte, uint64(lengthOldValue))
	metering.FreeGas(freeGas)
	logStorage.Trace("storage deleted", "key", key)
	return arwen.StorageDeleted, nil
}

func (context *storageContext) storageUnchanged(length int, usedCache bool) (arwen.StorageStatus, error) {
	useGas := context.computeGasForUnchangedValue(length, usedCache)
	context.host.Metering().UseGas(useGas)
	logStorage.Trace("storage set to identical value")
	return arwen.StorageUnchanged, nil
}

func (context *storageContext) computeGasForUnchangedValue(length int, usedCache bool) uint64 {
	metering := context.host.Metering()
	useGas := uint64(0)
	enableEpochsHandler := context.host.EnableEpochsHandler()
	if !usedCache || !enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled() {
		useGas = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	}
	return useGas
}

func (context *storageContext) getOldValue(storageUpdates map[string]*vmcommon.StorageUpdate, key []byte) ([]byte, bool, error) {
	var oldValue []byte
	var err error

	usedCache := true
	strKey := string(key)
	if update, ok := storageUpdates[strKey]; !ok {
		// if it's not in storageUpdates, GetStorageUnmetered() will use blockchain hook for sure
		oldValue, _, err = context.GetStorageUnmetered(key)
		if err != nil {
			return nil, false, err
		}
		storageUpdates[strKey] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   oldValue,
		}
		usedCache = false
	} else {
		oldValue = update.Data
	}
	return oldValue, usedCache, nil
}

func (context *storageContext) computeGasForKey(key []byte, usedCache bool) uint64 {
	metering := context.host.Metering()
	extraBytes := len(key) - arwen.AddressLen
	extraKeyLenGas := uint64(0)
	enableEpochsHandler := context.host.EnableEpochsHandler()
	if extraBytes > 0 &&
		(!usedCache || !enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled()) {
		extraKeyLenGas = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
	}
	return extraKeyLenGas
}

// UseGasForStorageLoad - single spot of gas consumption for storage load
func (context *storageContext) UseGasForStorageLoad(tracedFunctionName string, loadCost uint64, usedCache bool) {
	metering := context.host.Metering()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	if enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled() && usedCache {
		loadCost = metering.GasSchedule().ElrondAPICost.CachedStorageLoad
	}

	metering.UseGasAndAddTracedGas(tracedFunctionName, loadCost)
}

// IsUseDifferentGasCostFlagSet - getter for flag
func (context *storageContext) IsUseDifferentGasCostFlagSet() bool {
	enableEpochsHandler := context.host.EnableEpochsHandler()
	return enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled()
}

// IsInterfaceNil returns true if there is no value under the interface
func (context *storageContext) IsInterfaceNil() bool {
	return context == nil
}

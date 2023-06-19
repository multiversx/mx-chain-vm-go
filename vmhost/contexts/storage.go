package contexts

import (
	"bytes"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// TODO audit and verify consistency of all GetStorage*() methods

var logStorage = logger.GetOrCreate("vm/storage")

var _ vmhost.StorageContext = (*storageContext)(nil)

// VMStoragePrefix defines the VM prefix
const VMStoragePrefix = "VM@"

type storageContext struct {
	host                       vmhost.VMHost
	blockChainHook             vmcommon.BlockchainHook
	address                    []byte
	stateStack                 [][]byte
	protectedKeyPrefix         []byte
	vmProtectedKeyPrefix       []byte
	vmStorageProtectionEnabled bool
}

// NewStorageContext creates a new storageContext
func NewStorageContext(
	host vmhost.VMHost,
	blockChainHook vmcommon.BlockchainHook,
	protectedKeyPrefix []byte,
) (*storageContext, error) {
	if len(protectedKeyPrefix) == 0 {
		return nil, vmhost.ErrEmptyProtectedKeyPrefix
	}
	if check.IfNil(host) {
		return nil, vmhost.ErrNilVMHost
	}
	if check.IfNil(blockChainHook) {
		return nil, vmhost.ErrNilBlockChainHook
	}

	if check.IfNil(host) {
		return nil, vmhost.ErrNilVMHost
	}

	context := &storageContext{
		host:                       host,
		blockChainHook:             blockChainHook,
		stateStack:                 make([][]byte, 0),
		protectedKeyPrefix:         protectedKeyPrefix,
		vmProtectedKeyPrefix:       append(protectedKeyPrefix, []byte(VMStoragePrefix)...),
		vmStorageProtectionEnabled: true,
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
func (context *storageContext) GetStorage(key []byte) ([]byte, uint32, bool, error) {
	value, trieDepth, usedCache, err := context.GetStorageUnmetered(key)
	if err != nil {
		return nil, trieDepth, false, err
	}
	context.useExtraGasForKeyIfNeeded(key, usedCache)
	context.useGasForValueIfNeeded(value, usedCache)
	logStorage.Trace("get", "key", key, "value", value)

	return value, trieDepth, usedCache, nil
}

func (context *storageContext) useGasForValueIfNeeded(value []byte, usedCache bool) {
	metering := context.host.Metering()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	gasFlagSet := enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled()
	if !usedCache || !gasFlagSet {
		costPerByte := metering.GasSchedule().BaseOperationCost.DataCopyPerByte
		gasToUse := math.MulUint64(costPerByte, uint64(len(value)))
		// TODO replace UseGas with UseGasBounded
		metering.UseGas(gasToUse)
	}
}

func (context *storageContext) useExtraGasForKeyIfNeeded(key []byte, usedCache bool) {
	metering := context.host.Metering()
	extraBytes := len(key) - vmhost.AddressLen
	if extraBytes <= 0 {
		return
	}
	enableEpochsHandler := context.host.EnableEpochsHandler()
	gasFlagSet := enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled()
	if !gasFlagSet || !usedCache {
		gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
		// TODO replace UseGas with UseGasBounded
		metering.UseGas(gasToUse)
	}
}

// GetStorageFromAddress returns the data under the given key from the account mapped to the given address.
func (context *storageContext) GetStorageFromAddress(address []byte, key []byte) ([]byte, uint32, bool, error) {
	if !bytes.Equal(address, context.address) {
		userAcc, err := context.blockChainHook.GetUserAccount(address)
		if err != nil || check.IfNil(userAcc) {
			context.useExtraGasForKeyIfNeeded(key, false)
			return nil, 0, false, nil
		}

		metadata := vmcommon.CodeMetadataFromBytes(userAcc.GetCodeMetadata())
		if !metadata.Readable {
			context.useExtraGasForKeyIfNeeded(key, false)
			return nil, 0, false, nil
		}
	}

	return context.GetStorageFromAddressNoChecks(address, key)
}

// GetStorageFromAddressNoChecks same as GetStorageFromAddress but used internaly by vm, so no permissions checks are necessary
func (context *storageContext) GetStorageFromAddressNoChecks(address []byte, key []byte) ([]byte, uint32, bool, error) {
	// If the requested key is protected by the node, the stored value
	// could have been changed by a built-in function in the meantime, even if
	// contracts themselves cannot change protected values. Values stored under
	// protected keys must always be retrieved from the node, not from the cached
	// StorageUpdates.
	value, trieDepth, usedCache, err := context.getStorageFromAddressUnmetered(address, key)

	context.useExtraGasForKeyIfNeeded(key, usedCache)
	context.useGasForValueIfNeeded(value, usedCache)

	logStorage.Trace("get from address", "address", address, "key", key, "value", value)
	return value, trieDepth, usedCache, err
}

func (context *storageContext) getStorageFromAddressUnmetered(address []byte, key []byte) ([]byte, uint32, bool, error) {
	var value []byte
	var err error
	var trieDepth uint32

	enableEpochsHandler := context.host.EnableEpochsHandler()
	if context.isProtocolProtectedKey(key) && enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled() {
		value, trieDepth, err = context.readFromBlockchain(address, key)
		return value, trieDepth, false, err
	}

	storageUpdates := context.GetStorageUpdates(address)
	usedCache := true
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		value = storageUpdate.Data
	} else {
		value, trieDepth, err = context.readFromBlockchain(address, key)
		if err != nil {
			return nil, trieDepth, false, err
		}
		storageUpdates[string(key)] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   value,
		}
		usedCache = false
	}

	return value, trieDepth, usedCache, nil
}

func (context *storageContext) readFromBlockchain(address []byte, key []byte) ([]byte, uint32, error) {
	return context.blockChainHook.GetStorageData(address, key)
}

// GetStorageUnmetered returns the data under the given key.
func (context *storageContext) GetStorageUnmetered(key []byte) ([]byte, uint32, bool, error) {
	return context.getStorageFromAddressUnmetered(context.address, key)
}

// enableStorageProtection will prevent writing to protected keys
func (context *storageContext) enableStorageProtection() {
	context.vmStorageProtectionEnabled = true
}

// disableStorageProtection will prevent writing to protected keys
func (context *storageContext) disableStorageProtection() {
	context.vmStorageProtectionEnabled = false
}

func (context *storageContext) isVMProtectedKey(key []byte) bool {
	return bytes.HasPrefix(key, context.vmProtectedKeyPrefix)
}

func (context *storageContext) isProtocolProtectedKey(key []byte) bool {
	return bytes.HasPrefix(key, context.protectedKeyPrefix)
}

// SetProtectedStorage sets storage for timelocks and promises
func (context *storageContext) SetProtectedStorage(key []byte, value []byte) (vmhost.StorageStatus, error) {
	context.disableStorageProtection()
	defer context.enableStorageProtection()
	return context.SetStorage(key, value)
}

// SetStorage sets the given value at the given key.
func (context *storageContext) SetStorage(key []byte, value []byte) (vmhost.StorageStatus, error) {
	return context.setStorageToAddress(context.address, key, value)
}

// SetProtectedStorageToAddress sets the given value at the given key, for the specified address. This is only used internaly by vm!
func (context *storageContext) SetProtectedStorageToAddress(address []byte, key []byte, value []byte) (vmhost.StorageStatus, error) {
	context.disableStorageProtection()
	defer context.enableStorageProtection()
	return context.setStorageToAddress(address, key, value)
}

// SetProtectedStorageToAddressUnmetered sets the given value at the given key, for the specified address. This is only used internaly by vm!
// No gas cost involved, e.g. called by async.Save()
func (context *storageContext) SetProtectedStorageToAddressUnmetered(address []byte, key []byte, value []byte) (vmhost.StorageStatus, error) {
	context.disableStorageProtection()
	defer context.enableStorageProtection()
	return context.setStorageToAddressUnmetered(address, key, value)
}

func (context *storageContext) setStorageToAddress(address []byte, key []byte, value []byte) (vmhost.StorageStatus, error) {
	err := context.checkReservedAndProtection(key)
	if err != nil {
		return vmhost.StorageUnchanged, err
	}
	metering := context.host.Metering()

	length := len(value)

	storageUpdates := context.GetStorageUpdates(address)
	oldValue, usedCache, err := context.getOldValue(storageUpdates, key)
	if err != nil {
		return vmhost.StorageUnchanged, err
	}

	gasForKey := context.computeGasForKey(key, usedCache)
	err = metering.UseGasBounded(gasForKey)
	if err != nil {
		return vmhost.StorageUnchanged, err
	}

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

	err = metering.UseGasBounded(gasToUseForValue)
	if err != nil {
		return vmhost.StorageUnchanged, err
	}

	metering.FreeGas(gasToFreeForValue)

	logStorage.Trace("storage modified", "key", key, "value", value, "lengthDelta", newValueExtraLength)
	return vmhost.StorageModified, nil
}

func (context *storageContext) setStorageToAddressUnmetered(address []byte, key []byte, value []byte) (vmhost.StorageStatus, error) {
	err := context.checkReservedAndProtection(key)
	if err != nil {
		return vmhost.StorageUnchanged, err
	}

	storageUpdates := context.GetStorageUpdates(address)
	context.changeStorageUpdate(key, value, storageUpdates)

	logStorage.Trace("storage modified (unmetered)", "key", key, "value", value)
	return vmhost.StorageModified, nil
}

func (context *storageContext) checkReservedAndProtection(key []byte) error {
	if context.host.Runtime().ReadOnly() {
		logStorage.Trace("storage set", "error", "cannot set storage in readonly mode")
		return vmhost.ErrCannotWriteOnReadOnly
	}
	if !context.isVMProtectedKey(key) && context.isProtocolProtectedKey(key) {
		logStorage.Trace("storage set", "error", vmhost.ErrStoreReservedKey, "key", key)
		return vmhost.ErrStoreReservedKey
	}
	if context.isVMProtectedKey(key) && context.vmStorageProtectionEnabled {
		logStorage.Trace("storage set", "error", vmhost.ErrCannotWriteProtectedKey, "key", key)
		return vmhost.ErrCannotWriteProtectedKey
	}
	return nil
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

func (context *storageContext) storageAdded(length int, key []byte, value []byte) (vmhost.StorageStatus, error) {
	metering := context.host.Metering()
	useGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.StorePerByte, uint64(length))
	err := metering.UseGasBounded(useGas)
	if err != nil {
		return vmhost.StorageUnchanged, err
	}

	logStorage.Trace("storage added", "key", key, "value", value)
	return vmhost.StorageAdded, nil
}

func (context *storageContext) storageDeleted(lengthOldValue int, key []byte) (vmhost.StorageStatus, error) {
	metering := context.host.Metering()
	freeGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.ReleasePerByte, uint64(lengthOldValue))
	metering.FreeGas(freeGas)

	logStorage.Trace("storage deleted", "key", key)
	return vmhost.StorageDeleted, nil
}

func (context *storageContext) storageUnchanged(length int, usedCache bool) (vmhost.StorageStatus, error) {
	useGas := context.computeGasForUnchangedValue(length, usedCache)
	err := context.host.Metering().UseGasBounded(useGas)
	if err != nil {
		return vmhost.StorageUnchanged, err
	}

	logStorage.Trace("storage set to identical value")
	return vmhost.StorageUnchanged, nil
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
		oldValue, _, _, err = context.GetStorageUnmetered(key)
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
	extraBytes := len(key) - vmhost.AddressLen
	extraKeyLenGas := uint64(0)
	enableEpochsHandler := context.host.EnableEpochsHandler()
	if extraBytes > 0 &&
		(!usedCache || !enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled()) {
		extraKeyLenGas = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(extraBytes))
	}
	return extraKeyLenGas
}

// UseGasForStorageLoad - single spot of gas consumption for storage load
func (context *storageContext) UseGasForStorageLoad(tracedFunctionName string, trieDepth int64, staticGasCost uint64, usedCache bool) error {
	blockchainLoadCost, err := context.getBlockchainLoadCost(trieDepth, staticGasCost, usedCache)
	if err != nil {
		return err
	}

	return context.host.Metering().UseGasBoundedAndAddTracedGas(tracedFunctionName, blockchainLoadCost)
}

func (context *storageContext) getBlockchainLoadCost(trieDepth int64, staticGasCost uint64, usedCache bool) (uint64, error) {
	enableEpochsHandler := context.host.EnableEpochsHandler()
	if enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabled() && usedCache {
		return context.host.Metering().GasSchedule().BaseOpsAPICost.CachedStorageLoad, nil
	}

	return context.GetStorageLoadCost(trieDepth, staticGasCost)
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

// GetVmProtectedPrefix returns the VM protected prefix as byte slice
func (context *storageContext) GetVmProtectedPrefix(prefix string) []byte {
	return append(context.vmProtectedKeyPrefix, []byte(prefix)...)
}

// GetStorageLoadCost returns the gas cost for the storage load operation
func (context *storageContext) GetStorageLoadCost(trieDepth int64, staticGasCost uint64) (uint64, error) {
	if context.host.EnableEpochsHandler().IsDynamicGasCostForDataTrieStorageLoadEnabled() {
		return computeGasForStorageLoadBasedOnTrieDepth(
			trieDepth,
			context.host.Metering().GasSchedule().DynamicStorageLoad,
			staticGasCost,
		)
	}

	return staticGasCost, nil
}

func computeGasForStorageLoadBasedOnTrieDepth(trieDepth int64, coefficients config.DynamicStorageLoadCostCoefficients, staticGasCost uint64) (uint64, error) {
	overflowHandler := math.NewOverflowHandler()

	squaredTrieDepth := overflowHandler.MulInt64(trieDepth, trieDepth)                  // squaredTrieDepth = trieDepth * trieDepth
	quadraticTerm := overflowHandler.MulInt64(coefficients.Quadratic, squaredTrieDepth) // quadraticTerm = coefficients.Quadratic * trieDepth * trieDepth

	linearTerm := overflowHandler.MulInt64(coefficients.Linear, trieDepth) // linearTerm = coefficients.Linear * trieDepth

	firstSum := overflowHandler.AddInt64(quadraticTerm, linearTerm) // firstSum = coefficients.Quadratic * trieDepth * trieDepth + coefficients.Linear * trieDepth
	fx := overflowHandler.AddInt64(firstSum, coefficients.Constant) // fx = coefficients.Quadratic * trieDepth * trieDepth + coefficients.Linear * trieDepth + coefficients.Constant
	err := overflowHandler.Error()
	if err != nil {
		return 0, err
	}
	if fx < 0 {
		return 0, fmt.Errorf("invalid value for gas cost, quadratic coefficient = %v, linear coefficient = %v, constant coefficient = %v, trie depth = %v",
			coefficients.Quadratic, coefficients.Linear, coefficients.Constant, trieDepth)
	}

	if fx < int64(coefficients.MinGasCost) {
		logStorage.Error("invalid value for gas cost",
			"quadratic coefficient", coefficients.Quadratic,
			"linear coefficient", coefficients.Linear,
			"constant coefficient", coefficients.Constant,
			"trie depth", trieDepth,
			"min gas cost", coefficients.MinGasCost,
		)
		return staticGasCost, nil
	}

	return uint64(fx), nil
}

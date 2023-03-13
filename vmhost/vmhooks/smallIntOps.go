package vmhooks

import (
	"math/big"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

const (
	smallIntGetUnsignedArgumentName  = "smallIntGetUnsignedArgument"
	smallIntGetSignedArgumentName    = "smallIntGetSignedArgument"
	smallIntFinishUnsignedName       = "smallIntFinishUnsigned"
	smallIntFinishSignedName         = "smallIntFinishSigned"
	smallIntStorageStoreUnsignedName = "smallIntStorageStoreUnsigned"
	smallIntStorageStoreSignedName   = "smallIntStorageStoreSigned"
	smallIntStorageLoadUnsignedName  = "smallIntStorageLoadUnsigned"
	smallIntStorageLoadSignedName    = "smallIntStorageLoadSigned"
	int64getArgumentName             = "int64getArgument"
	int64storageStoreName            = "int64storageStore"
	int64storageLoadName             = "int64storageLoad"
	int64finishName                  = "int64finish"
)

// SmallIntGetUnsignedArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntGetUnsignedArgument(id int32) int64 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64GetArgument
	metering.UseGasAndAddTracedGas(smallIntGetUnsignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		_ = context.WithFault(vmhost.ErrArgIndexOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	arg := args[id]
	argBigInt := big.NewInt(0).SetBytes(arg)
	if !argBigInt.IsUint64() {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}
	return int64(argBigInt.Uint64())
}

// SmallIntGetSignedArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntGetSignedArgument(id int32) int64 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64GetArgument
	metering.UseGasAndAddTracedGas(smallIntGetSignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		_ = context.WithFault(vmhost.ErrArgIndexOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	arg := args[id]
	argBigInt := twos.SetBytes(big.NewInt(0), arg)
	if !argBigInt.IsInt64() {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}
	return argBigInt.Int64()
}

// SmallIntFinishUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntFinishUnsigned(value int64) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64Finish
	metering.UseGasAndAddTracedGas(smallIntFinishUnsignedName, gasToUse)

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	output.Finish(valueBytes)
}

// SmallIntFinishSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntFinishSigned(value int64) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64Finish
	metering.UseGasAndAddTracedGas(smallIntFinishSignedName, gasToUse)

	valueBytes := twos.ToBytes(big.NewInt(value))
	output.Finish(valueBytes)
}

// SmallIntStorageStoreUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageStoreUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength, value int64) int32 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64StorageStore
	metering.UseGasAndAddTracedGas(smallIntStorageStoreSignedName, gasToUse)

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

// SmallIntStorageStoreSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageStoreSigned(keyOffset executor.MemPtr, keyLength executor.MemLength, value int64) int32 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64StorageStore
	metering.UseGasAndAddTracedGas(smallIntStorageStoreSignedName, gasToUse)

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	valueBytes := twos.ToBytes(big.NewInt(value))
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

// SmallIntStorageLoadUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageLoadUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength) int64 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	err = storage.UseGasForStorageLoad(
		smallIntStorageLoadUnsignedName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.Int64StorageLoad,
		usedCache)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	valueBigInt := big.NewInt(0).SetBytes(data)
	if !valueBigInt.IsUint64() {
		_ = context.WithFault(vmhost.ErrStorageValueOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return int64(valueBigInt.Uint64())
}

// SmallIntStorageLoadSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageLoadSigned(keyOffset executor.MemPtr, keyLength executor.MemLength) int64 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	err = storage.UseGasForStorageLoad(
		smallIntStorageLoadSignedName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.Int64StorageLoad,
		usedCache)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	valueBigInt := twos.SetBytes(big.NewInt(0), data)
	if !valueBigInt.IsInt64() {
		_ = context.WithFault(vmhost.ErrStorageValueOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return valueBigInt.Int64()
}

// Int64getArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64getArgument(id int32) int64 {
	// backwards compatibility
	return context.SmallIntGetSignedArgument(id)
}

// Int64finish VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64finish(value int64) {
	// backwards compatibility
	context.SmallIntFinishSigned(value)
}

// Int64storageStore VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64storageStore(keyOffset executor.MemPtr, keyLength executor.MemLength, value int64) int32 {
	// backwards compatibility
	return context.SmallIntStorageStoreUnsigned(keyOffset, keyLength, value)
}

// Int64storageLoad VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64storageLoad(keyOffset executor.MemPtr, keyLength executor.MemLength) int64 {
	// backwards compatibility
	return context.SmallIntStorageLoadUnsigned(keyOffset, keyLength)
}

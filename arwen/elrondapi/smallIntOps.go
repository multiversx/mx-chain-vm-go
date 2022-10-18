package elrondapi

import (
	"math/big"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
	"github.com/ElrondNetwork/wasm-vm/arwen"
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
func (context *ElrondApi) SmallIntGetUnsignedArgument(id int32) int64 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64GetArgument
	metering.UseGasAndAddTracedGas(smallIntGetUnsignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		_ = context.WithFault(arwen.ErrArgIndexOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	arg := args[id]
	argBigInt := big.NewInt(0).SetBytes(arg)
	if !argBigInt.IsUint64() {
		_ = context.WithFault(arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}
	return int64(argBigInt.Uint64())
}

// SmallIntGetSignedArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) SmallIntGetSignedArgument(id int32) int64 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64GetArgument
	metering.UseGasAndAddTracedGas(smallIntGetSignedArgumentName, gasToUse)

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		_ = context.WithFault(arwen.ErrArgIndexOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	arg := args[id]
	argBigInt := twos.SetBytes(big.NewInt(0), arg)
	if !argBigInt.IsInt64() {
		_ = context.WithFault(arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}
	return argBigInt.Int64()
}

// SmallIntFinishUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) SmallIntFinishUnsigned(value int64) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64Finish
	metering.UseGasAndAddTracedGas(smallIntFinishUnsignedName, gasToUse)

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	output.Finish(valueBytes)
}

// SmallIntFinishSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) SmallIntFinishSigned(value int64) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64Finish
	metering.UseGasAndAddTracedGas(smallIntFinishSignedName, gasToUse)

	valueBytes := twos.ToBytes(big.NewInt(value))
	output.Finish(valueBytes)
}

// SmallIntStorageStoreUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) SmallIntStorageStoreUnsigned(keyOffset int32, keyLength int32, value int64) int32 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageStore
	metering.UseGasAndAddTracedGas(smallIntStorageStoreSignedName, gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if context.WithFault(err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

// SmallIntStorageStoreSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) SmallIntStorageStoreSigned(keyOffset int32, keyLength int32, value int64) int32 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageStore
	metering.UseGasAndAddTracedGas(smallIntStorageStoreSignedName, gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	valueBytes := twos.ToBytes(big.NewInt(value))
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if context.WithFault(err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

// SmallIntStorageLoadUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) SmallIntStorageLoadUnsigned(keyOffset int32, keyLength int32) int64 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	data, usedCache := storage.GetStorage(key)
	storage.UseGasForStorageLoad(smallIntStorageLoadUnsignedName, metering.GasSchedule().ElrondAPICost.Int64StorageLoad, usedCache)

	valueBigInt := big.NewInt(0).SetBytes(data)
	if !valueBigInt.IsUint64() {
		_ = context.WithFault(arwen.ErrStorageValueOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	return int64(valueBigInt.Uint64())
}

// SmallIntStorageLoadSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) SmallIntStorageLoadSigned(keyOffset int32, keyLength int32) int64 {
	runtime := context.GetRuntimeContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	data, usedCache := storage.GetStorage(key)
	storage.UseGasForStorageLoad(smallIntStorageLoadSignedName, metering.GasSchedule().ElrondAPICost.Int64StorageLoad, usedCache)

	valueBigInt := twos.SetBytes(big.NewInt(0), data)
	if !valueBigInt.IsInt64() {
		_ = context.WithFault(arwen.ErrStorageValueOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	return valueBigInt.Int64()
}

// Int64getArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) Int64getArgument(id int32) int64 {
	// backwards compatibility
	return context.SmallIntGetSignedArgument(id)
}

// Int64finish VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) Int64finish(value int64) {
	// backwards compatibility
	context.SmallIntFinishSigned(value)
}

// Int64storageStore VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) Int64storageStore(keyOffset int32, keyLength int32, value int64) int32 {
	// backwards compatibility
	return context.SmallIntStorageStoreUnsigned(keyOffset, keyLength, value)
}

// Int64storageLoad VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) Int64storageLoad(keyOffset int32, keyLength int32) int64 {
	// backwards compatibility
	return context.SmallIntStorageLoadUnsigned(keyOffset, keyLength)
}

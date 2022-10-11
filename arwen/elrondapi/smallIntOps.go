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

// SmallIntGetUnsignedArgument EIInterface implementation
func (context *EICallbacks) SmallIntGetUnsignedArgument(id int32) int64 {
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

// SmallIntGetSignedArgument EIInterface implementation
func (context *EICallbacks) SmallIntGetSignedArgument(id int32) int64 {
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

// SmallIntFinishUnsigned EIInterface implementation
func (context *EICallbacks) SmallIntFinishUnsigned(value int64) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64Finish
	metering.UseGasAndAddTracedGas(smallIntFinishUnsignedName, gasToUse)

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	output.Finish(valueBytes)
}

// SmallIntFinishSigned EIInterface implementation
func (context *EICallbacks) SmallIntFinishSigned(value int64) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64Finish
	metering.UseGasAndAddTracedGas(smallIntFinishSignedName, gasToUse)

	valueBytes := twos.ToBytes(big.NewInt(value))
	output.Finish(valueBytes)
}

// SmallIntStorageStoreUnsigned EIInterface implementation
func (context *EICallbacks) SmallIntStorageStoreUnsigned(keyOffset int32, keyLength int32, value int64) int32 {
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

// SmallIntStorageStoreSigned EIInterface implementation
func (context *EICallbacks) SmallIntStorageStoreSigned(keyOffset int32, keyLength int32, value int64) int32 {
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

// SmallIntStorageLoadUnsigned EIInterface implementation
func (context *EICallbacks) SmallIntStorageLoadUnsigned(keyOffset int32, keyLength int32) int64 {
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

// SmallIntStorageLoadSigned EIInterface implementation
func (context *EICallbacks) SmallIntStorageLoadSigned(keyOffset int32, keyLength int32) int64 {
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

// Int64getArgument EIInterface implementation
func (context *EICallbacks) Int64getArgument(id int32) int64 {
	// backwards compatibility
	return context.SmallIntGetSignedArgument(id)
}

// Int64finish EIInterface implementation
func (context *EICallbacks) Int64finish(value int64) {
	// backwards compatibility
	context.SmallIntFinishSigned(value)
}

// Int64storageStore EIInterface implementation
func (context *EICallbacks) Int64storageStore(keyOffset int32, keyLength int32, value int64) int32 {
	// backwards compatibility
	return context.SmallIntStorageStoreUnsigned(keyOffset, keyLength, value)
}

// Int64storageLoad EIInterface implementation
func (context *EICallbacks) Int64storageLoad(keyOffset int32, keyLength int32) int64 {
	// backwards compatibility
	return context.SmallIntStorageLoadUnsigned(keyOffset, keyLength)
}

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
func (context *VMHooksImpl) SmallIntGetUnsignedArgument(id int32) (int64, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(smallIntGetUnsignedArgumentName, gasToUse)
	if err != nil {
		return 1, err
	}

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		return 0, vmhost.ErrArgIndexOutOfRange
	}

	arg := args[id]
	argBigInt := big.NewInt(0).SetBytes(arg)
	if !argBigInt.IsUint64() {
		return 0, vmhost.ErrArgOutOfRange
	}
	return int64(argBigInt.Uint64()), nil
}

// SmallIntGetSignedArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntGetSignedArgument(id int32) (int64, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(smallIntGetSignedArgumentName, gasToUse)
	if err != nil {
		return 1, err
	}

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		return 0, vmhost.ErrArgIndexOutOfRange
	}

	arg := args[id]
	argBigInt := twos.SetBytes(big.NewInt(0), arg)
	if !argBigInt.IsInt64() {
		return 0, vmhost.ErrArgOutOfRange
	}
	return argBigInt.Int64(), nil
}

// SmallIntFinishUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntFinishUnsigned(value int64) error {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64Finish
	err := metering.UseGasBoundedAndAddTracedGas(smallIntFinishUnsignedName, gasToUse)
	if err != nil {
		return err
	}

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	output.Finish(valueBytes)
	return nil
}

// SmallIntFinishSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntFinishSigned(value int64) error {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64Finish
	err := metering.UseGasBoundedAndAddTracedGas(smallIntFinishSignedName, gasToUse)
	if err != nil {
		return err
	}

	valueBytes := twos.ToBytes(big.NewInt(value))
	output.Finish(valueBytes)
	return nil
}

// SmallIntStorageStoreUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageStoreUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength, value int64) (int32, error) {
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64StorageStore
	err := metering.UseGasBoundedAndAddTracedGas(smallIntStorageStoreSignedName, gasToUse)
	if err != nil {
		return -1, err
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if err != nil {
		return -1, err
	}

	return int32(storageStatus), nil
}

// SmallIntStorageStoreSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageStoreSigned(keyOffset executor.MemPtr, keyLength executor.MemLength, value int64) (int32, error) {
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64StorageStore
	err := metering.UseGasBoundedAndAddTracedGas(smallIntStorageStoreSignedName, gasToUse)
	if err != nil {
		return -1, err
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	valueBytes := twos.ToBytes(big.NewInt(value))
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if err != nil {
		return -1, err
	}

	return int32(storageStatus), nil
}

// SmallIntStorageLoadUnsigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageLoadUnsigned(keyOffset executor.MemPtr, keyLength executor.MemLength) (int64, error) {
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return 0, err
	}

	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return 0, err
	}

	err = storage.UseGasForStorageLoad(
		smallIntStorageLoadUnsignedName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.Int64StorageLoad,
		usedCache)
	if err != nil {
		return -1, err
	}

	valueBigInt := big.NewInt(0).SetBytes(data)
	if !valueBigInt.IsUint64() {
		return 0, vmhost.ErrStorageValueOutOfRange
	}

	return int64(valueBigInt.Uint64()), nil
}

// SmallIntStorageLoadSigned VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SmallIntStorageLoadSigned(keyOffset executor.MemPtr, keyLength executor.MemLength) (int64, error) {
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return 0, err
	}

	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return 0, err
	}

	err = storage.UseGasForStorageLoad(
		smallIntStorageLoadSignedName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.Int64StorageLoad,
		usedCache)
	if err != nil {
		return -1, err
	}

	valueBigInt := twos.SetBytes(big.NewInt(0), data)
	if !valueBigInt.IsInt64() {
		return 0, vmhost.ErrStorageValueOutOfRange
	}

	return valueBigInt.Int64(), nil
}

// Int64getArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64getArgument(id int32) (int64, error) {
	// backwards compatibility
	return context.SmallIntGetSignedArgument(id)
}

// Int64finish VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64finish(value int64) error {
	// backwards compatibility
	return context.SmallIntFinishSigned(value)
}

// Int64storageStore VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64storageStore(keyOffset executor.MemPtr, keyLength executor.MemLength, value int64) (int32, error) {
	// backwards compatibility
	return context.SmallIntStorageStoreUnsigned(keyOffset, keyLength, value)
}

// Int64storageLoad VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Int64storageLoad(keyOffset executor.MemPtr, keyLength executor.MemLength) (int64, error) {
	// backwards compatibility
	return context.SmallIntStorageLoadUnsigned(keyOffset, keyLength)
}

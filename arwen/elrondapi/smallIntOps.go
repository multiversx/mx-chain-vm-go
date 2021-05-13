package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern long long smallIntGetUnsignedArgument(void *context, int32_t id);
// extern long long smallIntGetSignedArgument(void *context, int32_t id);
//
// extern void smallIntFinishUnsigned(void* context, long long value);
// extern void smallIntFinishSigned(void* context, long long value);
//
// extern int32_t smallIntStorageStoreUnsigned(void *context, int32_t keyOffset, int32_t keyLength, long long value);
// extern int32_t smallIntStorageStoreSigned(void *context, int32_t keyOffset, int32_t keyLength, long long value);
// extern long long smallIntStorageLoadUnsigned(void *context, int32_t keyOffset, int32_t keyLength);
// extern long long smallIntStorageLoadSigned(void *context, int32_t keyOffset, int32_t keyLength);
//
// extern long long int64getArgument(void *context, int32_t id);
// extern int32_t int64storageStore(void *context, int32_t keyOffset, int32_t keyLength , long long value);
// extern long long int64storageLoad(void *context, int32_t keyOffset, int32_t keyLength );
// extern void int64finish(void* context, long long value);
//
import "C"

import (
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1.3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1.3/wasmer"
	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

// SmallIntImports creates a new wasmer.Imports populated with the small int (int64/uint64) API methods
func SmallIntImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("smallIntGetUnsignedArgument", smallIntGetUnsignedArgument, C.smallIntGetUnsignedArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("smallIntGetSignedArgument", smallIntGetSignedArgument, C.smallIntGetSignedArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("smallIntFinishUnsigned", smallIntFinishUnsigned, C.smallIntFinishUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("smallIntFinishSigned", smallIntFinishSigned, C.smallIntFinishSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("smallIntStorageStoreUnsigned", smallIntStorageStoreUnsigned, C.smallIntStorageStoreUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("smallIntStorageStoreSigned", smallIntStorageStoreSigned, C.smallIntStorageStoreSigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("smallIntStorageLoadUnsigned", smallIntStorageLoadUnsigned, C.smallIntStorageLoadUnsigned)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("smallIntStorageLoadSigned", smallIntStorageLoadSigned, C.smallIntStorageLoadSigned)
	if err != nil {
		return nil, err
	}

	// the last are just for backwards compatibility:

	imports, err = imports.Append("int64getArgument", int64getArgument, C.int64getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageStore", int64storageStore, C.int64storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageLoad", int64storageLoad, C.int64storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64finish", int64finish, C.int64finish)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export smallIntGetUnsignedArgument
func smallIntGetUnsignedArgument(context unsafe.Pointer, id int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		arwen.WithFault(arwen.ErrArgIndexOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	arg := args[id]
	argBigInt := big.NewInt(0).SetBytes(arg)
	if !argBigInt.IsUint64() {
		arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}
	return int64(argBigInt.Uint64())
}

//export smallIntGetSignedArgument
func smallIntGetSignedArgument(context unsafe.Pointer, id int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		arwen.WithFault(arwen.ErrArgIndexOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	arg := args[id]
	argBigInt := twos.SetBytes(big.NewInt(0), arg)
	if !argBigInt.IsInt64() {
		arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}
	return argBigInt.Int64()
}

//export smallIntFinishUnsigned
func smallIntFinishUnsigned(context unsafe.Pointer, value int64) {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64Finish
	metering.UseGas(gasToUse)

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	output.Finish(valueBytes)
}

//export smallIntFinishSigned
func smallIntFinishSigned(context unsafe.Pointer, value int64) {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64Finish
	metering.UseGas(gasToUse)

	valueBytes := twos.ToBytes(big.NewInt(value))
	output.Finish(valueBytes)
}

//export smallIntStorageStoreUnsigned
func smallIntStorageStoreUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, value int64) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	valueBytes := big.NewInt(0).SetUint64(uint64(value)).Bytes()
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export smallIntStorageStoreSigned
func smallIntStorageStoreSigned(context unsafe.Pointer, keyOffset int32, keyLength int32, value int64) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	valueBytes := twos.ToBytes(big.NewInt(value))
	storageStatus, err := storage.SetStorage(key, valueBytes)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export smallIntStorageLoadUnsigned
func smallIntStorageLoadUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	data := storage.GetStorage(key)
	valueBigInt := big.NewInt(0).SetBytes(data)
	if !valueBigInt.IsUint64() {
		arwen.WithFault(arwen.ErrStorageValueOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	return int64(valueBigInt.Uint64())
}

//export smallIntStorageLoadSigned
func smallIntStorageLoadSigned(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	data := storage.GetStorage(key)
	valueBigInt := twos.SetBytes(big.NewInt(0), data)
	if !valueBigInt.IsInt64() {
		arwen.WithFault(arwen.ErrStorageValueOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 0
	}

	return valueBigInt.Int64()
}

//export int64getArgument
func int64getArgument(context unsafe.Pointer, id int32) int64 {
	// backwards compatibility
	return smallIntGetSignedArgument(context, id)
}

//export int64finish
func int64finish(context unsafe.Pointer, value int64) {
	// backwards compatibility
	smallIntFinishSigned(context, value)
}

//export int64storageStore
func int64storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, value int64) int32 {
	// backwards compatibility
	return smallIntStorageStoreUnsigned(context, keyOffset, keyLength, value)
}

//export int64storageLoad
func int64storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	// backwards compatibility
	return smallIntStorageLoadUnsigned(context, keyOffset, keyLength)
}

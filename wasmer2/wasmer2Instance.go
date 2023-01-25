package wasmer2

// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

var _ = (executor.Instance)((*Wasmer2Instance)(nil))

// Wasmer2Instance represents a WebAssembly instance.
type Wasmer2Instance struct {
	// The underlying WebAssembly instance.
	cgoInstance *cWasmerInstanceT

	// The exported memory of a WebAssembly instance.
	memory Wasmer2Memory

	alreadyCleaned bool
}

func emptyInstance() *Wasmer2Instance {
	return &Wasmer2Instance{cgoInstance: nil}
}

func newInstance(c_instance *cWasmerInstanceT) (*Wasmer2Instance, error) {
	return &Wasmer2Instance{
		cgoInstance: c_instance,
		memory: Wasmer2Memory{
			cgoInstance: c_instance,
		},
	}, nil
}

// Clean cleans instance
func (instance *Wasmer2Instance) Clean() {
	if instance.alreadyCleaned {
		return
	}

	if instance.cgoInstance != nil {
		cWasmerInstanceDestroy(instance.cgoInstance)

		instance.alreadyCleaned = true
	}
}

// SetGasLimit sets the gas limit for the instance
func (instance *Wasmer2Instance) SetGasLimit(gasLimit uint64) {
	cWasmerInstanceSetGasLimit(instance.cgoInstance, gasLimit)
}

// SetPointsUsed sets the internal instance gas counter
func (instance *Wasmer2Instance) SetPointsUsed(points uint64) {
	cWasmerInstanceSetPointsUsed(instance.cgoInstance, points)
}

// GetPointsUsed returns the internal instance gas counter
func (instance *Wasmer2Instance) GetPointsUsed() uint64 {
	return cWasmerInstanceGetPointsUsed(instance.cgoInstance)
}

// SetBreakpointValue sets the breakpoint value for the instance
func (instance *Wasmer2Instance) SetBreakpointValue(value uint64) {
	cWasmerInstanceSetBreakpointValue(instance.cgoInstance, value)
}

// GetBreakpointValue returns the breakpoint value
func (instance *Wasmer2Instance) GetBreakpointValue() uint64 {
	return cWasmerInstanceGetBreakpointValue(instance.cgoInstance)
}

// Cache caches the instance
func (instance *Wasmer2Instance) Cache() ([]byte, error) {
	var cacheBytes *cUchar
	var cacheLen cUint32T

	var cacheResult = cWasmerInstanceCache(
		instance.cgoInstance,
		&cacheBytes,
		&cacheLen,
	)

	if cacheResult != cWasmerOk {
		return nil, ErrCachingFailed
	}

	goBytes := C.GoBytes(unsafe.Pointer(cacheBytes), C.int(cacheLen))

	C.free(unsafe.Pointer(cacheBytes))
	cacheBytes = nil
	return goBytes, nil
}

// IsFunctionImported returns true if the instance imports the specified function
func (instance *Wasmer2Instance) IsFunctionImported(name string) bool {
	return false
}

// CallFunction executes given function from loaded contract.
func (instance *Wasmer2Instance) CallFunction(functionName string) error {
	var wasmFunctionName = cCString(functionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	var callResult = cWasmerInstanceCall(
		instance.cgoInstance,
		wasmFunctionName,
	)

	if callResult != cWasmerOk {
		err := fmt.Errorf("failed to call the `%s` exported function", functionName)
		return newWrappedError(err)
	}

	return nil
}

// HasFunction checks if loaded contract has a function (endpoint) with given name.
func (instance *Wasmer2Instance) HasFunction(functionName string) bool {
	var wasmFunctionName = cCString(functionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	result := cWasmerInstanceHasFunction(
		instance.cgoInstance,
		wasmFunctionName,
	)

	return result == 1
}

// GetLastError returns the last error message if any, otherwise returns an error.
func (instance *Wasmer2Instance) getFunctionNamesConcat() (string, error) {
	var bufferLength = cWasmerInstanceExportedFunctionNamesLength(instance.cgoInstance)

	if bufferLength == 0 {
		return "", nil
	}

	var buffer = make([]cChar, bufferLength)
	var bufferPointer = (*cChar)(unsafe.Pointer(&buffer[0]))

	var result = cWasmerInstanceExportedFunctionNames(instance.cgoInstance, bufferPointer, bufferLength)

	if result == -1 {
		return "", errors.New("cannot read function names")
	}

	return cGoString(bufferPointer), nil
}

// GetFunctionNames returns a list of the function names exported by the contract.
func (instance *Wasmer2Instance) GetFunctionNames() []string {
	buffer, err := instance.getFunctionNamesConcat()
	if err != nil {
		return nil
	}
	return strings.Split(buffer, "|")
}

// ValidateFunctionArities checks that no function (endpoint) of the given contract has any parameters or returns any result.
// All arguments and results should be transferred via the import functions.
func (instance *Wasmer2Instance) ValidateFunctionArities() error {
	var result = cWasmerCheckSignatures(instance.cgoInstance)
	if result != cWasmerOk {
		return executor.ErrFunctionNonvoidSignature
	}
	return nil
}

// HasMemory checks whether the instance has at least one exported memory.
func (instance *Wasmer2Instance) HasMemory() bool {
	return true
}

// MemLoad returns the contents from the given offset of the WASM memory.
func (instance *Wasmer2Instance) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return executor.MemLoadFromMemory(&instance.memory, memPtr, length)
}

// MemStore stores the given data in the WASM memory at the given offset.
func (instance *Wasmer2Instance) MemStore(memPtr executor.MemPtr, data []byte) error {
	return executor.MemStoreToMemory(&instance.memory, memPtr, data)
}

// MemLength returns the length of the allocated memory. Only called directly in tests.
func (instance *Wasmer2Instance) MemLength() uint32 {
	return instance.memory.Length()
}

// MemGrow allocates more pages to the current memory. Only called directly in tests.
func (instance *Wasmer2Instance) MemGrow(pages uint32) error {
	return instance.memory.Grow(pages)
}

// MemDump yields the entire contents of the memory. Only used in tests.
func (instance *Wasmer2Instance) MemDump() []byte {
	return instance.memory.Data()
}

// Id returns an identifier for the instance, unique at runtime
func (instance *Wasmer2Instance) Id() string {
	return fmt.Sprintf("%p", instance.cgoInstance)
}

// Reset resets the instance memories and globals
func (instance *Wasmer2Instance) Reset() bool {
	if instance.alreadyCleaned {
		return false
	}

	result := cWasmerInstanceReset(instance.cgoInstance)
	return result == cWasmerOk
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *Wasmer2Instance) IsInterfaceNil() bool {
	return instance == nil
}

// SetVMHooksPtr sets the VM hooks pointer
func (instance *Wasmer2Instance) SetVMHooksPtr(vmHooksPtr uintptr) {
}

// GetVMHooksPtr returns the VM hooks pointer
func (instance *Wasmer2Instance) GetVMHooksPtr() uintptr {
	return uintptr(0)
}

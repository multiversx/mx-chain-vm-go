package wasmer2

// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.Instance = (*Wasmer2Instance)(nil)

var instanceIDCounter uint64

// Wasmer2Instance represents a WebAssembly instance.
type Wasmer2Instance struct {
	mut sync.RWMutex

	// The underlying WebAssembly instance.
	cgoInstance *cWasmerInstanceT

	// The exported memory of a WebAssembly instance.
	memory Wasmer2Memory

	AlreadyClean bool
	id           string
}

func emptyInstance() *Wasmer2Instance {
	id := atomic.AddUint64(&instanceIDCounter, 1)
	return &Wasmer2Instance{
		cgoInstance: nil,
		id:          fmt.Sprintf("empty-%d", id),
	}
}

// NewEmptyInstanceForTest creates an empty instance for testing
func NewEmptyInstanceForTest() *Wasmer2Instance {
	return emptyInstance()
}

func newInstance(c_instance *cWasmerInstanceT) (*Wasmer2Instance, error) {
	id := atomic.AddUint64(&instanceIDCounter, 1)
	return &Wasmer2Instance{
		cgoInstance: c_instance,
		memory: Wasmer2Memory{
			cgoInstance: c_instance,
		},
		id: fmt.Sprintf("%p-%d", c_instance, id),
	}, nil
}

// Clean cleans instance
func (instance *Wasmer2Instance) Clean() bool {
	instance.mut.Lock()
	defer instance.mut.Unlock()

	logWasmer2.Trace("cleaning instance", "id", instance.id)
	if instance.AlreadyClean {
		logWasmer2.Trace("clean: already cleaned instance", "id", instance.id)
		return false
	}

	if instance.cgoInstance != nil {
		cWasmerInstanceDestroy(instance.cgoInstance)

		instance.AlreadyClean = true
		instance.cgoInstance = nil
		instance.memory.cgoInstance = nil
		logWasmer2.Trace("cleaned instance", "id", instance.id)

		return true
	}

	return false
}

// IsAlreadyCleaned returns the internal field AlreadyClean
func (instance *Wasmer2Instance) IsAlreadyCleaned() bool {
	instance.mut.RLock()
	defer instance.mut.RUnlock()
	return instance.AlreadyClean
}

// Abandon marks the instance as cleaned without destroying CGO memory,
// avoiding crashes if the instance is still executing on an unresponsive CGO thread
func (instance *Wasmer2Instance) Abandon() {
	instance.mut.Lock()
	defer instance.mut.Unlock()
	instance.AlreadyClean = true
}

// SetGasLimit sets the gas limit for the instance
func (instance *Wasmer2Instance) SetGasLimit(gasLimit uint64) {
	instance.mut.RLock()
	defer instance.mut.RUnlock()
	if instance.cgoInstance != nil {
		cWasmerInstanceSetGasLimit(instance.cgoInstance, gasLimit)
	}
}

// SetPointsUsed sets the internal instance gas counter
func (instance *Wasmer2Instance) SetPointsUsed(points uint64) {
	instance.mut.RLock()
	defer instance.mut.RUnlock()
	if instance.cgoInstance != nil {
		cWasmerInstanceSetPointsUsed(instance.cgoInstance, points)
	}
}

// GetPointsUsed returns the internal instance gas counter
func (instance *Wasmer2Instance) GetPointsUsed() uint64 {
	instance.mut.RLock()
	defer instance.mut.RUnlock()
	if instance.cgoInstance != nil {
		return cWasmerInstanceGetPointsUsed(instance.cgoInstance)
	}
	return 0
}

// SetBreakpointValue sets the breakpoint value for the instance
func (instance *Wasmer2Instance) SetBreakpointValue(value uint64) {
	instance.mut.Lock()
	defer instance.mut.Unlock()

	if instance.AlreadyClean || instance.cgoInstance == nil {
		return
	}
	cWasmerInstanceSetBreakpointValue(instance.cgoInstance, value)
}

// GetBreakpointValue returns the breakpoint value
func (instance *Wasmer2Instance) GetBreakpointValue() uint64 {
	instance.mut.RLock()
	defer instance.mut.RUnlock()

	if instance.AlreadyClean || instance.cgoInstance == nil {
		return 0
	}
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
	var cImportName = cCString(name)
	defer cFree(unsafe.Pointer(cImportName))

	result := cWasmerInstanceHasImportedFunction(
		instance.cgoInstance,
		cImportName,
	)

	return result == 1
}

// CallFunction executes given function from loaded contract.
func (instance *Wasmer2Instance) CallFunction(functionName string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	instance.mut.RLock()
	if instance.AlreadyClean || instance.cgoInstance == nil {
		instance.mut.RUnlock()
		return errors.New("cannot run function on closed instance")
	}
	cgoInstance := instance.cgoInstance
	instance.mut.RUnlock()

	var wasmFunctionName = cCString(functionName)
	defer cFree(unsafe.Pointer(wasmFunctionName))

	var callResult = cWasmerInstanceCall(
		cgoInstance,
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

// exportedFunctionNamesConcat reads the exported function names out of an instance, as the
// single pipe-separated string the C API produces.
//
// Package-level rather than a method, because both instance types ask for it and the two-step
// buffer protocol is worth having once.
func exportedFunctionNamesConcat(cgoInstance *cWasmerInstanceT) (string, error) {
	var bufferLength = cWasmerInstanceExportedFunctionNamesLength(cgoInstance)

	if bufferLength == 0 {
		return "", nil
	}

	var buffer = make([]cChar, bufferLength)
	var bufferPointer = (*cChar)(unsafe.Pointer(&buffer[0]))

	var result = cWasmerInstanceExportedFunctionNames(cgoInstance, bufferPointer, bufferLength)

	if result == -1 {
		return "", errors.New("cannot read function names")
	}

	if result > 0 {
		strLen := result
		if buffer[result-1] == 0 {
			strLen = result - 1
		}
		return cGoStringN(bufferPointer, strLen), nil
	}

	return "", nil
}

// exportedFunctionNames reads the names of the functions an instance exports.
func exportedFunctionNames(cgoInstance *cWasmerInstanceT) ([]string, error) {
	concat, err := exportedFunctionNamesConcat(cgoInstance)
	if err != nil {
		return nil, err
	}

	return strings.Split(concat, "|"), nil
}

// GetFunctionNames returns a list of the function names exported by the contract.
func (instance *Wasmer2Instance) GetFunctionNames() []string {
	names, err := exportedFunctionNames(instance.cgoInstance)
	if err != nil {
		return nil
	}
	return names
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

// ID returns an identifier for the instance, unique at runtime
func (instance *Wasmer2Instance) ID() string {
	return instance.id
}

// Reset resets the instance memories and globals
func (instance *Wasmer2Instance) Reset() bool {
	instance.mut.Lock()
	defer instance.mut.Unlock()

	if instance.AlreadyClean || instance.cgoInstance == nil {
		logWasmer2.Trace("reset: already cleaned instance", "id", instance.id)
		return false
	}

	result := cWasmerInstanceReset(instance.cgoInstance)
	ok := result == cWasmerOk

	logWasmer2.Trace("reset: warm instance", "id", instance.id, "ok", ok)
	return ok
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

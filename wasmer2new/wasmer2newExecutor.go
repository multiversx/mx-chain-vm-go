package wasmer2new

// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
import "C"

import (
	"fmt"
	"runtime/cgo"
	"unsafe"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
)

var _ executor.ExecutorNew = (*Wasmer2Executor)(nil)

// Wasmer2Executor creates instances of the Wasmer 2.2 backend, through the new C API.
//
// It holds the runtime behind it and nothing else that execution depends on. Each instance it
// creates gets a cgo.Handle to its own hook bridge, which the library carries back in every
// hook context - so, as in package wasmer2 and in the Rust backend, an instance supplies its
// own identity and the executor keeps no record of which one is running.
type Wasmer2Executor struct {
	cExecutor *cExecutorT

	// runtime is the host behind this executor: supplied once, never replaced.
	runtime executor.ExecutorRuntime

	// opcodeConfig is the configuration cExecutor was built with, and opcodeCost the cost
	// table it was handed. See applyOpcodeConfig for why the C executor is built lazily.
	opcodeConfig executor.OpcodeConfig
	opcodeCost   OpcodeCost
}

// NewExecutor creates an executor around the given runtime.
//
// The library's own executor is not created here. Its configuration includes the opcode costs,
// which under the new interfaces travel per instance, so there is nothing to build it from
// until the first instance arrives; see applyOpcodeConfig.
func NewExecutor(executorRuntime executor.ExecutorRuntime) (*Wasmer2Executor, error) {
	if check.IfNil(executorRuntime) {
		return nil, ErrNilExecutorRuntime
	}

	logWasmer2New.Trace("wasmer2new executor created")

	return &Wasmer2Executor{runtime: executorRuntime}, nil
}

// NewInstance creates a new executor instance.
func (wasmerExecutor *Wasmer2Executor) NewInstance(
	wasmBytes []byte,
	compilationOptions *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {
	if len(wasmBytes) == 0 {
		return nil, wasmercommon.ErrInvalidBytecode
	}

	return wasmerExecutor.newInstance(compilationOptions,
		func(instanceOut **cInstanceT, options *C.struct_vm_exec_compilation_options_t, handle uintptr, buf errBuf) wasmercommon.Status {
			return cNewInstance(wasmerExecutor.cExecutor, instanceOut, wasmBytes, options, handle, buf)
		})
}

// NewInstanceFromCache restores an executor instance from cache.
//
// cacheBytes is what InstanceNew.Cache returned, for this same contract and this same
// library. It is mapped as executable code rather than validated, so it must come from the
// VM's own compiled-code store and never from an untrusted peer; the C API carries the same
// note. An artifact produced by a different Wasmer version is rejected rather than misread.
func (wasmerExecutor *Wasmer2Executor) NewInstanceFromCache(
	cacheBytes []byte,
	compilationOptions *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {
	if len(cacheBytes) == 0 {
		return nil, wasmercommon.ErrEmptyCacheArtifact
	}

	return wasmerExecutor.newInstance(compilationOptions,
		func(instanceOut **cInstanceT, options *C.struct_vm_exec_compilation_options_t, handle uintptr, buf errBuf) wasmercommon.Status {
			return cNewInstanceFromCache(wasmerExecutor.cExecutor, instanceOut, cacheBytes, options, handle, buf)
		})
}

// The half the two creation paths share. They differ only in which C entry point turns the
// input bytes into an instance - compile, or deserialize - and everything else, including the
// hooks handle and the name scan, is the same, and has to stay the same: an instance restored
// from cache is validated exactly like a freshly compiled one.
func (wasmerExecutor *Wasmer2Executor) newInstance(
	compilationOptions *executor.CompilationOptionsNew,
	create func(**cInstanceT, *C.struct_vm_exec_compilation_options_t, uintptr, errBuf) wasmercommon.Status,
) (executor.InstanceNew, error) {
	err := wasmerExecutor.applyOpcodeConfig(compilationOptions.OpcodeConfig)
	if err != nil {
		return nil, err
	}

	cOptions := C.struct_vm_exec_compilation_options_t{
		unmetered_locals:      C.uint64_t(compilationOptions.UnmeteredLocals),
		max_memory_grow:       C.uint64_t(compilationOptions.MaxMemoryGrow),
		max_memory_grow_delta: C.uint64_t(compilationOptions.MaxMemoryGrowDelta),
		opcode_trace:          C.bool(compilationOptions.OpcodeTrace),
	}

	// A cgo.Handle rather than a pointer, because the bridge is reached from C: cgo forbids
	// storing a Go pointer in C memory, and a handle is an integer index into a Go-side
	// table, which is exactly what that rule leaves available. It is released in Close, and
	// leaks the entry if Close is never called.
	handle := cgo.NewHandle(&instanceHooks{
		runtime:     wasmerExecutor.runtime,
		stateErrBuf: newErrBuf(),
	})

	buf := newErrBuf()
	var cInstance *cInstanceT
	status := create(&cInstance, &cOptions, uintptr(handle), buf)
	if err := buf.AsError(status, "could not create instance"); err != nil {
		handle.Delete()
		return nil, err
	}

	instance := &Wasmer2Instance{
		cInstance:   cInstance,
		hooksHandle: handle,
		id:          fmt.Sprintf("%p", cInstance),
	}

	err = instance.readFunctionNames()
	if err != nil {
		instance.Close()
		return nil, err
	}

	return instance, nil
}

// readFunctionNames captures the module's exported and imported names, once.
func (instance *Wasmer2Instance) readFunctionNames() error {
	exported, err := wasmercommon.ReadNames(
		func() int32 { return cInstanceExportedFunctionNamesLength(instance.cInstance) },
		func(dest []byte) int32 { return cInstanceExportedFunctionNames(instance.cInstance, dest) },
	)
	if err != nil {
		return err
	}

	imported, err := wasmercommon.ReadNames(
		func() int32 { return cInstanceImportedFunctionNamesLength(instance.cInstance) },
		func(dest []byte) int32 { return cInstanceImportedFunctionNames(instance.cInstance, dest) },
	)
	if err != nil {
		return err
	}

	instance.exportedFunctions = exported
	instance.importedFunctions = imported

	return nil
}

// applyOpcodeConfig makes sure the library's executor is the one this instance's compilation
// options ask for, building it the first time and rebuilding it if they ever change.
//
// Under the new interfaces the opcode config travels per instance. This C API still takes it
// per executor, with no setter - deliberately, and the header says why: moving it across the
// ABI is a coordinated change with the Go VM, and until then "the executor holds the config
// it is given and merges it into the options of every instance it creates".
//
// Rebuilding rather than refusing a change is behaviour-identical to what the VM already
// does: a gas schedule change clears the warm instance cache and rebuilds the executor. The
// instances already made stay valid, the hook table having been copied into each at creation.
func (wasmerExecutor *Wasmer2Executor) applyOpcodeConfig(opcodeConfig *executor.OpcodeConfig) error {
	if opcodeConfig == nil || opcodeConfig.OpcodeCosts == nil {
		// Rejected rather than passed on, because metering is always on under this API and a
		// missing table would mean a contract running for free.
		return wasmercommon.ErrMissingOpcodeCost
	}

	if wasmerExecutor.cExecutor != nil && wasmerExecutor.opcodeConfigMatches(opcodeConfig) {
		return nil
	}

	wasmerExecutor.opcodeConfig = *opcodeConfig
	wasmerExecutor.opcodeCost = *extractOpcodeCost(opcodeConfig.OpcodeCosts)

	// The config holds a pointer to the cost table, and cgo forbids handing C a Go pointer to
	// memory that itself contains a Go pointer, so both live in C memory for the call. The
	// executor reads them during the call and retains neither.
	cCost := (*cOpcodeCostT)(C.malloc(C.size_t(unsafe.Sizeof(cOpcodeCostT{}))))
	defer C.free(unsafe.Pointer(cCost))
	*cCost = *(*cOpcodeCostT)(unsafe.Pointer(&wasmerExecutor.opcodeCost))

	cConfig := (*C.struct_vm_exec_executor_config_t)(
		C.malloc(C.size_t(unsafe.Sizeof(C.struct_vm_exec_executor_config_t{}))))
	defer C.free(unsafe.Pointer(cConfig))
	cConfig.opcode_version = C.int32_t(opcodeConfig.OpcodeVersion)
	cConfig.opcode_cost = cCost

	buf := newErrBuf()
	var cExecutor *cExecutorT
	status := cNewExecutor(&cExecutor, cConfig, populateCgoFunctionPointers(), buf)
	if err := buf.AsError(status, "could not create executor"); err != nil {
		return err
	}

	wasmerExecutor.closeCExecutor()
	wasmerExecutor.cExecutor = cExecutor

	return nil
}

// opcodeConfigMatches reports whether the library's executor was built with this exact
// configuration. Compared by value, the cost table being nothing but integers.
func (wasmerExecutor *Wasmer2Executor) opcodeConfigMatches(opcodeConfig *executor.OpcodeConfig) bool {
	return wasmerExecutor.opcodeConfig.OpcodeVersion == opcodeConfig.OpcodeVersion &&
		wasmerExecutor.opcodeCost == *extractOpcodeCost(opcodeConfig.OpcodeCosts)
}

func (wasmerExecutor *Wasmer2Executor) closeCExecutor() {
	if wasmerExecutor.cExecutor == nil {
		return
	}

	cExecutorDestroy(wasmerExecutor.cExecutor)
	wasmerExecutor.cExecutor = nil
}

// Close releases the executor.
//
// Instances outlive it safely: the hook table is copied into each one at creation, and
// nothing else is shared, so the two can be closed in either order.
func (wasmerExecutor *Wasmer2Executor) Close() {
	wasmerExecutor.closeCExecutor()
}

// IsInterfaceNil returns true if underlying object is nil
func (wasmerExecutor *Wasmer2Executor) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}

// Tripwire for an opcode added to one side of the cost table and not the other. OpcodeCost is
// generated per backend package and vm_exec_opcode_cost_t is per library, so each backend has
// to assert the match against its own header; this is wasmer2new's.
var _ [unsafe.Sizeof(cOpcodeCostT{})]byte = [unsafe.Sizeof(OpcodeCost{})]byte{}

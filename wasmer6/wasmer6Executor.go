package wasmer6

// #include <stdlib.h>
// #include "./libvmexeccapi_w6.h"
import "C"

import (
	"fmt"
	"runtime/cgo"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// ExecutorConfig is everything the executor is configured with, in one struct.
//
// The legacy API configured an executor by a sequence of setters after construction
// (SetOpcodeConfig, SetRkyvSerializationEnabled, SetSIGSEGVPassthrough). This one takes its
// configuration once, which is both why there is nothing to mutate afterwards and why there
// is no window in which an executor exists but is not yet usable.
type ExecutorConfig struct {
	OpcodeVersion  executor.OpcodeVersion
	WASMOpcodeCost *executor.WASMOpcodeCost
}

// CompilationOptions is the per-instance half of the configuration.
//
// Four fields against the legacy seven. GasLimit moved to the call, where it belongs — it
// changes on every invocation while the compiled module does not — and metering and runtime
// breakpoints are no longer optional, so there is nothing to switch off.
//
// How this is reached from executor.CompilationOptions is a question for the interface
// convergence, not for the boundary; see PLAN-go-ffi.md, Phase 4b.
type CompilationOptions struct {
	UnmeteredLocals    uint64
	MaxMemoryGrow      uint64
	MaxMemoryGrowDelta uint64
	OpcodeTrace        bool
}

// Wasmer6Executor creates instances of the Wasmer 6 backend.
type Wasmer6Executor struct {
	cExecutor *cExecutorT

	// Kept so that the legacy Executor interface, whose NewInstanceWithOptions takes neither,
	// can still create instances; and so that SetOpcodeConfig can rebuild the C executor
	// rather than mutate a configuration this API deliberately makes immutable.
	config   ExecutorConfig
	provider executor.VMHooksProvider
}

// NewExecutor creates an executor from one configuration struct and the generated hook
// table.
func NewExecutor(config ExecutorConfig) (*Wasmer6Executor, error) {
	if config.WASMOpcodeCost == nil {
		return nil, ErrMissingOpcodeCost
	}

	buf := newErrBuf()

	// Both the config and the cost table live in C memory for the duration of the call.
	// They could not be Go-allocated: the config holds a pointer to the cost table, and
	// cgo forbids handing C a Go pointer to memory that itself contains a Go pointer.
	// The executor reads both during the call and retains neither, so they are freed here.
	cCost := (*cOpcodeCostT)(C.malloc(C.size_t(unsafe.Sizeof(cOpcodeCostT{}))))
	defer C.free(unsafe.Pointer(cCost))
	*cCost = *(*cOpcodeCostT)(unsafe.Pointer(executor.ExtractOpcodeCost(config.WASMOpcodeCost)))

	cConfig := (*C.struct_vm_exec_executor_config_t)(
		C.malloc(C.size_t(unsafe.Sizeof(C.struct_vm_exec_executor_config_t{}))))
	defer C.free(unsafe.Pointer(cConfig))
	cConfig.opcode_version = C.int32_t(config.OpcodeVersion)
	cConfig.opcode_cost = cCost

	var cExecutor *cExecutorT
	status := cNewExecutor(&cExecutor, cConfig, populateCgoFunctionPointers(), buf)
	if err := buf.asError(status, "could not create executor"); err != nil {
		return nil, err
	}

	logWasmer6.Trace("wasmer6 executor created")

	return &Wasmer6Executor{cExecutor: cExecutor, config: config}, nil
}

// NewInstance compiles and instantiates a contract.
//
// hooks is the VM hooks provider for this one instance. It is reached from every hook call
// the instance makes, and must stay usable until the instance is closed.
func (wasmerExecutor *Wasmer6Executor) NewInstance(
	contractCode []byte,
	options CompilationOptions,
	hooks executor.VMHooksProvider,
) (*Wasmer6Instance, error) {
	if len(contractCode) == 0 {
		return nil, ErrInvalidBytecode
	}

	cOptions := C.struct_vm_exec_compilation_options_t{
		unmetered_locals:      C.uint64_t(options.UnmeteredLocals),
		max_memory_grow:       C.uint64_t(options.MaxMemoryGrow),
		max_memory_grow_delta: C.uint64_t(options.MaxMemoryGrowDelta),
		opcode_trace:          C.bool(options.OpcodeTrace),
	}

	// A cgo.Handle rather than a pointer, because the hooks object is reached from C: cgo
	// forbids storing a Go pointer in C memory, and a handle is an integer index into a Go
	// side table, which is exactly what that rule leaves available. It is released in
	// Close, and leaks the entry if Close is never called.
	bridge := &instanceHooks{
		provider:    hooks,
		stateErrBuf: newErrBuf(),
	}
	handle := cgo.NewHandle(bridge)

	buf := newErrBuf()
	var cInstance *cInstanceT
	status := cNewInstance(
		wasmerExecutor.cExecutor,
		&cInstance,
		contractCode,
		&cOptions,
		uintptr(handle),
		buf,
	)
	if err := buf.asError(status, "could not create instance"); err != nil {
		handle.Delete()
		return nil, err
	}

	instance := &Wasmer6Instance{
		cInstance:   cInstance,
		hooksHandle: handle,
		hooks:       bridge,
		id:          fmt.Sprintf("%p", cInstance),
	}

	imported, err := instance.ImportedFunctionNames()
	if err != nil {
		instance.Close()
		return nil, err
	}
	instance.importedFunctions = make(map[string]struct{}, len(imported))
	for _, name := range imported {
		instance.importedFunctions[name] = struct{}{}
	}

	return instance, nil
}

// Close releases the executor.
//
// Instances outlive it safely: the hook table is copied into each one at creation, and
// nothing else is shared, so the two can be closed in either order.
func (wasmerExecutor *Wasmer6Executor) Close() {
	if wasmerExecutor.cExecutor == nil {
		return
	}

	cExecutorDestroy(wasmerExecutor.cExecutor)
	wasmerExecutor.cExecutor = nil
}

// Tripwire for an opcode added to one side of the cost table and not the other.
// executor.OpcodeCost is shared by every backend and vm_exec_opcode_cost_t is per-library, so
// each backend has to assert the match against its own header; this is wasmer6's.
var _ [unsafe.Sizeof(cOpcodeCostT{})]byte = [unsafe.Sizeof(executor.OpcodeCost{})]byte{}

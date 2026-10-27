package wasmer2

import (
	"runtime"
	"unsafe"

	"github.com/multiversx/mx-chain-core-go/core/check"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.ExecutorNew = (*Wasmer2ExecutorNew)(nil)
var _ executor.VMHooksLegacy = newInstanceVMHooks(nil, nil)

// Wasmer2ExecutorNew creates Wasmer 2.2 instances through the new interfaces.
//
// Counterpart of the Rust WasmerProdExecutor: the same backend, reached through the new
// boundary instead of the old one. It holds the runtime behind it and nothing else that
// execution depends on - the opcode costs and the memory limits arrive per instance, in the
// compilation options.
//
// As in the Rust version, each instance carries its own hooks, built around its own state.
// What the C API imposes is only where they are reached from: vm_exec_executor_set_vm_hooks_ptr
// is per executor, so the hooks of the instance that is running are installed in one slot
// here for the duration of its call. The executor keeps no record of which instances exist
// or which is running - that is the instance's own business, and the nesting of calls is the
// Go call stack.
type Wasmer2ExecutorNew struct {
	cgoExecutor *cWasmerExecutorT

	vmHookPointers *cWasmerVmHookPointers

	// runtime is the host behind this executor: supplied once, never replaced.
	runtime executor.ExecutorRuntime

	// vmHooks is the slot the C side reads on every hook, through vmHooksPtr. It holds the
	// hooks of the instance currently running. Its address is what crosses the boundary, so
	// the field must outlive every instance.
	vmHooks       executor.VMHooksLegacy
	vmHooksPtr    uintptr
	vmHooksPtrPtr unsafe.Pointer

	// opcodeCost is the cost table the C executor was last configured with. Kept alive here
	// because the C side is handed its address.
	opcodeCost *OpcodeCost
}

// CreateExecutorNew creates a new Wasmer 2.2 executor around the given runtime.
func CreateExecutorNew(executorRuntime executor.ExecutorRuntime) (*Wasmer2ExecutorNew, error) {
	if check.IfNil(executorRuntime) {
		return nil, ErrNilExecutorRuntime
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	vmHookPointers := populateCgoFunctionPointers()
	localPtr := uintptr(unsafe.Pointer(vmHookPointers))
	localPtrPtr := unsafe.Pointer(&localPtr)

	var cExecutor *cWasmerExecutorT

	if cWasmerNewExecutor(&cExecutor, localPtrPtr) != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	cWasmerForceInstallSighandlers()

	wasmerExecutor := &Wasmer2ExecutorNew{
		cgoExecutor:    cExecutor,
		vmHookPointers: vmHookPointers,
		runtime:        executorRuntime,
	}
	wasmerExecutor.initVMHooks()

	return wasmerExecutor, nil
}

// initVMHooks gives the C side the address of the slot it will read on every hook.
//
// The hooks are legacy-shaped because the functions this backend imports are plain C
// functions: they return a value and have nowhere to report that the hook wants execution to
// stop. The adapter applies that as a side effect instead.
//
// What goes into the slot to start with is a hooks object belonging to no instance. It
// cannot normally be reached, a hook only firing from inside a call that has installed its
// own, but the slot has to hold something, and a detached state reports itself as detached
// rather than dereferencing nothing.
func (wasmerExecutor *Wasmer2ExecutorNew) initVMHooks() {
	wasmerExecutor.vmHooks = newInstanceVMHooks(wasmerExecutor.runtime, nil)
	localPtr := uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
	wasmerExecutor.vmHooksPtr = localPtr
	wasmerExecutor.vmHooksPtrPtr = unsafe.Pointer(&localPtr)
	cWasmerExecutorContextDataSet(wasmerExecutor.cgoExecutor, wasmerExecutor.vmHooksPtrPtr)
}

// useVMHooks installs the hooks the C side will reach from now on, returning the ones that
// were there before, for the caller to put back.
func (wasmerExecutor *Wasmer2ExecutorNew) useVMHooks(vmHooks executor.VMHooksLegacy) executor.VMHooksLegacy {
	previous := wasmerExecutor.vmHooks
	wasmerExecutor.vmHooks = vmHooks

	return previous
}

// NewInstance creates a new executor instance.
func (wasmerExecutor *Wasmer2ExecutorNew) NewInstance(
	wasmBytes []byte,
	compilationOptions *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {
	if len(wasmBytes) == 0 {
		return nil, newWrappedError(ErrInvalidBytecode)
	}

	return wasmerExecutor.newInstance(compilationOptions,
		func(cInstance **cWasmerInstanceT, cOptions *cWasmerCompilationOptions) cWasmerResultT {
			return cWasmerInstantiateWithOptions(
				wasmerExecutor.cgoExecutor,
				cInstance,
				(*cUchar)(unsafe.Pointer(&wasmBytes[0])),
				cUint(len(wasmBytes)),
				cOptions,
			)
		})
}

// NewInstanceFromCache restores an executor instance from cache.
//
// cacheBytes is what InstanceNew.Cache returned, for this same contract and this same
// library. It is mapped as executable code rather than validated, so it must come from the
// VM's own compiled-code store and never from an untrusted peer.
func (wasmerExecutor *Wasmer2ExecutorNew) NewInstanceFromCache(
	cacheBytes []byte,
	compilationOptions *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {
	if len(cacheBytes) == 0 {
		return nil, newWrappedError(ErrEmptyCacheArtifact)
	}

	return wasmerExecutor.newInstance(compilationOptions,
		func(cInstance **cWasmerInstanceT, cOptions *cWasmerCompilationOptions) cWasmerResultT {
			return cWasmerInstanceFromCache(
				wasmerExecutor.cgoExecutor,
				cInstance,
				(*cUchar)(unsafe.Pointer(&cacheBytes[0])),
				cUint32T(len(cacheBytes)),
				cOptions,
			)
		})
}

// The half the two creation paths share. They differ only in which C entry point turns the
// input bytes into an instance - compile, or deserialize - and everything around that is the
// same, as it has to be: an instance restored from cache is configured exactly like a freshly
// compiled one.
func (wasmerExecutor *Wasmer2ExecutorNew) newInstance(
	compilationOptions *executor.CompilationOptionsNew,
	create func(**cWasmerInstanceT, *cWasmerCompilationOptions) cWasmerResultT,
) (executor.InstanceNew, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := wasmerExecutor.applyOpcodeConfig(compilationOptions.OpcodeConfig)
	if err != nil {
		return nil, err
	}

	cOptions := toCWasmerCompilationOptions(compilationOptions.ToLegacy())

	var cInstance *cWasmerInstanceT
	if create(&cInstance, &cOptions) != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	return newInstanceNew(cInstance, wasmerExecutor), nil
}

// applyOpcodeConfig hands the metering middleware the costs this instance is to be compiled
// with.
//
// Under the new interfaces the opcode config travels per instance, but the 2.2 C API still
// keeps it on the executor, behind a setter. Applied here, immediately before the module is
// compiled and on the same locked thread, so what the middleware reads is this instance's
// configuration. Nothing else writes it, which makes it a private copy in all but name - the
// same arrangement as the Rust backend's private_opcode_config.
func (wasmerExecutor *Wasmer2ExecutorNew) applyOpcodeConfig(opcodeConfig *executor.OpcodeConfig) error {
	if opcodeConfig == nil || opcodeConfig.OpcodeCosts == nil {
		// Tolerated, as it was under the old interface: tests that only compile and
		// instantiate a module never meter anything.
		return nil
	}

	wasmerExecutor.opcodeCost = extractOpcodeCost(opcodeConfig.OpcodeCosts)

	result := cWasmerExecutorSetOpcodeConfig(
		wasmerExecutor.cgoExecutor,
		(int32)(opcodeConfig.OpcodeVersion),
		(*cWasmerOpcodeCostT)(unsafe.Pointer(wasmerExecutor.opcodeCost)),
	)
	if result != cWasmerOk {
		return newWrappedError(ErrFailedInstantiation)
	}

	return nil
}

// FunctionNames returns the low-level function names provided to contracts.
//
// Not part of executor.ExecutorNew - under the new interfaces the set of hook names is the
// host's own business - but kept available, since this backend is what defines it.
func (wasmerExecutor *Wasmer2ExecutorNew) FunctionNames() vmcommon.FunctionNames {
	return functionNames
}

// IsInterfaceNil returns true if underlying object is nil
func (wasmerExecutor *Wasmer2ExecutorNew) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}

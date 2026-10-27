package wasmer2new

// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// getVMHooksFromHookContext builds the VM hooks that service this one host call, and finds
// the handler for an error they report.
//
// Field 0 of the hook context is the handle the instance was created with. The executor
// treats it as an opaque word and only hands it back; what it identifies is the instance's
// runtime. An integer handle rather than a pointer because cgo forbids storing a Go pointer
// in C memory, and because a handle stays valid under a moving collector.
//
// The hooks are built per call, not per instance, so that the instance state they reach
// through is the one belonging to this call.
func getVMHooksFromHookContext(ctx *C.vm_exec_hook_context_t) (executor.VMHooks, executor.VMHooksErrorHandler) {
	hooks := cgo.Handle(ctx.go_vm_hooks).Value().(*instanceHooks)
	state := &Wasmer2InstanceState{ctx: ctx, owner: hooks}

	vmHooks := hooks.runtime.VMHooksForCall(state)

	return vmHooks, errorHandlerOf(vmHooks)
}

// errorHandlerOf finds the error handler belonging to a hooks object.
//
// A type assertion rather than a second thing the runtime supplies, because the two roles are
// played by one object in the VM and executor.ExecutorRuntime hands over only the hooks. The
// same assertion the old C API's factory makes, moved to where the hooks are produced; see
// executor.VMHooksErrorHandler.
func errorHandlerOf(vmHooks executor.VMHooks) executor.VMHooksErrorHandler {
	if errorHandler, ok := vmHooks.(executor.VMHooksErrorHandler); ok {
		return errorHandler
	}

	return discardHookErrors{}
}

var _ executor.VMHooksErrorHandler = discardHookErrors{}

// discardHookErrors stands in for hooks that cannot absorb an error of their own.
//
// The early exit still travels to the executor and still stops the contract; what is lost is
// the VM-side record of why, which is the hooks object's own business. Logged rather than
// silent, because it means the VMOutput will carry no return message.
type discardHookErrors struct{}

// HandleHookError implements executor.VMHooksErrorHandler.
func (discardHookErrors) HandleHookError(err error) {
	logWasmer2New.Error("VM hooks cannot record a hook error", "err", err)
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}

// hookOk reports a hook that returned normally, with its own result value.
//
// The sentinel is what keeps "no early exit" distinct from early-exit code 0, which is a real
// code reserved for the legacy async call.
func hookOk(value int64) C.vm_exec_hook_result_t {
	return C.vm_exec_hook_result_t{
		value:      C.int64_t(value),
		early_exit: C.VM_EXEC_NO_EARLY_EXIT,
	}
}

// hookEarlyExit applies the error's side effects on the VM side and reports its code.
//
// Two steps rather than one, because they are two different concerns: HandleHookError records
// the return code and message in the VM's output state, while the code is what the executor
// needs in order to stop the instance. Under this executor the breakpoint that
// HandleHookError still sets is redundant - the early exit itself stops execution.
func hookEarlyExit(errorHandler executor.VMHooksErrorHandler, err error) C.vm_exec_hook_result_t {
	errorHandler.HandleHookError(err)

	return C.vm_exec_hook_result_t{
		value:      0,
		early_exit: C.uint64_t(executor.EarlyExitCodeOf(err)),
	}
}

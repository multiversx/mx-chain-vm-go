package wasmer3

// #include <stdlib.h>
// #include "./libvmexeccapi_new.h"
import "C"

import (
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// getVMHooksFromHookContext recovers the Go VM hooks object from field 0 of the hook context.
//
// The executor treats that field as an opaque word and only hands it back; what it holds is
// a pointer to an executor.VMHooks interface value, supplied when the instance was created.
// There is one such object per live instance.
func getVMHooksFromHookContext(ctx *C.vm_exec2_hook_context_t) executor.VMHooks {
	return *(*executor.VMHooks)(unsafe.Pointer(ctx.go_vm_hooks))
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}

// hookOk reports a hook that returned normally, with its own result value.
//
// The sentinel is what keeps "no early exit" distinct from early-exit code 0, which is a real
// code reserved for the legacy async call.
func hookOk(value int64) C.vm_exec2_hook_result_t {
	return C.vm_exec2_hook_result_t{
		value:      C.int64_t(value),
		early_exit: C.VM_EXEC2_NO_EARLY_EXIT,
	}
}

// hookEarlyExit applies the error's side effects on the VM side and reports its code.
//
// Two steps rather than one, because they are two different concerns: HandleHookError records
// the return code and message in the VM's output state, while the code is what the executor
// needs in order to stop the instance. Under this executor the breakpoint that
// HandleHookError still sets is redundant — the early exit itself stops execution — and it
// goes away with the breakpoint removal; see PLAN-go-breakpoint-removal.md.
func hookEarlyExit(vmHooks executor.VMHooks, err error) C.vm_exec2_hook_result_t {
	vmHooks.HandleHookError(err)

	return C.vm_exec2_hook_result_t{
		value:      0,
		early_exit: C.uint64_t(executor.EarlyExitCodeOf(err)),
	}
}

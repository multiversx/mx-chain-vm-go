package wasmer6

// #include <stdlib.h>
// #include "./libvmexeccapi_w6.h"
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// getVMHooksFromHookContext builds the VM hooks that service this one host call.
//
// Field 0 of the hook context is the handle the instance was created with. The executor
// treats it as an opaque word and only hands it back; what it identifies is the instance's
// VM hooks provider. An integer handle rather than a pointer because cgo forbids storing a
// Go pointer in C memory, and because a handle stays valid under a moving collector.
//
// The hooks are built per call, not per instance, so that the instance state they reach
// through is the one belonging to this call. See VMHooksProvider for why that matters.
func getVMHooksFromHookContext(ctx *C.vm_exec_w6_hook_context_t) executor.VMHooks {
	hooks := cgo.Handle(ctx.go_vm_hooks).Value().(*instanceHooks)
	state := &InstanceState{ctx: ctx, owner: hooks}

	// Published for the VM's instance-level memory and gas calls, which still arrive through
	// executor.Instance rather than through the state below. See instanceHooks.activeState.
	hooks.activeState = state

	return hooks.provider.VMHooksForCall(state)
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}

// hookOk reports a hook that returned normally, with its own result value.
//
// The sentinel is what keeps "no early exit" distinct from early-exit code 0, which is a real
// code reserved for the legacy async call.
func hookOk(value int64) C.vm_exec_w6_hook_result_t {
	return C.vm_exec_w6_hook_result_t{
		value:      C.int64_t(value),
		early_exit: C.VM_EXEC_W6_NO_EARLY_EXIT,
	}
}

// hookEarlyExit applies the error's side effects on the VM side and reports its code.
//
// Two steps rather than one, because they are two different concerns: HandleHookError records
// the return code and message in the VM's output state, while the code is what the executor
// needs in order to stop the instance. Under this executor the breakpoint that
// HandleHookError still sets is redundant — the early exit itself stops execution — and it
// goes away with the breakpoint removal; see PLAN-go-breakpoint-removal.md.
func hookEarlyExit(vmHooks executor.VMHooks, err error) C.vm_exec_w6_hook_result_t {
	vmHooks.HandleHookError(err)

	return C.vm_exec_w6_hook_result_t{
		value:      0,
		early_exit: C.uint64_t(executor.EarlyExitCodeOf(err)),
	}
}

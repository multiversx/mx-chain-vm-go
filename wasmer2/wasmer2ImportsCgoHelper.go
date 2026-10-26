package wasmer2

import (
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// getVMHooksFromContextRawPtr recovers the VM hooks from the executor's context data.
//
// The interface type here must match the one initVMHooks stored, since the cast reinterprets
// the interface value byte for byte.
func getVMHooksFromContextRawPtr(contextRawPtr unsafe.Pointer) executor.VMHooksLegacy {
	vmHooksPtrPtr := (*uintptr)(contextRawPtr)
	vmHooksPtr := *vmHooksPtrPtr
	return *(*executor.VMHooksLegacy)(unsafe.Pointer(vmHooksPtr))
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}

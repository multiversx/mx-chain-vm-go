package wasmer2

import (
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

func getVMHooksFromContextRawPtr(contextPtr unsafe.Pointer) executor.VMHooks {
	// instCtx := IntoInstanceContext(contextPtr)
	// var ptr = *(*uintptr)(instCtx.Data())
	// return *(*executor.VMHooks)(unsafe.Pointer(ptr))
	panic("getVMHooksFromContextRawPtr not yet implemented")
}

func funcPointer(f interface{}) *[0]byte {
	return (*[0]byte)(unsafe.Pointer(&f))
}

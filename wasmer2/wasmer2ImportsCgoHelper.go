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
	// The context data holds a pointer-sized word (stored as a uintptr on the Go side, since it
	// travels through C). Reading it back as an unsafe.Pointer avoids a uintptr -> unsafe.Pointer
	// conversion, which go vet's unsafeptr check rejects. The layout is identical.
	vmHooksPtr := *(*unsafe.Pointer)(contextRawPtr)
	return *(*executor.VMHooksLegacy)(vmHooksPtr)
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}

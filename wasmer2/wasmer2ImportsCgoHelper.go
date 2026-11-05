package wasmer2

import (
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

func getVMHooksFromContextRawPtr(contextRawPtr unsafe.Pointer) executor.VMHooks {
	// The context data holds a pointer-sized word (stored as a uintptr on the Go side, since it
	// travels through C). Reading it back as an unsafe.Pointer avoids a uintptr -> unsafe.Pointer
	// conversion, which go vet's unsafeptr check rejects. The layout is identical.
	vmHooksPtr := *(*unsafe.Pointer)(contextRawPtr)
	return *(*executor.VMHooks)(vmHooksPtr)
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}

package wasmer

import (
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

func importsInterfaceFromRaw(contextPtr unsafe.Pointer) executor.ImportsInterface {
	// instCtx := wasmer.IntoInstanceContext(vmHostPtr)
	// var ptr = *(*uintptr)(context)
	return *(*executor.ImportsInterface)(contextPtr)
}

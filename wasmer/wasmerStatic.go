package wasmer

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// SetRkyvSerializationEnabled enables or disables RKYV serialization of
// instances in Wasmer
func SetRkyvSerializationEnabled(enabled bool) {
	if enabled {
		cWasmerInstanceEnableRkyv()
	} else {
		cWasmerInstanceDisableRkyv()
	}
}

// SetSIGSEGVPassthrough instructs Wasmer to never register a handler for
// SIGSEGV. Only has effect if called before creating the first Wasmer instance
// since the process started. Calling this function after the first Wasmer
// instance will not unregister the signal handler set by Wasmer.
func SetSIGSEGVPassthrough() {
	cWasmerSetSIGSEGVPassthrough()
}

// SetOpcodeCosts sets gas costs globally for Wasmer.
func SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	cWasmerSetOpcodeCosts(opcodeCosts)
}

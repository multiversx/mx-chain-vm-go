package wasmer

import "github.com/multiversx/mx-chain-vm-go/executor"

// SetRkyvSerializationEnabled enables or disables RKYV serialization of instances in Wasmer.
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

// ForceInstallSighandlers triggers a forced installation of signal handlers in Wasmer 1
func ForceInstallSighandlers() {
	cWasmerForceInstallSighandlers()
}

// SetOpcodeCosts sets gas costs globally for Wasmer.
func SetOpcodeCosts(opcodeCosts *executor.WASMOpcodeCost) {
	opcodeCostsArray := toOpcodeCostsArray(opcodeCosts)
	cWasmerSetOpcodeCosts(&opcodeCostsArray)
}

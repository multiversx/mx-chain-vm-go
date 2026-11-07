package vmhooks

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// NewVMHooksLegacy creates VM hooks that absorb their own errors, exactly as the executor
// does when it runs a contract.
//
// For mock contracts and tests. They stand in for a contract, and a contract never sees the
// error channel: a hook that fails stops the execution, it does not hand back an error to
// inspect. Going through VMHooksImpl directly would make every call site deal with an error
// that a real contract would never have been given.
func NewVMHooksLegacy(host vmhost.VMHost) *executor.VMHooksLegacyAdapter {
	vmHooksImpl := NewVMHooksImpl(host)
	return executor.NewVMHooksLegacyAdapter(vmHooksImpl, vmHooksImpl)
}

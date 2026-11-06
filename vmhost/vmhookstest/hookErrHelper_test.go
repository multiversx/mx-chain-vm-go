package vmhookstest

import (
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// handleHookError applies the error returned by one of the package-level hook helpers the way
// the legacy adapter does for a hook called through the interface, so that tests observe the
// same failure side effects either way.
func handleHookError(host vmhost.VMHost, err error) {
	// Both halves, as a real backend does: the VM records why in its output state, and the
	// instance records the exit so that the call reports it.
	vmhooks.NewVMHooksImpl(host, vmhooks.InstanceStateOf(host)).HandleHookError(err)
	contextmock.GetMockInstance(host).HandleHookError(err)
}

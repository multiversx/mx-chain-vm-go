package vmhookstest

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// handleHookError applies the error returned by one of the package-level hook helpers the way
// the legacy adapter does for a hook called through the interface, so that tests observe the
// same failure side effects either way.
func handleHookError(host vmhost.VMHost, err error) {
	vmhooks.NewVMHooksImpl(host).HandleHookError(err)
}

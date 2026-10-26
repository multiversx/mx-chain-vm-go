package hostCoretest

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// handleHookErr applies the error returned by one of the package-level hook helpers the way
// the legacy adapter does for a hook called through the interface, so that tests observe the
// same failure side effects either way.
func handleHookErr(host vmhost.VMHost) func(int32, error) (int32, error) {
	return func(value int32, err error) (int32, error) {
		vmhooks.NewVMHooksImpl(host).HandleHookError(err)
		return value, err
	}
}

// handleHookError is the same, for a helper that only returns an error.
func handleHookError(host vmhost.VMHost, err error) {
	vmhooks.NewVMHooksImpl(host).HandleHookError(err)
}

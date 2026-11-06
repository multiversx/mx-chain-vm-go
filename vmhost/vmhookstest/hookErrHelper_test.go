package vmhookstest

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// handleHookErr applies the error returned by a directly-invoked VM hook the way the
// generated wrappers do, so that tests observe the same failure side effects.
func handleHookErr(host vmhost.VMHost) func(int32, error) (int32, error) {
	return func(value int32, err error) (int32, error) {
		vmhooks.NewVMHooksImpl(host).HandleHookError(err)
		return value, err
	}
}

func handleHookErr64(host vmhost.VMHost) func(int64, error) (int64, error) {
	return func(value int64, err error) (int64, error) {
		vmhooks.NewVMHooksImpl(host).HandleHookError(err)
		return value, err
	}
}

// handleHookError applies the error returned by a directly-invoked, error-only VM hook.
func handleHookError(host vmhost.VMHost, err error) {
	vmhooks.NewVMHooksImpl(host).HandleHookError(err)
}

package testcommon

import (
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// FailMockContract ends a mock contract the way a failed hook ends a real one.
//
// Two halves, as every backend does: the VM records in its output state why the contract
// stopped, and the instance records the exit so that its call reports it. A mock contract is
// Go code with no hook to return an error through, so it has to do both itself.
//
// This replaces the bare runtime.FailExecution that mock contracts used to call. That worked
// because FailExecution also set a breakpoint, which both stopped the contract and said how:
// in particular it recognised ErrNotEnoughGas and chose the out-of-gas breakpoint, which is
// why a mock contract that ran out of gas produced ReturnCode OutOfGas rather than
// ExecutionFailed. The classification now travels on the error, and reaches the instance
// through here.
func FailMockContract(host vmhost.VMHost, err error) {
	if err == nil {
		return
	}

	instance := mock.GetMockInstance(host)

	vmhooks.NewVMHooksImpl(host, instance).HandleHookError(err)
	instance.HandleHookError(err)
}

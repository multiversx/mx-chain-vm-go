package vmhooks

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// HandleHookError applies the side effects of an error returned by a VM hook.
// Called only from generated hook wrappers.
func (context *VMHooksImpl) HandleHookError(err error) {
	if err == nil {
		return
	}
	applyEarlyExit(context.host, executor.EarlyExitOf(err))
}

// applyEarlyExit records what the termination was in the VMOutput.
//
// It no longer performs the termination. That is the executor's, driven by the error the
// hook returns, which travels back out through the generated adapter. What is left here is
// the VMOutput bookkeeping each kind of exit needs.
func applyEarlyExit(host vmhost.VMHost, earlyExit *executor.EarlyExit) {
	if earlyExit.ConsumeAllGas {
		metering := host.Metering()
		_ = metering.UseGasBounded(metering.GasLeft())
	}

	runtime := host.Runtime()
	switch earlyExit.Code {
	case vmcommon.Ok:
		// A legacy async call: the contract suspended itself rather than failed, so there
		// is nothing to record. The host picks the registered call up once the call ends.
	case vmcommon.UserError:
		runtime.SignalUserError(earlyExit.Message)
	case vmcommon.OutOfGas:
		// Nothing to record: out of gas has no VMOutput state of its own, and the stopping
		// is the executor's, from the error this hook is about to return. This used to set
		// a breakpoint, which was how the contract was halted before hooks could report
		// their own exit.
	default:
		runtime.FailExecution(earlyExit.Cause)
	}
}

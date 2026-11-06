package vmhooks

import (
	"errors"

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
	applyEarlyExit(context.host, resolveEarlyExit(err))
}

// resolveEarlyExit classifies a hook error into the termination it requires.
// See the contract in the migration document: exactly one call site, never nil in, never nil out.
func resolveEarlyExit(err error) *executor.EarlyExit {
	var earlyExit *executor.EarlyExit
	if errors.As(err, &earlyExit) {
		return earlyExit
	}

	// The default: fail execution and burn the remaining gas. This is what the
	// 800+ plain-error failure sites used to do via FailExecution. The policy itself lives
	// in the executor package, so that the code an executor reports and the side effects
	// applied here cannot drift apart.
	return executor.ExitExecutionFailed(err)
}

// applyEarlyExit performs the termination. Each branch is exactly one of the three
// effect routines that exist today; no new behaviour is introduced here.
func applyEarlyExit(host vmhost.VMHost, earlyExit *executor.EarlyExit) {
	if earlyExit.ConsumeAllGas {
		metering := host.Metering()
		_ = metering.UseGasBounded(metering.GasLeft())
	}

	runtime := host.Runtime()
	switch earlyExit.Code {
	case vmcommon.UserError:
		runtime.SignalUserError(earlyExit.Message)
	case vmcommon.OutOfGas:
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
	default:
		runtime.FailExecution(earlyExit.Cause)
	}
}

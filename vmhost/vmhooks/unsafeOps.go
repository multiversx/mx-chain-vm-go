package vmhooks

const (
	activateUnsafeModeName   = "activateUnsafeMode"
	deactivateUnsafeModeName = "deactivateUnsafeMode"
)

func (context *VMHooksImpl) useGasForUnsafeActivation(traceString string) error {
	metering := context.GetMeteringContext()
	gasToUse := metering.GasSchedule().BaseOpsAPICost.Finish
	err := metering.UseGasBoundedAndAddTracedGas(traceString, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return err
	}

	return nil
}

// ActivateUnsafeMode VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ActivateUnsafeMode() {
	if err := context.useGasForUnsafeActivation(activateUnsafeModeName); err != nil {
		return
	}

	context.GetRuntimeContext().SetUnsafeMode(true)
}

// DeactivateUnsafeMode VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeactivateUnsafeMode() {
	if err := context.useGasForUnsafeActivation(deactivateUnsafeModeName); err != nil {
		return
	}

	context.GetRuntimeContext().SetUnsafeMode(false)
}

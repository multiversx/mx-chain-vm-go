package vmhooks

const (
	activateUnsafeModeName   = "activateUnsafeMode"
	deactivateUnsafeModeName = "deactivateUnsafeMode"
)

// ActivateUnsafeMode VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ActivateUnsafeMode() {
	metering := context.GetMeteringContext()
	err := metering.UseGasBoundedAndAddTracedGas(activateUnsafeModeName, 1)
	if err != nil {
		context.FailExecution(err)
		return
	}

	context.host.SetUnsafeMode(true)
}

// DeactivateUnsafeMode VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeactivateUnsafeMode() {
	metering := context.GetMeteringContext()
	err := metering.UseGasBoundedAndAddTracedGas(deactivateUnsafeModeName, 1)
	if err != nil {
		context.FailExecution(err)
		return
	}

	context.host.SetUnsafeMode(false)
}

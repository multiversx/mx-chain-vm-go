package vmhooks

import "github.com/multiversx/mx-chain-vm-go/vmhost"

const (
	activateUnsafeModeName       = "activateUnsafeMode"
	deactivateUnsafeModeName     = "deactivateUnsafeMode"
	managedGetNumErrorsName      = "managedGetNumErrors"
	managedGetErrorWithIndexName = "managedGetErrorWithIndex"
	managedGetLastErrorName      = "managedGetLastError"
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

// ManagedGetNumErrors VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetNumErrors() int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(managedGetNumErrorsName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return -1
	}

	allErrorsWrapper := runtime.GetAllErrors()
	if allErrorsWrapper == nil {
		return 0
	}

	wrappableErr, ok := allErrorsWrapper.(vmhost.WrappableError)
	if !ok {
		context.FailExecution(vmhost.ErrWrongType)
		return 1
	}

	return int32(len(wrappableErr.GetAllErrors()))
}

// ManagedGetErrorWithIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetErrorWithIndex(index int32, errorHandle int32) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(managedGetErrorWithIndexName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	allErrorsWrapper := runtime.GetAllErrors()
	if allErrorsWrapper == nil {
		context.FailExecutionConditionally(vmhost.ErrInvalidArgument)
		return
	}

	wrappableErr, ok := allErrorsWrapper.(vmhost.WrappableError)
	if !ok {
		context.FailExecution(vmhost.ErrWrongType) // Should not happen
		return
	}

	allErrors := wrappableErr.GetAllErrors()
	if index < 0 || int(index) >= len(allErrors) {
		context.FailExecutionConditionally(vmhost.ErrInvalidArgument)
		return
	}

	theError := allErrors[index]
	errorMessage := []byte(theError.Error())

	err = managedType.ConsumeGasForBytes(errorMessage)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(errorHandle, errorMessage)
}

// ManagedGetLastError VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetLastError(errorHandle int32) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(managedGetLastErrorName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	allErrorsWrapper := runtime.GetAllErrors()
	if allErrorsWrapper == nil {
		managedType.SetBytes(errorHandle, []byte{})
		return
	}

	wrappableErr, ok := allErrorsWrapper.(vmhost.WrappableError)
	if !ok {
		context.FailExecution(vmhost.ErrWrongType)
		return
	}

	lastError := wrappableErr.GetLastError()
	if lastError == nil {
		managedType.SetBytes(errorHandle, []byte{})
		return
	}

	errorMessage := []byte(lastError.Error())

	err = managedType.ConsumeGasForBytes(errorMessage)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(errorHandle, errorMessage)
}

package vmhooks

const (
	managedMapNewName      = "managedMapNew"
	managedMapPutName      = "managedMapPut"
	managedMapGetName      = "managedMapGet"
	managedMapRemoveName   = "managedMapRemove"
	managedMapContainsName = "managedMapContains"
)

// ManagedMapNew VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMapNew() int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedMapAPICost.ManagedMapNew
	err := metering.UseGasBoundedAndAddTracedGas(managedMapNewName, gasToUse)
	if context.WithFault(err, context.GetRuntimeContext().ManagedMapAPIErrorShouldFailExecution()) {
		return 1
	}

	return managedType.NewManagedMap()
}

// ManagedMapPut VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMapPut(mMapHandle int32, keyHandle int32, valueHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().ManagedMapAPICost.ManagedMapPut
	err := metering.UseGasBoundedAndAddTracedGas(managedMapPutName, gasToUse)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 1
	}

	err = managedType.ManagedMapPut(mMapHandle, keyHandle, valueHandle)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// ManagedMapGet VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMapGet(mMapHandle int32, keyHandle int32, outValueHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().ManagedMapAPICost.ManagedMapGet
	err := metering.UseGasBoundedAndAddTracedGas(managedMapGetName, gasToUse)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 1
	}

	err = managedType.ManagedMapGet(mMapHandle, keyHandle, outValueHandle)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// ManagedMapRemove VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMapRemove(mMapHandle int32, keyHandle int32, outValueHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().ManagedMapAPICost.ManagedMapRemove
	err := metering.UseGasBoundedAndAddTracedGas(managedMapRemoveName, gasToUse)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 1
	}

	err = managedType.ManagedMapRemove(mMapHandle, keyHandle, outValueHandle)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// ManagedMapContains VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMapContains(mMapHandle int32, keyHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().ManagedMapAPICost.ManagedMapContains
	err := metering.UseGasBoundedAndAddTracedGas(managedMapContainsName, gasToUse)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 2
	}

	foundValue, err := managedType.ManagedMapContains(mMapHandle, keyHandle)
	if context.WithFault(err, runtime.ManagedMapAPIErrorShouldFailExecution()) {
		return 2
	}

	if foundValue {
		return 1
	}

	return 0
}

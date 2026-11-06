package executor

// The adapter is generated, so its conformance is asserted here, next to the one hand-written
// piece it depends on.
var _ VMHooksLegacy = (*VMHooksLegacyAdapter)(nil)
var _ ExecutorRuntime = fixedVMHooksRuntime{}

// fixedVMHooksRuntime is the runtime of a caller that has one hooks object for good, and no
// instance to speak of.
//
// It serves VMHooksLegacyAdapter's other constituency: anything that simulates a contract -
// mock contracts, tests - which needs the adapter for the error handling alone, and reaches
// the hooks directly rather than through an instance. The state is ignored, there being
// none; the Rust ExecutorRuntimeDefault ignores it for the same reason.
type fixedVMHooksRuntime struct {
	vmHooks VMHooks
}

// VMHooksForCall returns the same hooks every time.
func (runtime fixedVMHooksRuntime) VMHooksForCall(_ InstanceState) VMHooks {
	return runtime.vmHooks
}

// IsInterfaceNil returns true if underlying object is nil
func (runtime fixedVMHooksRuntime) IsInterfaceNil() bool {
	return false
}

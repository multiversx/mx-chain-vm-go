package executor

// VMHooksErrorHandler applies the side effects of an error returned by a VM hook.
// It is implemented by the VM, and called only by generated hook wrappers.
type VMHooksErrorHandler interface {
	// HandleHookError applies the failure side effects for err. It is a no-op
	// when err is nil.
	HandleHookError(err error)
}

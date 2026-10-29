package executor

// VMHooksErrorHandler applies the side effects of an error returned by a VM hook.
//
// It is implemented by the VM, and called only by VMHooksLegacyAdapter, which is the one
// place that has an error it cannot return. It is deliberately not part of VMHooks: a hook
// that can report its error has no use for it.
//
// Counterpart of the Rust VMHooksSetEarlyExit, and like it, it lives next to the adapter
// rather than in a backend, because more than one caller adapts hooks this way.
type VMHooksErrorHandler interface {
	// HandleHookError applies the failure side effects for err. It is a no-op
	// when err is nil.
	HandleHookError(err error)
}

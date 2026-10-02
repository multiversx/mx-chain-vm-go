package executor

// VMHooksProvider hands out the VM hooks an executor calls into.
//
// There are two methods rather than one because the backends differ in where the instance
// state comes from, and that difference cannot be hidden:
//
//   - Wasmer 6 hands each host call a context from which the state is reachable, and only for
//     the duration of that call. Its hooks are therefore built per call, around that state.
//   - Wasmer 2.2 goes through a C API whose hooks pointer is per *executor*, not per instance,
//     so a hook has no way to say which instance called it. Its state is instance-lifetime
//     anyway, so the provider resolves it through the host's instance stack instead — which
//     is what the VM has always done.
//
// Exposing both rather than papering over the difference keeps each backend honest about the
// lifetime it actually offers. A provider must implement both; an executor calls the one that
// matches what it can supply.
type VMHooksProvider interface {
	// VMHooksForCall returns the hooks servicing one host call, reaching the instance through
	// the given call-scoped state. The result must not outlive the call.
	VMHooksForCall(state InstanceState) VMHooks

	// VMHooksForInstanceLifetime returns hooks that resolve the instance through the host
	// rather than being told which one. For a backend that cannot supply a call-scoped state.
	VMHooksForInstanceLifetime() VMHooks
}

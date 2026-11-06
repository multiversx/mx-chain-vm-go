package vmhooks

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// NewVMHooksLegacy creates VM hooks that absorb their own errors, exactly as the executor
// does when it runs a contract.
//
// For mock contracts and tests. They stand in for a contract, and a contract never sees the
// error channel: a hook that fails stops the execution, it does not hand back an error to
// inspect. Going through VMHooksImpl directly would make every call site deal with an error
// that a real contract would never have been given.
//
// Takes only the host, unlike NewVMHooksImpl: a caller here is standing in for a contract
// the host is already running, so the instance is whatever the host has on its runtime and
// the state follows from it. During a real call the executor supplies the state instead,
// because it is the side that knows whose call the hook is serving.
func NewVMHooksLegacy(host vmhost.VMHost) *executor.VMHooksLegacyAdapter {
	vmHooksImpl := NewVMHooksImpl(host, InstanceStateOf(host))
	return executor.NewVMHooksLegacyAdapter(vmHooksImpl, vmHooksImpl)
}

// InstanceStateOf is the per-call state of whatever instance the host is currently running.
//
// What the executor would have supplied, recovered from the host instead, for callers that
// reach the hooks from outside a call - mock contracts, and tests that drive a hook
// directly. The instance may be a mock or a real one: a mock instance is its own state,
// while a real one hands out a separate handle, which its backend exposes for this.
//
// Nil when there is no instance, or when the backend only lets a hook reach the state from
// inside a call, as wasmer2new does. The hooks treat nil as "no memory to reach".
func InstanceStateOf(host vmhost.VMHost) executor.InstanceState {
	instance := host.Runtime().GetInstance()

	type stateProvider interface {
		InstanceState() executor.InstanceState
	}

	if provider, ok := instance.(stateProvider); ok {
		return provider.InstanceState()
	}

	// A mock instance implements the state itself. Named structurally rather than by type,
	// so that this package does not reach for the mocks.
	if state, ok := instance.(executor.InstanceState); ok {
		return state
	}

	return nil
}

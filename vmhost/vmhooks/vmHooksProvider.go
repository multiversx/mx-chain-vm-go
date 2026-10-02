package vmhooks

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ = (executor.VMHooksProvider)((*vmHooksProvider)(nil))

// vmHooksProvider builds the VM's hooks for whichever backend is asking.
//
// It holds nothing but the host, so producing a hooks object is a two-word allocation. That is
// what makes it affordable to build them per host call, which is what a call-scoped backend
// needs.
type vmHooksProvider struct {
	host vmhost.VMHost
}

// NewVMHooksProvider creates the VM hooks provider for a host.
func NewVMHooksProvider(host vmhost.VMHost) executor.VMHooksProvider {
	return &vmHooksProvider{host: host}
}

// VMHooksForCall returns hooks bound to one host call's instance state.
func (provider *vmHooksProvider) VMHooksForCall(state executor.InstanceState) executor.VMHooks {
	return NewVMHooksImplForCall(provider.host, state)
}

// VMHooksForInstanceLifetime returns hooks that resolve the instance through the host.
func (provider *vmHooksProvider) VMHooksForInstanceLifetime() executor.VMHooks {
	return NewVMHooksImpl(provider.host)
}

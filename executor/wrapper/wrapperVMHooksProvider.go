package executorwrapper

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ = (executor.VMHooksProvider)((*WrapperVMHooksProvider)(nil))

// WrapperVMHooksProvider wraps a VM hooks provider so that every hooks object it hands out
// is itself wrapped for logging.
//
// The wrapping has to happen at the provider rather than once at construction, because a
// call-scoped backend asks for a fresh hooks object on every host call.
type WrapperVMHooksProvider struct {
	logger          ExecutorLogger
	wrappedProvider executor.VMHooksProvider
}

// NewWrapperVMHooksProvider creates a logging wrapper around a VM hooks provider.
func NewWrapperVMHooksProvider(
	logger ExecutorLogger,
	wrappedProvider executor.VMHooksProvider,
) *WrapperVMHooksProvider {
	return &WrapperVMHooksProvider{
		logger:          logger,
		wrappedProvider: wrappedProvider,
	}
}

// VMHooksForCall wraps the hooks servicing one host call.
func (provider *WrapperVMHooksProvider) VMHooksForCall(state executor.InstanceState) executor.VMHooks {
	return &WrapperVMHooks{
		logger:         provider.logger,
		wrappedVMHooks: provider.wrappedProvider.VMHooksForCall(state),
	}
}

// VMHooksForInstanceLifetime wraps the hooks that resolve the instance through the host.
func (provider *WrapperVMHooksProvider) VMHooksForInstanceLifetime() executor.VMHooks {
	return &WrapperVMHooks{
		logger:         provider.logger,
		wrappedVMHooks: provider.wrappedProvider.VMHooksForInstanceLifetime(),
	}
}

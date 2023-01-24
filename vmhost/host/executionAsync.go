package host

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	runtime := host.Runtime()
	async := host.Async()
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	legacyGroupID := arwen.LegacyAsyncCallGroupID
	legacyGroup, exists := async.GetCallGroup(legacyGroupID)
	if !exists {
		return arwen.ErrLegacyAsyncCallNotFound

	}

	if legacyGroup.IsComplete() {
		return arwen.ErrLegacyAsyncCallInvalid
	}

	return nil
}

package hostCore

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

func (host *vmHost) handleAsyncCallBreakpoint() error {
	async := host.Async()

	legacyGroupID := vmhost.LegacyAsyncCallGroupID
	legacyGroup, exists := async.GetCallGroup(legacyGroupID)
	if !exists {
		return vmhost.ErrLegacyAsyncCallNotFound

	}

	if legacyGroup.IsComplete() {
		return vmhost.ErrLegacyAsyncCallInvalid
	}

	return nil
}

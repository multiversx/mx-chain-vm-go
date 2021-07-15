package worldmock

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/check"
)

// EpochNotifierStub -
type EpochNotifierStub struct {
	CurrentEpochCalled          func() uint32
	RegisterNotifyHandlerCalled func(handler vmcommon.EpochSubscriberHandler)
}

// RegisterNotifyHandler -
func (ens *EpochNotifierStub) RegisterNotifyHandler(handler vmcommon.EpochSubscriberHandler) {
	if ens.RegisterNotifyHandlerCalled != nil {
		ens.RegisterNotifyHandlerCalled(handler)
	} else {
		if !check.IfNil(handler) {
			handler.EpochConfirmed(0, 0)
		}
	}
}

// IsInterfaceNil -
func (ens *EpochNotifierStub) IsInterfaceNil() bool {
	return ens == nil
}

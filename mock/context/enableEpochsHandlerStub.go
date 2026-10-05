package mock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.EnableEpochsHandler = (*DefaultEnableEpochsHandlerStub)(nil)

// DefaultEnableEpochsHandlerStub provides a default stub implementation of vmhost.EnableEpochsHandler
type DefaultEnableEpochsHandlerStub struct{}

// IsFlagDefined returns false
func (d *DefaultEnableEpochsHandlerStub) IsFlagDefined(_ core.EnableEpochFlag) bool {
	return false
}

// IsFlagEnabled returns false
func (d *DefaultEnableEpochsHandlerStub) IsFlagEnabled(_ core.EnableEpochFlag) bool {
	return false
}

// IsFlagEnabledInEpoch returns false
func (d *DefaultEnableEpochsHandlerStub) IsFlagEnabledInEpoch(_ core.EnableEpochFlag, _ uint32) bool {
	return false
}

// GetActivationEpoch returns 0
func (d *DefaultEnableEpochsHandlerStub) GetActivationEpoch(_ core.EnableEpochFlag) uint32 {
	return 0
}

// IsInterfaceNil returns whether the pointer is nil
func (d *DefaultEnableEpochsHandlerStub) IsInterfaceNil() bool {
	return d == nil
}

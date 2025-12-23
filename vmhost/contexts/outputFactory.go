package contexts

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type outputContextFactory struct{}

// NewOutputContextFactory create a new output context factory
func NewOutputContextFactory() *outputContextFactory {
	return &outputContextFactory{}
}

// CreateOutputContext create the output context
func (ocf *outputContextFactory) CreateOutputContext(host vmhost.VMHost) (vmhost.OutputContext, error) {
	return NewOutputContext(host)
}

// IsInterfaceNil returns true if underlying object is nil
func (ocf *outputContextFactory) IsInterfaceNil() bool {
	return ocf == nil
}

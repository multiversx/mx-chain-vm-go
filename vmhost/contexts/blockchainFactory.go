package contexts

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type blockchainContextFactory struct{}

// NewBlockchainContextFactory create a new blockchain context factory
func NewBlockchainContextFactory() *blockchainContextFactory {
	return &blockchainContextFactory{}
}

// CreateBlockchainContext create the blockchain context
func (bcf *blockchainContextFactory) CreateBlockchainContext(host vmhost.VMHost, blockChainHook vmcommon.BlockchainHook) (vmhost.BlockchainContext, error) {
	return NewBlockchainContext(host, blockChainHook)
}

// IsInterfaceNil returns true if underlying object is nil
func (bcf *blockchainContextFactory) IsInterfaceNil() bool {
	return bcf == nil
}

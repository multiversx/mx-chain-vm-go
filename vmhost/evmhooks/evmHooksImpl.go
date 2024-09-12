package evmhooks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

type EVMHooksImpl struct {
	host vmhost.VMHost
}

// NewEVMHooksImpl creates a new EVMHooksImpl instance.
func NewEVMHooksImpl(host vmhost.VMHost) *EVMHooksImpl {
	return &EVMHooksImpl{host: host}
}

// GetBlockchainContext returns the blockchain context
func (context *EVMHooksImpl) GetBlockchainContext() vmhost.BlockchainContext {
	return context.host.Blockchain()
}

// GetRuntimeContext returns the runtime context
func (context *EVMHooksImpl) GetRuntimeContext() vmhost.RuntimeContext {
	return context.host.Runtime()
}

// GetCryptoContext returns the crypto context
func (context *EVMHooksImpl) GetCryptoContext() crypto.VMCrypto {
	return context.host.Crypto()
}

// GetManagedTypesContext returns the big int context
func (context *EVMHooksImpl) GetManagedTypesContext() vmhost.ManagedTypesContext {
	return context.host.ManagedTypes()
}

// GetOutputContext returns the output context
func (context *EVMHooksImpl) GetOutputContext() vmhost.OutputContext {
	return context.host.Output()
}

// GetMeteringContext returns the metering context
func (context *EVMHooksImpl) GetMeteringContext() vmhost.MeteringContext {
	return context.host.Metering()
}

// GetStorageContext returns the storage context
func (context *EVMHooksImpl) GetStorageContext() vmhost.StorageContext {
	return context.host.Storage()
}

func (context *EVMHooksImpl) WithFault(err error) bool {
	return vmhooks.WithFaultAndHost(context.host, err, true)
}

func (context *EVMHooksImpl) toEVMAddress(address []byte) common.Address {
	addressResponse, err := context.GetBlockchainContext().RequestAddress(&vmcommon.AddressRequest{
		SourceAddress:       address,
		SourceIdentifier:    core.MVXAddressIdentifier,
		RequestedIdentifier: core.ETHAddressIdentifier,
		SaveOnGenerate:      true,
	})
	if err != nil {
		panic(err)
	}
	return common.BytesToAddress(addressResponse.RequestedAddress)
}

func (context *EVMHooksImpl) toMVXAddress(address common.Address) []byte {
	addressResponse, err := context.GetBlockchainContext().RequestAddress(&vmcommon.AddressRequest{
		SourceAddress:       address.Bytes(),
		SourceIdentifier:    core.ETHAddressIdentifier,
		RequestedIdentifier: core.MVXAddressIdentifier,
		SaveOnGenerate:      true,
	})
	if err != nil {
		panic(err)
	}
	return addressResponse.RequestedAddress
}

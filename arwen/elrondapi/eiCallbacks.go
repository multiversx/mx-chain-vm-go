package elrondapi

import (
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/crypto"
)

// EICallbacks is the VM structure that implements VMHooks,
// with all the hooks (callbacks) from the executor.
type EICallbacks struct {
	host arwen.VMHost
}

// NewEICallbacks creates a new EICallbacks instance.
func NewEICallbacks(host arwen.VMHost) *EICallbacks {
	return &EICallbacks{
		host: host,
	}
}

// GetVMHost returns the vm Context from the vm context map
func (context *EICallbacks) GetVMHost() arwen.VMHost {
	return context.host
}

// GetBlockchainContext returns the blockchain context
func (context *EICallbacks) GetBlockchainContext() arwen.BlockchainContext {
	return context.host.Blockchain()
}

// GetRuntimeContext returns the runtime context
func (context *EICallbacks) GetRuntimeContext() arwen.RuntimeContext {
	return context.host.Runtime()
}

// GetCryptoContext returns the crypto context
func (context *EICallbacks) GetCryptoContext() crypto.VMCrypto {
	return context.host.Crypto()
}

// GetManagedTypesContext returns the big int context
func (context *EICallbacks) GetManagedTypesContext() arwen.ManagedTypesContext {
	return context.host.ManagedTypes()
}

// GetOutputContext returns the output context
func (context *EICallbacks) GetOutputContext() arwen.OutputContext {
	return context.host.Output()
}

// GetMeteringContext returns the metering context
func (context *EICallbacks) GetMeteringContext() arwen.MeteringContext {
	return context.host.Metering()
}

// GetStorageContext returns the storage context
func (context *EICallbacks) GetStorageContext() arwen.StorageContext {
	return context.host.Storage()
}

// WithFault handles an error, taking into account whether it should completely
// fail the execution of a contract or not.
func (context *EICallbacks) WithFault(err error, failExecution bool) bool {
	return WithFaultAndHost(context.host, err, failExecution)
}

// WithFaultAndHost fails the execution with the provided error
func WithFaultAndHost(host arwen.VMHost, err error, failExecution bool) bool {
	if err == nil {
		return false
	}

	if failExecution {
		runtime := host.Runtime()
		metering := host.Metering()
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
	}

	return true
}

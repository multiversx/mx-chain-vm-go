package elrondapi

import (
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/crypto"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

//go:generate go run generate/cmd/eiGenMain.go

// ElrondApi is the VM structure that implements VMHooks,
// with all the hooks (callbacks) from the executor.
type ElrondApi struct {
	host arwen.VMHost
}

// NewElrondApi creates a new ElrondApi instance.
func NewElrondApi(host arwen.VMHost) *ElrondApi {
	return &ElrondApi{
		host: host,
	}
}

// MemLoad returns the contents from the given offset of the WASM memory.
func (context *ElrondApi) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return context.host.Runtime().GetInstance().MemLoad(memPtr, length)
}

// MemLoadMultiple returns multiple byte slices loaded from the WASM memory, starting at the given offset and having the provided lengths.
func (context *ElrondApi) MemLoadMultiple(memPtr executor.MemPtr, lengths []int32) ([][]byte, error) {
	if len(lengths) == 0 {
		return [][]byte{}, nil
	}

	results := make([][]byte, len(lengths))

	for i, length := range lengths {
		result, err := context.MemLoad(memPtr, length)
		if err != nil {
			return nil, err
		}

		results[i] = result
		memPtr = memPtr.Offset(length)
	}

	return results, nil
}

// GetVMHost returns the vm Context from the vm context map
func (context *ElrondApi) GetVMHost() arwen.VMHost {
	return context.host
}

// GetBlockchainContext returns the blockchain context
func (context *ElrondApi) GetBlockchainContext() arwen.BlockchainContext {
	return context.host.Blockchain()
}

// GetRuntimeContext returns the runtime context
func (context *ElrondApi) GetRuntimeContext() arwen.RuntimeContext {
	return context.host.Runtime()
}

// GetCryptoContext returns the crypto context
func (context *ElrondApi) GetCryptoContext() crypto.VMCrypto {
	return context.host.Crypto()
}

// GetManagedTypesContext returns the big int context
func (context *ElrondApi) GetManagedTypesContext() arwen.ManagedTypesContext {
	return context.host.ManagedTypes()
}

// GetOutputContext returns the output context
func (context *ElrondApi) GetOutputContext() arwen.OutputContext {
	return context.host.Output()
}

// GetMeteringContext returns the metering context
func (context *ElrondApi) GetMeteringContext() arwen.MeteringContext {
	return context.host.Metering()
}

// GetStorageContext returns the storage context
func (context *ElrondApi) GetStorageContext() arwen.StorageContext {
	return context.host.Storage()
}

// WithFault handles an error, taking into account whether it should completely
// fail the execution of a contract or not.
func (context *ElrondApi) WithFault(err error, failExecution bool) bool {
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

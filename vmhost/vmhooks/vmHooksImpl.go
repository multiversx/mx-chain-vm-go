package vmhooks

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

//go:generate go run generate/cmd/eiGenMain.go

// VMHooksImpl is the VM structure that implements VMHooks,
// with all the hooks (callbacks) from the executor.
//
// Built per host call, around the state of the instance making that call. It used to be
// built once per host and find the instance through the runtime's own notion of "current",
// which worked only because the two always agreed; the state now arrives from the executor,
// which is the side that actually knows whose call this is.
type VMHooksImpl struct {
	host vmhost.VMHost

	// instanceState is the executing instance, for the duration of one host call.
	//
	// Must not outlive the hook it was handed to - see executor.InstanceState. Nil when the
	// hooks were built outside a call, which is why every use of it is guarded: a mock
	// contract or a test may reach for hooks with no instance running, and those hooks have
	// no business touching memory.
	instanceState executor.InstanceState
}

// NewVMHooksImpl creates VM hooks servicing one host call, around the state of the instance
// making it.
func NewVMHooksImpl(host vmhost.VMHost, instanceState executor.InstanceState) *VMHooksImpl {
	return &VMHooksImpl{
		host:          host,
		instanceState: instanceState,
	}
}

// MemLoad returns the contents from the given offset of the WASM memory.
func (context *VMHooksImpl) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	if check.IfNil(context.instanceState) {
		return nil, vmhost.ErrNilInstance
	}

	return context.instanceState.MemoryLoadClamped(memPtr, length)
}

// MemLoadMultiple returns multiple byte slices loaded from the WASM memory, starting at the given offset and having the provided lengths.
func (context *VMHooksImpl) MemLoadMultiple(memPtr executor.MemPtr, lengths []int32) ([][]byte, error) {
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

// MemStore stores the given data in the WASM memory at the given offset.
func (context *VMHooksImpl) MemStore(memPtr executor.MemPtr, data []byte) error {
	if check.IfNil(context.instanceState) {
		return vmhost.ErrNilInstance
	}

	return context.instanceState.MemoryStore(memPtr, data)
}

// GetVMHost returns the vm Context from the vm context map
func (context *VMHooksImpl) GetVMHost() vmhost.VMHost {
	return context.host
}

// GetBlockchainContext returns the blockchain context
func (context *VMHooksImpl) GetBlockchainContext() vmhost.BlockchainContext {
	return context.host.Blockchain()
}

// GetRuntimeContext returns the runtime context
func (context *VMHooksImpl) GetRuntimeContext() vmhost.RuntimeContext {
	return context.host.Runtime()
}

// GetCryptoContext returns the crypto context
func (context *VMHooksImpl) GetCryptoContext() crypto.VMCrypto {
	return context.host.Crypto()
}

// GetManagedTypesContext returns the big int context
func (context *VMHooksImpl) GetManagedTypesContext() vmhost.ManagedTypesContext {
	return context.host.ManagedTypes()
}

// GetOutputContext returns the output context
func (context *VMHooksImpl) GetOutputContext() vmhost.OutputContext {
	return context.host.Output()
}

// GetMeteringContext returns the metering context
func (context *VMHooksImpl) GetMeteringContext() vmhost.MeteringContext {
	return context.host.Metering()
}

// GetStorageContext returns the storage context
func (context *VMHooksImpl) GetStorageContext() vmhost.StorageContext {
	return context.host.Storage()
}

// GetEnableEpochsHandler returns the enable epochs handler
func (context *VMHooksImpl) GetEnableEpochsHandler() vmhost.EnableEpochsHandler {
	return context.host.EnableEpochsHandler()
}

// FailExecution fails the execution with the provided error
func FailExecution(host vmhost.VMHost, err error) {
	if err == nil {
		return
	}

	runtime := host.Runtime()
	metering := host.Metering()
	_ = metering.UseGasBounded(metering.GasLeft())
	runtime.FailExecution(err)
}

package vmhooks

import (
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

//go:generate go run generate/cmd/eiGenMain.go

// VMHooksImpl is the VM structure that implements VMHooks,
// with all the hooks (callbacks) from the executor.
//
// Every hook that touches wasm memory does so through the state below, which is the one place
// the two backends' differing lifetimes meet. See executor.InstanceState.
type VMHooksImpl struct {
	host  vmhost.VMHost
	state executor.InstanceState
}

// NewVMHooksImpl creates hooks that reach the instance through the host's instance stack.
//
// This is the shape the legacy executor needs, because its hooks pointer is per executor and a
// hook cannot be told which instance called it. It is also what a test wants, since the host
// it is given already has an instance running.
func NewVMHooksImpl(host vmhost.VMHost) *VMHooksImpl {
	return &VMHooksImpl{
		host:  host,
		state: hostInstanceState{host: host},
	}
}

// NewVMHooksImplForCall creates hooks bound to the instance state of one host call.
//
// The state is taken as a parameter rather than kept in a mutable slot, which is what makes
// nesting safe by construction: a hook that re-enters the VM gets its own hooks object around
// its own state, and cannot clobber its caller's.
func NewVMHooksImplForCall(host vmhost.VMHost, state executor.InstanceState) *VMHooksImpl {
	return &VMHooksImpl{
		host:  host,
		state: state,
	}
}

// hostInstanceState answers "which instance" dynamically, from the host's instance stack.
//
// Sound only for a backend whose state stays valid for the instance's lifetime, which is the
// case for Wasmer 2.2. A call-scoped backend must supply its own state instead.
type hostInstanceState struct {
	host vmhost.VMHost
}

func (state hostInstanceState) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return state.host.Runtime().GetInstance().MemLoad(memPtr, length)
}

func (state hostInstanceState) MemStore(memPtr executor.MemPtr, data []byte) error {
	return state.host.Runtime().GetInstance().MemStore(memPtr, data)
}

// MemLoad returns the contents from the given offset of the WASM memory.
func (context *VMHooksImpl) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return context.state.MemLoad(memPtr, length)
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
	return context.state.MemStore(memPtr, data)
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

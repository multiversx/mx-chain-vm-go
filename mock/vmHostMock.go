package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/parsers"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ arwen.VMHost = (*VmHostMock)(nil)

type VmHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext arwen.BlockchainContext
	RuntimeContext    arwen.RuntimeContext
	AsyncContext      arwen.AsyncContext
	OutputContext     arwen.OutputContext
	MeteringContext   arwen.MeteringContext
	StorageContext    arwen.StorageContext
	BigIntContext     arwen.BigIntContext

	SCAPIMethods  *wasmer.Imports
	IsBuiltinFunc bool

	StoredInputs []*vmcommon.ContractCallInput

	VMOutputQueue    []*vmcommon.VMOutput
	VMOutputToReturn int
	Err              error
}

func (host *VmHostMock) Crypto() crypto.VMCrypto {
	return host.CryptoHook
}

func (host *VmHostMock) Blockchain() arwen.BlockchainContext {
	return host.BlockchainContext
}

func (host *VmHostMock) Runtime() arwen.RuntimeContext {
	return host.RuntimeContext
}

func (host *VmHostMock) Async() arwen.AsyncContext {
	return host.AsyncContext
}

func (host *VmHostMock) Output() arwen.OutputContext {
	return host.OutputContext
}

func (host *VmHostMock) Metering() arwen.MeteringContext {
	return host.MeteringContext
}

func (host *VmHostMock) Storage() arwen.StorageContext {
	return host.StorageContext
}

func (host *VmHostMock) BigInt() arwen.BigIntContext {
	return host.BigIntContext
}

func (host *VmHostMock) GetHostContexts() (
	arwen.BigIntContext,
	arwen.BlockchainContext,
	arwen.MeteringContext,
	arwen.OutputContext,
	arwen.RuntimeContext,
	arwen.AsyncContext,
	arwen.StorageContext,
) {
	return host.BigIntContext,
		host.BlockchainContext,
		host.MeteringContext,
		host.OutputContext,
		host.RuntimeContext,
		host.AsyncContext,
		host.StorageContext
}

func (host *VmHostMock) CallArgsParser() arwen.CallArgsParser {
	return parsers.NewCallArgsParser()
}

func (host *VmHostMock) IsArwenV2Enabled() bool {
	return true
}

func (host *VmHostMock) IsAheadOfTimeCompileEnabled() bool {
	return true
}

func (host *VmHostMock) IsDynamicGasLockingEnabled() bool {
	return true
}

func (host *VmHostMock) CreateNewContract(_ *vmcommon.ContractCreateInput) ([]byte, error) {
	return nil, host.Err
}

func (host *VmHostMock) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	return host.Err
}

func (host *VmHostMock) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if host.Err != nil {
		return nil, host.Err
	}
	host.StoreInput(input)
	return host.GetNextVMOutput(), nil
}

func (host *VmHostMock) EthereumCallData() []byte {
	return host.EthInput
}

func (host *VmHostMock) InitState() {
}

func (host *VmHostMock) PushState() {
}

func (host *VmHostMock) PopState() {
}

func (host *VmHostMock) ClearStateStack() {
}

func (host *VmHostMock) GetAPIMethods() *wasmer.Imports {
	return host.SCAPIMethods
}

func (host *VmHostMock) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

func (host *VmHostMock) IsBuiltinFunctionName(_ string) bool {
	return host.IsBuiltinFunc
}

func (host *VmHostMock) StoreInput(input *vmcommon.ContractCallInput) {
	if host.StoredInputs == nil {
		host.StoredInputs = make([]*vmcommon.ContractCallInput, 0)
	}
	host.StoredInputs = append(host.StoredInputs, input)
}

func (host *VmHostMock) EnqueueVMOutput(vmOutput *vmcommon.VMOutput) {
	if host.VMOutputQueue == nil {
		host.VMOutputQueue = make([]*vmcommon.VMOutput, 1)
		host.VMOutputQueue[0] = vmOutput
		host.VMOutputToReturn = 0
		return
	}

	host.VMOutputQueue = append(host.VMOutputQueue, vmOutput)
}

func (host *VmHostMock) GetNextVMOutput() *vmcommon.VMOutput {
	if host.VMOutputToReturn >= len(host.VMOutputQueue) {
		return nil
	}

	vmOutput := host.VMOutputQueue[host.VMOutputToReturn]
	host.VMOutputToReturn += 1
	return vmOutput
}

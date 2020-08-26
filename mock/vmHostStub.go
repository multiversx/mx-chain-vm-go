package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.VMHost = (*VmHostStub)(nil)

type VmHostStub struct {
	InitStateCalled       func()
	PushStateCalled       func()
	PopStateCalled        func()
	ClearStateStackCalled func()

	CryptoCalled                      func() crypto.VMCrypto
	BlockchainCalled                  func() arwen.BlockchainContext
	RuntimeCalled                     func() arwen.RuntimeContext
	BigIntCalled                      func() arwen.BigIntContext
	OutputCalled                      func() arwen.OutputContext
	MeteringCalled                    func() arwen.MeteringContext
	StorageCalled                     func() arwen.StorageContext
	CreateNewContractCalled           func(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContextCalled        func(input *vmcommon.ContractCallInput) (*arwen.AsyncContextInfo, error)
	ExecuteOnDestContextCalled        func(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, error)
	EthereumCallDataCalled            func() []byte
	GetAPIMethodsCalled               func() *wasmer.Imports
	GetProtocolBuiltinFunctionsCalled func() vmcommon.FunctionNames
	IsBuiltinFunctionNameCalled       func(functionName string) bool
}

func (vhs *VmHostStub) InitState() {
	if vhs.InitStateCalled != nil {
		vhs.InitStateCalled()
	}
}

func (vhs *VmHostStub) PushState() {
	if vhs.PushStateCalled != nil {
		vhs.PushStateCalled()
	}
}

func (vhs *VmHostStub) PopState() {
	if vhs.PopStateCalled != nil {
		vhs.PopStateCalled()
	}
}

func (vhs *VmHostStub) ClearStateStack() {
	if vhs.ClearStateStackCalled != nil {
		vhs.ClearStateStackCalled()
	}
}

func (vhs *VmHostStub) Crypto() crypto.VMCrypto {
	if vhs.CryptoCalled != nil {
		return vhs.CryptoCalled()
	}
	return nil
}

func (vhs *VmHostStub) Blockchain() arwen.BlockchainContext {
	if vhs.BlockchainCalled != nil {
		return vhs.BlockchainCalled()
	}
	return nil
}

func (vhs *VmHostStub) Runtime() arwen.RuntimeContext {
	if vhs.RuntimeCalled != nil {
		return vhs.RuntimeCalled()
	}
	return nil
}

func (vhs *VmHostStub) BigInt() arwen.BigIntContext {
	if vhs.BigIntCalled != nil {
		return vhs.BigIntCalled()
	}
	return nil
}

func (vhs *VmHostStub) Output() arwen.OutputContext {
	if vhs.OutputCalled != nil {
		return vhs.OutputCalled()
	}
	return nil
}

func (vhs *VmHostStub) Metering() arwen.MeteringContext {
	if vhs.MeteringCalled != nil {
		return vhs.MeteringCalled()
	}
	return nil
}

func (vhs *VmHostStub) Storage() arwen.StorageContext {
	if vhs.StorageCalled != nil {
		return vhs.StorageCalled()
	}
	return nil
}

func (vhs *VmHostStub) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	if vhs.CreateNewContractCalled != nil {
		return vhs.CreateNewContractCalled(input)
	}
	return nil, nil
}

func (vhs *VmHostStub) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*arwen.AsyncContextInfo, error) {
	if vhs.ExecuteOnSameContextCalled != nil {
		return vhs.ExecuteOnSameContextCalled(input)
	}
	return nil, nil
}

func (vhs *VmHostStub) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, error) {
	if vhs.ExecuteOnDestContextCalled != nil {
		return vhs.ExecuteOnDestContextCalled(input)
	}
	return nil, nil, nil
}

func (vhs *VmHostStub) EthereumCallData() []byte {
	if vhs.EthereumCallDataCalled != nil {
		return vhs.EthereumCallDataCalled()
	}
	return nil
}

func (vhs *VmHostStub) GetAPIMethods() *wasmer.Imports {
	if vhs.GetAPIMethodsCalled != nil {
		return vhs.GetAPIMethodsCalled()
	}
	return nil
}

func (vhs *VmHostStub) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	if vhs.GetProtocolBuiltinFunctionsCalled != nil {
		return vhs.GetProtocolBuiltinFunctionsCalled()
	}
	return make(vmcommon.FunctionNames)
}

func (vhs *VmHostStub) IsBuiltinFunctionName(functionName string) bool {
	if vhs.IsBuiltinFunctionNameCalled != nil {
		return vhs.IsBuiltinFunctionNameCalled(functionName)
	}
	return false
}

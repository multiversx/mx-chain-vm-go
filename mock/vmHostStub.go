package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.VMHost = (*VmHostStub)(nil)

type VmHostStub struct {
	InitStateCalled       func()
	PushStateCalled       func()
	PopStateCalled        func()
	ClearStateStackCalled func()

	CryptoCalled               func() vmcommon.CryptoHook
	BlockchainCalled           func() arwen.BlockchainContext
	RuntimeCalled              func() arwen.RuntimeContext
	BigIntCalled               func() arwen.BigIntContext
	OutputCalled               func() arwen.OutputContext
	MeteringCalled             func() arwen.MeteringContext
	StorageCalled              func() arwen.StorageContext
	CreateNewContractCalled    func(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContextCalled func(input *vmcommon.ContractCallInput) error
	ExecuteOnDestContextCalled func(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	EthereumCallDataCalled     func() []byte
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

func (vhs *VmHostStub) Crypto() vmcommon.CryptoHook {
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

func (vhs *VmHostStub) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	if vhs.ExecuteOnSameContextCalled != nil {
		return vhs.ExecuteOnSameContextCalled(input)
	}
	return nil
}

func (vhs *VmHostStub) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if vhs.ExecuteOnDestContextCalled != nil {
		return vhs.ExecuteOnDestContextCalled(input)
	}
	return nil, nil
}

func (vhs *VmHostStub) EthereumCallData() []byte {
	if vhs.EthereumCallDataCalled != nil {
		return vhs.EthereumCallDataCalled()
	}
	return nil
}

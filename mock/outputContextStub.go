package mock

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.OutputContext = (*OutputContextStub)(nil)

type OutputContextStub struct {
	InitStateCalled                   func()
	PushStateCalled                   func()
	PopSetActiveStateCalled           func()
	PopMergeActiveStateCalled         func()
	PopDiscardCalled                  func()
	ClearStateStackCalled             func()
	CopyTopOfStackToActiveStateCalled func()
	CensorVMOutputCalled              func()
	GetOutputAccountCalled            func(address []byte) (*vmcommon.OutputAccount, bool)
	WriteLogCalled                    func(address []byte, topics [][]byte, data []byte)
	TransferCalled                    func(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) error
	SelfDestructCalled                func(address []byte, beneficiary []byte)
	GetRefundCalled                   func() uint64
	SetRefundCalled                   func(refund uint64)
	ReturnCodeCalled                  func() vmcommon.ReturnCode
	SetReturnCodeCalled               func(returnCode vmcommon.ReturnCode)
	ReturnMessageCalled               func() string
	SetReturnMessageCalled            func(message string)
	ReturnDataCalled                  func() [][]byte
	ClearReturnDataCalled             func()
	FinishCalled                      func(data []byte)
	GetVMOutputCalled                 func() *vmcommon.VMOutput
	AddTxValueToAccountCalled         func(address []byte, value *big.Int)
	DeployCodeCalled                  func(input arwen.CodeDeployInput)
	CreateVMOutputInCaseOfErrorCalled func(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput
}

func (o *OutputContextStub) InitState() {
	if o.InitStateCalled != nil {
		o.InitStateCalled()
	}
}

func (o *OutputContextStub) PushState() {
	if o.PushStateCalled != nil {
		o.PushStateCalled()
	}
}

func (o *OutputContextStub) PopSetActiveState() {
	if o.PopSetActiveStateCalled != nil {
		o.PopSetActiveStateCalled()
	}
}

func (o *OutputContextStub) PopMergeActiveState() {
	if o.PopMergeActiveStateCalled != nil {
		o.PopMergeActiveStateCalled()
	}
}

func (o *OutputContextStub) PopDiscard() {
	if o.PopDiscardCalled != nil {
		o.PopDiscardCalled()
	}
}

func (o *OutputContextStub) ClearStateStack() {
	if o.ClearStateStackCalled != nil {
		o.ClearStateStackCalled()
	}
}

func (o *OutputContextStub) CopyTopOfStackToActiveState() {
	if o.CopyTopOfStackToActiveStateCalled != nil {
		o.CopyTopOfStackToActiveStateCalled()
	}
}

func (o *OutputContextStub) CensorVMOutput() {
	if o.CensorVMOutputCalled != nil {
		o.CensorVMOutputCalled()
	}
}

func (o *OutputContextStub) GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool) {
	if o.GetOutputAccountCalled != nil {
		return o.GetOutputAccountCalled(address)
	}
	return nil, false
}

func (o *OutputContextStub) WriteLog(address []byte, topics [][]byte, data []byte) {
	if o.WriteLogCalled != nil {
		o.WriteLogCalled(address, topics, data)
	}
}

func (o *OutputContextStub) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) error {
	if o.TransferCalled != nil {
		return o.TransferCalled(destination, sender, gasLimit, value, input)
	}

	return nil
}

func (o *OutputContextStub) SelfDestruct(address []byte, beneficiary []byte) {
	if o.SelfDestructCalled != nil {
		o.SelfDestructCalled(address, beneficiary)
	}
}

func (o *OutputContextStub) GetRefund() uint64 {
	if o.GetRefundCalled != nil {
		return o.GetRefundCalled()
	}
	return 0
}

func (o *OutputContextStub) SetRefund(refund uint64) {
	if o.SetRefundCalled != nil {
		o.SetRefundCalled(refund)
	}
}

func (o *OutputContextStub) ReturnCode() vmcommon.ReturnCode {
	if o.ReturnCodeCalled != nil {
		return o.ReturnCodeCalled()
	}
	return vmcommon.Ok
}

func (o *OutputContextStub) SetReturnCode(returnCode vmcommon.ReturnCode) {
	if o.SetReturnCodeCalled != nil {
		o.SetReturnCodeCalled(returnCode)
	}
}
func (o *OutputContextStub) ReturnMessage() string {
	if o.ReturnMessageCalled != nil {
		return o.ReturnMessageCalled()
	}
	return ""
}

func (o *OutputContextStub) SetReturnMessage(message string) {
	if o.SetReturnMessageCalled != nil {
		o.SetReturnMessageCalled(message)
	}
}

func (o *OutputContextStub) ReturnData() [][]byte {
	if o.ReturnDataCalled != nil {
		return o.ReturnDataCalled()
	}
	return [][]byte{}
}

func (o *OutputContextStub) ClearReturnData() {
	if o.ClearReturnDataCalled != nil {
		o.ClearReturnDataCalled()
	}
}

func (o *OutputContextStub) Finish(data []byte) {
	if o.FinishCalled != nil {
		o.FinishCalled(data)
	}
}

func (o *OutputContextStub) GetVMOutput() *vmcommon.VMOutput {
	if o.GetVMOutputCalled != nil {
		return o.GetVMOutputCalled()
	}
	return nil
}

func (o *OutputContextStub) AddTxValueToAccount(address []byte, value *big.Int) {
	if o.AddTxValueToAccountCalled != nil {
		o.AddTxValueToAccountCalled(address, value)
	}
}

func (o *OutputContextStub) DeployCode(input arwen.CodeDeployInput) {
	if o.DeployCodeCalled != nil {
		o.DeployCodeCalled(input)
	}
}

func (o *OutputContextStub) CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	if o.CreateVMOutputInCaseOfErrorCalled != nil {
		return o.CreateVMOutputInCaseOfErrorCalled(errCode, message)
	}
	return nil
}

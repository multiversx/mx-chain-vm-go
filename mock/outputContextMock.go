package mock

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.OutputContext = (*OutputContextMock)(nil)

type OutputContextMock struct {
	OutputStateMock    *vmcommon.VMOutput
	ReturnDataMock     [][]byte
	ReturnCodeMock     vmcommon.ReturnCode
	ReturnMessageMock  string
	GasRemaining       uint64
	GasRefund          *big.Int
	OutputAccounts     map[string]*vmcommon.OutputAccount
	DeletedAccounts    [][]byte
	TouchedAccounts    [][]byte
	Logs               []*vmcommon.LogEntry
	OutputAccountMock  *vmcommon.OutputAccount
	OutputAccountIsNew bool
	Err                error
}

func (o *OutputContextMock) InitState() {
}

func (o *OutputContextMock) NewVMOutputAccount(address []byte) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

func (o *OutputContextMock) NewVMOutputAccountFromMockAccount(account *Account) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        account.Address,
		Nonce:          account.Nonce,
		BalanceDelta:   account.BalanceDelta,
		Balance:        account.Balance,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

func (o *OutputContextMock) PushState() {
}

func (o *OutputContextMock) PopState() {
}

func (o *OutputContextMock) ClearStateStack() {
}

func (o *OutputContextMock) GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool) {
	return o.OutputAccountMock, o.OutputAccountIsNew
}

func (o *OutputContextMock) GetRefund() uint64 {
	return uint64(o.GasRefund.Int64())
}

func (o *OutputContextMock) SetRefund(refund uint64) {
	o.GasRefund = big.NewInt(int64(refund))
}

func (o *OutputContextMock) ReturnData() [][]byte {
	return o.ReturnDataMock
}

func (o *OutputContextMock) ReturnCode() vmcommon.ReturnCode {
	return o.ReturnCodeMock
}

func (o *OutputContextMock) SetReturnCode(returnCode vmcommon.ReturnCode) {
	o.ReturnCodeMock = returnCode
}

func (o *OutputContextMock) ReturnMessage() string {
	return o.ReturnMessageMock
}

func (o *OutputContextMock) SetReturnMessage(returnMessage string) {
	o.ReturnMessageMock = returnMessage
}

func (o *OutputContextMock) ClearReturnData() {
	o.ReturnDataMock = make([][]byte, 0)
}

func (o *OutputContextMock) SelfDestruct(_ []byte, _ []byte) {
	panic("not implemented")
}

func (o *OutputContextMock) Finish(data []byte) {
	if len(data) > 0 {
		o.ReturnDataMock = append(o.ReturnDataMock, data)
	}
}

func (o *OutputContextMock) FinishValue(value wasmer.Value) {
	if !value.IsVoid() {
		valueBytes := arwen.ConvertReturnValue(value)
		o.Finish(valueBytes)
	}
}

func (o *OutputContextMock) WriteLog(address []byte, topics [][]byte, data []byte) {
}

func (o *OutputContextMock) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) {
}

func (o *OutputContextMock) AddTxValueToAccount(address []byte, value *big.Int) {
}

func (o *OutputContextMock) GetVMOutput() *vmcommon.VMOutput {
	return o.OutputStateMock
}

func (o *OutputContextMock) DeployCode(address []byte, code []byte) {
}

func (o *OutputContextMock) CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	return o.OutputStateMock
}

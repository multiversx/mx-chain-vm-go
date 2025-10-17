package mock

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.OutputContext = (*OutputContextMock)(nil)

// OutputContextMock is used in tests to check the OutputContext interface method calls
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
	TransferResult     error
	CrtTransferIndex   uint32
}

// AddToActiveState mocked method
func (o *OutputContextMock) AddToActiveState(_ *vmcommon.VMOutput) {
}

// InitState mocked method
func (o *OutputContextMock) InitState() {
}

// NewVMOutputAccount mocked method
func (o *OutputContextMock) NewVMOutputAccount(address []byte) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

// NewVMOutputAccountFromMockAccount mocked method
func (o *OutputContextMock) NewVMOutputAccountFromMockAccount(account *worldmock.Account) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        account.Address,
		Nonce:          account.Nonce,
		BalanceDelta:   big.NewInt(0),
		Balance:        account.Balance,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

// PushState mocked method
func (o *OutputContextMock) PushState() {
}

// PopSetActiveState mocked method
func (o *OutputContextMock) PopSetActiveState() {
}

// PopMergeActiveState mocked method
func (o *OutputContextMock) PopMergeActiveState() {
}

// PopDiscard mocked method
func (o *OutputContextMock) PopDiscard() {
}

// ClearStateStack mocked method
func (o *OutputContextMock) ClearStateStack() {
}

// CopyTopOfStackToActiveState mocked method
func (o *OutputContextMock) CopyTopOfStackToActiveState() {
}

// CensorVMOutput mocked method
func (o *OutputContextMock) CensorVMOutput() {
}

// GetOutputAccounts mocked method
func (o *OutputContextMock) GetOutputAccounts() map[string]*vmcommon.OutputAccount {
	return o.OutputAccounts
}

// GetOutputAccount mocked method
func (o *OutputContextMock) GetOutputAccount(_ []byte) (*vmcommon.OutputAccount, bool) {
	return o.OutputAccountMock, o.OutputAccountIsNew
}

// DeleteAccount mocked method
func (o *OutputContextMock) DeleteAccount(_ []byte) {
}

// DeleteOutputAccount mocked method
func (o *OutputContextMock) DeleteOutputAccount(_ []byte) {
}

// GetRefund mocked method
func (o *OutputContextMock) GetRefund() uint64 {
	return uint64(o.GasRefund.Int64())
}

// SetRefund mocked method
func (o *OutputContextMock) SetRefund(refund uint64) {
	o.GasRefund = big.NewInt(int64(refund))
}

// ReturnData mocked method
func (o *OutputContextMock) ReturnData() [][]byte {
	return o.ReturnDataMock
}

// ReturnCode mocked method
func (o *OutputContextMock) ReturnCode() vmcommon.ReturnCode {
	return o.ReturnCodeMock
}

// SetReturnCode mocked method
func (o *OutputContextMock) SetReturnCode(returnCode vmcommon.ReturnCode) {
	o.ReturnCodeMock = returnCode
}

// ReturnMessage mocked method
func (o *OutputContextMock) ReturnMessage() string {
	return o.ReturnMessageMock
}

// SetReturnMessage mocked method
func (o *OutputContextMock) SetReturnMessage(returnMessage string) {
	o.ReturnMessageMock = returnMessage
}

// ClearReturnData mocked method
func (o *OutputContextMock) ClearReturnData() {
	o.ReturnDataMock = make([][]byte, 0)
}

// RemoveReturnData mocked method
func (o *OutputContextMock) RemoveReturnData(_ uint32) {
}

// Finish mocked method
func (o *OutputContextMock) Finish(data []byte) {
	o.ReturnDataMock = append(o.ReturnDataMock, data)
}

// PrependFinish mocked method
func (o *OutputContextMock) PrependFinish(data []byte) {
	o.ReturnDataMock = append([][]byte{data}, o.ReturnDataMock...)
}

// DeleteFirstReturnData mocked method
func (o *OutputContextMock) DeleteFirstReturnData() {
	if len(o.ReturnDataMock) > 0 {
		o.ReturnDataMock = o.ReturnDataMock[1:]
	}
}

// WriteLog mocked method
func (o *OutputContextMock) WriteLog(_ []byte, _ [][]byte, _ [][]byte) {}

// WriteLogWithIdentifier mocked method
func (o *OutputContextMock) WriteLogWithIdentifier(_ []byte, _ [][]byte, _ [][]byte, _ []byte) {
}

// TransferValueOnly mocked method
func (o *OutputContextMock) TransferValueOnly(_ []byte, _ []byte, _ *big.Int, _ bool) error {
	return o.TransferResult
}

// Transfer mocked method
func (o *OutputContextMock) Transfer(_ []byte, _ []byte, _ uint64, _ uint64, _ *big.Int, _ []byte, _ []byte, _ vm.CallType) error {
	return o.TransferResult
}

// TransferESDT mocked method
func (o *OutputContextMock) TransferESDT(_ *vmhost.ESDTTransfersArgs, _ *vmcommon.ContractCallInput) (uint64, error) {
	return 0, nil
}

// AddTxValueToAccount mocked method
func (o *OutputContextMock) AddTxValueToAccount(_ []byte, _ *big.Int) {
}

// GetVMOutput mocked method
func (o *OutputContextMock) GetVMOutput() *vmcommon.VMOutput {
	return o.OutputStateMock
}

// RemoveNonUpdatedStorage mocked method
func (o *OutputContextMock) RemoveNonUpdatedStorage() {
}

// DeployCode mocked method
func (o *OutputContextMock) DeployCode(_ vmhost.CodeDeployInput) {
}

// ChangeAccountCode mocked method
func (o *OutputContextMock) ChangeAccountCode(_ []byte, _ []byte) {
}

// SetIsCreatedInTransactionFlag mocked method
func (o *OutputContextMock) SetIsCreatedInTransactionFlag(_ []byte) {
}

// CreateVMOutputInCaseOfError mocked method
func (o *OutputContextMock) CreateVMOutputInCaseOfError(_ error) *vmcommon.VMOutput {
	return o.OutputStateMock
}

// GetCurrentTotalUsedGas mocked method
func (o *OutputContextMock) GetCurrentTotalUsedGas() (uint64, bool) {
	return 0, false
}

// NextOutputTransferIndex mocked method
func (o *OutputContextMock) NextOutputTransferIndex() uint32 {
	o.CrtTransferIndex++
	return o.CrtTransferIndex
}

// NextOutputTransferIndex mocked method
func (o *OutputContextMock) GetCrtTransferIndex() uint32 {
	return o.CrtTransferIndex
}

// NextOutputTransferIndex mocked method
func (o *OutputContextMock) SetCrtTransferIndex(index uint32) {
	o.CrtTransferIndex = index
}

// IsInterfaceNil mocked method
func (o *OutputContextMock) IsInterfaceNil() bool {
	return false
}

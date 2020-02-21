package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.RuntimeContext = (*RuntimeContextMock)(nil)

type RuntimeContextMock struct {
	Err                    error
	VmInput                *vmcommon.VMInput
	SCAddress              []byte
	CallFunction           string
	VmType                 []byte
	ReadOnlyFlag           bool
	CurrentBreakpointValue arwen.BreakpointValue
	PointsUsed             uint64
	InstanceCtxId          int
	MemLoadResult          []byte
	ArgParserMock          arwen.ArgumentsParser
	FailCryptoAPI          bool
	FailElrondAPI          bool
	FailBigIntAPI          bool
	AsyncCallInfo          *arwen.AsyncCallInfo
}

func (r *RuntimeContextMock) InitState() {
}

func (r *RuntimeContextMock) CreateWasmerInstance(contract []byte, gasLimit uint64) error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (r *RuntimeContextMock) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
}

func (r *RuntimeContextMock) PushState() {
}

func (r *RuntimeContextMock) PopState() {
}

func (r *RuntimeContextMock) ClearStateStack() {
}

func (r *RuntimeContextMock) PushInstance() {
}

func (r *RuntimeContextMock) PopInstance() {
}

func (r *RuntimeContextMock) ClearInstanceStack() {
}

func (r *RuntimeContextMock) GetVMType() []byte {
	return r.VmType
}

func (r *RuntimeContextMock) GetVMInput() *vmcommon.VMInput {
	return r.VmInput
}

func (r *RuntimeContextMock) SetVMInput(vmInput *vmcommon.VMInput) {
	r.VmInput = vmInput
}

func (r *RuntimeContextMock) GetSCAddress() []byte {
	return r.SCAddress
}

func (r *RuntimeContextMock) SetSCAddress(scAddress []byte) {
	r.SCAddress = scAddress
}

func (r *RuntimeContextMock) Function() string {
	return r.CallFunction
}

func (r *RuntimeContextMock) Arguments() [][]byte {
	return r.VmInput.Arguments
}

func (r *RuntimeContextMock) SignalExit(_ int) {
}

func (r *RuntimeContextMock) SignalUserError(_ string) {
}

func (r *RuntimeContextMock) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
}

func (r *RuntimeContextMock) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return r.CurrentBreakpointValue
}

func (r *RuntimeContextMock) VerifyContractCode() error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (r *RuntimeContextMock) ArgParser() arwen.ArgumentsParser {
	return r.ArgParserMock
}

func (r *RuntimeContextMock) GetPointsUsed() uint64 {
	return r.PointsUsed
}

func (r *RuntimeContextMock) SetPointsUsed(gasPoints uint64) {
	r.PointsUsed = gasPoints
}

func (r *RuntimeContextMock) ReadOnly() bool {
	return r.ReadOnlyFlag
}

func (r *RuntimeContextMock) SetReadOnly(readOnly bool) {
	r.ReadOnlyFlag = readOnly
}

func (r *RuntimeContextMock) SetInstanceContextId(id int) {
	r.InstanceCtxId = id
}

func (r *RuntimeContextMock) SetInstanceContext(instCtx *wasmer.InstanceContext) {
}

func (r *RuntimeContextMock) GetInstanceContext() *wasmer.InstanceContext {
	return nil
}

func (r *RuntimeContextMock) GetInstanceExports() wasmer.ExportsMap {
	return nil
}

func (r *RuntimeContextMock) CleanInstance() {
}

func (r *RuntimeContextMock) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return nil, nil
}

func (r *RuntimeContextMock) GetInitFunction() wasmer.ExportedFunctionCallback {
	return nil
}

func (r *RuntimeContextMock) MemLoad(offset int32, length int32) ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return r.MemLoadResult, nil
}

func (r *RuntimeContextMock) MemStore(offset int32, data []byte) error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (r *RuntimeContextMock) ElrondAPIErrorShouldFailExecution() bool {
	return r.FailElrondAPI
}

func (r *RuntimeContextMock) CryptoAPIErrorShouldFailExecution() bool {
	return r.FailCryptoAPI
}

func (r *RuntimeContextMock) BigIntAPIErrorShouldFailExecution() bool {
	return r.FailBigIntAPI
}

func (r *RuntimeContextMock) FailExecution(err error) {
}

func (r *RuntimeContextMock) GetAsyncCallInfo() *arwen.AsyncCallInfo {
	return r.AsyncCallInfo
}

func (r *RuntimeContextMock) SetAsyncCallInfo(asyncCallInfo *arwen.AsyncCallInfo) {
	r.AsyncCallInfo = asyncCallInfo
}

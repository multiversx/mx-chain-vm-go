package mock

import (
	"errors"
	"fmt"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// memoryMockPageSize matches MemoryMock's own page size, and is what MemoryGrow reports
// page counts in.
const memoryMockPageSize = 65536

var _ executor.InstanceNew = (*InstanceMock)(nil)
var _ executor.InstanceState = (*InstanceMock)(nil)

type mockMethod func() *InstanceMock

// InstanceMock is a mock for Wasmer instances; it allows creating mock smart
// contracts within tests, without needing actual WASM smart contracts.
type InstanceMock struct {
	Code          []byte
	Exports       map[string]string
	DefaultErrors map[string]error
	Methods       map[string]mockMethod
	Points        uint64
	Data          executor.VMHooks
	GasLimit      uint64
	Memory        executor.Memory
	Host          vmhost.VMHost
	T             testing.TB
	Address       []byte
	AlreadyClean  bool

	// earlyExit is how a mock contract ends other than by returning: out of gas, or any
	// other exit that leaves nothing behind in the VMOutput for Call to read.
	//
	// The same slot Wasmer2InstanceNew keeps, for the same reason. A real contract's exit
	// travels out through the hook that raised it; a mock contract is Go code with no hook
	// to return through, so it puts the exit here.
	earlyExit *executor.EarlyExit
}

// NewInstanceMock creates a new InstanceMock
func NewInstanceMock(code []byte) *InstanceMock {
	return &InstanceMock{
		Code:          code,
		Exports:       make(map[string]string),
		DefaultErrors: make(map[string]error),
		Methods:       make(map[string]mockMethod),
		Points:        0,
		Data:          nil,
		GasLimit:      0,
		Memory:        NewMemoryMock(),
		AlreadyClean:  false,
	}
}

// AddMockMethod adds the provided function as a mocked method to the instance under the specified name.
func (instance *InstanceMock) AddMockMethod(name string, method mockMethod) {
	instance.AddMockMethodWithError(name, method, nil)
}

// AddMockMethodWithError adds the provided function as a mocked method to the instance under the specified name and returns an error
func (instance *InstanceMock) AddMockMethodWithError(name string, method mockMethod, err error) {
	instance.Methods[name] = method
	instance.DefaultErrors[name] = err
	instance.Exports[name] = name
}

var _ executor.VMHooksErrorHandler = (*InstanceMock)(nil)

// HandleHookError ends the mock contract because one of its hooks failed.
//
// The instance's half of what a backend does with a failed hook: a real one halts the
// contract, this one records the exit for Call to report. The VM's half - the return code
// and message in the output state - is applied by whoever has the hooks, exactly as on the
// real path.
func (instance *InstanceMock) HandleHookError(err error) {
	if instance == nil || err == nil {
		return
	}

	if instance.earlyExit == nil {
		instance.earlyExit = executor.EarlyExitOf(err)
	}
}

// EarlyExit is the exit the mock contract has recorded so far, or nil if it has not ended
// other than by returning. For tests that assert on how a hook ended.
func (instance *InstanceMock) EarlyExit() *executor.EarlyExit {
	return instance.earlyExit
}

// SetEarlyExit records how the mock contract is ending.
//
// What a real contract does by returning an error from a VM hook. Mock contracts reach for
// this for the exits that record nothing elsewhere - out of gas above all; an exit that sets
// a return code on the output, as FailExecution and SignalUserError do, is picked up by Call
// without needing this.
func (instance *InstanceMock) SetEarlyExit(earlyExit *executor.EarlyExit) {
	if instance == nil {
		return
	}

	instance.earlyExit = earlyExit
}

// Call runs a mocked contract method and reports how it ended.
//
// A mock contract has no hook to return its exit through, so the exit is recovered from what
// it left behind, in the order a real backend would see it: the explicit exit first, then
// any legacy async call it registered, then whatever it recorded in the VMOutput.
func (instance *InstanceMock) Call(funcName string, pointsLimit uint64) executor.InstanceCallResult {
	method, found := instance.Methods[funcName]
	if !found {
		return executor.CallResultFunctionNotFound()
	}

	instance.earlyExit = nil
	instance.GasLimit = pointsLimit

	returnCodeBefore := instance.Host.Output().ReturnCode()
	method()

	if instance.earlyExit != nil {
		return executor.CallResultEarlyExit(instance.earlyExit, uint64(instance.earlyExit.Code))
	}

	if instance.hasPendingLegacyAsyncCall() {
		exit := vmhost.ExitAsyncCall()
		return executor.CallResultEarlyExit(exit, uint64(exit.Code))
	}

	// FailExecution and SignalUserError record a return code and nothing else. Comparing
	// against the code before the call is what keeps a failure left behind by a nested
	// execution from being read as this contract's own.
	returnCode := instance.Host.Output().ReturnCode()
	if returnCode != returnCodeBefore && returnCode != vmcommon.Ok {
		exit := &executor.EarlyExit{
			Code:    returnCode,
			Message: instance.Host.Output().ReturnMessage(),
			Cause:   errors.New(instance.Host.Output().ReturnMessage()),
		}
		return executor.CallResultEarlyExit(exit, uint64(returnCode))
	}

	if err := instance.DefaultErrors[funcName]; err != nil {
		return executor.CallResultRuntimeError(err)
	}

	return executor.CallResultOk()
}

// hasPendingLegacyAsyncCall reports whether the contract registered a legacy async call it
// is now waiting on, which is exactly the condition the host validates after such an exit.
func (instance *InstanceMock) hasPendingLegacyAsyncCall() bool {
	group, exists := instance.Host.Async().GetCallGroup(vmhost.LegacyAsyncCallGroupID)
	return exists && !group.IsComplete()
}

// GetPointsUsedNew mocked method
func (instance *InstanceMock) GetPointsUsedNew() (uint64, error) {
	return instance.Points, nil
}

// SetPointsUsedNew mocked method
func (instance *InstanceMock) SetPointsUsedNew(points uint64) error {
	instance.Points = points
	return nil
}

// Cache mocked method
func (instance *InstanceMock) Cache() ([]byte, error) {
	return instance.Code, nil
}

// Close mocked method
func (instance *InstanceMock) Close() {
	instance.AlreadyClean = true
}

// IsAlreadyCleaned mocked method
func (instance *InstanceMock) IsAlreadyCleaned() bool {
	return instance.AlreadyClean
}

// Abandon mocked method
func (instance *InstanceMock) Abandon() {
	instance.AlreadyClean = true
}

// ResetNew mocked method
func (instance *InstanceMock) ResetNew() error {
	instance.Points = 0
	instance.earlyExit = nil
	return nil
}

// RequestStop mocked method. Nothing can be running on a mock instance when this is called:
// a mock contract is Go code on the calling goroutine, so by the time a watchdog could ask,
// the call has returned.
func (instance *InstanceMock) RequestStop() {
}

// HasFunction mocked method
func (instance *InstanceMock) HasFunction(name string) bool {
	_, has := instance.Methods[name]
	return has
}

// GetExportedFunctionNames mocked method
func (instance *InstanceMock) GetExportedFunctionNames() []string {
	var functionNames []string
	for functionName := range instance.Exports {
		functionNames = append(functionNames, functionName)
	}
	return functionNames
}

// IsFunctionImported mocked method.
//
// A mock contract imports whatever it exports: it has one name table and no notion of a
// separate import section.
func (instance *InstanceMock) IsFunctionImported(name string) bool {
	_, ok := instance.Exports[name]
	return ok
}

// CheckSignatures mocked method
func (instance *InstanceMock) CheckSignatures() bool {
	return true
}

// MemoryLoad is the strict load from the mocked memory.
func (instance *InstanceMock) MemoryLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return executor.MemoryLoadFromMemory(instance.Memory, memPtr, length)
}

// MemoryLoadClamped is the truncating load from the mocked memory.
func (instance *InstanceMock) MemoryLoadClamped(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return executor.MemLoadFromMemory(instance.Memory, memPtr, length)
}

// MemoryStore stores the given data in the mocked memory at the given offset.
func (instance *InstanceMock) MemoryStore(memPtr executor.MemPtr, data []byte) error {
	return executor.MemStoreToMemory(instance.Memory, memPtr, data)
}

// MemoryLength returns the length of the mocked memory.
func (instance *InstanceMock) MemoryLength() (uint64, error) {
	return uint64(instance.Memory.Length()), nil
}

// MemoryGrow allocates more pages to the mocked memory.
func (instance *InstanceMock) MemoryGrow(byNumPages uint32) (uint32, error) {
	pagesBefore := instance.Memory.Length() / memoryMockPageSize

	err := instance.Memory.Grow(byNumPages)
	if err != nil {
		return 0, err
	}

	return pagesBefore, nil
}

// MemDump yields the entire contents of the memory. Only used in tests.
func (instance *InstanceMock) MemDump() []byte {
	return instance.Memory.Data()
}

// GetMockInstance gets the mock instance from the runtime of the provided host, or nil if
// no contract is running.
//
// Nil rather than a panic, because tests legitimately call VM hooks with no instance at all -
// a hook that only touches the managed types does not need one - and those tests have to be
// able to ask for the state to pass in. Every method here is safe on a nil receiver, and the
// hooks refuse any memory operation when their state is nil.
func GetMockInstance(host vmhost.VMHost) *InstanceMock {
	instance, _ := host.Runtime().GetInstance().(*InstanceMock)
	return instance
}

// ID returns an identifier for the instance, unique at runtime
func (instance *InstanceMock) ID() string {
	return fmt.Sprintf("%p", instance)
}

// IsInterfaceNil mocked method
func (instance *InstanceMock) IsInterfaceNil() bool {
	return instance == nil
}

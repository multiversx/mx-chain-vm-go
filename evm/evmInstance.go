package evm

import (
	"fmt"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	interpreter "github.com/multiversx/mx-chain-vm-go/evm/interpreter"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const NoBreakpoint = uint64(vmhost.BreakpointNone)

var _ executor.Instance = (*EVMInstance)(nil)

// EVMInstance represents a EVM instance.
type EVMInstance struct {
	evmExecutor *EVMExecutor
	evm         *interpreter.EVM
	options     executor.CompilationOptions

	isCompiled bool
	wasCleaned bool
	code       []byte
	gasUsed    uint64
	breakpoint uint64
}

func newInstance(
	evmExecutor *EVMExecutor,
	isCompiled bool,
	code []byte,
	options executor.CompilationOptions,
) (*EVMInstance, error) {
	instance := &EVMInstance{
		evmExecutor: evmExecutor,
		options:     options,

		isCompiled: isCompiled,
		wasCleaned: false,
		code:       code,
		gasUsed:    0,
		breakpoint: NoBreakpoint,
	}

	instance.resetEVM()
	return instance, nil
}

// Clean cleans instance
func (instance *EVMInstance) Clean() bool {
	logEVM.Trace("clean: start", "id", instance.ID())
	if instance.wasCleaned {
		logEVM.Trace("clean: was cleaned", "id", instance.ID())
		return false
	}

	instance.evm = nil
	instance.evmExecutor = nil
	instance.options = executor.CompilationOptions{}

	instance.isCompiled = false
	instance.wasCleaned = true
	instance.code = []byte{}
	instance.gasUsed = 0
	instance.breakpoint = NoBreakpoint

	logEVM.Trace("clean: end", "id", instance.ID())
	return true
}

// IsAlreadyCleaned returns the internal field AlreadyClean
func (instance *EVMInstance) IsAlreadyCleaned() bool {
	return instance.wasCleaned
}

// SetGasLimit sets the gas limit for the instance
func (instance *EVMInstance) SetGasLimit(uint64) {
}

// SetPointsUsed sets the internal instance gas counter
func (instance *EVMInstance) SetPointsUsed(points uint64) {
	if !instance.options.Metering {
		return
	}

	instance.gasUsed = points
}

// GetPointsUsed returns the internal instance gas counter
func (instance *EVMInstance) GetPointsUsed() uint64 {
	return instance.gasUsed
}

// SetBreakpointValue sets the breakpoint value for the instance
func (instance *EVMInstance) SetBreakpointValue(value uint64) {
	if !instance.options.RuntimeBreakpoints {
		return
	}

	instance.breakpoint = value
	if value != NoBreakpoint {
		instance.evm.Cancel()
	}
}

// GetBreakpointValue returns the breakpoint value
func (instance *EVMInstance) GetBreakpointValue() uint64 {
	return instance.breakpoint
}

// HasCompiledCode specifies if the code is compiled
func (instance *EVMInstance) HasCompiledCode() bool {
	return instance.isCompiled
}

// Cache caches the instance
func (instance *EVMInstance) Cache() ([]byte, error) {
	if !instance.isCompiled {
		return nil, ErrCodeNotCompiled
	}
	return instance.code, nil
}

// IsFunctionImported returns true if the instance imports the specified function
func (instance *EVMInstance) IsFunctionImported(string) bool {
	return false
}

// CallFunction executes given function from loaded contract.
func (instance *EVMInstance) CallFunction(functionName string) error {
	err := instance.prepareAddress(functionName)
	if err != nil {
		return err
	}

	contract, input := instance.prepareContract(functionName)
	returnData, err := instance.evm.Interpreter().Run(contract, input, instance.evmExecutor.evmHooks.ReadOnly())
	if err != nil {
		return err
	}
	if instance.evm.Cancelled() {
		return ErrExecutionAborted
	}

	instance.consumeOutput(functionName, returnData)
	return nil
}

// HasFunction checks if loaded contract has a function (endpoint) with given name.
func (instance *EVMInstance) HasFunction(functionName string) bool {
	switch functionName {
	case vmhost.InitFunctionName:
		return true
	case vmhost.UpgradeFunctionName, vmhost.DeleteFunctionName, vmhost.CallbackFunctionName, vmhost.ContractsUpgradeFunctionName:
		return false
	default:
		return parsers.EVMSelectorSize == len(functionNameToSelector(functionName))
	}
}

// GetFunctionNames returns a list of the function names exported by the contract.
func (instance *EVMInstance) GetFunctionNames() []string {
	return []string{}
}

// ValidateFunctionArities checks that no function (endpoint) of the given contract has any parameters or returns any result.
// All arguments and results should be transferred via the import functions.
func (instance *EVMInstance) ValidateFunctionArities() error {
	return nil
}

// HasMemory checks whether the instance has at least one exported memory.
func (instance *EVMInstance) HasMemory() bool {
	return true
}

// MemLoad returns the contents from the given offset of the EVM memory.
func (instance *EVMInstance) MemLoad(executor.MemPtr, executor.MemLength) ([]byte, error) {
	return nil, ErrActionNotSupported
}

// MemStore stores the given data in the EVM memory at the given offset.
func (instance *EVMInstance) MemStore(executor.MemPtr, []byte) error {
	return ErrActionNotSupported
}

// MemLength returns the length of the allocated memory. Only called directly in tests.
func (instance *EVMInstance) MemLength() uint32 {
	return 0
}

// MemGrow allocates more pages to the current memory. Only called directly in tests.
func (instance *EVMInstance) MemGrow(uint32) error {
	return ErrActionNotSupported
}

// MemDump yields the entire contents of the memory. Only used in tests.
func (instance *EVMInstance) MemDump() []byte {
	return []byte{}
}

// ID Id returns an identifier for the instance, unique at runtime
func (instance *EVMInstance) ID() string {
	return fmt.Sprintf("%p", instance)
}

// Reset resets the instance memories and globals
func (instance *EVMInstance) Reset() bool {
	if instance.wasCleaned {
		logEVM.Trace("reset: was cleaned", "id", instance.ID())
		return false
	}

	instance.resetEVM()
	instance.gasUsed = 0
	instance.breakpoint = NoBreakpoint

	logEVM.Trace("reset: warm instance", "id", instance.ID())
	return true
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *EVMInstance) IsInterfaceNil() bool {
	return instance == nil
}

// SetVMHooksPtr sets the VM hooks pointer
func (instance *EVMInstance) SetVMHooksPtr(uintptr) {
}

// GetVMHooksPtr returns the VM hooks pointer
func (instance *EVMInstance) GetVMHooksPtr() uintptr {
	return uintptr(0)
}

package wasmer2

import (
	"unsafe"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.Executor = (*Wasmer2Executor)(nil)

// Wasmer2Executor oversees the creation of Wasmer instances and execution.
type Wasmer2Executor struct {
	cgoExecutor *cWasmerExecutorT

	vmHookPointers *cWasmerVmHookPointers
	vmHooks        executor.VMHooks
	vmHooksPtr     uintptr
	vmHooksPtrPtr  unsafe.Pointer

	opcodeCost *OpcodeCost
}

// CreateExecutor creates a new wasmer executor.
func CreateExecutor() (*Wasmer2Executor, error) {
	vmHookPointers := populateCgoFunctionPointers()
	localPtr := uintptr(unsafe.Pointer(vmHookPointers))
	localPtrPtr := unsafe.Pointer(&localPtr)

	var cExecutor *cWasmerExecutorT

	var result = cWasmerNewExecutor(
		&cExecutor,
		localPtrPtr,
	)

	if result != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	cWasmerForceInstallSighandlers()

	wasmerExecutor := &Wasmer2Executor{
		cgoExecutor:    cExecutor,
		vmHookPointers: vmHookPointers,
	}

	return wasmerExecutor, nil
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (wasmerExecutor *Wasmer2Executor) SetOpcodeCosts(wasmOps *executor.WASMOpcodeCost) {
	// extract only wasmer2 opcodes
	wasmerExecutor.opcodeCost = wasmerExecutor.extractOpcodeCost(wasmOps)
	cWasmerExecutorSetOpcodeCost(
		wasmerExecutor.cgoExecutor,
		(*cWasmerOpcodeCostT)(unsafe.Pointer(wasmerExecutor.opcodeCost)),
	)
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (wasmerExecutor *Wasmer2Executor) SetRkyvSerializationEnabled(_ bool) {
}

// SetSIGSEGVPassthrough controls a Wasmer flag.
func (wasmerExecutor *Wasmer2Executor) SetSIGSEGVPassthrough() {
}

func (wasmerExecutor *Wasmer2Executor) FunctionNames() vmcommon.FunctionNames {
	return functionNames
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (wasmerExecutor *Wasmer2Executor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	var cInstance *cWasmerInstanceT

	if len(contractCode) == 0 {
		return nil, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var compileResult = cWasmerInstantiateWithOptions(
		wasmerExecutor.cgoExecutor,
		&cInstance,
		(*cUchar)(unsafe.Pointer(&contractCode[0])),
		cUint(len(contractCode)),
		(*cWasmerCompilationOptions)(cOptions),
	)

	if compileResult != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	return newInstance(cInstance)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (wasmerExecutor *Wasmer2Executor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	var cInstance *cWasmerInstanceT

	if len(compiledCode) == 0 {
		return nil, newWrappedError(ErrInvalidBytecode)
	}

	cOptions := unsafe.Pointer(&options)
	var compileResult = cWasmerInstanceFromCache(
		wasmerExecutor.cgoExecutor,
		&cInstance,
		(*cUchar)(unsafe.Pointer(&compiledCode[0])),
		cUint32T(len(compiledCode)),
		(*cWasmerCompilationOptions)(cOptions),
	)

	if compileResult != cWasmerOk {
		return nil, newWrappedError(ErrFailedInstantiation)
	}

	return newInstance(cInstance)
}

// IsInterfaceNil returns true if underlying object is nil
func (wasmerExecutor *Wasmer2Executor) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}

// InitVMHooks inits the VM hooks
func (wasmerExecutor *Wasmer2Executor) initVMHooks(vmHooks executor.VMHooks) {
	wasmerExecutor.vmHooks = vmHooks
	localPtr := uintptr(unsafe.Pointer(&wasmerExecutor.vmHooks))
	wasmerExecutor.vmHooksPtr = localPtr
	wasmerExecutor.vmHooksPtrPtr = unsafe.Pointer(&localPtr)
	cWasmerExecutorContextDataSet(wasmerExecutor.cgoExecutor, wasmerExecutor.vmHooksPtrPtr)
}

func (wasmerExecutor *Wasmer2Executor) extractOpcodeCost(wasmOps *executor.WASMOpcodeCost) *OpcodeCost {
	return &OpcodeCost{
		Block:              wasmOps.Block,
		Br:                 wasmOps.Br,
		BrIf:               wasmOps.BrIf,
		BrTable:            wasmOps.BrTable,
		Call:               wasmOps.Call,
		CallIndirect:       wasmOps.CallIndirect,
		Catch:              wasmOps.Catch,
		CatchAll:           wasmOps.CatchAll,
		Delegate:           wasmOps.Delegate,
		Drop:               wasmOps.Drop,
		Else:               wasmOps.Else,
		End:                wasmOps.End,
		GlobalGet:          wasmOps.GlobalGet,
		GlobalSet:          wasmOps.GlobalSet,
		I32Add:             wasmOps.I32Add,
		I32And:             wasmOps.I32And,
		I32Clz:             wasmOps.I32Clz,
		I32Const:           wasmOps.I32Const,
		I32Ctz:             wasmOps.I32Ctz,
		I32DivS:            wasmOps.I32DivS,
		I32DivU:            wasmOps.I32DivU,
		I32Eq:              wasmOps.I32Eq,
		I32Eqz:             wasmOps.I32Eqz,
		I32Extend16S:       wasmOps.I32Extend16S,
		I32Extend8S:        wasmOps.I32Extend8S,
		I32GeS:             wasmOps.I32GeS,
		I32GeU:             wasmOps.I32GeU,
		I32GtS:             wasmOps.I32GtS,
		I32GtU:             wasmOps.I32GtU,
		I32LeS:             wasmOps.I32LeS,
		I32LeU:             wasmOps.I32LeU,
		I32Load:            wasmOps.I32Load,
		I32Load16S:         wasmOps.I32Load16S,
		I32Load16U:         wasmOps.I32Load16U,
		I32Load8S:          wasmOps.I32Load8S,
		I32Load8U:          wasmOps.I32Load8U,
		I32LtS:             wasmOps.I32LtS,
		I32LtU:             wasmOps.I32LtU,
		I32Mul:             wasmOps.I32Mul,
		I32Ne:              wasmOps.I32Ne,
		I32Or:              wasmOps.I32Or,
		I32Popcnt:          wasmOps.I32Popcnt,
		I32RemS:            wasmOps.I32RemS,
		I32RemU:            wasmOps.I32RemU,
		I32Rotl:            wasmOps.I32Rotl,
		I32Rotr:            wasmOps.I32Rotr,
		I32Shl:             wasmOps.I32Shl,
		I32ShrS:            wasmOps.I32ShrS,
		I32ShrU:            wasmOps.I32ShrU,
		I32Store:           wasmOps.I32Store,
		I32Store16:         wasmOps.I32Store16,
		I32Store8:          wasmOps.I32Store8,
		I32Sub:             wasmOps.I32Sub,
		I32WrapI64:         wasmOps.I32WrapI64,
		I32Xor:             wasmOps.I32Xor,
		I64Add:             wasmOps.I64Add,
		I64And:             wasmOps.I64And,
		I64Clz:             wasmOps.I64Clz,
		I64Const:           wasmOps.I64Const,
		I64Ctz:             wasmOps.I64Ctz,
		I64DivS:            wasmOps.I64DivS,
		I64DivU:            wasmOps.I64DivU,
		I64Eq:              wasmOps.I64Eq,
		I64Eqz:             wasmOps.I64Eqz,
		I64Extend16S:       wasmOps.I64Extend16S,
		I64Extend32S:       wasmOps.I64Extend32S,
		I64Extend8S:        wasmOps.I64Extend8S,
		I64ExtendI32S:      wasmOps.I64ExtendI32S,
		I64ExtendI32U:      wasmOps.I64ExtendI32U,
		I64GeS:             wasmOps.I64GeS,
		I64GeU:             wasmOps.I64GeU,
		I64GtS:             wasmOps.I64GtS,
		I64GtU:             wasmOps.I64GtU,
		I64LeS:             wasmOps.I64LeS,
		I64LeU:             wasmOps.I64LeU,
		I64Load:            wasmOps.I64Load,
		I64Load16S:         wasmOps.I64Load16S,
		I64Load16U:         wasmOps.I64Load16U,
		I64Load32S:         wasmOps.I64Load32S,
		I64Load32U:         wasmOps.I64Load32U,
		I64Load8S:          wasmOps.I64Load8S,
		I64Load8U:          wasmOps.I64Load8U,
		I64LtS:             wasmOps.I64LtS,
		I64LtU:             wasmOps.I64LtU,
		I64Mul:             wasmOps.I64Mul,
		I64Ne:              wasmOps.I64Ne,
		I64Or:              wasmOps.I64Or,
		I64Popcnt:          wasmOps.I64Popcnt,
		I64RemS:            wasmOps.I64RemS,
		I64RemU:            wasmOps.I64RemU,
		I64Rotl:            wasmOps.I64Rotl,
		I64Rotr:            wasmOps.I64Rotr,
		I64Shl:             wasmOps.I64Shl,
		I64ShrS:            wasmOps.I64ShrS,
		I64ShrU:            wasmOps.I64ShrU,
		I64Store:           wasmOps.I64Store,
		I64Store16:         wasmOps.I64Store16,
		I64Store32:         wasmOps.I64Store32,
		I64Store8:          wasmOps.I64Store8,
		I64Sub:             wasmOps.I64Sub,
		I64Xor:             wasmOps.I64Xor,
		If:                 wasmOps.If,
		LocalGet:           wasmOps.LocalGet,
		LocalSet:           wasmOps.LocalSet,
		LocalTee:           wasmOps.LocalTee,
		LocalAllocate:      wasmOps.LocalAllocate,
		Loop:               wasmOps.Loop,
		MemoryGrow:         wasmOps.MemoryGrow,
		MemorySize:         wasmOps.MemorySize,
		Nop:                wasmOps.Nop,
		RefFunc:            wasmOps.RefFunc,
		RefIsNull:          wasmOps.RefIsNull,
		RefNull:            wasmOps.RefNull,
		Rethrow:            wasmOps.Rethrow,
		Return:             wasmOps.Return,
		ReturnCall:         wasmOps.ReturnCall,
		ReturnCallIndirect: wasmOps.ReturnCallIndirect,
		Select:             wasmOps.Select,
		TableGet:           wasmOps.TableGet,
		TableGrow:          wasmOps.TableGrow,
		TableInit:          wasmOps.TableInit,
		TableSet:           wasmOps.TableSet,
		TableSize:          wasmOps.TableSize,
		Throw:              wasmOps.Throw,
		Try:                wasmOps.Try,
		TypedSelect:        wasmOps.TypedSelect,
		Unreachable:        wasmOps.Unreachable,
		Unwind:             wasmOps.Unwind,
	}
}

package subcontexts

import (
	"fmt"
	"math/big"
	"unsafe"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Runtime struct {
	host            arwen.VMContext
	instance        *wasmer.Instance
	instanceContext *wasmer.InstanceContext
	vmInput         *vmcommon.VMInput
	scAddress       []byte
	callFunction    string
	vmType          []byte
	readOnly        bool
	ethInput        []byte
	stateStack      []*Runtime
}

func NewRuntimeSubcontext(host arwen.VMContext, blockChainHook vmcommon.BlockchainHook) (*Runtime, error) {
	runtime := &Runtime{
		host:           host,
		stateStack:     make([]*Runtime, 0),
	}

	runtime.InitState()

	return runtime, nil
}

func (runtime *Runtime) InitState() {
	runtime.vmInput = &vmcommon.VMInput{}
	runtime.scAddress = make([]byte, 0)
	runtime.callFunction = ""
	runtime.ethInput = nil
	runtime.readOnly = false
}

func (runtime *Runtime) CreateWasmerInstance(contract []byte) error {
	var err error
	runtime.instance, err = wasmer.NewMeteredInstance(contract, runtime.vmInput.GasProvided)
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	return err
}

func (runtime *Runtime) InitializeFromInput(input *vmcommon.ContractCallInput) error {
	runtime.vmInput = &input.VMInput
	runtime.scAddress = input.RecipientAddr
	runtime.callFunction = input.Function

	contract, err := runtime.host.Blockchain().GetCode(runtime.scAddress)
	gasProvided := runtime.vmInput.GasProvided
	instance, err := wasmer.NewMeteredInstance(contract, gasProvided)
	if err != nil {
		return err
	}
	runtime.instance = instance
	return nil
}

func (runtime *Runtime) PushState() {
	newState := &Runtime{
		vmInput:      runtime.vmInput,
		scAddress:    runtime.scAddress,
		callFunction: runtime.callFunction,
		readOnly:     runtime.readOnly,
	}

	runtime.stateStack = append(runtime.stateStack, newState)
}

func (runtime *Runtime) PopState() error {
	stateStackLen := len(runtime.stateStack)
	if stateStackLen < 1 {
		return StateStackUnderflow
	}

	prevState := runtime.stateStack[stateStackLen-1]
	runtime.stateStack = runtime.stateStack[:stateStackLen-1]

	runtime.vmInput = prevState.vmInput
	runtime.scAddress = prevState.scAddress
	runtime.callFunction = prevState.callFunction
	runtime.readOnly = prevState.readOnly

	return nil
}

func (runtime *Runtime) LoadFromStateCopy(otherRuntime *Runtime) {
	runtime.vmInput = otherRuntime.vmInput
	runtime.scAddress = otherRuntime.scAddress
	runtime.callFunction = otherRuntime.callFunction
	runtime.readOnly = otherRuntime.readOnly
}

func (runtime *Runtime) GetVMInput() *vmcommon.VMInput {
	return runtime.vmInput
}

func (runtime *Runtime) GetSCAddress() []byte {
	return runtime.scAddress
}

func (runtime *Runtime) Function() string {
	return runtime.callFunction
}

func (runtime *Runtime) Arguments() [][]byte {
	return runtime.vmInput.Arguments
}

func (runtime *Runtime) SignalUserError() {
	// SignalUserError() remains in Runtime, and won't be moved into Output,
	// because there will be extra handling added here later, which requires
	// information from Runtime (e.g. runtime breakpoints)
	runtime.host.Output().SetReturnCode(vmcommon.UserError)
}

func (runtime *Runtime) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	runtime.instance.SetBreakpointValue(uint64(value))
}

func (runtime *Runtime) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return arwen.BreakpointValue(runtime.instance.GetBreakpointValue())
}

func (runtime *Runtime) GetPointsUsed() uint64 {
	return runtime.instance.GetPointsUsed()
}

func (runtime *Runtime) SetPointsUsed(gasPoints uint64) {
	runtime.instance.SetPointsUsed(gasPoints)
}

func (runtime *Runtime) CallData() []byte {
	panic("not implemented")
}

func (runtime *Runtime) ReadOnly() bool {
	return runtime.readOnly
}

func (runtime *Runtime) SetReadOnly(readOnly bool) {
	runtime.readOnly = readOnly
}

func (runtime *Runtime) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	panic("not implemented")
}

func (runtime *Runtime) ExecuteOnDestContext(input *vmcommon.ContractCallInput) error {
	panic("not implemented")
}

func (runtime *Runtime) SetInstanceContextId(id int) {
	runtime.instance.SetContextData(unsafe.Pointer(&id))
}

func (runtime *Runtime) SetInstanceContext(instCtx *wasmer.InstanceContext) {
	runtime.instanceContext = instCtx
}

func (runtime *Runtime) GetInstanceContext() *wasmer.InstanceContext {
	return runtime.instanceContext
}

func (runtime *Runtime) MemStore(offset int32, data []byte) error {
	memory := runtime.instanceContext.Memory()
	return arwen.StoreBytes(memory, offset, data)
}

func (runtime *Runtime) MemLoad(offset int32, length int32) ([]byte, error) {
	memory := runtime.instanceContext.Memory()
	return arwen.LoadBytes(memory, offset, length)
}

func (runtime *Runtime) Clean() {
	runtime.instance.Clean()
}

func (runtime *Runtime) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	if runtime.readOnly {
		return nil, ErrInvalidCallOnReadOnlyMode
	}

	currVmInput := runtime.vmInput
	currScAddress := runtime.scAddress
	currCallFunction := runtime.callFunction

	defer func() {
		runtime.vmInput = currVmInput
		runtime.scAddress = currScAddress
		runtime.callFunction = currCallFunction
	}()

	runtime.vmInput = &input.VMInput
	nonce := runtime.host.Blockchain().GetNonce(input.CallerAddr)
	address, err := runtime.host.Blockchain().NewAddress(input.CallerAddr, nonce, runtime.vmType)
	if err != nil {
		return nil, err
	}

	runtime.host.Output().Transfer(address, input.CallerAddr, 0, input.CallValue, nil)
	runtime.host.Blockchain().IncreaseNonce(input.CallerAddr)
	runtime.scAddress = address

	totalGasConsumed := input.GasProvided
	defer func() {
		runtime.host.Metering().UseGas(totalGasConsumed)
	}()

	gasLeft, err := runtime.deductInitialCodeCost(
		input.GasProvided,
		input.ContractCode,
		0, // create cost was elrady taken care of. as it is different for ethereum and elrond
		runtime.host.Metering().GasSchedule().BaseOperationCost.StorePerByte,
	)
	if err != nil {
		return nil, err
	}

	newInstance, err := wasmer.NewMeteredInstance(input.ContractCode, gasLeft)
	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		return nil, err
	}

	idContext := arwen.AddHostContext(runtime.host)
	oldInstance := runtime.instance
	runtime.instance = newInstance
	defer func() {
		runtime.instance = oldInstance
		newInstance.Clean()
		arwen.RemoveHostContext(idContext)
	}()

	runtime.instance.SetContextData(unsafe.Pointer(&idContext))

	initCalled, result, err := runtime.callInitFunction()
	if err != nil {
		return nil, err
	}

	if initCalled {
		runtime.host.Output().Finish(result)
	}

	outputAccounts := runtime.host.Output().GetOutputAccounts()
	newSCAcc, ok := outputAccounts[string(address)]
	if !ok {
		outputAccounts[string(address)] = &vmcommon.OutputAccount{
			Address:        address,
			Nonce:          0,
			BalanceDelta:   big.NewInt(0),
			StorageUpdates: nil,
			Code:           input.ContractCode,
		}
	} else {
		newSCAcc.Code = input.ContractCode
	}

	totalGasConsumed = input.GasProvided - gasLeft - newInstance.GetPointsUsed()

	return address, nil
}

func (runtime *Runtime) deductInitialCodeCost(
	gasProvided uint64,
	code []byte,
	baseCost uint64,
	costPerByte uint64,
) (uint64, error) {
	codeLength := uint64(len(code))
	codeCost := codeLength * costPerByte
	initialCost := baseCost + codeCost

	if initialCost > gasProvided {
		return 0, ErrNotEnoughGas
	}

	return gasProvided - initialCost, nil
}

func (runtime *Runtime) callInitFunction() (bool, []byte, error) {
	init, ok := runtime.instance.Exports[arwen.InitFunctionName]

	if !ok {
		init, ok = runtime.instance.Exports[arwen.InitFunctionNameEth]
	}

	if !ok {
		// There's no initialization function, don't do anything.
		return false, nil, nil
	}

	out, err := init()
	if err != nil {
		fmt.Println("arwen.callInitFunction() error:", err.Error())
		return true, nil, err
	}

	convertedResult := arwen.ConvertReturnValue(out)
	result := convertedResult.Bytes()
	return true, result, nil
}


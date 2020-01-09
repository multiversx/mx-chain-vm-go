package context

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"unsafe"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/context/subcontexts"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/ethapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)


// vmContext implements HostContext interface.
type vmContext struct {
	BigIntContainer
	blockChainHook vmcommon.BlockchainHook
	cryptoHook     vmcommon.CryptoHook
	instance       *wasmer.Instance

	vmInput vmcommon.VMInput

	vmType       []byte
	callFunction string
	scAddress    []byte

	selfDestruct  map[string][]byte
	ethInput      []byte
	blockGasLimit uint64
	refund        uint64

	gasCostConfig *config.GasCost

	// -- refactored subcontexts
	blockchainSubcontext arwen.BlockchainSubcontext
	runtimeSubcontext    arwen.RuntimeSubcontext
	outputSubcontext     arwen.OutputSubcontext
	meteringSubcontext   arwen.MeteringSubcontext
	storageSubcontext    arwen.StorageSubcontext
	bigIntSubcontext     arwen.BigIntSubcontext
}

func NewArwenVM(
	blockChainHook vmcommon.BlockchainHook,
	cryptoHook vmcommon.CryptoHook,
	vmType []byte,
	blockGasLimit uint64,
	gasSchedule map[string]map[string]uint64,
) (*vmContext, error) {

	imports, err := elrondapi.ElrondEImports()
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.BigIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = ethapi.EthereumImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = crypto.CryptoImports(imports)
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(gasSchedule)
	if err != nil {
		return nil, err
	}

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()


	host := &vmContext{
		BigIntContainer: NewBigIntContainer(),
		blockChainHook:  blockChainHook,
		cryptoHook:      cryptoHook,
		vmType:          vmType,
		blockGasLimit:   blockGasLimit,
		gasCostConfig:   gasCostConfig,

		meteringSubcontext:   nil,
		runtimeSubcontext:    nil,
		blockchainSubcontext: nil,
		storageSubcontext: nil,
	}

	host.blockchainSubcontext, err = subcontexts.NewBlockchainSubcontext(blockChainHook, host)
	if err != nil {
		return nil, err
	}

	host.runtimeSubcontext, err = subcontexts.NewRuntimeSubcontext(blockChainHook)
	if err != nil {
		return nil, err
	}

	host.meteringSubcontext, err = subcontexts.NewMeteringSubcontext(gasSchedule, blockGasLimit, host)
	if err != nil {
		return nil, err
	}

	host.outputSubcontext, err = subcontexts.NewOutputSubcontext(host)
	if err != nil {
		return nil, err
	}

	host.storageSubcontext, err = subcontexts.NewStorageSubcontext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.bigIntSubcontext, err = subcontexts.NewBigIntSubcontext()
	if err != nil {
		return nil, err
	}

	host.initInternalValues()

	err = wasmer.SetImports(imports)
	if err != nil {
		return nil, err
	}
	wasmer.SetOpcodeCosts(&opcodeCosts)

	return host, nil
}

func (host *vmContext) Crypto() vmcommon.CryptoHook {
	return host.cryptoHook
}

func (host *vmContext) Blockchain() arwen.BlockchainSubcontext {
	return host.blockchainSubcontext
}

func (host *vmContext) Runtime() arwen.RuntimeSubcontext {
	return host.runtimeSubcontext
}

func (host *vmContext) Output() arwen.OutputSubcontext {
	return host.outputSubcontext
}

func (host *vmContext) Metering() arwen.MeteringSubcontext {
	return host.meteringSubcontext
}

func (host *vmContext) Storage() arwen.StorageSubcontext {
	return host.storageSubcontext
}

func (host *vmContext) BigInt() arwen.BigIntSubcontext {
	return host.bigIntSubcontext
}

func (host *vmContext) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	try := func() {
		vmOutput, err = host.doRunSmartContractCreate(input)
	}

	catch := func(caught error) {
		err = caught
	}

	arwen.TryCatch(try, catch, "arwen.RunSmartContractCreate")
	return
}

func (host *vmContext) doRunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	host.initInternalValues()
	host.vmInput = input.VMInput

	nonce, err := host.blockChainHook.GetNonce(input.CallerAddr)
	if err != nil {
		return nil, err
	}

	if nonce > 0 {
		nonce -= 1
	}

	address, err := host.blockChainHook.NewAddress(input.CallerAddr, nonce, host.vmType)
	if err != nil {
		return nil, err
	}

	host.scAddress = address
	host.addTxValueToSmartContract(input.CallValue, address)

	host.vmInput.GasProvided, err = host.deductInitialCodeCost(
		input.GasProvided,
		input.ContractCode,
		host.GasSchedule().ElrondAPICost.CreateContract,
		host.GasSchedule().BaseOperationCost.StorePerByte,
	)
	if err != nil {
		return host.createVMOutputInCaseOfError(vmcommon.OutOfGas), nil
	}

	err = host.createMeteredWasmerInstance(input.ContractCode)

	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid), nil
	}

	idContext := arwen.AddHostContext(host)
	defer func() {
		arwen.RemoveHostContext(idContext)
		host.instance.Clean()
	}()

	host.instance.SetContextData(unsafe.Pointer(&idContext))

	_, result, err := host.callInitFunction()
	if err != nil {
		return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature), nil
	}

	outputAccounts := host.Output().GetOutputAccounts()
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

	vmOutput := host.createVMOutput(result)

	return vmOutput, err
}

func (host *vmContext) deductInitialCodeCost(
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

func (host *vmContext) callInitFunction() (bool, []byte, error) {
	init, ok := host.instance.Exports[arwen.InitFunctionName]

	if !ok {
		init, ok = host.instance.Exports[arwen.InitFunctionNameEth]
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

func (host *vmContext) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	try := func() {
		vmOutput, err = host.doRunSmartContractCall(input)
	}

	catch := func(caught error) {
		err = caught
	}

	arwen.TryCatch(try, catch, "arwen.RunSmartContractCall")
	return
}

func (host *vmContext) doRunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.initInternalValues()
	host.vmInput = input.VMInput
	host.scAddress = input.RecipientAddr
	host.callFunction = input.Function

	host.addTxValueToSmartContract(input.CallValue, input.RecipientAddr)

	contract, err := host.GetCode(host.scAddress)
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid), nil
	}

	host.vmInput.GasProvided, err = host.deductInitialCodeCost(
		input.GasProvided,
		contract,
		0,
		host.GasSchedule().BaseOperationCost.CompilePerByte,
	)
	if err != nil {
		return host.createVMOutputInCaseOfError(vmcommon.OutOfGas), nil
	}

	err = host.createMeteredWasmerInstance(contract)

	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid), nil
	}

	idContext := arwen.AddHostContext(host)

	defer func() {
		host.instance.Clean()
		arwen.RemoveHostContext(idContext)
	}()

	host.instance.SetContextData(unsafe.Pointer(&idContext))

	if host.isInitFunctionCalled() {
		fmt.Println("arwen Error", ErrInitFuncCalledInRun.Error())
		return host.createVMOutputInCaseOfError(vmcommon.UserError), nil
	}

	function, err := host.getFunctionToCall()
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.FunctionNotFound), nil
	}

	result, err := function()
	if err != nil {
		strError, _ := wasmer.GetLastError()

		fmt.Println("arwen Error", err.Error(), strError)
		return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature), nil
	}

	returnCode := host.Output().ReturnCode()
	if returnCode != vmcommon.Ok {
		// user error: signalError()
		return host.createVMOutputInCaseOfError(returnCode), nil
	}

	convertedResult := arwen.ConvertReturnValue(result)
	vmOutput := host.createVMOutput(convertedResult.Bytes())

	return vmOutput, nil
}

func (host *vmContext) createMeteredWasmerInstance(contract []byte) error {
	var err error
	host.instance, err = wasmer.NewMeteredInstance(contract, host.vmInput.GasProvided)
	host.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	return err
}

func (host *vmContext) isInitFunctionCalled() bool {
	return host.callFunction == arwen.InitFunctionName || host.callFunction == arwen.InitFunctionNameEth
}

func (host *vmContext) createVMOutputInCaseOfError(errCode vmcommon.ReturnCode) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{GasRemaining: 0, GasRefund: big.NewInt(0)}
	vmOutput.ReturnCode = errCode
	return vmOutput
}

func (host *vmContext) getFunctionToCall() (func(...interface{}) (wasmer.Value, error), error) {
	exports := host.instance.Exports
	function, ok := exports[host.callFunction]
	if !ok {
		function, ok = exports["main"]
	}

	if !ok {
		return nil, ErrFuncNotFound
	}

	return function, nil
}

func (host *vmContext) initInternalValues() {
	host.BigInt().InitState()
	host.selfDestruct = make(map[string][]byte)
	host.vmInput = vmcommon.VMInput{}
	host.outputAccounts = make(map[string]*vmcommon.OutputAccount, 0)
	host.scAddress = make([]byte, 0)
	host.callFunction = ""
	host.returnData = nil
	host.returnCode = vmcommon.Ok
	host.ethInput = nil
	host.readOnly = false
	host.refund = 0
}

func (host *vmContext) addTxValueToSmartContract(value *big.Int, scAddress []byte) {
	destAcc, ok := host.outputAccounts[string(scAddress)]
	if !ok {
		destAcc = &vmcommon.OutputAccount{
			Address:      scAddress,
			BalanceDelta: big.NewInt(0),
		}
		host.outputAccounts[string(destAcc.Address)] = destAcc
	}

	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

func (host *vmContext) CallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmContext) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	if host.Runtime().ReadOnly() {
		return nil, ErrInvalidCallOnReadOnlyMode
	}

	currVmInput := host.vmInput
	currScAddress := host.scAddress
	currCallFunction := host.callFunction

	defer func() {
		host.vmInput = currVmInput
		host.scAddress = currScAddress
		host.callFunction = currCallFunction
	}()

	host.vmInput = input.VMInput
	nonce := host.Blockchain().GetNonce(input.CallerAddr)
	address, err := host.blockChainHook.NewAddress(input.CallerAddr, nonce, host.vmType)
	if err != nil {
		return nil, err
	}

	host.Output().Transfer(address, input.CallerAddr, 0, input.CallValue, nil)
	host.Blockchain().IncreaseNonce(input.CallerAddr)
	host.scAddress = address

	totalGasConsumed := input.GasProvided
	defer func() {
		host.Metering().UseGas(totalGasConsumed)
	}()

	gasLeft, err := host.deductInitialCodeCost(
		input.GasProvided,
		input.ContractCode,
		0, // create cost was elrady taken care of. as it is different for ethereum and elrond
		host.Metering().GasSchedule().BaseOperationCost.StorePerByte,
	)
	if err != nil {
		return nil, err
	}

	newInstance, err := wasmer.NewMeteredInstance(input.ContractCode, gasLeft)
	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	oldInstance := host.instance
	host.instance = newInstance
	defer func() {
		host.instance = oldInstance
		newInstance.Clean()
		arwen.RemoveHostContext(idContext)
	}()

	host.instance.SetContextData(unsafe.Pointer(&idContext))

	initCalled, result, err := host.callInitFunction()
	if err != nil {
		return nil, err
	}

	if initCalled {
		host.Output().Finish(result)
	}

	outputAccounts := host.Output().GetOutputAccounts()
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

func (host *vmContext) execute(input *vmcommon.ContractCallInput) error {
	contract, err := host.Blockchain().GetCode(host.scAddress)
	if err != nil {
		return err
	}

	totalGasConsumed := input.GasProvided

	defer func() {
		host.Metering().UseGas(totalGasConsumed)
	}()

	gasLeft, err := host.deductInitialCodeCost(
		input.GasProvided,
		contract,
		0, // create cost was elrady taken care of. as it is different for ethereum and elrond
		host.Metering().GasSchedule().BaseOperationCost.StorePerByte,
	)
	if err != nil {
		return err
	}

	newInstance, err := wasmer.NewMeteredInstance(contract, gasLeft)
	if err != nil {
		host.Metering().UseGas(input.GasProvided)
		return err
	}

	idContext := arwen.AddHostContext(host)
	oldInstance := host.instance
	host.instance = newInstance
	host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointNone)
	defer func() {
		host.instance = oldInstance
		newInstance.Clean()
		arwen.RemoveHostContext(idContext)
	}()

	newInstance.SetContextData(unsafe.Pointer(&idContext))

	if host.isInitFunctionCalled() {
		return ErrInitFuncCalledInRun
	}

	function, ok := newInstance.Exports[host.callFunction]
	if !ok {
		return ErrFuncNotFound
	}

	result, err := function()
	if err != nil {
		return ErrFunctionRunError
	}

	if host.Output().ReturnCode() != vmcommon.Ok {
		return ErrReturnCodeNotOk
	}

	convertedResult := arwen.ConvertReturnValue(result)
	host.Output().Finish(convertedResult.Bytes())

	totalGasConsumed = input.GasProvided - gasLeft - newInstance.GetPointsUsed()

	return nil
}

func (host *vmContext) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	currVmInput := host.vmInput
	currScAddress := host.scAddress
	currCallFunction := host.callFunction

	host.vmInput = input.VMInput
	host.scAddress = input.RecipientAddr
	host.callFunction = input.Function

	err := host.execute(input)

	host.vmInput = currVmInput
	host.scAddress = currScAddress
	host.callFunction = currCallFunction

	return err
}

func (host *vmContext) PushState() {
	host.bigIntSubcontext.PushState()
	host.runtimeSubcontext.PushState()
	host.outputSubcontext.PushState()
}

func (host *vmContext) PopState() error {
	err := host.bigIntSubcontext.PopState()
	if err != nil {
		return err
	}

	err = host.runtimeSubcontext.PopState()
	if err != nil {
		return err
	}
	
	err = host.outputSubcontext.PopState()
	if err != nil {
		return err
	}

	return nil
}

func (host *vmContext) ExecuteOnDestContext(input *vmcommon.ContractCallInput) error {
	host.PushState()

	var err error
	defer func() {
		popErr := host.PopState()
		if popErr != nil {
			err = popErr
		}
	}()

	host.vmInput = input.VMInput
	host.scAddress = input.RecipientAddr
	host.callFunction = input.Function

	host.initInternalValues()
	err = host.execute(input)

	return err
}

// The first four bytes is the method selector. The rest of the input data are method arguments in chunks of 32 bytes.
// The method selector is the kecccak256 hash of the method signature.
func (host *vmContext) createETHCallInput() []byte {
	newInput := make([]byte, 0)

	if len(host.callFunction) > 0 {
		hashOfFunction, err := host.cryptoHook.Keccak256([]byte(host.callFunction))
		if err != nil {
			return nil
		}

		newInput = append(newInput, hashOfFunction[0:4]...)
	}

	for _, arg := range host.vmInput.Arguments {
		paddedArg := make([]byte, arwen.ArgumentLenEth)
		copy(paddedArg[arwen.ArgumentLenEth-len(arg):], arg)
		newInput = append(newInput, paddedArg...)
	}

	return newInput
}

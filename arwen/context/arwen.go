package context

import (
	"fmt"
	"math/big"

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

	ethInput []byte

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

	host := &vmContext{
		blockChainHook:       blockChainHook,
		cryptoHook:           cryptoHook,
		meteringSubcontext:   nil,
		runtimeSubcontext:    nil,
		blockchainSubcontext: nil,
		storageSubcontext:    nil,
		bigIntSubcontext:     nil,
	}

	var err error
	
	host.blockchainSubcontext, err = subcontexts.NewBlockchainSubcontext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeSubcontext, err = subcontexts.NewRuntimeSubcontext(host, blockChainHook, vmType)
	if err != nil {
		return nil, err
	}

	host.meteringSubcontext, err = subcontexts.NewMeteringSubcontext(host, gasSchedule, blockGasLimit)
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

	err = wasmer.SetImports(imports)
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(gasSchedule)
	if err != nil {
		return nil, err
	}

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.InitState()

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
	host.InitState()

	blockchain := host.Blockchain()
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return nil, err
	}

	runtime.SetSCAddress(address)
	output.AddTxValueToAccount(address, input.CallValue)

	runtime.GetVMInput().GasProvided, err = host.deductInitialCodeCost(
		input.GasProvided,
		input.ContractCode,
		metering.GasSchedule().ElrondAPICost.CreateContract,
		metering.GasSchedule().BaseOperationCost.StorePerByte,
	)
	if err != nil {
		return host.createVMOutputInCaseOfError(vmcommon.OutOfGas, err.Error()), nil
	}

	err = runtime.CreateWasmerInstance(input.ContractCode)

	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid, err.Error()), nil
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextId(idContext)
	defer func() {
		arwen.RemoveHostContext(idContext)
		runtime.CleanInstance()
	}()

	_, result, err := host.callInitFunction()
	if err != nil {
		return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature, err.Error()), nil
	}

	output.DeployCode(address, input.ContractCode)

	vmOutput := output.CreateVMOutput(result)

	return vmOutput, err
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
	host.InitState()
	runtime := host.Runtime()
	output := host.Output()
	blockchain := host.Blockchain()

	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)

	contract, err := blockchain.GetCode(runtime.GetSCAddress())
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid, err.Error()), nil
	}

	runtime.GetVMInput().GasProvided, err = host.deductInitialCodeCost(
		input.GasProvided,
		contract,
		0,
		host.Metering().GasSchedule().BaseOperationCost.CompilePerByte,
	)
	if err != nil {
		return host.createVMOutputInCaseOfError(vmcommon.OutOfGas, err.Error()), nil
	}

	err = runtime.CreateWasmerInstance(contract)
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid, err.Error()), nil
	}

	idContext := arwen.AddHostContext(host)
	runtime.SetInstanceContextId(idContext)
	defer func() {
		runtime.CleanInstance()
		arwen.RemoveHostContext(idContext)
	}()

	if host.isInitFunctionBeingCalled() {
		fmt.Println("arwen Error", ErrInitFuncCalledInRun.Error())
		return host.createVMOutputInCaseOfError(vmcommon.UserError, err.Error()), nil
	}

	function, err := runtime.GetFunctionToCall()
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.FunctionNotFound, err.Error()), nil
	}

	result, err := function()
	if err != nil {
		strError, _ := wasmer.GetLastError()

		fmt.Println("arwen Error", err.Error(), strError)
		return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature, err.Error()), nil
	}

	returnCode := output.ReturnCode()
	if returnCode != vmcommon.Ok {
		// user error: signalError()
		return host.createVMOutputInCaseOfError(returnCode, ""), nil
	}

	vmOutput := output.CreateVMOutputFromValue(result)

	return vmOutput, nil
}

func (host *vmContext) isInitFunctionBeingCalled() bool {
	functionName := host.Runtime().Function()
	return functionName == arwen.InitFunctionName || functionName == arwen.InitFunctionNameEth
}

func (host *vmContext) createVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{GasRemaining: 0, GasRefund: big.NewInt(0)}
	vmOutput.ReturnCode = errCode
	vmOutput.ReturnMessage = message
	return vmOutput
}

func (host *vmContext) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()
	metering := host.Metering()
	output := host.Output()

	if runtime.ReadOnly() {
		return nil, ErrInvalidCallOnReadOnlyMode
	}

	runtime.PushState()

	defer func() {
		runtime.PopState()
	}()

	runtime.SetVMInput(&input.VMInput)
	address, err := blockchain.NewAddress(input.CallerAddr)
	if err != nil {
		return nil, err
	}

	output.Transfer(address, input.CallerAddr, 0, input.CallValue, nil)
	blockchain.IncreaseNonce(input.CallerAddr)
	runtime.SetSCAddress(address)

	totalGasConsumed := input.GasProvided
	defer func() {
		metering.UseGas(totalGasConsumed)
	}()

	gasLeft, err := host.deductInitialCodeCost(
		input.GasProvided,
		input.ContractCode,
		0, // create cost was elrady taken care of. as it is different for ethereum and elrond
		metering.GasSchedule().BaseOperationCost.StorePerByte,
	)
	if err != nil {
		return nil, err
	}

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	defer func() {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
	}()

	err = runtime.CreateWasmerInstanceWithGasLimit(input.ContractCode, gasLeft)
	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		// TODO use all gas here?
		return nil, err
	}

	runtime.SetInstanceContextId(idContext)

	initCalled, result, err := host.callInitFunction()
	if err != nil {
		return nil, err
	}

	if initCalled {
		output.Finish(result)
	}

	outputAccounts := output.GetOutputAccounts()
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

	totalGasConsumed = input.GasProvided - gasLeft - runtime.GetPointsUsed()

	return address, nil
}

func (host *vmContext) InitState() {
	host.BigInt().InitState()
	host.Output().InitState()
	host.Runtime().InitState()
}

func (host *vmContext) CallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmContext) execute(input *vmcommon.ContractCallInput) error {
	runtime := host.Runtime()
	contract, err := host.Blockchain().GetCode(host.Runtime().GetSCAddress())
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

	idContext := arwen.AddHostContext(host)
	runtime.PushInstance()

	defer func() {
		runtime.PopInstance()
		arwen.RemoveHostContext(idContext)
	}()

	newInstance, err := wasmer.NewMeteredInstance(contract, gasLeft)
	if err != nil {
		host.Metering().UseGas(input.GasProvided)
		return err
	}

	runtime.SetInstanceContextId(idContext)

	if host.isInitFunctionBeingCalled() {
		return ErrInitFuncCalledInRun
	}

	exports := runtime.GetInstanceExports()
	functionName := runtime.Function()
	function, ok := exports[functionName]
	if !ok {
		return subcontexts.ErrFuncNotFound
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
	host.Runtime().PushState()
	host.Runtime().InitStateFromContractCallInput(input)

	err := host.execute(input)

	host.Runtime().PopState()
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

	host.InitState()

	host.Runtime().InitStateFromContractCallInput(input)
	err = host.execute(input)

	return err
}

func (host *vmContext) callInitFunction() (bool, []byte, error) {
	exports := host.Runtime().GetInstanceExports()
	init, ok := exports[arwen.InitFunctionName]

	if !ok {
		init, ok = exports[arwen.InitFunctionNameEth]
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

// The first four bytes is the method selector. The rest of the input data are method arguments in chunks of 32 bytes.
// The method selector is the kecccak256 hash of the method signature.
func (host *vmContext) createETHCallInput() []byte {
	newInput := make([]byte, 0)

	function := host.Runtime().Function()
	if len(function) > 0 {
		hashOfFunction, err := host.cryptoHook.Keccak256([]byte(function))
		if err != nil {
			return nil
		}

		newInput = append(newInput, hashOfFunction[0:4]...)
	}

	for _, arg := range host.Runtime().Arguments() {
		paddedArg := make([]byte, arwen.ArgumentLenEth)
		copy(paddedArg[arwen.ArgumentLenEth-len(arg):], arg)
		newInput = append(newInput, paddedArg...)
	}

	return newInput
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

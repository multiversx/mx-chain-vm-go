package host

import (
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/contexts"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/ethapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var log = logger.GetOrCreate("arwen/host")

var MaximumWasmerInstanceCount = uint64(10)

// TryFunction corresponds to the try() part of a try / catch block
type TryFunction func()

// CatchFunction corresponds to the catch() part of a try / catch block
type CatchFunction func(error)

// vmHost implements HostContext interface.
type vmHost struct {
	blockChainHook vmcommon.BlockchainHook
	cryptoHook     vmcommon.CryptoHook

	ethInput []byte

	blockchainContext arwen.BlockchainContext
	runtimeContext    arwen.RuntimeContext
	outputContext     arwen.OutputContext
	meteringContext   arwen.MeteringContext
	storageContext    arwen.StorageContext
	bigIntContext     arwen.BigIntContext

	scAPIMethods             *wasmer.Imports
	protocolBuiltinFunctions vmcommon.FunctionNames
}

// NewArwenVM creates a new Arwen vmHost
func NewArwenVM(
	blockChainHook vmcommon.BlockchainHook,
	cryptoHook vmcommon.CryptoHook,
	hostParameters *arwen.VMHostParameters,
) (*vmHost, error) {

	host := &vmHost{
		blockChainHook:           blockChainHook,
		cryptoHook:               cryptoHook,
		meteringContext:          nil,
		runtimeContext:           nil,
		blockchainContext:        nil,
		storageContext:           nil,
		bigIntContext:            nil,
		scAPIMethods:             nil,
		protocolBuiltinFunctions: hostParameters.ProtocolBuiltinFunctions,
	}

	var err error

	imports, err := elrondapi.ElrondEIImports()
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

	host.scAPIMethods = imports

	host.blockchainContext, err = contexts.NewBlockchainContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeContext, err = contexts.NewRuntimeContext(host, hostParameters.VMType)
	if err != nil {
		return nil, err
	}

	host.meteringContext, err = contexts.NewMeteringContext(host, hostParameters.GasSchedule, hostParameters.BlockGasLimit)
	if err != nil {
		return nil, err
	}

	host.outputContext, err = contexts.NewOutputContext(host)
	if err != nil {
		return nil, err
	}

	host.storageContext, err = contexts.NewStorageContext(host, blockChainHook, hostParameters.ElrondProtectedKeyPrefix)
	if err != nil {
		return nil, err
	}

	host.bigIntContext, err = contexts.NewBigIntContext()
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(hostParameters.GasSchedule)
	if err != nil {
		return nil, err
	}

	host.runtimeContext.SetMaxInstanceCount(MaximumWasmerInstanceCount)

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.InitState()

	return host, nil
}

func (host *vmHost) Crypto() vmcommon.CryptoHook {
	return host.cryptoHook
}

func (host *vmHost) Blockchain() arwen.BlockchainContext {
	return host.blockchainContext
}

func (host *vmHost) Runtime() arwen.RuntimeContext {
	return host.runtimeContext
}

func (host *vmHost) Output() arwen.OutputContext {
	return host.outputContext
}

func (host *vmHost) Metering() arwen.MeteringContext {
	return host.meteringContext
}

func (host *vmHost) Storage() arwen.StorageContext {
	return host.storageContext
}

func (host *vmHost) BigInt() arwen.BigIntContext {
	return host.bigIntContext
}

func (host *vmHost) GetContexts() (
	arwen.BigIntContext,
	arwen.BlockchainContext,
	arwen.MeteringContext,
	arwen.OutputContext,
	arwen.RuntimeContext,
	arwen.StorageContext,
) {
	return host.bigIntContext,
		host.blockchainContext,
		host.meteringContext,
		host.outputContext,
		host.runtimeContext,
		host.storageContext
}

func (host *vmHost) InitState() {
	host.ClearContextStateStack()
	host.bigIntContext.InitState()
	host.outputContext.InitState()
	host.runtimeContext.InitState()
	host.storageContext.InitState()
	host.ethInput = nil
}

func (host *vmHost) ClearContextStateStack() {
	host.bigIntContext.ClearStateStack()
	host.outputContext.ClearStateStack()
	host.runtimeContext.ClearStateStack()
	host.storageContext.ClearStateStack()
}

func (host *vmHost) Clean() {
	host.runtimeContext.CleanInstance()
	arwen.RemoveAllHostContexts()
}

func (host *vmHost) GetAPIMethods() *wasmer.Imports {
	return host.scAPIMethods
}

func (host *vmHost) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return host.protocolBuiltinFunctions
}

func (host *vmHost) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	log.Trace("RunSmartContractCreate begin", "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	try := func() {
		vmOutput = host.doRunSmartContractCreate(input)
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCreate", "error", err)
	}

	ok := TryCatch(try, catch, "arwen.RunSmartContractCreate")
	if ok {
		log.Trace("RunSmartContractCreate end", "returnCode", vmOutput.ReturnCode, "returnMessage", vmOutput.ReturnMessage)
	}

	return
}

func (host *vmHost) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	log.Trace("RunSmartContractCall begin", "function", input.Function)

	tryUpgrade := func() {
		vmOutput = host.doRunSmartContractUpgrade(input)
	}

	tryCall := func() {
		vmOutput = host.doRunSmartContractCall(input)
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCall", "error", err)
	}

	var ok bool
	isUpgrade := input.Function == arwen.UpgradeFunctionName
	if isUpgrade {
		ok = TryCatch(tryUpgrade, catch, "arwen.RunSmartContractUpgrade")
	} else {
		ok = TryCatch(tryCall, catch, "arwen.RunSmartContractCall")
	}

	if ok {
		log.Trace("RunSmartContractCall end", "returnCode", vmOutput.ReturnCode, "returnMessage", vmOutput.ReturnMessage)
	}

	return
}

// TryCatch simulates a try/catch block using golang's recover() functionality
func TryCatch(try TryFunction, catch CatchFunction, catchFallbackMessage string) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%s, panic: %v", catchFallbackMessage, r)
			}

			catch(err)
			ok = false
		}
	}()

	try()
	ok = true
	return
}

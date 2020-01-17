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
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

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
}

func NewArwenVM(
	blockChainHook vmcommon.BlockchainHook,
	cryptoHook vmcommon.CryptoHook,
	vmType []byte,
	blockGasLimit uint64,
	gasSchedule map[string]map[string]uint64,
) (*vmHost, error) {

	host := &vmHost{
		blockChainHook:    blockChainHook,
		cryptoHook:        cryptoHook,
		meteringContext:   nil,
		runtimeContext:    nil,
		blockchainContext: nil,
		storageContext:    nil,
		bigIntContext:     nil,
	}

	var err error

	host.blockchainContext, err = contexts.NewBlockchainContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeContext, err = contexts.NewRuntimeContext(host, blockChainHook, vmType)
	if err != nil {
		return nil, err
	}

	host.meteringContext, err = contexts.NewMeteringContext(host, gasSchedule, blockGasLimit)
	if err != nil {
		return nil, err
	}

	host.outputContext, err = contexts.NewOutputContext(host)
	if err != nil {
		return nil, err
	}

	host.storageContext, err = contexts.NewStorageContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.bigIntContext, err = contexts.NewBigIntContext()
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

func (host *vmHost) InitState() {
	host.BigInt().InitState()
	host.Output().InitState()
	host.Runtime().InitState()
	host.ethInput = nil
}

func (host *vmHost) PushState() {
	host.bigIntContext.PushState()
	host.runtimeContext.PushState()
	host.outputContext.PushState()
}

func (host *vmHost) PopState() error {
	err := host.bigIntContext.PopState()
	if err != nil {
		return err
	}

	err = host.runtimeContext.PopState()
	if err != nil {
		return err
	}

	err = host.outputContext.PopState()
	if err != nil {
		return err
	}

	return nil
}

func (host *vmHost) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	try := func() {
		vmOutput = host.doRunSmartContractCreate(input)
	}

	catch := func(caught error) {
		err = caught
	}

	TryCatch(try, catch, "arwen.RunSmartContractCreate")
	return
}

func (host *vmHost) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	try := func() {
		vmOutput = host.doRunSmartContractCall(input)
	}

	catch := func(caught error) {
		err = caught
	}

	TryCatch(try, catch, "arwen.RunSmartContractCall")
	return
}

// TryCatch simulates a try/catch block using golang's recover() functionality
func TryCatch(try TryFunction, catch CatchFunction, catchFallbackMessage string) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%s, panic: %v", catchFallbackMessage, r)
			}

			catch(err)
		}
	}()

	try()
}

package context

import (
	"fmt"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/context/subcontexts"
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

// vmContext implements HostContext interface.
type vmContext struct {
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

func (host *vmContext) InitState() {
	host.BigInt().InitState()
	host.Output().InitState()
	host.Runtime().InitState()
	host.ethInput = nil
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

func (host *vmContext) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	try := func() {
		vmOutput, err = host.doRunSmartContractCreate(input)
	}

	catch := func(caught error) {
		err = caught
	}

	TryCatch(try, catch, "arwen.RunSmartContractCreate")
	return
}

func (host *vmContext) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	try := func() {
		vmOutput, err = host.doRunSmartContractCall(input)
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

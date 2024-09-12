package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	interpreter "github.com/multiversx/mx-chain-vm-go/evm/interpreter"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"math/big"
)

func (instance *EVMInstance) resetEVM() {
	instance.evm = interpreter.NewEVM(
		instance.buildBlockContext(),
		instance.buildTxContext(),
		interpreter.CreateEVMStateDB(instance.evmExecutor.evmHooks),
		instance.buildChainConfig(),
		instance.evmExecutor.gasConfig,
		&instance.evmExecutor.instructionSet,
	)
}

func (instance *EVMInstance) prepareAddress(functionName string) error {
	if functionName != vmhost.InitFunctionName {
		return nil
	}

	return instance.evmExecutor.evmHooks.SaveAliasAddress()
}

func (instance *EVMInstance) prepareContract(functionName string) (*interpreter.Contract, []byte) {
	evmHooks := instance.evmExecutor.evmHooks
	input, code := instance.prepareInputAndCode(functionName)

	contractAddress := evmHooks.ContractAddress()
	contract := interpreter.NewContract(
		interpreter.AccountRef(evmHooks.CallerAddress()),
		interpreter.AccountRef(contractAddress),
		evmHooks.CallValue(),
	)
	contract.Code = code
	contract.CodeHash = evmHooks.CodeHash()
	return contract, input
}

func (instance *EVMInstance) consumeOutput(functionName string, returnData []byte) {
	evmHooks := instance.evmExecutor.evmHooks

	switch functionName {
	case vmhost.InitFunctionName:
		instance.isCompiled = true
		instance.code = returnData
		evmHooks.FinishCreate(returnData)
	default:
		evmHooks.Finish(returnData)
	}
}

func (instance *EVMInstance) prepareInputAndCode(functionName string) ([]byte, []byte) {
	code := instance.code
	input := instance.flattenInput()

	switch functionName {
	case vmhost.InitFunctionName:
		return nil, append(code, input...)
	default:
		selector := functionNameToSelector(functionName)
		return append(selector, input...), code
	}
}

func (instance *EVMInstance) flattenInput() []byte {
	var input []byte
	arguments := instance.evmExecutor.evmHooks.Arguments()
	for _, slice := range arguments {
		input = append(input, slice...)
	}
	return input
}

func (instance *EVMInstance) buildBlockContext() interpreter.BlockContext {
	evmHooks := instance.evmExecutor.evmHooks
	return interpreter.BlockContext{
		GetHash:     evmHooks.GetHash,
		Coinbase:    common.Address{},
		GasLimit:    evmHooks.BlockGasLimit(),
		BlockNumber: evmHooks.BlockNumber(),
		Time:        evmHooks.Time(),
		Difficulty:  new(big.Int),
		BaseFee:     new(big.Int),
		BlobBaseFee: new(big.Int),
		Random:      evmHooks.Random(),
	}
}

func (instance *EVMInstance) buildTxContext() interpreter.TxContext {
	evmHooks := instance.evmExecutor.evmHooks
	return interpreter.TxContext{
		Origin:     evmHooks.Origin(),
		GasPrice:   evmHooks.GasPrice(),
		BlobHashes: make([]common.Hash, 0),
		BlobFeeCap: new(big.Int),
	}
}

func (instance *EVMInstance) buildChainConfig() *params.ChainConfig {
	evmHooks := instance.evmExecutor.evmHooks
	return &params.ChainConfig{
		ChainID: evmHooks.ChainID(),
	}
}

func functionNameToSelector(functionName string) []byte {
	return common.Hex2Bytes(functionName)
}

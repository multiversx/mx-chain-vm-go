package evmhooks

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

const evmCallExpectedReturnDataSize = 1

const evmCallExpectedReturnDataPosition = 0

const nonEvmCallFunctionNameInputPosition = 0

const nonEvmCallArgumentsInputPositionStart = 1

func (context *EVMHooksImpl) IsSmartContractAddress(address common.Address) bool {
	return core.IsSmartContractAddress(context.toMVXAddress(address))
}

func (context *EVMHooksImpl) Create(code []byte, gas uint64, value *uint256.Int) ([]byte, common.Address, error) {
	aliasAddress, err := context.createEthereumContractAddress(context.ContractMvxAddress())
	if err != nil {
		return nil, common.Address{}, err
	}

	return context.createContract(code, gas, value, aliasAddress)
}

func (context *EVMHooksImpl) Create2(code []byte, gas uint64, value *uint256.Int, salt *uint256.Int) ([]byte, common.Address, error) {
	aliasAddress := crypto.CreateAddress2(context.ContractAddress(), salt.Bytes32(), crypto.Keccak256Hash(code).Bytes())
	return context.createContract(code, gas, value, aliasAddress)
}

func (context *EVMHooksImpl) Call(address common.Address, value *uint256.Int, input []byte, gas uint64) ([]byte, error) {
	destination := context.toMVXAddress(address)
	isNonEvmCall := context.isNonEvmCall(destination)
	functionName, arguments, err := parseInput(input, isNonEvmCall)
	if err != nil {
		return nil, err
	}

	returnDataLength := context.returnDataLength()
	_, err = vmhooks.ExecuteOnDestContextUnmetered(
		context.host,
		int64(gas),
		value.ToBig(),
		functionName,
		destination,
		arguments,
		0,
	)
	if err != nil {
		return nil, err
	}

	return context.parseReturnData(isNonEvmCall, returnDataLength)
}

func (context *EVMHooksImpl) StaticCall(address common.Address, input []byte, gas uint64) ([]byte, error) {
	destination := context.toMVXAddress(address)
	isNonEvmCall := context.isNonEvmCall(destination)
	functionName, arguments, err := parseInput(input, isNonEvmCall)
	if err != nil {
		return nil, err
	}

	returnDataLength := context.returnDataLength()
	_, err = vmhooks.ExecuteReadOnlyUnmetered(
		context.host,
		int64(gas),
		functionName,
		context.toMVXAddress(address),
		arguments,
		0,
	)
	if err != nil {
		return nil, err
	}

	return context.parseReturnData(isNonEvmCall, returnDataLength)
}

func (context *EVMHooksImpl) CallCode(address common.Address, value *uint256.Int, input []byte, gas uint64) ([]byte, error) {
	sender := context.ContractMvxAddress()
	return context.executeOnSameContext(sender, sender, address, value, input, gas, true)
}

func (context *EVMHooksImpl) DelegateCall(address common.Address, input []byte, gas uint64) ([]byte, error) {
	sender := context.CallerMvxAddress()
	receiver := context.ContractMvxAddress()
	return context.executeOnSameContext(sender, receiver, address, context.CallValue(), input, gas, false)
}

func (context *EVMHooksImpl) createContract(code []byte, gas uint64, value *uint256.Int, aliasAddress common.Address) ([]byte, common.Address, error) {
	sender := context.ContractMvxAddress()
	returnDataLength := context.returnDataLength()
	_, err := vmhooks.CreateContractWithAddress(
		sender,
		[][]byte{},
		value.ToBig(),
		int64(gas),
		code,
		nil,
		context.host,
		vmhooks.CreateContract,
		aliasAddress.Bytes(),
	)
	if err != nil {
		return nil, common.Address{}, err
	}

	returnData, err := context.parseReturnData(false, returnDataLength)
	if err != nil {
		return nil, common.Address{}, err
	}

	return returnData, aliasAddress, err
}

func (context *EVMHooksImpl) executeOnSameContext(sender []byte, receiver []byte, codeAddress common.Address, value *uint256.Int, input []byte, gas uint64, doTransfer bool) ([]byte, error) {
	returnDataLength := context.returnDataLength()
	functionName, arguments, err := parseInput(input, false)
	if err != nil {
		return nil, err
	}

	_, err = vmhooks.ExecuteOnSameContextUnmetered(
		context.host,
		int64(gas),
		value.ToBig(),
		functionName,
		receiver,
		arguments,
		sender,
		0,
		context.toMVXAddress(codeAddress),
		doTransfer,
	)
	if err != nil {
		return nil, err
	}

	return context.parseReturnData(false, returnDataLength)
}

func (context *EVMHooksImpl) isNonEvmCall(destination []byte) bool {
	return context.host.IsOutOfVMFunctionExecution(&vmcommon.ContractCallInput{RecipientAddr: destination})
}

func (context *EVMHooksImpl) returnDataLength() int {
	return len(context.GetOutputContext().ReturnData())
}

func (context *EVMHooksImpl) extractReturnData(previousReturnDataLength int) [][]byte {
	returnData := context.GetOutputContext().ReturnData()
	context.GetOutputContext().ClearReturnData()
	if len(returnData) > previousReturnDataLength {
		return returnData[previousReturnDataLength:]
	}
	return nil
}

func (context *EVMHooksImpl) parseReturnData(isNonEvmCall bool, oldLength int) ([]byte, error) {
	returnData := context.extractReturnData(oldLength)
	if returnData == nil {
		return nil, nil
	}
	if isNonEvmCall {
		return lengthPrefixEncode(returnData), nil
	}
	if len(returnData) != evmCallExpectedReturnDataSize {
		return nil, ErrInvalidReturnDataSize
	}
	return returnData[evmCallExpectedReturnDataPosition], nil
}

func parseInput(input []byte, isNonEvmCall bool) ([]byte, [][]byte, error) {
	if isNonEvmCall {
		decodedArguments, err := lengthPrefixDecode(input)
		if err != nil {
			return nil, nil, err
		}
		return decodedArguments[nonEvmCallFunctionNameInputPosition], decodedArguments[nonEvmCallArgumentsInputPositionStart:], nil
	}

	functionName, arguments, err := parsers.ParseEthereumCallInput(input)
	if err != nil {
		return nil, nil, err
	}

	functionNameHex := []byte(hex.EncodeToString(functionName))
	if len(arguments) > 0 {
		return functionNameHex, [][]byte{arguments}, nil
	}
	return functionNameHex, nil, nil
}

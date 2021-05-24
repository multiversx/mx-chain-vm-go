package contracts

import (
	"encoding/binary"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/elrondapi"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/testcommon"
)

// ExecESDTTransferAndCallParentMock is an exposed mock contract method
func ExecESDTTransferAndCallParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("execESDTTransferAndCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 3 {
			host.Runtime().SignalUserError("need 3 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.CallerAddr = host.Runtime().GetSCAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.ESDTTestTokenName,
			big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
			arguments[2],
		}
		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		_, _, err := host.ExecuteOnDestContext(input)
		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

// ExecESDTTransferWithAPICall is an exposed mock contract method
func ExecESDTTransferWithAPICall(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("execESDTTransferWithAPICall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 3 {
			host.Runtime().SignalUserError("need 3 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.CallerAddr = host.Runtime().GetSCAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.ESDTTestTokenName,
			big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
			arguments[2],
		}
		input.RecipientAddr = arguments[0]

		runtime := host.Runtime()

		recipientAddr := arguments[0]
		offsetDest := int32(0)
		runtime.MemStore(int32(offsetDest), recipientAddr)

		offsetTokenID := offsetDest + int32(len(recipientAddr))
		tokenLen := int32(len(test.ESDTTestTokenName))
		runtime.MemStore(offsetTokenID, test.ESDTTestTokenName)

		value := big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes()
		offsetValue := offsetTokenID + tokenLen
		valueLen := int32(arwen.BalanceLen)
		value = arwen.PadBytesLeft(value, arwen.BalanceLen)
		runtime.MemStore(int32(offsetValue), value)

		functionName := arguments[1]
		offsetFunction := offsetValue + valueLen
		funcNameLen := int32(len(functionName))
		runtime.MemStore(offsetFunction, functionName)

		noOfArguments := 1
		argumentsLengths := []uint32{uint32(len(arguments[2]))}

		offsetArgumentsLength := offsetFunction + funcNameLen
		argumentsLengthsAsBytes := make([]byte, noOfArguments*4)
		binary.LittleEndian.PutUint32(argumentsLengthsAsBytes, argumentsLengths[0])
		runtime.MemStore(int32(offsetArgumentsLength), argumentsLengthsAsBytes)

		argumentsData := make([]byte, argumentsLengths[0])
		copy(argumentsData, arguments[2])
		offsetArgumentsData := offsetArgumentsLength + int32(len(argumentsLengthsAsBytes))
		runtime.MemStore(offsetArgumentsData, argumentsData)

		elrondapi.TransferESDTNFTExecuteWithHost(host,
			offsetDest,
			offsetTokenID,
			tokenLen,
			offsetValue,
			0,
			int64(testConfig.GasProvidedToChild),
			offsetFunction,
			funcNameLen,
			1,
			offsetArgumentsLength,
			offsetArgumentsData)

		return instance
	})
}

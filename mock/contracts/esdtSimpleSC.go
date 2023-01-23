package contracts

import (
	"fmt"
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	mock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooks"
)

// ExecESDTTransferAndCallChild is an exposed mock contract method
func ExecESDTTransferAndCallChild(instanceMock *mock.InstanceMock, config interface{}) {
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
		input.CallerAddr = host.Runtime().GetContextAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.ESDTTestTokenName,
			big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
			arguments[2],
		}
		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		returnValue := ExecuteOnDestContextInMockContracts(host, input)
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
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
		input.CallerAddr = host.Runtime().GetContextAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.ESDTTestTokenName,
			big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
			arguments[2],
		}
		input.RecipientAddr = arguments[0]

		functionName := arguments[1]
		args := [][]byte{arguments[2]}

		transfer := &vmcommon.ESDTTransfer{
			ESDTValue:      big.NewInt(int64(testConfig.ESDTTokensToTransfer)),
			ESDTTokenName:  test.ESDTTestTokenName,
			ESDTTokenType:  0,
			ESDTTokenNonce: 0,
		}

		vmhooks.TransferESDTNFTExecuteWithTypedArgs(
			host,
			input.RecipientAddr,
			[]*vmcommon.ESDTTransfer{transfer},
			int64(testConfig.GasProvidedToChild),
			functionName,
			args)

		return instance
	})
}

// ExecESDTTransferAndAsyncCallChild is an exposed mock contract method
func ExecESDTTransferAndAsyncCallChild(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("execESDTTransferAndAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 3 {
			host.Runtime().SignalUserError("need 3 arguments")
			return instance
		}

		functionToCallOnChild := arguments[2]

		receiver := arguments[0]
		builtInFunction := arguments[1]

		callData := txDataBuilder.NewBuilder()
		// function to be called on child
		callData.Func(string(builtInFunction))
		callData.Bytes(test.ESDTTestTokenName)
		callData.Bytes(big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes())
		callData.Bytes(functionToCallOnChild)

		value := big.NewInt(0).Bytes()

		err := host.Runtime().ExecuteAsyncCall(receiver, callData.ToBytes(), value)

		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

// ExecESDTTransferInAsyncCall is an exposed mock contract method
func ExecESDTTransferInAsyncCall(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("esdtTransferInAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 1 {
			host.Runtime().SignalUserError("need 1 arguments")
			return instance
		}

		receiver := arguments[0]

		callData := txDataBuilder.NewBuilder()
		callData.Func("ESDTTransfer")
		callData.Bytes(test.ESDTTestTokenName)
		callData.Bytes(big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes())

		value := big.NewInt(0).Bytes()

		err := host.Runtime().ExecuteAsyncCall(receiver, callData.ToBytes(), value)

		if err != nil {
			host.Runtime().FailExecution(err)
			return instance
		}

		return instance
	})
}

// EvilCallback is an exposed mock contract method
func EvilCallback(instanceMock *mock.InstanceMock, _ interface{}) {
	// testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		retVal := vmhooks.ExecuteOnDestContextByCallerWithTypedArgs(
			host,
			int64(host.Metering().GasLeft()),
			big.NewInt(0),
			[]byte("wasteGas"),
			test.ChildAddress, // owned by UserAddress2 (the CallserAddr of this callback)
			[][]byte{},
		)

		if retVal != 0 {
			host.Runtime().SignalUserError("execution by caller failed")
			return instance
		}

		return instance
	})
}

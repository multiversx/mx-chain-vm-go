package contracts

import (
	"fmt"
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/multiversx/mx-chain-vm-go/executor"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
)

// ExecESDTTransferAndCallChild is an exposed mock contract method
func ExecESDTTransferAndCallChild(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("execESDTTransferAndCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

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
	instanceMock.AddMockMethod("execESDTTransferWithAPICall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

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
	instanceMock.AddMockMethod("execESDTTransferAndAsyncCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		arguments := host.Runtime().Arguments()
		if len(arguments) != 4 {
			host.Runtime().SignalUserError("need 4 arguments")
			return instance
		}

		receiver := arguments[0]
		builtInFunction := arguments[1]
		functionToCallOnChild := arguments[2]
		asyncCallType := arguments[3]

		callData := txDataBuilder.NewBuilder()
		// function to be called on child
		callData.Func(string(builtInFunction))
		callData.Bytes(test.ESDTTestTokenName)
		callData.Bytes(big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes())
		callData.Bytes(functionToCallOnChild)
		callData.Bytes(asyncCallType)

		value := big.NewInt(0).Bytes()

		if asyncCallType[0] == 0 {
			err = host.Async().RegisterLegacyAsyncCall(receiver, callData.ToBytes(), value)
		} else {
			callbackName := "callBack"
			if host.Runtime().ValidateCallbackName(callbackName) == executor.ErrFuncNotFound {
				callbackName = ""
			}
			err = host.Async().RegisterAsyncCall("testGroup", &vmhost.AsyncCall{
				Status:          vmhost.AsyncCallPending,
				Destination:     receiver,
				Data:            callData.ToBytes(),
				ValueBytes:      value,
				SuccessCallback: callbackName,
				ErrorCallback:   callbackName,
				GasLimit:        testConfig.GasProvidedToChild,
				GasLocked:       testConfig.GasToLock,
			})
		}

		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

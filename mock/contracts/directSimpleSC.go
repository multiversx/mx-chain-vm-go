package contracts

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/arwen/elrondapimeta"
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
)

// WasteGasChildMock is an exposed mock contract method
func WasteGasChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("wasteGas", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByChild))
}

// FailChildMock is an exposed mock contract method
func FailChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("fail", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Runtime().FailExecution(errors.New("forced fail"))
		return instance
	})
}

// FailChildAndBurnESDTMock is an exposed mock contract method
func FailChildAndBurnESDTMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("failAndBurn", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		runtime := host.Runtime()

		input := test.DefaultTestContractCallInput()
		input.CallerAddr = runtime.GetContextAddress()
		input.GasProvided = runtime.GetVMInput().GasProvided / 2
		input.Arguments = [][]byte{
			test.ESDTTestTokenName,
			runtime.Arguments()[0],
		}
		input.RecipientAddr = host.Runtime().GetContextAddress()
		input.Function = "ESDTLocalBurn"

		returnValue := ExecuteOnDestContextInMockContracts(host, input)
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("Return value %d", returnValue))
		}

		host.Runtime().FailExecution(errors.New("forced fail"))
		return instance
	})
}

// ExecOnSameCtxParentMock is an exposed mock contract method
func ExecOnSameCtxParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("execOnSameCtx", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		argsPerCall := 3
		arguments := host.Runtime().Arguments()
		if len(arguments)%argsPerCall != 0 {
			host.Runtime().SignalUserError("need 3 arguments per individual call")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address

		for callIndex := 0; callIndex < len(arguments); callIndex += argsPerCall {
			input.RecipientAddr = arguments[callIndex+0]
			input.Function = string(arguments[callIndex+1])
			numCalls := big.NewInt(0).SetBytes(arguments[callIndex+2]).Uint64()

			for i := uint64(0); i < numCalls; i++ {
				returnValue := ExecuteOnSameContextInMockContracts(host, input)
				if returnValue != 0 {
					host.Runtime().FailExecution(fmt.Errorf("Return value %d", returnValue))
				}
			}
		}

		return instance
	})
}

// ExecOnDestCtxParentMock is an exposed mock contract method
func ExecOnDestCtxParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("execOnDestCtx", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		argsPerCall := 3
		arguments := host.Runtime().Arguments()
		if len(arguments)%argsPerCall != 0 {
			host.Runtime().SignalUserError("need 3 arguments per individual call")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address

		for callIndex := 0; callIndex < len(arguments); callIndex += argsPerCall {
			input.RecipientAddr = arguments[callIndex+0]
			input.Function = string(arguments[callIndex+1])
			numCalls := big.NewInt(0).SetBytes(arguments[callIndex+2]).Uint64()

			for i := uint64(0); i < numCalls; i++ {
				returnValue := ExecuteOnDestContextInMockContracts(host, input)
				if returnValue != 0 {
					host.Runtime().FailExecution(fmt.Errorf("Return value %d", returnValue))
				}
			}
		}

		return instance
	})
}

// ExecOnDestCtxSingleCallParentMock is an exposed mock contract method
func ExecOnDestCtxSingleCallParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("execOnDestCtxSingleCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 2 {
			host.Runtime().SignalUserError("need 2 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address

		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		returnValue := ExecuteOnDestContextInMockContracts(host, input)
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("Return value %d", returnValue))
		}

		return instance
	})
}

// WasteGasParentMock is an exposed mock contract method
func WasteGasParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("wasteGas", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByParent))
}

const (
	esdtOnCallbackSuccess int = iota
	esdtOnCallbackWrongNumOfArgs
	esdtOnCallbackFail
	esdtOnCallbackNewAsync
)

// ESDTTransferToParentMock is an exposed mock contract method
func ESDTTransferToParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	esdtTransferToParentMock(instanceMock, testConfig, esdtOnCallbackSuccess)
}

// ESDTTransferToParentWrongESDTArgsNumberMock is an exposed mock contract method
func ESDTTransferToParentWrongESDTArgsNumberMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	esdtTransferToParentMock(instanceMock, testConfig, esdtOnCallbackWrongNumOfArgs)
}

// ESDTTransferToParentCallbackWillFail is an exposed mock contract method
func ESDTTransferToParentCallbackWillFail(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	esdtTransferToParentMock(instanceMock, testConfig, esdtOnCallbackFail)
}

// ESDTTransferToParentAndNewAsyncFromCallbackMock is an exposed mock contract method
func ESDTTransferToParentAndNewAsyncFromCallbackMock(instanceMock *mock.InstanceMock, config interface{}) {
	esdtTransferToParentMock(instanceMock, config, esdtOnCallbackNewAsync)
}

func esdtTransferToParentMock(instanceMock *mock.InstanceMock, config interface{}, behavior int) {
	instanceMock.AddMockMethod("transferESDTToParent", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByChild)

		callData := txDataBuilder.NewBuilder()
		callData.Func(string("ESDTTransfer"))
		callData.Bytes(test.ESDTTestTokenName)
		callData.Bytes(big.NewInt(int64(testConfig.CallbackESDTTokensToTransfer)).Bytes())

		switch behavior {
		case esdtOnCallbackSuccess:
			host.Output().Finish([]byte("success"))
		case esdtOnCallbackWrongNumOfArgs:
			callData.Bytes([]byte{})
		case esdtOnCallbackFail:
			host.Output().Finish([]byte("fail"))
		case esdtOnCallbackNewAsync:
			host.Output().Finish([]byte("new_async"))
			host.Output().Finish(host.Runtime().GetContextAddress())
			host.Output().Finish([]byte("wasteGas"))
		}

		value := big.NewInt(0).Bytes()

		arguments := host.Runtime().Arguments()
		asyncCallType := arguments[0]

		async := host.Async()
		var err error
		if asyncCallType[0] == 0 {
			err = async.RegisterLegacyAsyncCall(test.ParentAddress, callData.ToBytes(), value)
		} else {
			callbackName := "callBack"
			if host.Runtime().ValidateCallbackName(callbackName) == elrondapimeta.ErrFuncNotFound {
				callbackName = ""
			}
			err = host.Async().RegisterAsyncCall("testGroup", &arwen.AsyncCall{
				Status:          arwen.AsyncCallPending,
				Destination:     test.ParentAddress,
				Data:            callData.ToBytes(),
				ValueBytes:      value,
				SuccessCallback: callbackName,
				ErrorCallback:   callbackName,
				GasLimit:        testConfig.GasProvidedToChild / 2,
				GasLocked:       testConfig.GasToLock,
			})
		}

		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

var TestStorageValue1 = []byte{1, 2, 3, 4}
var TestStorageValue2 = []byte{1, 2, 3}
var TestStorageValue3 = []byte{1, 2}
var TestStorageValue4 = []byte{1}

// ParentSetStorageMock is an exposed mock contract method
func ParentSetStorageMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("parentSetStorage", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Storage().SetStorage(test.ParentKeyA, TestStorageValue1) // add
		host.Storage().SetStorage(test.ParentKeyA, TestStorageValue2) // delete
		host.Storage().SetStorage(test.ParentKeyB, TestStorageValue2) // add
		host.Storage().SetStorage(test.ParentKeyB, TestStorageValue3) // delete

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address
		input.RecipientAddr = test.ChildAddress
		input.Function = "childSetStorage"

		arguments := host.Runtime().Arguments()
		var returnValue int32
		if bytes.Equal(arguments[0], []byte{0}) {
			returnValue = ExecuteOnSameContextInMockContracts(host, input)
		} else {
			returnValue = ExecuteOnDestContextInMockContracts(host, input)
		}
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("Return value %d", returnValue))
		}

		return instance
	})
}

// ChildSetStorageMock is an exposed mock contract method
func ChildSetStorageMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("childSetStorage", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Storage().SetStorage(test.ChildKey, TestStorageValue2)  // add
		host.Storage().SetStorage(test.ChildKey, TestStorageValue1)  // add
		host.Storage().SetStorage(test.ChildKeyB, TestStorageValue1) // add
		host.Storage().SetStorage(test.ChildKeyB, TestStorageValue4) // delete
		return instance
	})
}

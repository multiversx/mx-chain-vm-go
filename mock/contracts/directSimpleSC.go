package contracts

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// WasteGasChildMock is an exposed mock contract method
func WasteGasChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("wasteGas", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByChild))
}

// ReportOriginalCaller is an exposed mock contract method
func ReportOriginalCaller(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("reportOriginalCaller", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		err := host.Metering().UseGasBounded(testConfig.GasUsedByChild)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		}

		originalCaller := host.Runtime().GetOriginalCallerAddress()
		host.Output().Finish(originalCaller)
		return instance
	})
}

// FailChildMock is an exposed mock contract method
func FailChildMock(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("fail", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Runtime().FailExecution(errors.New("forced fail"))
		return instance
	})
}

// FailChildAndBurnESDTMock is an exposed mock contract method
func FailChildAndBurnESDTMock(instanceMock *mock.InstanceMock, _ interface{}) {
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
			host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
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
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
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
					host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
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
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
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
					host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
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
		if len(arguments) < 2 {
			host.Runtime().SignalUserError("needs at least 2 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = instance.Address

		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		if len(arguments) > 2 {
			input.Arguments = arguments[2:]
		}

		returnValue := ExecuteOnDestContextInMockContracts(host, input)
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
		}

		originalCaller := host.Runtime().GetOriginalCallerAddress()
		host.Output().Finish(originalCaller)

		return instance
	})
}

// WasteGasParentMock is an exposed mock contract method
func WasteGasParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("wasteGas", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByParent))
}

// InitFunctionMock is the exposed init function
func InitFunctionMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod(vmhost.InitFunctionName, func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Output().Finish([]byte(vmhost.InitFunctionName))
		return instance
	})
}

// InitFunctionMock is the exposed upgrade function
func UpgradeFunctionMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod(vmhost.ContractsUpgradeFunctionName, func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Output().Finish([]byte(vmhost.UpgradeFunctionName))
		return instance
	})
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

		switch behavior {
		case esdtOnCallbackSuccess:
			host.Output().Finish([]byte("success"))
		case esdtOnCallbackFail:
			host.Output().Finish([]byte("fail"))
		case esdtOnCallbackNewAsync:
			host.Output().Finish([]byte("new_async"))
			host.Output().Finish(host.Runtime().GetContextAddress())
			host.Output().Finish([]byte("wasteGas"))
		}

		arguments := host.Runtime().Arguments()
		numberOfBackTransfers := uint64(big.NewInt(0).SetBytes(arguments[0]).Int64())

		var err error
		for numCallbacks := uint64(0); numCallbacks < numberOfBackTransfers; numCallbacks++ {
			transfer := &vmcommon.ESDTTransfer{
				ESDTValue:      big.NewInt(int64(testConfig.CallbackESDTTokensToTransfer)),
				ESDTTokenName:  test.ESDTTestTokenName,
				ESDTTokenType:  0,
				ESDTTokenNonce: 0,
			}

			ret := vmhooks.TransferESDTNFTExecuteWithTypedArgs(
				host,
				test.ParentAddress,
				[]*vmcommon.ESDTTransfer{transfer},
				int64(testConfig.GasProvidedToChild),
				nil,
				nil)
			if ret != 0 {
				host.Runtime().FailExecution(fmt.Errorf("Transfer ESDT failed"))
			}

		}

		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

// test variables
var (
	TestStorageValue1 = []byte{1, 2, 3, 4}
	TestStorageValue2 = []byte{1, 2, 3}
	TestStorageValue3 = []byte{1, 2}
	TestStorageValue4 = []byte{1}
)

// ParentSetStorageMock is an exposed mock contract method
func ParentSetStorageMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("parentSetStorage", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		_, _ = host.Storage().SetStorage(test.ParentKeyA, TestStorageValue1) // add
		_, _ = host.Storage().SetStorage(test.ParentKeyA, TestStorageValue2) // delete
		_, _ = host.Storage().SetStorage(test.ParentKeyB, TestStorageValue2) // add
		_, _ = host.Storage().SetStorage(test.ParentKeyB, TestStorageValue3) // delete

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
			host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
		}

		return instance
	})
}

// ChildSetStorageMock is an exposed mock contract method
func ChildSetStorageMock(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("childSetStorage", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		_, _ = host.Storage().SetStorage(test.ChildKey, TestStorageValue2)  // add
		_, _ = host.Storage().SetStorage(test.ChildKey, TestStorageValue1)  // add
		_, _ = host.Storage().SetStorage(test.ChildKeyB, TestStorageValue1) // add
		_, _ = host.Storage().SetStorage(test.ChildKeyB, TestStorageValue4) // delete
		return instance
	})
}

// SimpleChildSetStorageMock is an exposed mock contract method
func SimpleChildSetStorageMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("simpleChildSetStorage", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		_, _ = host.Storage().SetStorage(test.ChildKey, test.ChildData)
		return instance
	})
}

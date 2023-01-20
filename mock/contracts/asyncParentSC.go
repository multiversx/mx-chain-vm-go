package contracts

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/arwen/elrondapi"
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
)

// test variables
var (
	AsyncChildFunction = "transferToThirdParty"
	AsyncChildData     = " there"
	ExpectedClosure    = []byte{4, 5, 6}
)

// PerformAsyncCallParentMock is an exposed mock contract method
func PerformAsyncCallParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("performAsyncCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		_, _ = host.Storage().SetStorage(test.ParentKeyA, test.ParentDataA)
		_, _ = host.Storage().SetStorage(test.ParentKeyB, test.ParentDataB)
		host.Output().Finish(test.ParentFinishA)
		host.Output().Finish(test.ParentFinishB)

		scAddress := host.Runtime().GetContextAddress()
		transferValue := big.NewInt(testConfig.TransferToThirdParty)
		err = host.Output().Transfer(testConfig.GetThirdPartyAddress(), scAddress, 0, 0, transferValue, nil, []byte("hello"), 0)
		if err != nil {
			host.Runtime().SignalUserError(err.Error())
			return instance
		}

		_, err = host.Async().GetCallbackClosure()
		if err != arwen.ErrAsyncNoCallbackForClosure {
			host.Runtime().SignalUserError("should not be a closure")
			return instance
		}

		err = RegisterAsyncCallToChild(host, testConfig, host.Runtime().Arguments())
		if err != nil {
			host.Runtime().SignalUserError(err.Error())
			return instance
		}

		return instance

	})
}

// RegisterAsyncCallToChild is resued also in some tests before async context serialization
func RegisterAsyncCallToChild(host arwen.VMHost, config interface{}, arguments [][]byte) error {
	testConfig := config.(*test.TestConfig)
	callData := txDataBuilder.NewBuilder()
	callData.Func(AsyncChildFunction)
	callData.Int64(testConfig.TransferToThirdParty)
	callData.Str(AsyncChildData)
	callData.Bytes(arguments[0])

	value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()
	async := host.Async()
	if !testConfig.IsLegacyAsync {
		err := host.Async().RegisterAsyncCall("testGroup", &arwen.AsyncCall{
			Status:          arwen.AsyncCallPending,
			Destination:     testConfig.GetChildAddress(),
			Data:            callData.ToBytes(),
			ValueBytes:      value,
			SuccessCallback: testConfig.SuccessCallback,
			ErrorCallback:   testConfig.ErrorCallback,
			GasLimit:        testConfig.GasProvidedToChild,
			GasLocked:       testConfig.GasToLock,
			CallbackClosure: ExpectedClosure,
		})
		if err != nil {
			return err
		}
		return nil
	} else {
		return async.RegisterLegacyAsyncCall(testConfig.GetChildAddress(), callData.ToBytes(), value)
	}
}

// SimpleCallbackMock is an exposed mock contract method
func SimpleCallbackMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		arguments := host.Runtime().Arguments()

		err := host.Metering().UseGasBounded(testConfig.GasUsedByCallback)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		if string(arguments[1]) == "fail" {
			host.Runtime().SignalUserError("callback failed intentionally")
			return instance
		}

		if string(arguments[1]) == "new_async" {
			destination := arguments[2]
			function := arguments[3]
			err = host.Async().RegisterAsyncCall("testGroup", &arwen.AsyncCall{
				Status:          arwen.AsyncCallPending,
				Destination:     destination,
				Data:            function,
				ValueBytes:      big.NewInt(0).Bytes(),
				SuccessCallback: "callBack",
				ErrorCallback:   "callBack",
				GasLimit:        testConfig.GasProvidedToChild,
				GasLocked:       150,
			})
			if err != nil {
				host.Runtime().FailExecution(err)
			}
		}

		return instance
	})
}

// CallBackParentMock is an exposed mock contract method
func CallBackParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	callbackFunc := func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		arguments := host.Runtime().Arguments()

		if !testConfig.IsLegacyAsync {
			managed := host.ManagedTypes()
			closureHandle := managed.NewManagedBuffer()
			elrondapi.GetCallbackClosureWithHost(host, closureHandle)
			closure, err := managed.GetBytes(closureHandle)
			if err != nil {
				host.Runtime().SignalUserError("can't get closure")
				return instance
			}
			if err != nil || !bytes.Equal(closure, ExpectedClosure) {
				host.Runtime().SignalUserError("can't get closure")
				return instance
			}
		}

		err := host.Metering().UseGasBounded(testConfig.GasUsedByCallback)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		if len(arguments) < 2 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		loadedData, _, err := host.Storage().GetStorage(test.ParentKeyB)
		if err != nil {
			host.Runtime().FailExecution(err)
			return instance
		}

		status := bytes.Compare(loadedData, test.ParentDataB)
		if status != 0 {
			status = 1
		}

		if len(arguments) >= 4 {
			err := handleParentBehaviorArgument(host, big.NewInt(0).SetBytes(arguments[1]))
			if err != nil {
				return instance
			}
		}
		err = handleTransferToVault(host, arguments)
		if err != nil {
			host.Runtime().SignalUserError(err.Error())
			return instance
		}

		finishResult(host, status)

		_, _ = host.Storage().SetStorage(test.CallbackKey, test.CallbackData)

		return instance
	}

	instanceMock.AddMockMethod("callBack", callbackFunc)
	instanceMock.AddMockMethod("myCallBack", callbackFunc)
}

// CallbackWithOnSameContext is an exposed mock contract method
func CallbackWithOnSameContext(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		retVal := elrondapi.ExecuteOnSameContextWithTypedArgs(
			host,
			int64(host.Metering().GasLeft()),
			big.NewInt(0),
			[]byte("executedOnSameContextByCallback"),
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

func handleParentBehaviorArgument(host arwen.VMHost, behavior *big.Int) error {
	if behavior.Cmp(big.NewInt(3)) == 0 {
		host.Runtime().SignalUserError("callBack error")
		return errors.New("behavior / parent error")
	}
	if behavior.Cmp(big.NewInt(4)) == 0 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}

	behaviorBytes := behavior.Bytes()
	if len(behaviorBytes) == 0 {
		behaviorBytes = []byte{0}
	}
	host.Output().Finish(behaviorBytes)

	return nil
}

func mustTransferToVault(arguments [][]byte) bool {
	vault := "vault"
	numArgs := len(arguments)
	if numArgs == 3 {
		if string(arguments[2]) == vault {
			return false
		}
	}

	if numArgs == 4 {
		if string(arguments[3]) == vault {
			return false
		}
	}

	return true
}

func handleTransferToVault(host arwen.VMHost, arguments [][]byte) error {
	err := error(nil)
	if mustTransferToVault(arguments) {
		valueToTransfer := big.NewInt(4)
		err = host.Output().Transfer(test.VaultAddress, host.Runtime().GetContextAddress(), 0, 0, valueToTransfer, nil, arguments[1], 0)
	}

	return err
}

func finishResult(host arwen.VMHost, result int) {
	outputContext := host.Output()
	if result == 0 {
		outputContext.Finish([]byte("succ"))
	}
	if result == 1 {
		outputContext.Finish([]byte("fail"))
	}
	if result != 0 && result != 1 {
		outputContext.Finish([]byte("unkn"))
	}
}

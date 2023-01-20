package contracts

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	mock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooks"
	"github.com/stretchr/testify/require"
)

var AsyncChildFunction = "transferToThirdParty"
var AsyncChildData = " there"

// PerformAsyncCallParentMock is an exposed mock contract method
func PerformAsyncCallParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("performAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(testConfig.GasUsedByParent)

		_, _ = host.Storage().SetStorage(test.ParentKeyA, test.ParentDataA)
		_, _ = host.Storage().SetStorage(test.ParentKeyB, test.ParentDataB)
		host.Output().Finish(test.ParentFinishA)
		host.Output().Finish(test.ParentFinishB)

		scAddress := host.Runtime().GetContextAddress()
		transferValue := big.NewInt(testConfig.TransferToThirdParty)
		err := host.Output().Transfer(test.ThirdPartyAddress, scAddress, 0, 0, transferValue, []byte("hello"), 0)
		require.Nil(t, err)

		arguments := host.Runtime().Arguments()

		callData := txDataBuilder.NewBuilder()
		// function to be called on child
		callData.Func(AsyncChildFunction)
		// value to send to third party
		callData.Int64(testConfig.TransferToThirdParty)
		// data for child -> third party tx
		callData.Str(AsyncChildData)
		// behavior param for child
		callData.Bytes(arguments[0])

		// amount to transfer from parent to child
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		err = host.Runtime().ExecuteAsyncCall(test.ChildAddress, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance

	})
}

// SimpleCallbackMock is an exposed mock contract method
func SimpleCallbackMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.GasUsedByCallback)

		if string(arguments[1]) == "fail" {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		return instance
	})
}

// CallBackParentMock is an exposed mock contract method
func CallBackParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.GasUsedByCallback)

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
		require.Nil(t, err)

		finishResult(host, status)

		return instance
	})
}

// CallbackWithOnSameContext is an exposed mock contract method
func CallbackWithOnSameContext(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		retVal := vmhooks.ExecuteOnSameContextWithTypedArgs(
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

func handleParentBehaviorArgument(host vmhost.VMHost, behavior *big.Int) error {
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

func handleTransferToVault(host vmhost.VMHost, arguments [][]byte) error {
	err := error(nil)
	if mustTransferToVault(arguments) {
		valueToTransfer := big.NewInt(4)
		err = host.Output().Transfer(test.VaultAddress, host.Runtime().GetContextAddress(), 0, 0, valueToTransfer, arguments[1], 0)
	}

	return err
}

func finishResult(host vmhost.VMHost, result int) {
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

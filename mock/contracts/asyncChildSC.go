package contracts

import (
	"errors"
	"math/big"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
)

func TransferToAsyncParentOnCallbackChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("transferToThirdParty", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		host.Metering().UseGas(testConfig.GasUsedByChild)

		runtime := host.Runtime()
		output := host.Output()

		vmInput := runtime.GetVMInput()
		scAddress := host.Runtime().GetContextAddress()
		arguments := host.Runtime().Arguments()

		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])

		output.Transfer(
			vmInput.CallerAddr,
			scAddress,
			0,
			0,
			valueToTransfer,
			nil,
			nil,
			vm.DirectCall)
		return instance
	})
}

// TransferToThirdPartyAsyncChildMock is an exposed mock contract method
func TransferToThirdPartyAsyncChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("transferToThirdParty", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		metering := host.Metering()
		err := metering.UseGasBounded(testConfig.GasUsedByChild)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		arguments := host.Runtime().Arguments()
		outputContext := host.Output()

		if len(arguments) != 3 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		behavior := byte(0)
		if len(arguments[2]) != 0 {
			behavior = arguments[2][0]
		}
		err = handleChildBehaviorArgument(host, behavior)
		if err != nil {
			return instance
		}

		scAddress := host.Runtime().GetContextAddress()
		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])
		err = outputContext.Transfer(
			testConfig.GetThirdPartyAddress(),
			scAddress,
			0,
			0,
			valueToTransfer,
			nil,
			arguments[1],
			0)
		if err != nil {
			host.Runtime().SignalUserError(err.Error())
			return instance
		}

		outputContext.Finish([]byte("thirdparty"))

		valueToTransfer = big.NewInt(testConfig.TransferToVault)
		err = outputContext.Transfer(
			testConfig.GetVaultAddress(),
			scAddress,
			0,
			0,
			valueToTransfer,
			nil,
			[]byte{},
			0)
		if err != nil {
			host.Runtime().SignalUserError(err.Error())
			return instance
		}

		outputContext.Finish([]byte("vault"))

		host.Storage().SetStorage(test.ChildKey, test.ChildData)

		return instance
	})
}

// ExecutedOnSameContextByCallback is an exposed mock contract method
func ExecutedOnSameContextByCallback(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("executedOnSameContextByCallback", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Storage().SetStorage(test.ParentKeyB, test.ParentDataA)
		return instance
	})
}

func handleChildBehaviorArgument(host arwen.VMHost, behavior byte) error {
	if behavior == 1 {
		host.Runtime().SignalUserError("child error")
		return errors.New("behavior / child error")
	}
	if behavior == 2 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}
	host.Output().Finish([]byte{behavior})
	return nil
}

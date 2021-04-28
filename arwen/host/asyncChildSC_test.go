package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/stretchr/testify/require"
)

func createTestAsyncChildContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock, testConfig *asyncCallTestConfig) {
	childInstance := imb.CreateAndStoreInstanceMock(t, host, childAddress, testConfig.childBalance)
	addDummyMethodsToInstanceMock(childInstance, gasUsedByChild)
	addAsyncChildMethodsToInstanceMock(childInstance, testConfig)
}

func addAsyncChildMethodsToInstanceMock(instanceMock *mock.InstanceMock, testConfig *asyncCallTestConfig) {
	instanceMock.AddMockMethod("transferToThirdParty", func() {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T

		host.Metering().UseGas(testConfig.gasUsedByChild)

		arguments := host.Runtime().Arguments()
		outputContext := host.Output()

		if len(arguments) != 3 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return
		}

		handleChildBehaviorArgument(host, arguments[2][0])

		scAddress := host.Runtime().GetSCAddress()
		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])
		err := outputContext.Transfer(
			thirdPartyAddress,
			scAddress,
			0,
			0,
			valueToTransfer,
			arguments[1],
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("thirdparty"))

		valueToTransfer = big.NewInt(testConfig.transferFromChildToVault)
		err = outputContext.Transfer(
			vaultAddress,
			scAddress,
			0,
			0,
			valueToTransfer,
			[]byte{},
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("vault"))

		host.Storage().SetStorage(childKey, childData)
	})
}

func handleChildBehaviorArgument(host arwen.VMHost, behavior byte) {
	if behavior == 1 {
		host.Runtime().SignalUserError("child error")
	}
	if behavior == 2 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}

	host.Output().Finish([]byte{behavior})
}

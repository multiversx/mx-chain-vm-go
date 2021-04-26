package host

import (
	"math/big"
	"testing"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/stretchr/testify/require"
)

func createTestAsyncChildContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock) {
	gasUsedByChild := uint64(200)

	childInstance := imb.CreateAndStoreInstanceMock(t, host, childAddress, 0)
	addDummyMethodsToInstanceMock(childInstance, gasUsedByChild)
	addAsyncChildMethodsToInstanceMock(childInstance, gasUsedByChild)
}

func addAsyncChildMethodsToInstanceMock(instance *mock.InstanceMock, gasPerCall uint64) {

	t := instance.T

	handleBehaviorArgument := func(behavior byte) {
		host := instance.Host
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

	instance.AddMockMethod("transferToThirdParty", func() {
		host := instance.Host

		host.Metering().UseGas(gasPerCall)

		arguments := host.Runtime().Arguments()
		outputContext := host.Output()

		if len(arguments) != 3 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return
		}

		handleBehaviorArgument(arguments[1][0])

		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])
		err := outputContext.Transfer(thirdPartyAddress, host.Runtime().GetSCAddress(), 0, 0, valueToTransfer, arguments[1], 0)
		require.Nil(t, err)
		outputContext.Finish([]byte("thirdparty"))

		valueToTransfer = big.NewInt(4)
		err = outputContext.Transfer(vaultAddress, host.Runtime().GetSCAddress(), 0, 0, valueToTransfer, []byte{}, 0)
		require.Nil(t, err)
		outputContext.Finish([]byte("vault"))

		host.Storage().SetStorage(childKey, childData)
	})
}

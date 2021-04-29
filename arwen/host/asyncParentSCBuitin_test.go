package host

import (
	"math/big"
	"testing"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func createTestAsyncBuiltinParentContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock, testConfig *asyncCallBaseTestConfig) {
	parentInstance := imb.CreateAndStoreInstanceMock(t, host, parentAddress, testConfig.parentBalance)
	addAsyncBuiltinParentMethodsToInstanceMock(parentInstance, testConfig)
}

func addAsyncBuiltinParentMethodsToInstanceMock(instanceMock *mock.InstanceMock, testConfig *asyncCallBaseTestConfig) {
	input := DefaultTestContractCallInput()
	input.GasProvided = testConfig.gasProvidedToChild

	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.transferFromParentToChild).Bytes()

		host.Metering().UseGas(testConfig.gasUsedByParent)

		destinationForBuildInCall := host.Runtime().GetSCAddress()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.Bytes(destinationForBuildInCall)
		callData.Bytes(arguments[2])

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})

	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.gasUsedByCallback)
		return instance
	})
}

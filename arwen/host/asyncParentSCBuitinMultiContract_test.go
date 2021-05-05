package host

import (
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func forwardAsyncCallMultiContractParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncBuiltInCallTestConfig)
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
}

func callBackMultiContractParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncBuiltInCallTestConfig)
	instanceMock.AddMockMethod("callBack", simpleWasteGasMockMethod(instanceMock, testConfig.gasUsedByCallback))
}

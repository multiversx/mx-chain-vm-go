package host

import (
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func childFunctionAsyncChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncBuiltInCallTestConfig)
	instanceMock.AddMockMethod("childFunction", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.gasUsedByChild)

		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.transferFromChildToParent).Bytes()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

func callBackAsyncChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncBuiltInCallTestConfig)
	instanceMock.AddMockMethod("callBack", simpleWasteGasMockMethod(instanceMock, testConfig.gasUsedByCallback))
}

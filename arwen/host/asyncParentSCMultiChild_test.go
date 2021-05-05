package host

import (
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func forwardAsyncCallMultiChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncCallMultiChildTestConfig)
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.transferFromParentToChild).Bytes()

		host.Metering().UseGas(testConfig.gasUsedByParent)

		for childCall := 0; childCall < testConfig.childCalls; childCall++ {
			callData := txDataBuilder.NewBuilder()
			callData.Func(function)
			// recursiveChildCalls
			callData.BigInt(big.NewInt(1))

			err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
			require.Nil(t, err)
		}

		return instance

	})
}

func callBackMultiChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncCallMultiChildTestConfig)
	instanceMock.AddMockMethod("callBack", simpleWasteGasMockMethod(instanceMock, testConfig.gasUsedByCallback))
}

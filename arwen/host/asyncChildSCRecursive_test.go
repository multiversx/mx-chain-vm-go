package host

import (
	"math/big"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func recursiveAsyncCallRecursiveChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncCallBaseTestConfig)
	instanceMock.AddMockMethod("recursiveAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.gasUsedByChild)

		recursiveChildCalls := big.NewInt(0).SetBytes(arguments[0])
		recursiveChildCalls.Sub(recursiveChildCalls, big.NewInt(1))
		if recursiveChildCalls.Int64() == 0 {
			return instance
		}

		destination := host.Runtime().GetSCAddress()
		function := string("recursiveAsyncCall")
		value := big.NewInt(testConfig.transferFromParentToChild).Bytes()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.BigInt(recursiveChildCalls)

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

func callBackRecursiveChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncCallBaseTestConfig)
	instanceMock.AddMockMethod("callBack", simpleWasteGasMockMethod(instanceMock, testConfig.gasUsedByCallback))
}

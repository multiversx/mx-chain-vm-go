package contracts

import (
	"math/big"

	mock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
	test "github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallRecursiveParentMock is an exposed mock contract method
func ForwardAsyncCallRecursiveParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallRecursiveTestConfig)
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		host.Metering().UseGas(testConfig.GasUsedByParent)

		// only one child call by default
		recursiveChildCalls := big.NewInt(1)
		if len(arguments) > 2 {
			recursiveChildCalls = big.NewInt(0).SetBytes(arguments[2])
		}

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.BigInt(recursiveChildCalls)

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

// CallBackRecursiveParentMock is an exposed mock contract method
func CallBackRecursiveParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallRecursiveTestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

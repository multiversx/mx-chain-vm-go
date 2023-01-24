package contracts

import (
	"math/big"

	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallRecursiveParentMock is an exposed mock contract method
func ForwardAsyncCallRecursiveParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		// only one child call by default
		recursiveChildCalls := big.NewInt(1)
		if len(arguments) > 2 {
			recursiveChildCalls = big.NewInt(0).SetBytes(arguments[2])
		}

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.BigInt(recursiveChildCalls)

		async := host.Async()
		err = async.RegisterLegacyAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

// CallBackRecursiveParentMock is an exposed mock contract method
func CallBackRecursiveParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

package contracts

import (
	"math/big"

	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	mock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallParentBuiltinMock is an exposed mock contract method
func ForwardAsyncCallParentBuiltinMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallBaseTestConfig)
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		host.Metering().UseGas(testConfig.GasUsedByParent)

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

// CallBackParentBuiltinMock is an exposed mock contract method
func CallBackParentBuiltinMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallBaseTestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

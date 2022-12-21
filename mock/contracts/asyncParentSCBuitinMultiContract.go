package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	mock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
	test "github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallMultiContractParentMock is an exposed mock contract method
func ForwardAsyncCallMultiContractParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncBuiltInCallTestConfig)
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		host.Metering().UseGas(testConfig.GasUsedByParent)

		destinationForBuiltInCall := host.Runtime().GetContextAddress()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.Bytes(destinationForBuiltInCall)
		callData.Bytes(arguments[2])

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

// CallBackMultiContractParentMock is an exposed mock contract method
func CallBackMultiContractParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncBuiltInCallTestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

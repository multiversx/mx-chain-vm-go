package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallMultiContractParentMock is an exposed mock contract method
func ForwardAsyncCallMultiContractParentMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		destinationForBuiltInCall := host.Runtime().GetSCAddress()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.Bytes(destinationForBuiltInCall)
		callData.Bytes(arguments[2])

		async := host.Async()
		err = async.RegisterLegacyAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

// CallBackMultiContractParentMock is an exposed mock contract method
func CallBackMultiContractParentMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

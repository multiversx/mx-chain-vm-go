package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func childFunctionAsyncChildMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("childFunction", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		err := host.Metering().UseGasBounded(testConfig.GasUsedByChild)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		destination := arguments[0]
		function := string(arguments[1])
		value := big.NewInt(testConfig.TransferFromChildToParent).Bytes()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)

		async := host.Async()
		err = async.RegisterLegacyAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

func callBackAsyncChildMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

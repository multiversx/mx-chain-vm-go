package contracts

import (
	"math/big"

	mock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
	test "github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

// RecursiveAsyncCallRecursiveChildMock is an exposed mock contract method
func RecursiveAsyncCallRecursiveChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallBaseTestConfig)
	instanceMock.AddMockMethod("recursiveAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.GasUsedByChild)

		recursiveChildCalls := big.NewInt(0).SetBytes(arguments[0]).Uint64()
		recursiveChildCalls = recursiveChildCalls - 1
		if recursiveChildCalls == 0 {
			return instance
		}

		destination := host.Runtime().GetContextAddress()
		function := string("recursiveAsyncCall")
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		callData := txDataBuilder.NewBuilder()
		callData.Func(function)
		callData.BigInt(big.NewInt(int64(recursiveChildCalls)))

		err := host.Runtime().ExecuteAsyncCall(destination, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance
	})
}

// CallBackRecursiveChildMock is an exposed mock contract method
func CallBackRecursiveChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallBaseTestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

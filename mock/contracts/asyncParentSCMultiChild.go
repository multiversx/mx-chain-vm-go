package contracts

import (
	"math/big"

	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	"github.com/multiversx/wasm-vm/arwen"
	mock "github.com/multiversx/wasm-vm/mock/context"
	test "github.com/multiversx/wasm-vm/testcommon"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallMultiChildMock is an exposed mock contract method
func ForwardAsyncCallMultiChildMock(instanceMock *mock.InstanceMock, config interface{}) {
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
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		for childCall := 0; childCall < testConfig.ChildCalls; childCall++ {
			callData := txDataBuilder.NewBuilder()
			callData.Func(function)
			// recursiveChildCalls
			callData.BigInt(big.NewInt(1))
			// child will return this
			callData.BigInt(big.NewInt(int64(childCall)))

			async := host.Async()
			err := async.RegisterAsyncCall("myAsyncGroup", &arwen.AsyncCall{
				Status:          arwen.AsyncCallPending,
				Destination:     destination,
				Data:            callData.ToBytes(),
				ValueBytes:      value,
				GasLimit:        uint64(300),
				SuccessCallback: "callBack",
				ErrorCallback:   "callBack",
			})
			require.Nil(t, err)
		}

		return instance

	})
}

// CallBackMultiChildMock is an exposed mock contract method
func CallBackMultiChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}

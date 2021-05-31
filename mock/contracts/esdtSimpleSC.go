package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen/elrondapi"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
)

// ExecESDTTransferAndCallParentMock is an exposed mock contract method
func ExecESDTTransferAndCallParentMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("execESDTTransferAndCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 3 {
			host.Runtime().SignalUserError("need 3 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.CallerAddr = host.Runtime().GetSCAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.ESDTTestTokenName,
			big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
			arguments[2],
		}
		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])

		_, _, err := host.ExecuteOnDestContext(input)
		if err != nil {
			host.Runtime().FailExecution(err)
		}

		return instance
	})
}

// ExecESDTTransferWithAPICall is an exposed mock contract method
func ExecESDTTransferWithAPICall(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(DirectCallGasTestConfig)
	instanceMock.AddMockMethod("execESDTTransferWithAPICall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 3 {
			host.Runtime().SignalUserError("need 3 arguments")
			return instance
		}

		input := test.DefaultTestContractCallInput()
		input.CallerAddr = host.Runtime().GetSCAddress()
		input.GasProvided = testConfig.GasProvidedToChild
		input.Arguments = [][]byte{
			test.ESDTTestTokenName,
			big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
			arguments[2],
		}
		input.RecipientAddr = arguments[0]

		elrondapi.TransferESDTNFTExecute(
			host,
			input.RecipientAddr,
			test.ESDTTestTokenName,
			testConfig.ESDTTokensToTransfer,
			arguments[1],
			[][]byte{arguments[2]},
			testConfig.GasProvidedToChild)

		return instance
	})
}

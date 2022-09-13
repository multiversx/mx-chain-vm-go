package contracts

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/wasm-vm/arwen/elrondapi"
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	"github.com/ElrondNetwork/wasm-vm/testcommon"
)

// TransferAndExecuteFuncName -
var TransferAndExecuteFuncName = "transferAndExecute"

// TransferAndExecuteReturnData -
var TransferAndExecuteReturnData = []byte{1, 2, 3}

// TransferAndExecute is an exposed mock contract method
func TransferAndExecute(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod(TransferAndExecuteFuncName, func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		noOfTransfers := int(big.NewInt(0).SetBytes(arguments[0]).Int64())

		for transfer := 0; transfer < noOfTransfers; transfer++ {
			elrondapi.TransferValueExecuteWithTypedArgs(host,
				GetChildAddressForTransfer(transfer),
				big.NewInt(testConfig.TransferFromParentToChild),
				int64(testConfig.GasProvidedToChild),
				big.NewInt(int64(transfer)).Bytes(), // transfer data
				[][]byte{},
			)
		}

		host.Output().Finish([]byte(TransferAndExecuteReturnData))

		return instance
	})
}

// GetChildAddressForTransfer -
func GetChildAddressForTransfer(transfer int) []byte {
	return testcommon.MakeTestSCAddress(fmt.Sprintf("childSC-%d", transfer))
}

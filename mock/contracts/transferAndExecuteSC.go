package contracts

import (
	"encoding/hex"
	"fmt"
	"math/big"

	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
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

		_ = host.Metering().UseGasBounded(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		noOfTransfers := int(big.NewInt(0).SetBytes(arguments[0]).Int64())

		for transfer := 0; transfer < noOfTransfers; transfer++ {
			vmhooks.TransferValueExecuteWithTypedArgs(host,
				GetChildAddressForTransfer(transfer),
				big.NewInt(testConfig.TransferFromParentToChild),
				int64(testConfig.GasProvidedToChild),
				big.NewInt(int64(transfer)).Bytes(), // transfer data
				[][]byte{},
			)
		}

		host.Output().Finish(TransferAndExecuteReturnData)

		return instance
	})
}

func TransferEGLDToParent(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("transferEGLDToParent", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		_ = host.Metering().UseGasBounded(testConfig.GasUsedByChild)

		vmhooks.TransferValueExecuteWithTypedArgs(host,
			test.ParentAddress,
			big.NewInt(testConfig.ChildBalance/2),
			0,
			[]byte{}, // transfer data
			[][]byte{},
		)

		return instance
	})
}

func TransferAndExecuteWithBuiltIn(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("transferAndExecuteWithBuiltIn", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		_ = host.Metering().UseGasBounded(testConfig.GasUsedByChild)

		transferString := "ESDTTransfer@" + hex.EncodeToString([]byte("XYY-ABCDEF")) + "@" + hex.EncodeToString(big.NewInt(100000).Bytes())
		vmhooks.TransferValueExecuteWithTypedArgs(host,
			test.UserAddress,
			big.NewInt(0),
			1,
			[]byte(transferString),
			[][]byte{},
		)

		return instance
	})
}

// GetChildAddressForTransfer -
func GetChildAddressForTransfer(transfer int) []byte {
	return testcommon.MakeTestSCAddress(fmt.Sprintf("childSC-%d", transfer))
}

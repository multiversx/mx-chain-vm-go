package contracts

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// BackTransfer_ParentCallsChild is an exposed mock contract method
func BackTransfer_ParentCallsChild(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("callChild", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		storedResult := []byte("ok")

		testConfig := config.(*test.TestConfig)
		input := test.DefaultTestContractCallInput()
		input.GasProvided = testConfig.GasProvidedToChild
		input.CallerAddr = test.ParentAddress
		input.RecipientAddr = test.ChildAddress
		input.Function = "childFunction"
		returnValue := ExecuteOnDestContextInMockContracts(host, input)
		if returnValue != 0 {
			host.Runtime().FailExecution(fmt.Errorf("return value %d", returnValue))
		}
		managedTypes := host.ManagedTypes()
		esdtTransfers, egld := managedTypes.GetBackTransfers()
		if len(esdtTransfers) != 1 {
			host.Runtime().FailExecution(fmt.Errorf("found esdt transfers %d", len(esdtTransfers)))
			storedResult = []byte("err")
		}
		if !bytes.Equal(test.ESDTTestTokenName, esdtTransfers[0].ESDTTokenName) {
			host.Runtime().FailExecution(fmt.Errorf("invalid token name %s", string(esdtTransfers[0].ESDTTokenName)))
			storedResult = []byte("err")
		}
		if big.NewInt(0).SetUint64(testConfig.ESDTTokensToTransfer).Cmp(esdtTransfers[0].ESDTValue) != 0 {
			host.Runtime().FailExecution(fmt.Errorf("invalid token value %d", esdtTransfers[0].ESDTValue.Uint64()))
			storedResult = []byte("err")
		}
		if egld.Cmp(big.NewInt(testConfig.TransferFromChildToParent)) != 0 {
			host.Runtime().FailExecution(fmt.Errorf("invalid egld value %d", egld))
			storedResult = []byte("err")
		}

		host.Storage().SetStorage(test.ParentKeyA, storedResult)

		return instance
	})
}

// BackTransfer_ChildMakesAsync is an exposed mock contract method
func BackTransfer_ChildMakesAsync(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("childFunction", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		testConfig := config.(*test.TestConfig)

		callData := txDataBuilder.NewBuilder()
		callData.Func("calledByAsync")
		callData.Int64(0)

		err := host.Async().RegisterAsyncCall("testGroup", &vmhost.AsyncCall{
			Status:          vmhost.AsyncCallPending,
			Destination:     test.NephewAddress,
			Data:            callData.ToBytes(),
			ValueBytes:      big.NewInt(0).Bytes(),
			SuccessCallback: testConfig.SuccessCallback,
			ErrorCallback:   testConfig.ErrorCallback,
			GasLimit:        uint64(300),
			GasLocked:       testConfig.GasToLock,
			CallbackClosure: nil,
		})
		if err != nil {
			host.Runtime().FailExecution(err)
		}
		return instance
	})
}

// BackTransfer_ChildCallback is an exposed mock contract method
func BackTransfer_ChildCallback(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("myCallback", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		testConfig := config.(*test.TestConfig)

		valueBytes := big.NewInt(testConfig.TransferFromChildToParent).Bytes()
		err := host.Output().Transfer(
			test.ParentAddress,
			test.ChildAddress, 0, 0, big.NewInt(0).SetBytes(valueBytes), nil, []byte{}, vm.DirectCall)
		if err != nil {
			host.Runtime().FailExecution(err)
		}

		transfer := &vmcommon.ESDTTransfer{
			ESDTValue:      big.NewInt(int64(testConfig.ESDTTokensToTransfer)),
			ESDTTokenName:  test.ESDTTestTokenName,
			ESDTTokenType:  0,
			ESDTTokenNonce: 0,
		}

		ret := vmhooks.TransferESDTNFTExecuteWithTypedArgs(
			host,
			test.ParentAddress,
			[]*vmcommon.ESDTTransfer{transfer},
			int64(testConfig.GasProvidedToChild),
			nil,
			nil)
		if ret != 0 {
			host.Runtime().FailExecution(fmt.Errorf("Transfer ESDT failed"))
		}

		return instance
	})
}

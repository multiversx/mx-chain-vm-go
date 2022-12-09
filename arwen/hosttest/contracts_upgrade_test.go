package hosttest

import (
	"testing"

	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	contextmock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
	test "github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
	testcommon "github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
)

func TestUpgrade(t *testing.T) {
	arwen.SetLoggingForTests()
	childAddress := test.MakeTestSCAddress("destAddress")
	childCode := test.GetTestSCCode("init-simple", "../../")

	numMockERC20s := 0

	childContract := test.CreateInstanceContract(childAddress).
		WithCode(childCode).
		WithBalance(1000)
	upgraderContract := test.CreateInstanceContract(test.ParentAddress).
		WithCode(test.GetTestSCCode("upgrader", "../../")).
		WithBalance(1000)
	mockERC20s := createUniqueERC20Contracts(numMockERC20s)

	mockContracts := make([]*test.InstanceTestSmartContract, 2+numMockERC20s)
	copy(mockContracts, mockERC20s)
	mockContracts[numMockERC20s] = childContract
	mockContracts[numMockERC20s+1] = upgraderContract

	testCase := test.BuildInstanceCallTest(t).
		WithContracts(mockContracts...)

	testCase.WithInput(test.CreateTestContractCallInputBuilder().
		WithRecipientAddr(test.ParentAddress).
		WithFunction("dummy").
		WithGasProvided(10_000).
		Build())

	testCase.AndAssertResultsWithoutReset(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
		verify.Ok().
			ReturnData(
				[]byte("dummy text"),
			)
	})

	testCase.WithInput(test.CreateTestContractCallInputBuilder().
		WithRecipientAddr(childAddress).
		WithFunction("dummy").
		WithGasProvided(10_000).
		Build())

	testCase.AndAssertResultsWithoutReset(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
		verify.Ok().
			ReturnData(
				[]byte("dummy text"),
			)
	})

	for i := 0; i < numMockERC20s; i++ {
		testCase.WithInput(createTransferInput(i))
		testCase.AndAssertResultsWithoutReset(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.UserError()
		})
	}

	testCase.WithInput(test.CreateTestContractCallInputBuilder().
		WithRecipientAddr(test.ParentAddress).
		WithFunction("upgradeChildContract").
		WithArguments(childAddress, childCode).
		WithGasProvided(1_000_000).
		Build())

	testCase.AndAssertResults(func(host arwen.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
		verify.Ok()
	})
}

func createUniqueERC20Contracts(nContracts int) []*test.InstanceTestSmartContract {
	mockContracts := make([]*test.InstanceTestSmartContract, nContracts)
	originalCode := testcommon.GetTestSCCode("erc20", "../../")
	for i := 0; i < nContracts; i++ {
		address := createAddress(i)
		modifiedCode := make([]byte, len(originalCode))
		copy(modifiedCode, originalCode)
		modifyERC20BytecodeWithCustomTransferEvent(modifiedCode, []byte{byte(i)})
		mockContracts[i] = test.CreateInstanceContract(address).WithCode(modifiedCode)
	}

	return mockContracts

}

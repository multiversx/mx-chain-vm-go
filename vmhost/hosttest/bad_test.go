package hostCoretest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-vm-go/executor"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/stretchr/testify/require"
)

func TestBadContract_NoPanic_Memoryfault(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("memoryFault").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrorAndInfo(vmhost.ErrExecutionFailed.Error(), "memoryFault")
		})
}

func TestBadContract_NoPanic_DivideByZero(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("divideByZero").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestBadContract_NoPanic_BadGetOwner1(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badGetOwner1").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrors(vmhost.ErrBadBounds.Error())
		})
}

func TestBadContract_NoPanic_BadBigIntStorageStore1(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badBigIntStorageStore1").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestBadContract_NoPanic_BadWriteLog1(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badWriteLog1").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrors("negative length")
		})
}

func TestBadContract_NoPanic_BadWriteLog2(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badWriteLog2").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrors("negative length")
		})
}

func TestBadContract_NoPanic_BadWriteLog3(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badWriteLog3").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestBadContract_NoPanic_BadWriteLog4(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badWriteLog4").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrors("mem load: bad bounds")
		})
}

func TestBadContract_NoPanic_BadGetBlockHash1(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badGetBlockHash1").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrors(vmhost.ErrExecutionFailed.Error())
		})
}

func TestBadContract_NoPanic_BadGetBlockHash2(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badGetBlockHash2").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestBadContract_NoPanic_BadGetBlockHash3(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("badGetBlockHash3").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrors(vmhost.ErrExecutionFailed.Error())
		})
}

func TestBadContract_NoPanic_BadRecursive(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-misc", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(10000000).
			WithFunction("badRecursive").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				ExecutionFailed().
				HasRuntimeErrors(vmhost.ErrExecutionFailed.Error())
		})
}

func TestBadContract_NoPanic_NonExistingFunction(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-empty", "../../")).
				WithBalance(1000)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("thisDoesNotExist").
			Build()).
		WithWasmerSIGSEGVPassthrough(false).
		AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.
				FunctionNotFound().
				HasRuntimeErrorAndInfo(executor.ErrInvalidFunction.Error(), "thisDoesNotExist")
		})
}

func TestBadContractExtra_LongIntLoop_Wasmer1(t *testing.T) {
	if os.Getenv("VMEXECUTOR") != "wasmer1" {
		t.Skip("Skipping test")
	}

	testBadContractExtraLongIntLoop(t, wasmer.ExecutorFactory())
}

func TestBadContractExtra_LongIntLoop_Wasmer2(t *testing.T) {
	testBadContractExtraLongIntLoop(t, wasmer2.ExecutorFactory())
}

func testBadContractExtraLongIntLoop(t *testing.T, executorFactory executor.ExecutorAbstractFactory) {
	testCase := test.BuildInstanceCallTest(t).WithContracts(
		test.CreateInstanceContract(test.ParentAddress).
			WithCode(test.GetTestSCCode("bad-extra", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(test.GasProvided).
			WithFunction("bigLoop").
			Build()).
		WithExecutorFactory(executorFactory).
		WithWasmerSIGSEGVPassthrough(false)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		testCase.AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.OutOfGas()
		})
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		require.FailNow(t, "test timed out")
	}
}

func TestBadContractExtra_NoPanic_BadRecursive(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	testCase := test.BuildInstanceCallTest(t).
		WithWasmerSIGSEGVPassthrough(false).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("bad-recursive", "../../")).
				WithBalance(1000)).
		WithExecutorFactory(wasmer2.ExecutorFactory())

	input := test.CreateTestContractCallInputBuilder().
		WithRecipientAddr(test.ParentAddress).
		WithGasProvided(10000000).
		WithFunction("badRecursive").
		Build()

	repetitions := 25_000

	for i := 0; i < repetitions; i++ {
		testCase.
			WithInput(input).
			AndAssertResultsWithoutReset(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
				verify.ReturnMessage("execution failed")
				verify.ExecutionFailed()
			})
	}
}

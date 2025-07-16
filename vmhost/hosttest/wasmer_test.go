package hostCoretest

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-vm-go/executor"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWASMGlobals_NoGlobals(t *testing.T) {
	value := int64(42)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/noglobals", "noglobals", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("getnumber").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData(
				big.NewInt(value).Bytes(),
			)
		})
}

func TestWASMGlobals_SingleMutable(t *testing.T) {
	value := int64(66561)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/single-mutable", "single-mutable", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("getglobal").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData(
				big.NewInt(value).Bytes(),
			)
		})
}

func TestWASMGlobals_ImportedGlobal(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/imported-global", "imported-global", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("get_imported_global").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMGlobals_MultipleMutables_WithReset(t *testing.T) {
	testCase := test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/multiple-mutable", "multiple-mutable", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("increment_globals").
			Build())

	assertFunc := func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
		verify.Ok().ReturnData(
			[]byte{},
			[]byte{2},
			[]byte{3},
			[]byte{5},
			[]byte{7})
	}

	testCase.AndAssertResultsWithoutReset(assertFunc)
	testCase.AndAssertResultsWithoutReset(assertFunc)
}

func TestWASMGlobals_SingleImmutable(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/single-immutable", "single-immutable", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("getglobal").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWasmMemories_DeployWithoutMemory(t *testing.T) {
	test.BuildInstanceCreatorTest(t).
		WithInput(test.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithContractCode(test.GetTestSCCodeModule("wasmbacking/memoryless", "memoryless", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_NoPages(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-no-pages", "mem-no-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_NoMaxPages(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-no-max-pages", "mem-no-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_SinglePage(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-single-page", "mem-single-page", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_MultiplePages(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-multiple-pages", "mem-multiple-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_MultipleMaxPages(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-multiple-max-pages", "mem-multiple-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_ExceededPages(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-exceeded-pages", "mem-exceeded-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_ExceededMaxPages(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-exceeded-max-pages", "mem-exceeded-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_MinPagesGreaterThanMaxPages(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-min-pages-greater-than-max-pages", "mem-min-pages-greater-than-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_MultipleMemories(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/multiple-memories", "multiple-memories", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_ResetContent(t *testing.T) {
	testCase := test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-content", "mem-content", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build())

	keyword := "ok"
	keywordOffset := 1024

	assertFunc := func(host vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
		verify.Ok().ReturnData([]byte(keyword))
		instance := extractSingleTrackedInstanceFromHost(verify.T, host)
		require.NotNil(verify.T, instance)
		memory := instance.MemDump()
		require.Len(verify.T, memory, int(1*vmhost.WASMPageSize))
		require.Equal(verify.T, keyword, string(memory[keywordOffset:keywordOffset+len(keyword)]))
	}

	testCase.AndAssertResultsWithoutReset(assertFunc)
	testCase.AndAssertResultsWithoutReset(assertFunc)
}

func TestWASMMemories_ResetDataInitializers(t *testing.T) {
	testCase := test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-data-initializer", "mem-data-initializer", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build())

	keyword := "ok"
	keywordOffset := 1024

	assertFunc := func(host vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
		verify.Ok().ReturnData([]byte(keyword))

		instance := extractSingleTrackedInstanceFromHost(verify.T, host)
		require.NotNil(verify.T, instance)

		memory := instance.MemDump()
		require.Len(verify.T, memory, int(1*vmhost.WASMPageSize))
		require.Equal(verify.T, keyword, string(memory[keywordOffset:keywordOffset+len(keyword)]))
	}

	testCase.AndAssertResultsWithoutReset(assertFunc)
	testCase.AndAssertResultsWithoutReset(assertFunc)
}

func TestWASMMemories_WithGrow(t *testing.T) {
	testCase := test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-grow", "mem-grow", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build())

	assertFunc := func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
		verify.Ok().ReturnData(
			big.NewInt(6).Bytes(),
		)
	}

	for i := 0; i < 10; i++ {
		testCase.AndAssertResultsWithoutReset(assertFunc)
	}
}

func TestWASMCreateAndCall(t *testing.T) {
	deployInput := test.CreateTestContractCreateInputBuilder().
		WithGasProvided(100000).
		WithContractCode(test.GetTestSCCode("counter", "../../")).
		WithCallerAddr(test.UserAddress).
		Build()

	world := worldmock.NewMockWorld()
	host := test.NewTestHostBuilder(t).
		WithBlockchainHook(world).
		WithBuiltinFunctions().
		Build()
	world.NewAddressMocks = append(world.NewAddressMocks, &worldmock.NewAddressMock{
		CreatorAddress: test.UserAddress,
		CreatorNonce:   0,
		NewAddress:     test.ParentAddress,
	})
	world.AcctMap.CreateAccount(test.UserAddress, world)
	vmOutput, err := host.RunSmartContractCreate(deployInput)
	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok()
	err = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
	assert.Nil(t, err)

	input := test.CreateTestContractCallInputBuilder().
		WithGasProvided(100000).
		WithFunction("increment").
		WithCallerAddr(test.UserAddress).
		WithRecipientAddr(test.ParentAddress).
		Build()

	for i := 0; i < 10; i++ {
		vmOutput, err = host.RunSmartContractCall(input)
		verify = test.NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok()
	}
}

func TestWASMMiddleware_GlobalsUnreachable(t *testing.T) {
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/middleware-globals", "middleware-globals", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("getglobal").
			Build()).
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func extractSingleTrackedInstanceFromHost(tb testing.TB, host vmhost.VMHost) executor.Instance {
	trackedInstances := host.Runtime().GetInstanceTracker().TrackedInstances()
	require.Len(tb, trackedInstances, 1)

	for _, inst := range trackedInstances {
		return inst
	}

	return nil
}

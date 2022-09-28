package hosttest

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	contextmock "github.com/ElrondNetwork/wasm-vm/mock/context"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
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
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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

	assertFunc := func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
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
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-no-pages", "mem-no-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_NoMaxPages(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-no-max-pages", "mem-no-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_SinglePage(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-single-page", "mem-single-page", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_MultiplePages(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-multiple-pages", "mem-multiple-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_MultipleMaxPages(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-multiple-max-pages", "mem-multiple-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok().ReturnData([]byte{42})
		})
}

func TestWASMMemories_ExceededPages(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-exceeded-pages", "mem-exceeded-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_ExceededMaxPages(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-exceeded-max-pages", "mem-exceeded-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_MinPagesGreaterThanMaxPages(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/mem-min-pages-greater-than-max-pages", "mem-min-pages-greater-than-max-pages", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

func TestWASMMemories_MultipleMemories(t *testing.T) {
	arwen.SetLoggingForTests()
	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmbacking/multiple-memories", "multiple-memories", "../../"))).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithGasProvided(100000).
			WithFunction("main").
			Build()).
		AndAssertResults(func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

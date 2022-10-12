package hosttest

import (
	"testing"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	contextmock "github.com/ElrondNetwork/wasm-vm/mock/context"
	"github.com/ElrondNetwork/wasm-vm/testcommon"
)

func TestForbiddenOps_BulkAndSIMD(t *testing.T) {
	wasmModules := []string{"data-drop", "memory-init", "memory-fill", "memory-copy", "simd"}

	for _, moduleName := range wasmModules {
		testCase := testcommon.BuildInstanceCallTest(t).
			WithContracts(
				testcommon.CreateInstanceContract(testcommon.ParentAddress).
					WithCode(testcommon.GetTestSCCodeModule("forbidden-opcodes/"+moduleName, moduleName, "../../"))).
			WithInput(testcommon.CreateTestContractCallInputBuilder().
				WithGasProvided(100000).
				WithFunction("main").
				Build())

		assertResults := func(_ arwen.VMHost, _ *contextmock.BlockchainHookStub, verify *testcommon.VMOutputVerifier) {
			verify.ContractInvalid()
		}

		testCase.AndAssertResults(assertResults)
	}
}

func TestForbiddenOps_FloatingPoints(t *testing.T) {
	testcommon.BuildInstanceCreatorTest(t).
		WithInput(testcommon.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000).
			WithCallValue(88).
			WithArguments([]byte{2}).
			WithContractCode(testcommon.GetTestSCCode("num-with-fp", "../../")).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(func(_ *contextmock.BlockchainHookStub, verify *testcommon.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}

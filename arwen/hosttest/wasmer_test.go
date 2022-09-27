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
				WithCode(test.GetTestSCCodeModule("wasmglobals/noglobals", "noglobals", "../../"))).
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
	value := int64(66560)

	test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCodeModule("wasmglobals/single-mutable", "single-mutable", "../../"))).
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

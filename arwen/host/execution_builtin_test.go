package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestExecution_ExecuteOnSameContext_BuiltinFunctions(t *testing.T) {
	code := GetTestSCCode("exec-same-ctx-builtin", "../../")
	scBalance := big.NewInt(1000)

	host, stubBlockchainHook := DefaultTestArwenForCall(t, code, scBalance)
	stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
	host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress

	// Run function testBuiltins1
	input.Function = "testBuiltins1"
	input.GasProvided = 100000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	require.NotNil(t, vmOutput)
	expectedVMOutput := expectedVMOutput_SameCtx_BuiltinFunctions_1(code)
	require.Equal(t, expectedVMOutput, vmOutput)

	// Run function testBuiltins2
	input.Function = "testBuiltins2"
	input.GasProvided = 100000

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)

	require.NotNil(t, vmOutput)
	expectedVMOutput = expectedVMOutput_SameCtx_BuiltinFunctions_2(code)
	require.Equal(t, expectedVMOutput, vmOutput)

	// Run function testBuiltins3
	input.Function = "testBuiltins3"
	input.GasProvided = 100000

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)

	require.NotNil(t, vmOutput)
	expectedVMOutput = expectedVMOutput_SameCtx_BuiltinFunctions_3(code)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_AsyncCall_BuiltinFails(t *testing.T) {
	code := GetTestSCCode("async-call-builtin", "../../")
	scBalance := big.NewInt(1000)

	host, stubBlockchainHook := DefaultTestArwenForCall(t, code, scBalance)
	stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
	host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "performAsyncCallToBuiltin"
	input.Arguments = [][]byte{{1}}
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{[]byte("hello"), {4}}, vmOutput.ReturnData)
}

func TestExecution_AsyncCall_CallBackFailsBeforeExecution(t *testing.T) {
	config.AsyncCallbackGasLockForTests = uint64(2)

	code := GetTestSCCode("async-call-builtin", "../../")
	scBalance := big.NewInt(1000)

	host, stubBlockchainHook := DefaultTestArwenForCall(t, code, scBalance)
	stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
	host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "performAsyncCallToBuiltin"
	input.Arguments = [][]byte{{1}}
	input.GasProvided = 1000000
	input.CurrentTxHash = []byte("txhash")

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{[]byte("hello"), []byte("out of gas"), []byte("txhash")}, vmOutput.ReturnData)
}

func dummyProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	outPutAccounts := make(map[string]*vmcommon.OutputAccount)
	outPutAccounts[string(parentAddress)] = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(0),
		Address:      parentAddress}

	if input.Function == "builtinClaim" {
		outPutAccounts[string(parentAddress)].BalanceDelta = big.NewInt(42)
		return &vmcommon.VMOutput{GasRemaining: 400, OutputAccounts: outPutAccounts}, nil
	}
	if input.Function == "builtinDoSomething" {
		return &vmcommon.VMOutput{OutputAccounts: outPutAccounts, GasRemaining: input.GasProvided}, nil
	}
	if input.Function == "builtinFail" {
		return &vmcommon.VMOutput{
			GasRemaining:  0,
			GasRefund:     big.NewInt(0),
			ReturnCode:    vmcommon.UserError,
			ReturnMessage: "whatdidyoudo",
		}, nil
	}

	return nil, arwen.ErrFuncNotFound
}

func getDummyBuiltinFunctionNames() vmcommon.FunctionNames {
	names := make(vmcommon.FunctionNames)

	var empty struct{}
	names["builtinClaim"] = empty
	names["builtinDoSomething"] = empty
	names["builtinFail"] = empty
	return names
}

func TestESDT_SimpleTransferFromSC(t *testing.T) {
}

package host

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
	"github.com/stretchr/testify/require"
)

var ESDTTransferGasCost = uint64(1)
var ESDTTestTokenName = []byte("TT")
var ESDTTestTokenKey = worldmock.MakeTokenKey(ESDTTestTokenName, 0)

func TestExecution_ExecuteOnDestContext_ESDTTransferWithoutExecute(t *testing.T) {
	code := GetTestSCCodeModule("exec-dest-ctx-esdt/basic", "basic", "../../")
	scBalance := big.NewInt(1000)
	host, world := defaultTestArwenForCallWithWorldMock(t, code, scBalance)

	tokenKey := worldmock.MakeTokenKey(ESDTTestTokenName, 0)
	err := world.BuiltinFuncs.SetTokenData(parentAddress, tokenKey, &esdt.ESDigitalToken{
		Value: big.NewInt(100),
		Type:  uint32(core.Fungible),
	})
	require.Nil(t, err)

	input := DefaultTestContractCallInput()
	input.Function = "basic_transfer"
	input.GasProvided = 100000
	input.ESDTTokenName = ESDTTestTokenName
	input.ESDTValue = big.NewInt(16)

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok()
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Claim(t *testing.T) {
	parentGasUsed := uint64(1988)
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(gasProvided).
			withFunction("callBuiltinClaim").
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(parentAddress, 42).
				GasUsed(parentAddress, parentGasUsed).
				GasRemaining(gasProvided - parentGasUsed).
				ReturnData([]byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_DoSomething(t *testing.T) {
	parentGasUsed := uint64(1992)
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(gasProvided).
			withFunction("callBuiltinDoSomething").
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, big.NewInt(0).Sub(arwen.One, arwen.One).Int64()).
				GasUsed(parentAddress, parentGasUsed).
				GasRemaining(gasProvided - parentGasUsed).
				ReturnData([]byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Nonexistent(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(gasProvided).
			withFunction("callNonexistingBuiltin").
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				ReturnMessage(arwen.ErrFuncNotFound.Error()).
				GasRemaining(0)
		})
}

func TestExecution_ExecuteOnDestContext_MockBuiltinFunctions_Fail(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-builtin", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(gasProvided).
			withFunction("callBuiltinFail").
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				ReturnMessage("whatdidyoudo").
				GasRemaining(0)
		})
}

func TestExecution_AsyncCall_MockBuiltinFails(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("async-call-builtin", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(gasProvided).
			withFunction("performAsyncCallToBuiltin").
			withArguments([]byte{1}).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte("hello"), []byte{10})
		})
}

func TestESDT_GettersAPI(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exchange", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(gasProvided).
			withFunction("validateGetters").
			withESDTValue(big.NewInt(5)).
			withESDTTokenName(ESDTTestTokenName).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.ProcessBuiltInFunctionCalled = dummyProcessBuiltInFunction
			host.protocolBuiltinFunctions = getDummyBuiltinFunctionNames()
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestESDT_GettersAPI_ExecuteAfterBuiltinCall(t *testing.T) {
	host, world := defaultTestArwenWithWorldMock(t)

	initialESDTTokenBalance := uint64(1000)

	// Deploy the "parent" contract, which will call the exchange; the actual
	// code of the contract is not important, because the exchange will be called
	// by the "parent" using a manual call to host.ExecuteOnDestContext().
	dummyCode := GetTestSCCode("init-simple", "../../")
	parentAccount := world.AcctMap.CreateSmartContractAccount(userAddress, parentAddress, dummyCode)
	parentAccount.SetTokenBalanceUint64(ESDTTestTokenKey, initialESDTTokenBalance)

	// Deploy the exchange contract, which will receive ESDT and verify that it
	// can see the received token amount and token name.
	exchangeAddress := MakeTestSCAddress("exchange")
	exchangeCode := GetTestSCCode("exchange", "../../")
	exchange := world.AcctMap.CreateSmartContractAccount(userAddress, exchangeAddress, exchangeCode)
	exchange.Balance = big.NewInt(1000)

	host.InitState()

	// Prepare Arwen to appear as if the parent contract is being executed
	input := DefaultTestContractCallInput()
	host.Runtime().InitStateFromContractCallInput(input)
	_ = host.Runtime().StartWasmerInstance(dummyCode, input.GasProvided, true)
	err := host.Output().TransferValueOnly(input.RecipientAddr, input.CallerAddr, input.CallValue, false)
	require.Nil(t, err)

	// Transfer ESDT to the exchange and call its "validateGetters" method
	esdtValue := int64(5)
	input.CallerAddr = parentAddress
	input.RecipientAddr = exchangeAddress
	input.Function = core.BuiltInFunctionESDTTransfer
	input.GasProvided = 10000
	input.Arguments = [][]byte{
		ESDTTestTokenName,
		big.NewInt(esdtValue).Bytes(),
		[]byte("validateGetters"),
	}

	vmOutput, asyncInfo, err := host.ExecuteOnDestContext(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok()

	require.Zero(t, len(asyncInfo.AsyncContextMap))

	parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(ESDTTestTokenKey)
	require.Equal(t, initialESDTTokenBalance-uint64(esdtValue), parentESDTBalance)

	host.Clean()
}

func dummyProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	outputAccounts := make(map[string]*vmcommon.OutputAccount)
	outputAccounts[string(parentAddress)] = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(0),
		Address:      parentAddress}

	if input.Function == "builtinClaim" {
		outputAccounts[string(parentAddress)].BalanceDelta = big.NewInt(42)
		return &vmcommon.VMOutput{
			GasRemaining:   400 + input.GasLocked,
			OutputAccounts: outputAccounts,
		}, nil
	}
	if input.Function == "builtinDoSomething" {
		return &vmcommon.VMOutput{
			GasRemaining:   400 + input.GasLocked,
			OutputAccounts: outputAccounts,
		}, nil
	}
	if input.Function == "builtinFail" {
		return nil, errors.New("whatdidyoudo")
	}
	if input.Function == core.BuiltInFunctionESDTTransfer {
		vmOutput := &vmcommon.VMOutput{
			GasRemaining: 0,
		}
		function := string(input.Arguments[2])
		esdtTransferTxData := function
		for _, arg := range input.Arguments[3:] {
			esdtTransferTxData += "@" + hex.EncodeToString(arg)
		}
		outTransfer := vmcommon.OutputTransfer{
			Value:         big.NewInt(0),
			GasLimit:      input.GasProvided - ESDTTransferGasCost + input.GasLocked,
			Data:          []byte(esdtTransferTxData),
			CallType:      vmcommon.AsynchronousCall,
			SenderAddress: input.CallerAddr,
		}
		vmOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
		vmOutput.OutputAccounts[string(input.RecipientAddr)] = &vmcommon.OutputAccount{
			Address:         input.RecipientAddr,
			OutputTransfers: []vmcommon.OutputTransfer{outTransfer},
		}
		// TODO when ESDT token balance querying is implemented, ensure the
		// transfers that happen here are persisted in the mock accounts
		return vmOutput, nil
	}

	return nil, arwen.ErrFuncNotFound
}

func getDummyBuiltinFunctionNames() vmcommon.FunctionNames {
	names := make(vmcommon.FunctionNames)

	var empty struct{}
	names["builtinClaim"] = empty
	names["builtinDoSomething"] = empty
	names["builtinFail"] = empty
	names[core.BuiltInFunctionESDTTransfer] = empty

	return names
}

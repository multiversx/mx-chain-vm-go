package testcommon

import (
	"testing"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/wasm-vm/arwen/host"
	"github.com/ElrondNetwork/wasm-vm/arwen/mock"
	"github.com/ElrondNetwork/wasm-vm/config"
	"github.com/ElrondNetwork/wasm-vm/executor"
	contextmock "github.com/ElrondNetwork/wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
	"github.com/stretchr/testify/require"
)

// TestCreateTemplateConfig holds the data to build a contract creation test
type TestCreateTemplateConfig struct {
	t                        *testing.T
	address                  []byte
	input                    *vmcommon.ContractCreateInput
	setup                    func(arwen.VMHost, *contextmock.BlockchainHookStub)
	assertResults            func(*contextmock.BlockchainHookStub, *VMOutputVerifier)
	host                     arwen.VMHost
	gasSchedule              config.GasScheduleMap
	wasmerSIGSEGVPassthrough bool
	executorFactory          executor.ExecutorFactory
	stubAccountInitialNonce  uint64
	blockchainHookStub       *contextmock.BlockchainHookStub
}

// BuildInstanceCreatorTest starts the building process for a contract creation test
func BuildInstanceCreatorTest(t *testing.T) *TestCreateTemplateConfig {
	return &TestCreateTemplateConfig{
		t:                        t,
		setup:                    func(arwen.VMHost, *contextmock.BlockchainHookStub) {},
		gasSchedule:              config.MakeGasMapForTests(),
		wasmerSIGSEGVPassthrough: true,
		executorFactory:          wasmer.ExecutorFactory(),
		stubAccountInitialNonce:  24,
	}
}

// WithExecutor allows caller to choose the Executor type.
func (callerTest *TestCreateTemplateConfig) WithExecutor(executorFactory executor.ExecutorFactory) *TestCreateTemplateConfig {
	callerTest.executorFactory = executorFactory
	return callerTest
}

// WithInput provides the ContractCreateInput for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithInput(input *vmcommon.ContractCreateInput) *TestCreateTemplateConfig {
	callerTest.input = input
	return callerTest
}

// WithAddress provides the address for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithAddress(address []byte) *TestCreateTemplateConfig {
	callerTest.address = address
	return callerTest
}

// WithSetup provides the setup function for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithSetup(setup func(arwen.VMHost, *contextmock.BlockchainHookStub)) *TestCreateTemplateConfig {
	callerTest.setup = setup
	return callerTest
}

// AndAssertResults provides the function that will aserts the results
func (callerTest *TestCreateTemplateConfig) AndAssertResults(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTest(true)
}

// AndAssertResults provides the function that will aserts the results
func (callerTest *TestCreateTemplateConfig) AndAssertResultsWithoutReset(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTest(false)
}

func (callerTest *TestCreateTemplateConfig) runTest(reset bool) {
	if callerTest.blockchainHookStub == nil {
		callerTest.blockchainHookStub = callerTest.createBlockchainMock()
	}
	if callerTest.host == nil {
		callerTest.host = callerTest.createTestArwenVM()
		callerTest.setup(callerTest.host, callerTest.blockchainHookStub)
	}
	defer func() {
		if reset {
			callerTest.host.Reset()
		}
	}()

	vmOutput, err := callerTest.host.RunSmartContractCreate(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(callerTest.blockchainHookStub, verify)
}

func (callerTest *TestCreateTemplateConfig) createBlockchainMock() *contextmock.BlockchainHookStub {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return callerTest.address, nil
	}
	return stubBlockchainHook
}

func (callerTest *TestCreateTemplateConfig) createTestArwenVM() arwen.VMHost {
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	host, err := arwenHost.NewArwenVM(
		callerTest.blockchainHookStub,
		callerTest.executorFactory,
		&arwen.VMHostParameters{
			VMType:                   DefaultVMType,
			BlockGasLimit:            uint64(1000),
			GasSchedule:              callerTest.gasSchedule,
			BuiltInFuncContainer:     builtInFunctions.NewBuiltInFunctionContainer(),
			ElrondProtectedKeyPrefix: []byte("ELROND"),
			ESDTTransferParser:       esdtTransferParser,
			EpochNotifier:            &mock.EpochNotifierStub{},
			EnableEpochsHandler: &worldmock.EnableEpochsHandlerStub{
				IsStorageAPICostOptimizationFlagEnabledField:         true,
				IsMultiESDTTransferFixOnCallBackFlagEnabledField:     true,
				IsFixOOGReturnCodeFlagEnabledField:                   true,
				IsRemoveNonUpdatedStorageFlagEnabledField:            true,
				IsCreateNFTThroughExecByCallerFlagEnabledField:       true,
				IsManagedCryptoAPIsFlagEnabledField:                  true,
				IsFailExecutionOnEveryAPIErrorFlagEnabledField:       true,
				IsRefactorContextFlagEnabledField:                    true,
				IsCheckCorrectTokenIDForTransferRoleFlagEnabledField: true,
				IsDisableExecByCallerFlagEnabledField:                true,
				IsESDTTransferRoleFlagEnabledField:                   true,
				IsSendAlwaysFlagEnabledField:                         true,
				IsGlobalMintBurnFlagEnabledField:                     true,
				IsCheckFunctionArgumentFlagEnabledField:              true,
				IsCheckExecuteOnReadOnlyFlagEnabledField:             true,
			},
			WasmerSIGSEGVPassthrough: callerTest.wasmerSIGSEGVPassthrough,
		})
	require.Nil(callerTest.t, err)
	require.NotNil(callerTest.t, host)

	return host
}

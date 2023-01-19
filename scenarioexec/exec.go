package scenarioexec

import (
	"fmt"

	logger "github.com/multiversx/mx-chain-logger-go"
	vmi "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	worldhook "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	gasSchedules "github.com/multiversx/mx-chain-vm-v1_4-go/scenarioexec/gasSchedules"
	mc "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/controller"
	er "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/expression/reconstructor"
	fr "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/fileresolver"
	mj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/model"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/mock"
)

var log = logger.GetOrCreate("vm/scenarios")

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// VMTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with VM.
type VMTestExecutor struct {
	World             *worldhook.MockWorld
	vm                vmi.VMExecutionHandler
	vmHost            vmhost.VMHost
	checkGas          bool
	scenarioTraceGas  []bool
	fileResolver      fr.FileResolver
	exprReconstructor er.ExprReconstructor
}

var _ mc.TestExecutor = (*VMTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*VMTestExecutor)(nil)

// NewVMTestExecutor prepares a new VMTestExecutor instance.
func NewVMTestExecutor() (*VMTestExecutor, error) {
	world := worldhook.NewMockWorld()

	return &VMTestExecutor{
		World:             world,
		vm:                nil,
		checkGas:          true,
		scenarioTraceGas:  make([]bool, 0),
		fileResolver:      nil,
		exprReconstructor: er.ExprReconstructor{},
	}, nil
}

// InitVM will initialize the VM and the builtin function container.
// Does nothing if the VM is already initialized.
func (ae *VMTestExecutor) InitVM(scenGasSchedule mj.GasSchedule) error {
	if ae.vm != nil {
		return nil
	}

	gasSchedule, err := ae.gasScheduleMapFromScenarios(scenGasSchedule)
	if err != nil {
		return err
	}

	err = ae.World.InitBuiltinFunctions(gasSchedule)
	if err != nil {
		return err
	}

	blockGasLimit := uint64(10000000)
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldhook.WorldMarshalizer)
	vm, err := hostCore.NewVMHost(ae.World, &vmhost.VMHostParameters{
		VMType:               TestVMType,
		BlockGasLimit:        blockGasLimit,
		GasSchedule:          gasSchedule,
		BuiltInFuncContainer: ae.World.BuiltinFuncs.Container,
		ProtectedKeyPrefix:   []byte(ProtectedKeyPrefix),
		ESDTTransferParser:   esdtTransferParser,
		EpochNotifier:        &mock.EpochNotifierStub{},
		EnableEpochsHandler: &mock.EnableEpochsHandlerStub{
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
			IsGlobalMintBurnFlagEnabledField:                     true,
			IsTransferToMetaFlagEnabledField:                     true,
			IsCheckFrozenCollectionFlagEnabledField:              true,
			IsFixAsyncCallbackCheckFlagEnabledField:              true,
			IsESDTNFTImprovementV1FlagEnabledField:               true,
			IsSaveToSystemAccountFlagEnabledField:                true,
			IsValueLengthCheckFlagEnabledField:                   true,
			IsSCDeployFlagEnabledField:                           true,
			IsRepairCallbackFlagEnabledField:                     true,
			IsAheadOfTimeGasUsageFlagEnabledField:                true,
			IsCheckFunctionArgumentFlagEnabledField:              true,
			IsCheckExecuteOnReadOnlyFlagEnabledField:             true,
		},
		WasmerSIGSEGVPassthrough: false,
		Hasher:                   worldhook.DefaultHasher,
	})
	if err != nil {
		return err
	}

	ae.vm = vm
	ae.vmHost = vm
	return nil
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *VMTestExecutor) GetVM() vmi.VMExecutionHandler {
	return ae.vm
}

// GetVMHost returns the VM host instance in the test executor.
func (ae *VMTestExecutor) GetVMHost() vmhost.VMHost {
	return ae.vmHost
}

func (ae *VMTestExecutor) gasScheduleMapFromScenarios(scenGasSchedule mj.GasSchedule) (config.GasScheduleMap, error) {
	switch scenGasSchedule {
	case mj.GasScheduleDefault:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	case mj.GasScheduleDummy:
		return config.MakeGasMapForTests(), nil
	case mj.GasScheduleV3:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV3())
	case mj.GasScheduleV4:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	default:
		return nil, fmt.Errorf("unknown scenario GasSchedule: %d", scenGasSchedule)
	}
}

// PeekTraceGas -
func (ae *VMTestExecutor) PeekTraceGas() bool {
	length := len(ae.scenarioTraceGas)
	if length != 0 {
		return ae.scenarioTraceGas[length-1]
	}
	return false
}

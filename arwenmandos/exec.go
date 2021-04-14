package arwenmandos

import (
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	vr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/valuereconstructor"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var log = logger.GetOrCreate("arwen/mandos")

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// ArwenTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with Arwen.
type ArwenTestExecutor struct {
	fileResolver            fr.FileResolver
	World                   *worldhook.MockWorld
	vm                      vmi.VMExecutionHandler
	checkGas                bool
	mandosGasScheduleLoaded bool
	valueReconstructor      vr.ValueReconstructor
}

var _ mc.TestExecutor = (*ArwenTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*ArwenTestExecutor)(nil)

// NewArwenTestExecutor prepares a new ArwenTestExecutor instance.
func NewArwenTestExecutor() (*ArwenTestExecutor, error) {
	world := worldhook.NewMockWorld()

	gasScheduleMap := config.MakeGasMapForTests()
	err := world.InitBuiltinFunctions(gasScheduleMap)
	if err != nil {
		return nil, err
	}

	blockGasLimit := uint64(10000000)
	vm, err := arwenHost.NewArwenVM(world, &arwen.VMHostParameters{
		VMType:                   TestVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              gasScheduleMap,
		ProtocolBuiltinFunctions: world.GetBuiltinFunctionNames(),
		ElrondProtectedKeyPrefix: []byte(ElrondProtectedKeyPrefix),
	})
	if err != nil {
		return nil, err
	}

	return &ArwenTestExecutor{
		fileResolver:            nil,
		World:                   world,
		vm:                      vm,
		checkGas:                true,
		mandosGasScheduleLoaded: false,
		valueReconstructor:      vr.ValueReconstructor{},
	}, nil
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *ArwenTestExecutor) GetVM() vmi.VMExecutionHandler {
	return ae.vm
}

func gasScheduleMapFromMandos(mandosGasSchedule mj.GasSchedule) (config.GasScheduleMap, error) {
	switch mandosGasSchedule {
	case mj.GasScheduleDefault:
		return arwenHost.LoadGasScheduleConfig("../../arwenmandos/gasSchedules/gasScheduleV2.toml")
	case mj.GasScheduleDummy:
		return config.MakeGasMapForTests(), nil
	case mj.GasScheduleV1:
		return arwenHost.LoadGasScheduleConfig("../../arwenmandos/gasSchedules/gasScheduleV1.toml")
	case mj.GasScheduleV2:
		return arwenHost.LoadGasScheduleConfig("../../arwenmandos/gasSchedules/gasScheduleV2.toml")
	default:
		return nil, fmt.Errorf("unknown mandos GasSchedule: %d", mandosGasSchedule)
	}
}

// SetMandosGasSchedule updates the gas costs based on the mandos scenario config
// only changes the gas schedule once,
// this prevents subsequent gasSchedule declarations in externalSteps to overwrite
func (ae *ArwenTestExecutor) SetMandosGasSchedule(newGasSchedule mj.GasSchedule) error {
	if ae.mandosGasScheduleLoaded {
		return nil
	}
	ae.mandosGasScheduleLoaded = true
	gasSchedule, err := gasScheduleMapFromMandos(newGasSchedule)
	if err != nil {
		return err
	}
	ae.vm.GasScheduleChange(gasSchedule)
	return nil
}

package arwenmandos

import (
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// ArwenTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with Arwen.
type ArwenTestExecutor struct {
	fileResolver fr.FileResolver
	World        *worldhook.MockWorld
	vm           vmi.VMExecutionHandler
	checkGas     bool
}

var _ mc.TestExecutor = (*ArwenTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*ArwenTestExecutor)(nil)

// NewArwenTestExecutor prepares a new ArwenTestExecutor instance.
func NewArwenTestExecutor() (*ArwenTestExecutor, error) {
	world := worldhook.NewMockWorld()

	blockGasLimit := uint64(10000000)
	vm, err := arwenHost.NewArwenVM(world, &arwen.VMHostParameters{
		VMType:                   TestVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              config.MakeGasMapForTests(),
		ProtocolBuiltinFunctions: world.GetBuiltinFunctionNames(),
		ElrondProtectedKeyPrefix: []byte(ElrondProtectedKeyPrefix),
	})
	if err != nil {
		return nil, err
	}
	return &ArwenTestExecutor{
		fileResolver: nil,
		World:        world,
		vm:           vm,
		checkGas:     true,
	}, nil
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *ArwenTestExecutor) GetVM() vmi.VMExecutionHandler {
	return ae.vm
}

func gasScheduleMapFromMandos(mandosGasSchedule mj.GasSchedule) (config.GasScheduleMap, error) {
	switch mandosGasSchedule {
	case mj.GasScheduleDefault:
		return config.MakeGasMapForTests(), nil // TODO: change to v2 after all tests pass
		// return arwenHost.LoadGasScheduleConfig("../../arwenmandos/gasSchedules/gasScheduleV2.toml")
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

// SetGasSchedule updates the gas costs based on the mandos scenario config.
func (ae *ArwenTestExecutor) setGasSchedule(mandosGasSchedule mj.GasSchedule) error {
	gasSchedule, err := gasScheduleMapFromMandos(mandosGasSchedule)
	if err != nil {
		return err
	}
	ae.vm.GasScheduleChange(gasSchedule)
	return nil
}

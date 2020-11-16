package arwenmandos

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/test/mock-hook-blockchain"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/test/test-util/mandos/controller"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/test/test-util/mandos/json/fileresolver"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// ArwenTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with Arwen.
type ArwenTestExecutor struct {
	fileResolver fr.FileResolver
	World        *worldhook.BlockchainHookMock
	vm           vmi.VMExecutionHandler
	checkGas     bool
}

var _ mc.TestExecutor = (*ArwenTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*ArwenTestExecutor)(nil)

// NewArwenTestExecutor prepares a new ArwenTestExecutor instance.
func NewArwenTestExecutor() (*ArwenTestExecutor, error) {
	world := worldhook.NewMock()
	world.EnableMockAddressGeneration()

	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMapForTests()
	vm, err := arwenHost.NewArwenVM(world, &arwen.VMHostParameters{
		VMType:                   TestVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              gasSchedule,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
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

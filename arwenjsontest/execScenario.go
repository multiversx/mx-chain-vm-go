package arwenjsontest

import (
	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	cryptohook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-crypto"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// ArwenScenarioExecutor parses, interprets and executes .scen.json test scenarios.
type ArwenScenarioExecutor struct {
	world    *worldhook.BlockchainHookMock
	vm       vmi.VMExecutionHandler
	checkGas bool
}

// NewArwenScenarioExecutor prepares a new ArwenScenarioExecutor instance.
func NewArwenScenarioExecutor() (*ArwenScenarioExecutor, error) {
	world := worldhook.NewMock()
	world.EnableMockAddressGeneration()

	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)
	vm, err := arwenHost.NewArwenVM(world, cryptohook.KryptoHookMockInstance, &arwen.VMHostParameters{
		VMType:                   TestVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              gasSchedule,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
	})
	if err != nil {
		return nil, err
	}
	return &ArwenScenarioExecutor{
		world:    world,
		vm:       vm,
		checkGas: false,
	}, nil
}

// Run executes an individual test.
func (se *ArwenScenarioExecutor) Run(scenario *ij.Scenario) error {
	// reset world
	se.world.Clear()

	txIndex := 0
	for _, generalStep := range scenario.Steps {
		switch step := generalStep.(type) {
		case *ij.SetStateStep:
			for _, acct := range step.Accounts {
				se.world.AcctMap.PutAccount(convertAccount(acct))
			}
			se.world.Blockhashes = ij.JSONBytesValues(step.BlockHashes)
		case *ij.CheckStateStep:
			err := checkAccounts(step.CheckAccounts, se.world)
			if err != nil {
				return err
			}
		case *ij.TxStep:
			// execute tx
			output, err := executeTx(step.Tx, se.world, se.vm)
			if err != nil {
				return err
			}

			// check results
			checkGas := se.checkGas && scenario.CheckGas && step.ExpectedResult.CheckGas
			err = checkTxResults(txIndex, step.ExpectedResult, checkGas, output)
			if err != nil {
				return err
			}
			txIndex++
		}

	}

	return nil
}

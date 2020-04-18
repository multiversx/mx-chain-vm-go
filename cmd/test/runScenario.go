package main

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

type arwenScenarioExecutor struct {
	world    *worldhook.BlockchainHookMock
	vm       vmi.VMExecutionHandler
	checkGas bool
}

func newArwenScenarioExecutor() *arwenScenarioExecutor {
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
		panic(err)
	}
	return &arwenScenarioExecutor{
		world:    world,
		vm:       vm,
		checkGas: false,
	}
}

// Run executes an individual test.
func (se *arwenScenarioExecutor) Run(scenario *ij.Scenario) error {
	world := se.world
	vm := se.vm

	// reset world
	world.Clear()

	txIndex := 0
	for _, generalStep := range scenario.Steps {
		switch step := generalStep.(type) {
		case *ij.SetStateStep:
			for _, acct := range step.Accounts {
				world.AcctMap.PutAccount(convertAccount(acct))
			}
			world.Blockhashes = step.BlockHashes
		case *ij.CheckStateStep:
			err := checkAccounts(step.CheckAccounts, world)
			if err != nil {
				return err
			}
		case *ij.TxStep:
			// execute tx
			output, err := executeTx(step.Tx, world, vm)
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

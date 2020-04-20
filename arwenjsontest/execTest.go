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

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

const ignoreGas = true
const ignoreAllLogs = false

// ArwenTestExecutor parses, interprets and executes .test.json tests.
type ArwenTestExecutor struct {
	world    *worldhook.BlockchainHookMock
	vm       vmi.VMExecutionHandler
	checkGas bool
}

// NewArwenTestExecutor prepares a new ArwenTestExecutor instance.
func NewArwenTestExecutor() *ArwenTestExecutor {
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
	return &ArwenTestExecutor{
		world:    world,
		vm:       vm,
		checkGas: false,
	}
}

// Run executes an individual test.
func (te *ArwenTestExecutor) Run(test *ij.Test) error {
	world := te.world
	vm := te.vm

	// reset world
	world.Clear()
	world.Blockhashes = ij.JSONBytesValues(test.BlockHashes)

	for _, acct := range test.Pre {
		world.AcctMap.PutAccount(convertAccount(acct))
	}

	//spew.Dump(world.AcctMap)

	for _, block := range test.Blocks {
		for txIndex, tx := range block.Transactions {
			//fmt.Printf("%d\n", txIndex)
			output, err := executeTx(tx, world, vm)
			if err != nil {
				return err
			}

			blResult := block.Results[txIndex]

			// check results
			checkGas := te.checkGas && test.CheckGas && blResult.CheckGas
			err = checkTxResults(txIndex, blResult, checkGas, output)
			if err != nil {
				return err
			}
		}
	}

	return checkAccounts(test.PostState, world)
}

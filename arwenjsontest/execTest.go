package arwenjsontest

import (
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// Run executes an individual test.
func (ae *ArwenTestExecutor) ExecuteTest(test *ij.Test) error {
	// reset world
	ae.world.Clear()
	ae.world.Blockhashes = ij.JSONBytesValues(test.BlockHashes)

	for _, acct := range test.Pre {
		ae.world.AcctMap.PutAccount(convertAccount(acct))
	}

	for _, block := range test.Blocks {
		for txIndex, tx := range block.Transactions {
			output, err := ae.executeTx(tx)
			if err != nil {
				return err
			}

			blResult := block.Results[txIndex]

			// check results
			checkGas := ae.checkGas && test.CheckGas && blResult.CheckGas
			err = checkTxResults(txIndex, blResult, checkGas, output)
			if err != nil {
				return err
			}
		}
	}

	return checkAccounts(test.PostState, ae.world)
}

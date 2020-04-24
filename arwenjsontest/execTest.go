package arwenjsontest

import (
	"fmt"

	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// ExecuteTest executes an individual test.
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
			txName := fmt.Sprintf("%d", txIndex)
			err = checkTxResults(txName, blResult, test.CheckGas, output)
			if err != nil {
				return err
			}
		}
	}

	return checkAccounts(test.PostState, ae.world)
}

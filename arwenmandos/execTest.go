package arwenmandos

import (
	"fmt"

	mj "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/json/model"
)

// ExecuteTest executes an individual test.
func (ae *ArwenTestExecutor) ExecuteTest(test *mj.Test) error {
	// reset world
	ae.World.Clear()
	ae.World.Blockhashes = mj.JSONBytesFromStringValues(test.BlockHashes)

	for _, acct := range test.Pre {
		ae.World.AcctMap.PutAccount(convertAccount(acct))
	}

	for _, block := range test.Blocks {
		for txIndex, tx := range block.Transactions {
			txName := fmt.Sprintf("%d", txIndex)

			// execute
			output, err := ae.executeTx(txName, tx)
			if err != nil {
				return err
			}

			blResult := block.Results[txIndex]

			// check results
			err = checkTxResults(txName, blResult, test.CheckGas, output)
			if err != nil {
				return err
			}
		}
	}

	return checkAccounts(test.PostState, ae.World)
}

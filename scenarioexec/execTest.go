package scenarioexec

import (
	"fmt"

	mj "github.com/multiversx/mx-chain-vm-go/scenarios/model"
)

// ExecuteTest executes an individual test.
func (ae *ArwenTestExecutor) ExecuteTest(test *mj.Test) error {
	// reset world
	ae.World.Clear()
	ae.World.Blockhashes = test.BlockHashes.ToValues()

	for _, acct := range test.Pre {
		account, err := convertAccount(acct, ae.World)
		if err != nil {
			return err
		}

		ae.World.AcctMap.PutAccount(account)
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
			err = ae.checkTxResults(txName, blResult, test.CheckGas, output)
			if err != nil {
				return err
			}
		}
	}

	baseErrMsg := "Legacy test check: "
	err := ae.checkAccounts(baseErrMsg, test.PostState)
	ae.Close()
	return err
}

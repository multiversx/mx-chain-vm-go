package arwenjsontest

import (
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// ExecuteScenario executes an individual test.
func (ae *ArwenTestExecutor) ExecuteScenario(scenario *ij.Scenario) error {
	// reset world
	ae.world.Clear()

	txIndex := 0
	for _, generalStep := range scenario.Steps {
		switch step := generalStep.(type) {
		case *ij.SetStateStep:
			for _, acct := range step.Accounts {
				ae.world.AcctMap.PutAccount(convertAccount(acct))
			}
			ae.world.Blockhashes = ij.JSONBytesValues(step.BlockHashes)
			ae.world.NewAddressMocks = convertNewAddressMocks(step.NewAddressMocks)
		case *ij.CheckStateStep:
			err := checkAccounts(step.CheckAccounts, ae.world)
			if err != nil {
				return err
			}
		case *ij.TxStep:
			// execute tx
			output, err := ae.executeTx(step.TxIdent, step.Tx)
			if err != nil {
				return err
			}

			// check results
			err = checkTxResults(step.TxIdent, step.ExpectedResult, scenario.CheckGas, output)
			if err != nil {
				return err
			}
			txIndex++
		}

	}

	return nil
}

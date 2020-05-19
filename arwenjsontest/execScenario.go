package arwenjsontest

import (
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// ExecuteScenario executes an individual test.
func (ae *ArwenTestExecutor) ExecuteScenario(scenario *ij.Scenario) error {
	// reset world
	ae.World.Clear()
	ae.checkGas = scenario.CheckGas

	txIndex := 0
	for _, generalStep := range scenario.Steps {
		err := ae.ExecuteStep(generalStep)
		if err != nil {
			return err
		}

		txIndex++
	}

	return nil
}

// ExecuteStep executes a single scenario step and updates mock state.
func (ae *ArwenTestExecutor) ExecuteStep(generalStep ij.Step) error {
	switch step := generalStep.(type) {
	case *ij.SetStateStep:
		for _, acct := range step.Accounts {
			ae.World.AcctMap.PutAccount(convertAccount(acct))
		}
		ae.World.PreviousBlockInfo = convertBlockInfo(step.PreviousBlockInfo)
		ae.World.CurrentBlockInfo = convertBlockInfo(step.CurrentBlockInfo)
		ae.World.Blockhashes = ij.JSONBytesValues(step.BlockHashes)
		ae.World.NewAddressMocks = convertNewAddressMocks(step.NewAddressMocks)
	case *ij.CheckStateStep:
		err := checkAccounts(step.CheckAccounts, ae.World)
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
		if step.ExpectedResult != nil {
			err = checkTxResults(step.TxIdent, step.ExpectedResult, ae.checkGas, output)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

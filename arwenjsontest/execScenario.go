package arwenjsontest

import (
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// ExecuteScenario executes an individual test.
func (ae *ArwenTestExecutor) ExecuteScenario(scenario *ij.Scenario, fileResolver ij.FileResolver) error {
	ae.fileResolver = fileResolver

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
	case *ij.ExternalStepsStep:
		externalStepsRunner := controller.NewScenarioRunner(
			ae,
			ae.fileResolver,
		)
		extAbsPth := ae.fileResolver.ResolveAbsolutePath(step.Path)
		err := externalStepsRunner.RunSingleJSONScenario(extAbsPth)
		if err != nil {
			return err
		}
	case *ij.SetStateStep:
		// append accounts
		for _, acct := range step.Accounts {
			ae.World.AcctMap.PutAccount(convertAccount(acct))
		}

		// replace block info
		ae.World.PreviousBlockInfo = convertBlockInfo(step.PreviousBlockInfo)
		ae.World.CurrentBlockInfo = convertBlockInfo(step.CurrentBlockInfo)
		ae.World.Blockhashes = ij.JSONBytesValues(step.BlockHashes)

		// append NewAddressMocks
		addressMocksToAdd := convertNewAddressMocks(step.NewAddressMocks)
		ae.World.NewAddressMocks = append(ae.World.NewAddressMocks, addressMocksToAdd...)
	case *ij.CheckStateStep:
		err := checkAccounts(step.CheckAccounts, ae.World)
		if err != nil {
			return err
		}
	case *ij.TxStep:
		// execute tx
		_, err := ae.ExecuteTxStep(step)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExecuteTxStep executes a tx step and updates mock state.
func (ae *ArwenTestExecutor) ExecuteTxStep(txStep *ij.TxStep) (*vmi.VMOutput, error) {
	output, err := ae.executeTx(txStep.TxIdent, txStep.Tx)
	if err != nil {
		return nil, err
	}

	// check results
	if txStep.ExpectedResult != nil {
		err = checkTxResults(txStep.TxIdent, txStep.ExpectedResult, ae.checkGas, output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

package arwenmandos

import (
	mc "github.com/ElrondNetwork/arwen-wasm-vm/test/test-util/mandos/controller"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/test/test-util/mandos/json/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/test/test-util/mandos/json/model"
	vmi "github.com/ElrondNetwork/elrond-go/core/vm-common"
)

// Reset clears state/world.
// Is called in RunAllJSONScenariosInDirectory, but not in RunSingleJSONScenario.
func (ae *ArwenTestExecutor) Reset() {
	ae.World.Clear()
}

// ExecuteScenario executes an individual test.
func (ae *ArwenTestExecutor) ExecuteScenario(scenario *mj.Scenario, fileResolver fr.FileResolver) error {
	ae.fileResolver = fileResolver
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
func (ae *ArwenTestExecutor) ExecuteStep(generalStep mj.Step) error {
	switch step := generalStep.(type) {
	case *mj.ExternalStepsStep:
		fileResolverBackup := ae.fileResolver
		externalStepsRunner := mc.NewScenarioRunner(
			ae,
			ae.fileResolver.Clone(),
		)
		extAbsPth := ae.fileResolver.ResolveAbsolutePath(step.Path)
		err := externalStepsRunner.RunSingleJSONScenario(extAbsPth)
		if err != nil {
			return err
		}
		ae.fileResolver = fileResolverBackup
	case *mj.SetStateStep:
		// append accounts
		for _, acct := range step.Accounts {
			ae.World.AcctMap.PutAccount(convertAccount(acct))
		}

		// replace block info
		ae.World.PreviousBlockInfo = convertBlockInfo(step.PreviousBlockInfo)
		ae.World.CurrentBlockInfo = convertBlockInfo(step.CurrentBlockInfo)
		ae.World.Blockhashes = mj.JSONBytesFromStringValues(step.BlockHashes)

		// append NewAddressMocks
		addressMocksToAdd := convertNewAddressMocks(step.NewAddressMocks)
		ae.World.NewAddressMocks = append(ae.World.NewAddressMocks, addressMocksToAdd...)
	case *mj.CheckStateStep:
		err := checkAccounts(step.CheckAccounts, ae.World)
		if err != nil {
			return err
		}
	case *mj.DumpStateStep:
		err := dumpWorld(ae.World)
		if err != nil {
			return err
		}
	case *mj.TxStep:
		// execute tx
		_, err := ae.ExecuteTxStep(step)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExecuteTxStep executes a tx step and updates mock state.
func (ae *ArwenTestExecutor) ExecuteTxStep(txStep *mj.TxStep) (*vmi.VMOutput, error) {
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

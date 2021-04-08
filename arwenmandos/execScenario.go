package arwenmandos

import (
	"fmt"
	"sort"

	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
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
	err := ae.setGasSchedule(scenario.GasSchedule)
	if err != nil {
		return err
	}

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

func (ae *ArwenTestExecutor) ExecuteStep(generalStep mj.Step) error {
	err := error(nil)

	switch step := generalStep.(type) {
	case *mj.ExternalStepsStep:
		err = ae.ExecuteExternalStep(step)
	case *mj.SetStateStep:
		ae.ExecuteSetStateStep(step)
	case *mj.CheckStateStep:
		log.Trace("CheckStateStep", "comment", step.Comment)
		err = ae.checkStateStep(step)
	case *mj.TxStep:
		_, err = ae.ExecuteTxStep(step)
	case *mj.DumpStateStep:
		err = ae.DumpWorld()
	}

	return err
}

// ExecuteExternalStep executes an external step referenced by the scenario.
func (ae *ArwenTestExecutor) ExecuteExternalStep(step *mj.ExternalStepsStep) error {
	log.Trace("ExternalStepsStep", "path", step.Path)
	if len(step.Comment) > 0 {
		log.Trace("ExternalStepsStep", "comment", step.Comment)
	}

	fileResolverBackup := ae.fileResolver
	clonedFileResolver := ae.fileResolver.Clone()
	externalStepsRunner := mc.NewScenarioRunner(ae, clonedFileResolver)

	extAbsPth := ae.fileResolver.ResolveAbsolutePath(step.Path)
	err := externalStepsRunner.RunSingleJSONScenario(extAbsPth)
	if err != nil {
		return err
	}

	ae.fileResolver = fileResolverBackup

	return nil
}

// ExecuteSetStateStep executes a SetStateStep.
func (ae *ArwenTestExecutor) ExecuteSetStateStep(step *mj.SetStateStep) {
	if len(step.Comment) > 0 {
		log.Trace("SetStateStep", "comment", step.Comment)
	}

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
}

// ExecuteTxStep executes a TxStep.
func (ae *ArwenTestExecutor) ExecuteTxStep(step *mj.TxStep) (*vmi.VMOutput, error) {
	log.Trace("ExecuteTxStep", "id", step.TxIdent)
	if len(step.Comment) > 0 {
		log.Trace("ExecuteTxStep", "comment", step.Comment)
	}

	output, err := ae.executeTx(step.TxIdent, step.Tx)
	if err != nil {
		return nil, err
	}

	// check results
	if step.ExpectedResult != nil {
		err = checkTxResults(step.TxIdent, step.ExpectedResult, ae.checkGas, output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

// DumpWorld prints the state of the MockWorld to stdout.
func (ae *ArwenTestExecutor) DumpWorld() error {
	fmt.Print("world state dump:\n")

	for addr, account := range ae.World.AcctMap {
		fmt.Printf("\t%s\n", byteArrayPretty([]byte(addr)))
		fmt.Printf("\t\tnonce: %d\n", account.Nonce)
		fmt.Printf("\t\tbalance: %d\n", account.Balance)

		if len(account.Storage) > 0 {
			var keys []string
			for key := range account.Storage {
				keys = append(keys, key)
			}

			fmt.Print("\t\tstorage:\n")
			sort.Strings(keys)
			for _, key := range keys {
				value := account.Storage[key]
				if len(value) > 0 {
					fmt.Printf("\t\t\t%s => %s\n", byteArrayPretty([]byte(key)), byteArrayPretty(value))
				}
			}
		}
	}

	return nil
}

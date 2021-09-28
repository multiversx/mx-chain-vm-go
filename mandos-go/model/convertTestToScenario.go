package mandosjsonmodel

import (
	"errors"
	"fmt"
)

// ConvertTestToScenario converts the old test format to the new scenario format for tests.
func ConvertTestToScenario(top []*Test) (*Scenario, error) {
	if len(top) != 1 {
		return nil, errors.New("only one test per test file supported")
	}

	test := top[0]

	scenario := &Scenario{
		CheckGas: test.CheckGas,
	}

	if len(test.Blocks) != 1 {
		return nil, errors.New("only one block per test file supported")
	}
	block := test.Blocks[0]

	scenario.Steps = append(scenario.Steps, &SetStateStep{
		Accounts:    test.Pre,
		BlockHashes: test.BlockHashes,
	})

	if len(block.Transactions) != len(block.Results) {
		return nil, errors.New("transactions must match results")
	}
	for txIndex, tx := range block.Transactions {
		scenario.Steps = append(scenario.Steps, &TxStep{
			TxIdent:        fmt.Sprintf("%d", txIndex+1),
			Tx:             tx,
			ExpectedResult: block.Results[txIndex],
		})
	}

	scenario.Steps = append(scenario.Steps, &CheckStateStep{
		CheckAccounts: test.PostState,
	})

	return scenario, nil
}

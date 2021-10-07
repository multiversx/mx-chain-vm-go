package elrondgo_exporter

import (
	"errors"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwenmandos"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/controller"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/esdtconvert"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
)

var errFirstStepMustSetState = errors.New("first step must be of type SetState")

var errNoStepsProvided = errors.New("no steps were provided")

var errStepIsNotTxStep = errors.New("step is not deploy or scCall")

func GetAccountsAndTransactionsFromMandos(mandosTestPath string) (accounts []*testAccount, scAccounts []*testAccount, txs []*transaction, err error) {
	scenario, err := getScenario(mandosTestPath)
	if err != nil {
		return nil, nil, nil, err
	}
	steps := scenario.Steps
	accounts, scAccounts, txs, err = getAccountsAndTransactionsFromSteps(steps)
	if err != nil {
		return nil, nil, nil, err
	}
	return accounts, scAccounts, txs, nil
}

func setAccounts(setStateStep *mj.SetStateStep) (userAccounts []*testAccount, scAccounts []*testAccount, err error) {
	accounts := make([]*testAccount, 0)
	scAccounts = make([]*testAccount, 0)
	// append accounts
	for _, mandosAccount := range setStateStep.Accounts {
		if mandosAccount.Code.Value != nil {
			account, err := convertMandosToTestAccount(mandosAccount)
			if err != nil {
				return nil, nil, err
			}
			accounts = append(accounts, account)
		}
	}
	for _, mandosNewAccount := range setStateStep.NewAddressMocks {
		newScAccount := NewTestAccount().WithAddress(mandosNewAccount.NewAddress.Value)
		scAccounts = append(scAccounts, newScAccount)
	}
	return accounts, scAccounts, nil
}

func getScenario(testPath string) (scenario *mj.Scenario, err error) {
	executor, err := arwenmandos.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}
	runner := mc.NewScenarioRunner(
		executor,
		mc.NewDefaultFileResolver(),
	)
	scenario, err = runner.ParseTestToScenario(testPath)
	if err != nil {
		return nil, err
	}
	return scenario, err
}

func getAccountsAndTransactionsFromSteps(steps []mj.Step) (accounts []*testAccount, scAccounts []*testAccount, txs []*transaction, err error) {
	if len(steps) == 0 {
		return nil, nil, nil, errNoStepsProvided
	}
	if !stepIsSetState(steps[0]) {
		return nil, nil, nil, errFirstStepMustSetState
	}

	switch step := steps[0].(type) {
	case *mj.SetStateStep:
		accounts, scAccounts, err = setAccounts(step)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	txs = make([]*transaction, 0)

	for i := 1; i < len(steps); i++ {
		switch txStep := steps[i].(type) {
		case *mj.TxStep:
			switch txStep.StepTypeName() {
			case "scDeploy":
				arguments := getArguments(txStep.Tx.Arguments)
				tx := CreateDeployTransaction(
					txStep.Tx.Code.Original,
					arguments,
					txStep.Tx.Nonce.Value,
					txStep.Tx.Value.Value,
					txStep.Tx.From.Value,
					txStep.Tx.GasLimit.Value,
					txStep.Tx.GasPrice.Value,
				)
				txs = append(txs, tx)
			case "scCall":
				arguments := getArguments(txStep.Tx.Arguments)
				tx := CreateTransaction(
					txStep.Tx.Function,
					arguments,
					txStep.Tx.Nonce.Value,
					txStep.Tx.Value.Value,
					txStep.Tx.ESDTValue,
					txStep.Tx.From.Value,
					txStep.Tx.To.Value,
					txStep.Tx.GasLimit.Value,
					txStep.Tx.GasPrice.Value,
				)
				txs = append(txs, tx)
			}
		default:
			return nil, nil, nil, errStepIsNotTxStep
		}
	}
	return accounts, scAccounts, txs, nil
}

func convertMandosToTestAccount(mandosAcc *mj.Account) (*testAccount, error) {
	if len(mandosAcc.Address.Value) != 32 {
		return nil, errors.New("bad test: account address should be 32 bytes long")
	}

	storage := make(map[string][]byte)
	for _, stkvp := range mandosAcc.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}
	esdtconvert.WriteESDTToStorage(mandosAcc.ESDTData, storage)
	account := SetNewAccount(mandosAcc.Nonce.Value, mandosAcc.Address.Value, mandosAcc.Balance.Value, storage)
	return account, nil
}

func getArguments(args []mj.JSONBytesFromTree) [][]byte {
	arguments := make([][]byte, len(args))
	for i := 0; i < len(args); i++ {
		arguments[i] = append(arguments[i], args[i].Value...)
	}
	return arguments
}

func stepIsSetState(step mj.Step) bool {
	return step.StepTypeName() == "setState"
}

package elrondgo_exporter

import (
	"errors"
	"math/big"

	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/controller"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/esdtconvert"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
)

var errFirstStepMustSetState = errors.New("first step must be of type SetState")

var errNoStepsProvided = errors.New("no steps were provided")

var errStepIsNotTxStep = errors.New("step is not scCall")

var errTxStepIsNotScCall = errors.New("txStep is not scCall")

var errScAccountMustHaveOwner = errors.New("scAccount must have owner")

var okStatus = big.NewInt(0)

var ScAddressPrefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 5, 0}

var ScAddressPrefixLength = 10

func GetAccountsAndTransactionsFromMandos(mandosTestPath string) (accounts []*TestAccount, txs []*Transaction, deployTxs []*Transaction, err error) {
	scenario, err := getScenario(mandosTestPath)
	if err != nil {
		return nil, nil, nil, err
	}
	steps := scenario.Steps
	accounts, txs, deployTxs, err = getAccountsAndTransactionsFromSteps(steps)
	if err != nil {
		return nil, nil, nil, err
	}
	return accounts, txs, deployTxs, nil
}

func setAccounts(setStateStep *mj.SetStateStep) (accounts []*TestAccount, err error) {
	accounts = make([]*TestAccount, 0)
	// append accounts
	for _, mandosAccount := range setStateStep.Accounts {
		account, err := convertMandosToTestAccount(mandosAccount)
		if err != nil {
			return nil, err
		}
		if len(account.code) > 0 {
			account.address = append(ScAddressPrefix, account.address[ScAddressPrefixLength:]...)
		}
		accounts = append(accounts, account)
	}
	for _, newMandosAddressMock := range setStateStep.NewAddressMocks {
		scAddress := newMandosAddressMock.NewAddress.Value
		ownerAddress := newMandosAddressMock.CreatorAddress.Value
		account := SetNewAccount(0, scAddress, big.NewInt(0), make(map[string][]byte), make([]byte, 0), ownerAddress)
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func getScenario(testPath string) (scenario *mj.Scenario, err error) {
	scenario, err = mc.ParseMandosScenarioDefaultParser(testPath)
	if err != nil {
		return nil, err
	}
	return scenario, err
}

func getAccountsAndTransactionsFromSteps(steps []mj.Step) (accounts []*TestAccount, txs []*Transaction, deployTxs []*Transaction, err error) {
	if len(steps) == 0 {
		return nil, nil, nil, errNoStepsProvided
	}
	if !stepIsSetState(steps[0]) && !stepIsExternalStep(steps[0]) {
		return nil, nil, nil, errFirstStepMustSetState
	}

	txs = make([]*Transaction, 0)
	deployTxs = make([]*Transaction, 0)
	accounts = make([]*TestAccount, 0)

	for i := 0; i < len(steps); i++ {
		switch step := steps[i].(type) {
		case *mj.SetStateStep:

			setStepAccounts, err := setAccounts(step)
			if err != nil {
				return nil, nil, nil, err
			}
			accounts = append(accounts, setStepAccounts...)

		case *mj.TxStep:
			arguments := getArguments(step.Tx.Arguments)
			if step.ExpectedResult.Status.Value.Cmp(okStatus) == 0 {
				switch step.StepTypeName() {
				case "scCall":
					tx := CreateTransaction(
						step.Tx.Function,
						arguments,
						step.Tx.Nonce.Value,
						step.Tx.EGLDValue.Value,
						step.Tx.ESDTValue,
						step.Tx.From.Value,
						append(ScAddressPrefix, step.Tx.To.Value[ScAddressPrefixLength:]...),
						step.Tx.GasLimit.Value,
						step.Tx.GasPrice.Value,
					)
					txs = append(txs, tx)
				case "scDeploy":
					deployTx := CreateDeployTransaction(
						arguments,
						step.Tx.Code.Value,
						step.Tx.From.Value,
						step.Tx.GasLimit.Value,
						step.Tx.GasPrice.Value,
					)
					deployTxs = append(deployTxs, deployTx)
				default:
					return nil, nil, nil, errTxStepIsNotScCall
				}
			}
		case *mj.ExternalStepsStep:
			externalStepAccounts, externalStepTransactions, externalDeployTransactions, err := GetAccountsAndTransactionsFromMandos(step.Path)
			if err != nil {
				return nil, nil, nil, err
			}
			accounts = append(accounts, externalStepAccounts...)
			txs = append(txs, externalStepTransactions...)
			deployTxs = append(deployTxs, externalDeployTransactions...)
		default:
			return nil, nil, nil, errStepIsNotTxStep
		}
	}
	return accounts, txs, deployTxs, nil
}

func convertMandosToTestAccount(mandosAcc *mj.Account) (*TestAccount, error) {
	if len(mandosAcc.Address.Value) != 32 {
		return nil, errors.New("bad test: account address should be 32 bytes long")
	}
	storage := make(map[string][]byte)
	for _, stkvp := range mandosAcc.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}
	esdtconvert.WriteESDTToStorage(mandosAcc.ESDTData, storage)
	account := SetNewAccount(mandosAcc.Nonce.Value, mandosAcc.Address.Value, mandosAcc.Balance.Value, storage, mandosAcc.Code.Value, mandosAcc.Owner.Value)

	if len(account.code) != 0 && len(account.ownerAddress) == 0 {
		return nil, errScAccountMustHaveOwner
	}

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

func stepIsExternalStep(step mj.Step) bool {
	return step.StepTypeName() == "externalSteps"
}

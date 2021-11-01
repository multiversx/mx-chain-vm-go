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

var errScAccountMustHaveOwner = errors.New("scAccount must have owner")

var okStatus = big.NewInt(0)

var ScAddressPrefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 5, 0}

var ScAddressPrefixLength = 10

var benchmarkTxIdent = "benchmark"

var InvalidBenchmarkTxPos = -1

var minimumAcceptedGasPrice = uint64(1)

func GetAccountsAndTransactionsFromMandos(mandosTestPath string) (accounts []*TestAccount, deployedAccounts []*TestAccount, txs []*Transaction, deployTxs []*Transaction, benchmarkTxPos int, err error) {
	scenario, err := getScenario(mandosTestPath)
	if err != nil {
		return nil, nil, nil, nil, -1, err
	}
	steps := scenario.Steps
	accounts, deployedAccounts, txs, deployTxs, benchmarkTxPos, err = getAccountsAndTransactionsFromSteps(steps)
	if err != nil {
		return nil, nil, nil, nil, -1, err
	}
	return accounts, deployedAccounts, txs, deployTxs, benchmarkTxPos, nil
}

func getScenario(testPath string) (scenario *mj.Scenario, err error) {
	scenario, err = mc.ParseMandosScenarioDefaultParser(testPath)
	if err != nil {
		return nil, err
	}
	return scenario, err
}

func getAccountsAndTransactionsFromSteps(steps []mj.Step) (accounts []*TestAccount, deployedAccounts []*TestAccount, txs []*Transaction, deployTxs []*Transaction, benchmarkTxPos int, err error) {
	benchmarkTxPos = -1

	if len(steps) == 0 {
		return nil, nil, nil, nil, InvalidBenchmarkTxPos, errNoStepsProvided
	}
	if !stepIsSetState(steps[0]) && !stepIsExternalStep(steps[0]) {
		return nil, nil, nil, nil, InvalidBenchmarkTxPos, errFirstStepMustSetState
	}

	txs = make([]*Transaction, 0)
	deployTxs = make([]*Transaction, 0)
	accounts = make([]*TestAccount, 0)
	deployedAccounts = make([]*TestAccount, 0)

	for i := 0; i < len(steps); i++ {
		switch step := steps[i].(type) {
		case *mj.SetStateStep:

			setStepAccounts, setStepDeployedAccounts, err := setAccounts(step)
			if err != nil {
				return nil, nil, nil, nil, InvalidBenchmarkTxPos, err
			}
			accounts = append(accounts, setStepAccounts...)
			deployedAccounts = append(deployedAccounts, setStepDeployedAccounts...)

		case *mj.TxStep:
			if step.ExpectedResult.Status.Value.Cmp(okStatus) == 0 {

				if step.Tx.GasPrice.Value == 0 {
					step.Tx.GasPrice.Value = minimumAcceptedGasPrice
				}
				arguments := getArguments(step.Tx.Arguments)
				switch step.StepTypeName() {
				case "scCall":
					if txIdRequieresBenchmark(step.TxIdent) && benchmarkTxPosIsNotSet(benchmarkTxPos) {
						benchmarkTxPos = len(txs)
					}
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
						step.Tx.Code.Original,
						step.Tx.From.Value,
						step.Tx.GasLimit.Value,
						step.Tx.GasPrice.Value,
					)
					deployTxs = append(deployTxs, deployTx)
				default:
					steps = append(steps[:i], steps[i+1:]...)
					i--
				}
			}
		case *mj.ExternalStepsStep:
			externalStepAccounts, externalStepDeployedAccounts, externalStepTransactions, externalDeployTransactions, externalBenchmarkTxPos, err := GetAccountsAndTransactionsFromMandos(step.Path)
			if err != nil {
				return nil, nil, nil, nil, InvalidBenchmarkTxPos, err
			}
			if benchmarkTxPosIsNotSet(benchmarkTxPos) {
				benchmarkTxPos = externalBenchmarkTxPos
			}
			accounts = append(accounts, externalStepAccounts...)
			deployedAccounts = append(deployedAccounts, externalStepDeployedAccounts...)
			txs = append(txs, externalStepTransactions...)
			deployTxs = append(deployTxs, externalDeployTransactions...)
		default:
			steps = append(steps[:i], steps[i+1:]...)
			i--
		}
	}
	return accounts, deployedAccounts, txs, deployTxs, benchmarkTxPos, nil
}

func setAccounts(setStateStep *mj.SetStateStep) (accounts []*TestAccount, deployedAccounts []*TestAccount, err error) {
	accounts = make([]*TestAccount, 0)
	deployedAccounts = make([]*TestAccount, 0)
	for _, mandosAccount := range setStateStep.Accounts {
		account, err := convertMandosToTestAccount(mandosAccount)
		if err != nil {
			return nil, nil, err
		}
		if len(account.code) > 0 {
			account.address = append(ScAddressPrefix, account.address[ScAddressPrefixLength:]...)
		}
		accounts = append(accounts, account)
	}
	for _, newMandosAddressMock := range setStateStep.NewAddressMocks {
		scAddress := append(ScAddressPrefix, newMandosAddressMock.NewAddress.Value[ScAddressPrefixLength:]...)
		ownerAddress := newMandosAddressMock.CreatorAddress.Value
		account := SetNewAccount(0, scAddress, big.NewInt(0), make(map[string][]byte), make([]byte, 0), ownerAddress)
		deployedAccounts = append(deployedAccounts, account)
	}

	return accounts, deployedAccounts, nil
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

func benchmarkTxPosIsNotSet(benchmarkTxPos int) bool {
	return benchmarkTxPos == -1
}

func txIdRequieresBenchmark(txIdent string) bool {
	return txIdent == benchmarkTxIdent
}

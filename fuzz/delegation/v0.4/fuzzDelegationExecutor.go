//nolint:all
package delegation

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	vmi "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	am "github.com/multiversx/mx-chain-vm-v1_4-go/scenarioexec"
	fr "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/fileresolver"
	mjparse "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/json/parse"
	mjwrite "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/json/write"
	mj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/model"
	worldhook "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
)

type fuzzDelegationExecutor struct {
	vmTestExecutor *am.VMTestExecutor
	world             *worldhook.MockWorld
	vm                vmi.VMExecutionHandler
	parser      mjparse.Parser
	txIndex           int

	serviceFee                  int
	numBlocksBeforeForceUnstake int
	numBlocksBeforeUnbond       int
	numDelegators               int
	stakePerNode                *big.Int
	ownerAddress                []byte
	delegationContractAddress   []byte
	auctionMockAddress          []byte
	faucetAddress               []byte
	withdrawTargetAddress       []byte
	stakePurchaseForwardAddress []byte
	numNodes                    int
	totalStakeAdded             *big.Int
	totalStakeWithdrawn         *big.Int
	totalRewards                *big.Int
	generatedScenario           *mj.Scenario
}

func newFuzzDelegationExecutor(fileResolver fr.FileResolver) (*fuzzDelegationExecutor, error) {
	vmTestExecutor, err := am.NewVMTestExecutor()
	if err != nil {
		return nil, err
	}
	parser := mjparse.NewParser(fileResolver)
	return &fuzzDelegationExecutor{
		vmTestExecutor:   vmTestExecutor,
		world:               vmTestExecutor.World,
		vm:                  vmTestExecutor.GetVM(),
		parser:        parser,
		txIndex:             0,
		numNodes:            0,
		totalStakeAdded:     big.NewInt(0),
		totalStakeWithdrawn: big.NewInt(0),
		totalRewards:        big.NewInt(0),
		generatedScenario: &mj.Scenario{
			Name: "fuzz generated",
		},
	}, nil
}

func (pfe *fuzzDelegationExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}

type fuzzDelegationExecutorInitArgs struct {
	serviceFee                  int
	ownerMinStake               int
	numBlocksBeforeForceUnstake int
	numBlocksBeforeUnbond       int
	numDelegators               int
	stakePerNode                *big.Int
	numGenesisNodes             int
}

func (pfe *fuzzDelegationExecutor) addStep(step mj.Step) {
	pfe.generatedScenario.Steps = append(pfe.generatedScenario.Steps, step)
}

func (pfe *fuzzDelegationExecutor) saveGeneratedScenario() {
	vmHost := pfe.vm.(vmhost.VMHost)
	vmHost.Reset()
	serialized := mjwrite.ScenarioToJSONString(pfe.generatedScenario)

	err := ioutil.WriteFile("fuzz_gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (pfe *fuzzDelegationExecutor) nextTxIndex() int {
	pfe.txIndex++
	return pfe.txIndex
}

func (pfe *fuzzDelegationExecutor) getContractBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.delegationContractAddress)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getDelegatorBalance(delegIndex int) *big.Int {
	delegAddr := pfe.delegatorAddress(delegIndex)
	acct := pfe.world.AcctMap.GetAccount(delegAddr)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getAllDelegatorsBalance() *big.Int {
	totalDelegatorBalance := big.NewInt(0)
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		balance := pfe.getDelegatorBalance(delegatorIdx)
		totalDelegatorBalance.Add(totalDelegatorBalance, balance)
	}
	return totalDelegatorBalance
}

func (pfe *fuzzDelegationExecutor) getAuctionBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.auctionMockAddress)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) getWithdrawTargetBalance() *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.withdrawTargetAddress)
	return acct.Balance
}

func (pfe *fuzzDelegationExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.parser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}
	pfe.addStep(step)
	return pfe.vmTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzDelegationExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := pfe.parser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}
	pfe.addStep(step)
	txStep, isTx := step.(*mj.TxStep)
	if !isTx {
		return nil, errors.New("tx step expected")
	}
	return pfe.vmTestExecutor.ExecuteTxStep(txStep)
}

func (pfe *fuzzDelegationExecutor) querySingleResult(funcName string, args string) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "%s",
			"arguments": [
				%s
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "",
			"logs": "*",
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.nextTxIndex(),
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
		funcName,
		args,
	))
	if err != nil {
		return nil, err
	}

	result := big.NewInt(0).SetBytes(output.ReturnData[0])
	pfe.log("query: %s -> %d", funcName, result)
	return result, nil
}

func (pfe *fuzzDelegationExecutor) simpleQuery(funcName string) (*big.Int, error) {
	return pfe.querySingleResult(funcName, "")
}

func (pfe *fuzzDelegationExecutor) delegatorQuery(funcName string, delegIndex int) (*big.Int, error) {
	delegAddr := fmt.Sprintf(`"''%s"`, string(pfe.delegatorAddress(delegIndex)))
	return pfe.querySingleResult(funcName, delegAddr)
}

func (pfe *fuzzDelegationExecutor) increaseBlockNonce(nonceDelta int) error {
	curentBlockNonce := uint64(0)
	if pfe.world.CurrentBlockInfo != nil {
		curentBlockNonce = pfe.world.CurrentBlockInfo.BlockNonce
	}

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "%d - increase block nonce",
		"currentBlockInfo": {
			"blockNonce": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		curentBlockNonce+uint64(nonceDelta),
	))
	if err != nil {
		return err
	}
	pfe.log("block nonce: %d ---> %d", curentBlockNonce, curentBlockNonce+uint64(nonceDelta))
	return nil
}

func (pfe *fuzzDelegationExecutor) dumpState() error {
	return pfe.executeStep(`
	{
		"step": "dumpState"
	}`)
}

func (pfe *fuzzDelegationExecutor) delegatorAddress(delegIndex int) []byte {
	if delegIndex == 0 {
		return pfe.ownerAddress
	}
	return []byte(fmt.Sprintf("delegator %5d               s1", delegIndex))
}

func blsKey(index int) string {
	return fmt.Sprintf(
		"bls key %5d ..................................................................................",
		index)
}

func blsSignature(index int) string {
	return fmt.Sprintf(
		"bls key signature %5d ........",
		index)
}

func blsKeySignatureArgsString(startIndex, numNodes int) string {
	var blsKeyArgs []string
	for i := startIndex; i < startIndex+numNodes; i++ {
		blsKeyArgs = append(blsKeyArgs, "\"''"+blsKey(i)+"\"")
		blsKeyArgs = append(blsKeyArgs, "\"''"+blsSignature(i)+"\"")
	}
	return strings.Join(blsKeyArgs, ",")
}

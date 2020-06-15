package delegation

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	am "github.com/ElrondNetwork/arwen-wasm-vm/arwenmandos"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	mj "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/json/model"
	mjparse "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/json/parse"
	mjwrite "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/json/write"
)

type fuzzDelegationExecutor struct {
	arwenTestExecutor *am.ArwenTestExecutor
	world             *worldhook.BlockchainHookMock
	vm                vmi.VMExecutionHandler
	mandosParser      mjparse.Parser

	initialDelegatorBalance     *big.Int
	serviceFee                  int
	numBlocksBeforeForceUnstake int
	numBlocksBeforeUnbond       int
	numDelegators               int
	stakePerNode                *big.Int
	ownerAddress                []byte
	delegationContractAddress   []byte
	auctionMockAddress          []byte
	numNodes                    int
	generatedScenario           *mj.Scenario
}

func newFuzzDelegationExecutor(fileResolver mjparse.FileResolver) (*fuzzDelegationExecutor, error) {
	arwenTestExecutor, err := am.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}
	parser := mjparse.Parser{
		FileResolver: fileResolver,
	}
	return &fuzzDelegationExecutor{
		arwenTestExecutor: arwenTestExecutor,
		world:             arwenTestExecutor.World,
		vm:                arwenTestExecutor.GetVM(),
		mandosParser:      parser,
		numNodes:          0,
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
	numBlocksBeforeForceUnstake int
	numBlocksBeforeUnbond       int
	numDelegators               int
	stakePerNode                *big.Int
}

func (pfe *fuzzDelegationExecutor) addStep(step mj.Step) {
	pfe.generatedScenario.Steps = append(pfe.generatedScenario.Steps, step)
	pfe.saveGeneratedScenario()
}

func (pfe *fuzzDelegationExecutor) saveGeneratedScenario() {
	serialized := mjwrite.ScenarioToJSONString(pfe.generatedScenario)

	err := ioutil.WriteFile("fuzz_gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (pfe *fuzzDelegationExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}
	pfe.addStep(step)
	return pfe.arwenTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzDelegationExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}
	pfe.addStep(step)
	txStep, isTx := step.(*mj.TxStep)
	if !isTx {
		return nil, errors.New("tx step expected")
	}
	return pfe.arwenTestExecutor.ExecuteTxStep(txStep)
}

func (pfe *fuzzDelegationExecutor) simpleQuery(funcName string) (*big.Int, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-simpleQuery-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "%s",
			"arguments": [],
			"gasLimit": "100,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
		funcName,
	))
	if err != nil {
		return nil, err
	}

	result := big.NewInt(0).SetBytes(output.ReturnData[0])
	pfe.log("simpleQuery: %s -> %d", funcName, result)
	return result, nil
}

func (pfe *fuzzDelegationExecutor) init(args *fuzzDelegationExecutorInitArgs) error {
	pfe.serviceFee = args.serviceFee
	pfe.numBlocksBeforeForceUnstake = args.numBlocksBeforeForceUnstake
	pfe.numBlocksBeforeUnbond = args.numBlocksBeforeUnbond
	pfe.numDelegators = args.numDelegators
	pfe.stakePerNode = args.stakePerNode
	pfe.initialDelegatorBalance, _ = big.NewInt(0).SetString("1000000000000000", 10)

	pfe.world.Clear()

	pfe.ownerAddress = []byte("fuzz_owner_addr_______________s1")
	pfe.delegationContractAddress = []byte("fuzz_sc_delegation_addr_______s1")
	pfe.auctionMockAddress = []byte("fuzz_sc_auction_mock_addr_____s1")

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"''%s": {
				"nonce": "0",
				"balance": "1,000,000,000",
				"storage": {},
				"code": ""
			},
			"''%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {
					"''stake_per_node": "%d"
				},
				"code": "file:auction-mock.wasm"
			}
		},
		"newAddresses": [
			{
				"creatorAddress": "''%s",
				"creatorNonce": "0",
				"newAddress": "''%s"
			}
		]
	}`,
		string(pfe.ownerAddress),
		string(pfe.auctionMockAddress),
		pfe.stakePerNode,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
	))
	if err != nil {
		return err
	}

	// delegators
	for i := 0; i < args.numDelegators; i++ {
		err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"''%s": {
					"nonce": "0",
					"balance": "%d",
					"storage": {},
					"code": ""
				}
			}
		}`,
			string(delegatorAddress(i)),
			pfe.initialDelegatorBalance,
		))
		if err != nil {
			return err
		}
	}

	// deploy delegation
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scDeploy",
		"txId": "-deploy-",
		"tx": {
			"from": "''%s",
			"value": "0",
			"contractCode": "file:delegation.wasm",
			"arguments": [
				"''%s",
				"%d",
				"%d",
				"%d"
			],
			"gasLimit": "100,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.auctionMockAddress),
		args.serviceFee,
		args.numBlocksBeforeForceUnstake,
		args.numBlocksBeforeUnbond,
	))
	if err != nil {
		return err
	}

	err = pfe.setStakePerNode(args.stakePerNode)
	if err != nil {
		return err
	}

	pfe.log("init ok")
	return nil
}

func (pfe *fuzzDelegationExecutor) setStakePerNode(stakePerNode *big.Int) error {
	pfe.log("setStakePerNode: %d", stakePerNode)
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-nr-nodes-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "setStakePerNode",
			"arguments": [
				"%d"
			],
			"gasLimit": "100,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
		stakePerNode,
	))
	return err
}

func (pfe *fuzzDelegationExecutor) enableAutoActivation() error {
	pfe.log("enableAutoActivation")
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-enable-auto-activation-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "enableAutoActivation",
			"arguments": [],
			"gasLimit": "100,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
	))
	return err
}

func (pfe *fuzzDelegationExecutor) addNodes(numNodesToAdd int) error {
	pfe.log("addNodes %d -> %d", numNodesToAdd, pfe.numNodes+numNodesToAdd)

	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-add-nodes-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "addNodes",
			"arguments": [
				%s
			],
			"gasLimit": "1,000,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
		blsKeySignatureArgsString(pfe.numNodes, numNodesToAdd),
	))
	pfe.numNodes += numNodesToAdd
	return err
}

func (pfe *fuzzDelegationExecutor) stake(delegIndex int, amount *big.Int) error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-stake-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "%d",
			"function": "stake",
			"arguments": [],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": "*",
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(delegatorAddress(delegIndex)),
		string(pfe.delegationContractAddress),
		amount,
	))
	pfe.log("stake, delegator: %d, amount: %d", delegIndex, amount)
	return err
}

func (pfe *fuzzDelegationExecutor) withdrawInactiveStake(delegIndex int, amount *big.Int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-unstake-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "withdrawInactiveStake",
			"arguments": [
				"%d"
			],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		}
	}`,
		string(delegatorAddress(delegIndex)),
		string(pfe.delegationContractAddress),
		amount,
	))
	if output.ReturnCode == vmi.Ok {
		pfe.log("unstake, delegator: %d, amount: %d", delegIndex, amount)
	} else {
		pfe.log("unstake, delegator: %d, amount: %d, fail, %s", delegIndex, amount, output.ReturnMessage)
	}
	return err
}

func delegatorAddress(delegIndex int) []byte {
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

func blsKeyArgsString(numNodes int) string {
	var blsKeyArgs []string
	for i := 0; i < numNodes; i++ {
		blsKey := fmt.Sprintf(
			"bls key %5d ..................................................................................",
			i)
		blsKeyArg := "\"''" + blsKey + "\""
		blsKeyArgs = append(blsKeyArgs, blsKeyArg)
	}
	return strings.Join(blsKeyArgs, ",")
}

func blsSignatureArgsString(numNodes int) string {
	var blsSigArgs []string
	for i := 0; i < numNodes; i++ {
		blsSig := fmt.Sprintf(
			"bls key signature %5d ........",
			i)
		blsSigArg := "\"''" + blsSig + "\""
		blsSigArgs = append(blsSigArgs, blsSigArg)
	}
	return strings.Join(blsSigArgs, ",")
}

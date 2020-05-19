package delegation

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ajt "github.com/ElrondNetwork/arwen-wasm-vm/arwenjsontest"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

type fuzzDelegationExecutor struct {
	arwenTestExecutor *ajt.ArwenTestExecutor
	world             *worldhook.BlockchainHookMock
	vm                vmi.VMExecutionHandler
	mandosParser      ij.Parser

	initialDelegatorBalance   *big.Int
	nodeShare                 int
	timeBeforeForceUnstake    int
	numDelegators             int
	numNodes                  int
	stakePerNode              *big.Int
	ownerAddress              []byte
	delegationContractAddress []byte
	auctionMockAddress        []byte
	expectedStake             *big.Int
	active                    bool
}

func newFuzzDelegationExecutor(fileResolver ij.FileResolver) (*fuzzDelegationExecutor, error) {
	arwenTestExecutor, err := ajt.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}
	parser := ij.Parser{
		FileResolver: fileResolver,
	}
	return &fuzzDelegationExecutor{
		arwenTestExecutor: arwenTestExecutor,
		world:             arwenTestExecutor.World,
		vm:                arwenTestExecutor.GetVM(),
		mandosParser:      parser,
		active:            false,
	}, nil
}

func (pfe *fuzzDelegationExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}

type fuzzDelegationExecutorInitArgs struct {
	nodeShare              int
	timeBeforeForceUnstake int
	numDelegators          int
	numNodes               int
	stakePerNode           *big.Int
}

func (pfe *fuzzDelegationExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}
	return pfe.arwenTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzDelegationExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}
	txStep, isTx := step.(*ij.TxStep)
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
	pfe.nodeShare = args.nodeShare
	pfe.timeBeforeForceUnstake = args.timeBeforeForceUnstake
	pfe.numDelegators = args.numDelegators
	pfe.numNodes = args.numNodes
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
				"storage": {},
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
				"%d",
				"''%s",
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
		args.nodeShare,
		string(pfe.auctionMockAddress),
		args.timeBeforeForceUnstake,
	))
	if err != nil {
		return err
	}

	err = pfe.setNumNodes(args.numNodes)
	if err != nil {
		return err
	}

	err = pfe.setStakePerNode(args.stakePerNode)
	if err != nil {
		return err
	}

	err = pfe.setBlsKeys(args.numNodes)
	if err != nil {
		return err
	}

	pfe.expectedStake, err = pfe.simpleQuery("getExpectedStake")
	if err != nil {
		panic(err)
	}

	pfe.log("init ok")
	return nil
}

func (pfe *fuzzDelegationExecutor) setNumNodes(numNodes int) error {
	pfe.log("setNumNodes: %d", numNodes)
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-nr-nodes-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "setNumNodes",
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
		numNodes,
	))
	return err
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

func (pfe *fuzzDelegationExecutor) setBlsKeys(numNodes int) error {
	pfe.log("setBlsKeys")

	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-set-bls-keys-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "setBlsKeys",
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
		blsKeyArgsString(numNodes),
	))
	return err
}

func (pfe *fuzzDelegationExecutor) tryStake(delegIndex int, amount *big.Int) error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
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
		}
	}`,
		string(delegatorAddress(delegIndex)),
		string(pfe.delegationContractAddress),
		amount,
	))
	if output.ReturnCode == vmi.Ok {
		pfe.log("try stake, delegator: %d, amount: %d, ok", delegIndex, amount)
	} else {
		pfe.log("try stake, delegator: %d, amount: %d, fail", delegIndex, amount)
	}
	return err
}

func (pfe *fuzzDelegationExecutor) stakeTheRest(delegIndex int) error {
	filledStake, err := pfe.simpleQuery("getFilledStake")
	if err != nil {
		return err
	}
	expectedStake, err := pfe.simpleQuery("getExpectedStake")
	if err != nil {
		return err
	}
	unfilledStake := big.NewInt(0).Sub(expectedStake, filledStake)

	pfe.log("stake the rest, delegator: %d, amount: %d ", delegIndex, unfilledStake)
	_, err = pfe.executeTxStep(fmt.Sprintf(`
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
		unfilledStake,
	))
	return err
}

func (pfe *fuzzDelegationExecutor) maybeActivate() error {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-activate-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "activate",
			"arguments": [
				%s
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
		blsSignatureArgsString(pfe.numNodes),
	))
	if output.ReturnCode == vmi.Ok {
		pfe.active = true
		pfe.log("try activate, ok")
	} else {
		pfe.log("try activate, fail, %s", output.ReturnMessage)
	}
	return err
}

func (pfe *fuzzDelegationExecutor) mustActivate() error {
	pfe.log("must activate")
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-activate-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "activate",
			"arguments": [
				%s
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [
                    {
                        "address": "''%s",
                        "identifier": "0x0000000000000000000000000000000000000000000000000000000000000003",
                        "topics": [],
                        "data": ""
                    }
                ],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
		blsSignatureArgsString(pfe.numNodes),
		string(pfe.delegationContractAddress),
	))
	pfe.active = true
	return err
}

func delegatorAddress(delegIndex int) []byte {
	return []byte(fmt.Sprintf("delegator %5d               s1", delegIndex))
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

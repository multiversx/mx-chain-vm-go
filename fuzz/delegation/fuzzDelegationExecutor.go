package delegation

import (
	"fmt"
	"math/big"
	"strings"

	ajt "github.com/ElrondNetwork/arwen-wasm-vm/arwenjsontest"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

type fuzzDelegationExecutor struct {
	arwenTestExecutor         *ajt.ArwenTestExecutor
	world                     *worldhook.BlockchainHookMock
	vm                        vmi.VMExecutionHandler
	fileResolver              ij.FileResolver
	initArgs                  *fuzzDelegationExecutorInitArgs
	ownerAddress              []byte
	delegationContractAddress []byte
	auctionMockAddress        []byte
	expectedStake             *big.Int
}

func newFuzzDelegationExecutor(fileResolver ij.FileResolver) (*fuzzDelegationExecutor, error) {
	arwenTestExecutor, err := ajt.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}
	return &fuzzDelegationExecutor{
		arwenTestExecutor: arwenTestExecutor,
		world:             arwenTestExecutor.World,
		vm:                arwenTestExecutor.GetVM(),
		fileResolver:      fileResolver,
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
	p := ij.Parser{
		FileResolver: pfe.fileResolver,
	}
	step, err := p.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}
	return pfe.arwenTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzDelegationExecutor) simpleQuery(funcName string) (*big.Int, error) {
	query1 := &vmi.ContractCallInput{
		RecipientAddr: pfe.delegationContractAddress,
		Function:      funcName,
		VMInput: vmi.VMInput{
			CallerAddr:  pfe.ownerAddress,
			Arguments:   [][]byte{},
			CallValue:   big.NewInt(0),
			GasPrice:    0,
			GasProvided: 1000000,
		},
	}

	queryOutput, err := pfe.vm.RunSmartContractCall(query1)
	if err != nil {
		return nil, err
	}

	result := big.NewInt(0).SetBytes(queryOutput.ReturnData[0])
	pfe.log("simpleQuery: %s -> %d", funcName, result)
	return result, nil
}

func (pfe *fuzzDelegationExecutor) init(args *fuzzDelegationExecutorInitArgs) error {
	pfe.initArgs = args
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
		pfe.world.AcctMap.PutAccount(&worldhook.Account{
			Exists:  true,
			Address: delegatorAddress(i),
			Nonce:   0,
			Balance: big.NewInt(0x100000000),
			Storage: make(map[string][]byte),
			Code:    []byte{},
		})
	}

	// deploy delegation
	err = pfe.executeStep(fmt.Sprintf(`
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
	return pfe.executeStep(fmt.Sprintf(`
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
}

func (pfe *fuzzDelegationExecutor) setStakePerNode(stakePerNode *big.Int) error {
	pfe.log("setStakePerNode: %d", stakePerNode)
	return pfe.executeStep(fmt.Sprintf(`
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
}

func (pfe *fuzzDelegationExecutor) setBlsKeys(numNodes int) error {
	pfe.log("setBlsKeys")
	var blsKeyArgs []string
	for i := 0; i < numNodes; i++ {
		blsKeyArg := "\"''" + string(blsKey(i)) + "\""
		blsKeyArgs = append(blsKeyArgs, blsKeyArg)
	}

	return pfe.executeStep(fmt.Sprintf(`
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
		strings.Join(blsKeyArgs, ","),
	))
}

func (pfe *fuzzDelegationExecutor) stake(delegIndex int, amount *big.Int) error {
	pfe.log("stake, delegator: %d, amount: %d ", delegIndex, amount)
	return pfe.executeStep(fmt.Sprintf(`
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

	return pfe.stake(delegIndex, unfilledStake)
}

func (pfe *fuzzDelegationExecutor) activate() error {
	pfe.log("activate")
	return pfe.executeStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "-activate-",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "stake",
			"arguments": [],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
	))
}

func delegatorAddress(delegIndex int) []byte {
	return []byte(fmt.Sprintf("delegator %5d               s1", delegIndex))
}

func blsKey(index int) []byte {
	return []byte(fmt.Sprintf(
		"bls key %5d ..................................................................................",
		index))
}

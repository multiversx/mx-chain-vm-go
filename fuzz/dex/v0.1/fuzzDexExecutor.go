package dex

import (
	"bytes"
	"errors"
	"fmt"
	am "github.com/ElrondNetwork/arwen-wasm-vm/arwenmandos"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/fileresolver"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	mjparse "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/parse"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/write"
	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"io/ioutil"
)

type fuzzDexExecutorInitArgs struct {
	wegldTokenId				string
	numUsers					int
	numTokens					int
	numEvents					int
	removeLiquidityProb			float32
	addLiquidityProb			float32
	swapProb					float32
	queryPairsProb				float32
	stakeProb					float32
	unstakeProb					float32
	unbondProb					float32
	increaseEpochProb			float32
	removeLiquidityMaxValue		int
	addLiquidityMaxValue 		int
	swapMaxValue 				int
	stakeMaxValue				int
	unstakeMaxValue				int
	unbondMaxValue				int
	blockEpochIncrease			int
	tokensCheckFrequency		int
}

type StakeInfo struct {
	user						string
	value						int64
	lpToken						string
}

type UnstakeInfo struct {
	user						string
	value						int64
	lpToken						string
}

type fuzzDexExecutor struct {
	arwenTestExecutor 			*am.ArwenTestExecutor
	world             			*worldhook.MockWorld
	vm                			vmi.VMExecutionHandler
	mandosParser      			mjparse.Parser
	txIndex           			int

	wegldTokenId				string
	ownerAddress				[]byte
	routerAddress				[]byte
	stakingAddress				[]byte
	numUsers					int
	numTokens					int
	numEvents					int
	removeLiquidityProb			float32
	addLiquidityProb			float32
	swapProb					float32
	queryPairsProb				float32
	stakeProb					float32
	unstakeProb					float32
	unbondProb					float32
	increaseEpochProb			float32
	removeLiquidityMaxValue		int
	addLiquidityMaxValue 		int
	swapMaxValue 				int
	stakeMaxValue				int
	unstakeMaxValue				int
	unbondMaxValue				int
	blockEpochIncrease			int
	tokensCheckFrequency		int
	currentStakeTokenNonce		int
	stakers						map[int]StakeInfo
	currentUnstakeTokenNonce	int
	unstakers					map[int]UnstakeInfo
	generatedScenario           *mj.Scenario
}

type eventsStatistics struct {
	swapFixedInputHits			int
	swapFixedInputMisses		int

	swapFixedOutputHits			int
	swapFixedOutputMisses		int

	addLiquidityHits			int
	addLiquidityMisses			int
	addLiquidityPriceChecks		int

	removeLiquidityHits			int
	removeLiquidityMisses		int
	removeLiquidityPriceChecks	int

	queryPairsHits				int
	queryPairsMisses			int

	stakeHits					int
	stakeMisses					int

	unstakeHits					int
	unstakeMisses				int
	unstakeWithRewards			int

	unbondHits					int
	unbondMisses				int
}

func newFuzzDexExecutor(fileResolver fr.FileResolver) (*fuzzDexExecutor, error) {
	arwenTestExecutor, err := am.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}

	parser := mjparse.NewParser(fileResolver)

	return &fuzzDexExecutor{
		arwenTestExecutor: arwenTestExecutor,
		world:             arwenTestExecutor.World,
		vm:                arwenTestExecutor.GetVM(),
		mandosParser:      parser,
		txIndex:           0,
		generatedScenario: &mj.Scenario{
			Name: "fuzz generated",
		},
	}, nil
}

func (pfe *fuzzDexExecutor) saveGeneratedScenario() {
	serialized := mjwrite.ScenarioToJSONString(pfe.generatedScenario)

	err := ioutil.WriteFile("fuzz_gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (pfe *fuzzDexExecutor) executeStep(stepSnippet string) error {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}

	pfe.addStep(step)
	return pfe.arwenTestExecutor.ExecuteStep(step)
}

func (pfe *fuzzDexExecutor) addStep(step mj.Step) {
	pfe.generatedScenario.Steps = append(pfe.generatedScenario.Steps, step)
}


func (pfe *fuzzDexExecutor) executeTxStep(stepSnippet string) (*vmi.VMOutput, error) {
	step, err := pfe.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}

	txStep, isTx := step.(*mj.TxStep)
	if !isTx {
		return nil, errors.New("tx step expected")
	}

	pfe.addStep(step)

	return pfe.arwenTestExecutor.ExecuteTxStep(txStep)
}

func (pfe *fuzzDexExecutor) log(info string, args ...interface{}) {
	fmt.Printf(info+"\n", args...)
}

func (pfe *fuzzDexExecutor) userAddress(userIndex int) []byte {
	return []byte(fmt.Sprintf("user%06d____________________s1", userIndex))
}

func (pfe *fuzzDexExecutor) tokenTicker(index int) string {
	return fmt.Sprintf("TOKEN-%06d", index)
}

func (pfe *fuzzDexExecutor) lpTokenTicker(index int) string {
	return fmt.Sprintf("LPTOK-%06d", index)
}


func (pfe *fuzzDexExecutor) fullOfEsdtWalletString() string {
	esdtString := ""
	for i := 1; i <= pfe.numTokens; i++ {
		esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.tokenTicker(i))
	}
	esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.wegldTokenId)
	for i := 1; i <= (pfe.numTokens * (pfe.numTokens + 1) / 2); i++ {
		esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.lpTokenTicker(i))
	}
	esdtString = esdtString[:len(esdtString)-1]
	return esdtString
}

func (pfe *fuzzDexExecutor) createPairs() error {
	for i := 1; i < pfe.numTokens; i++ {
		for j := i + 1; j <= pfe.numTokens; j++ {
			err := pfe.createPair(pfe.tokenTicker(i), pfe.tokenTicker(j))
			if err != nil {
				return err
			}
		}
	}
	for i := 1; i <= pfe.numTokens; i++ {
		err := pfe.createPair("WEGLD-abcdef", pfe.tokenTicker(i))
		if err != nil {
			return err
		}
	}

	return nil
}

func (pfe *fuzzDexExecutor) createPair(tokenA string, tokenB string) error {
	// deploy pair sc
	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "deploy-pair-contract",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "createPair",
			"arguments": [
				"str:%s",
				"str:%s"
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [ "*" ],
			"status": "0",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.routerAddress),
		tokenA,
		tokenB,
	))
	if err != nil {
		return err
	}

	rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
		"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
	if err != nil {
		return err
	}

	pairHexStr := "0x"
	for i := 0; i < len(rawResponse[0]); i++ {
		toAppend := fmt.Sprintf("%02x", rawResponse[0][i])
		pairHexStr += toAppend
	}

	// issue lp token for pair
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "issue-lp-token",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "5,000,000,000,000,000,000",
			"function": "issueLpToken",
			"arguments": [
				"%s",
				"0x4c49515544495459504f4f4c544f4b454e",
				"0x4c50544f4b454e"
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.routerAddress),
		pairHexStr,
	))
	if err != nil {
		return err
	}

	rawOutput, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "get_lp_token_identifier",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "getLpTokenIdentifier",
			"arguments": [],
			"gasLimit": "10,000,000",
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
		string(pfe.routerAddress),
		pairHexStr,
	))

	rawResponse = rawOutput.ReturnData
	lpTokenHexStr := "0x"
	for i := 0; i < len(rawResponse[0]); i++ {
		toAppend := fmt.Sprintf("%02x", rawResponse[0][i])
		lpTokenHexStr += toAppend
	}

	// set local roles for pair + lp token
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "deploy-pair-contract",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "setLocalRoles",
			"arguments": [
				"%s"
			],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.routerAddress),
		pairHexStr,
	))
	if err != nil {
		return err
	}

	return nil
}


func (pfe *fuzzDexExecutor) querySingleResult(from []byte, to []byte, funcName string, args string) ([][]byte, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%s",
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
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		funcName,
		string(from),
		string(to),
		funcName,
		args,
	))
	if err != nil {
		return [][]byte{}, err
	}

	return output.ReturnData, nil
}


func (pfe *fuzzDexExecutor) querySingleResultStringAddr(from []byte, to string, funcName string, args string) ([][]byte, error) {
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%s",
		"tx": {
			"from": "''%s",
			"to": "%s",
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
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		funcName,
		string(from),
		to,
		funcName,
		args,
	))
	if err != nil {
		return [][]byte{}, err
	}

	return output.ReturnData, nil
}


func (pfe *fuzzDexExecutor) setFeeOn() error {
	for i := 1; i <= pfe.numTokens; i++ {
		tokenA := pfe.wegldTokenId
		tokenB := pfe.tokenTicker(i)

		rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
			"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
		if err != nil {
			return err
		}

		pairHexStr := "0x"
		for i := 0; i < len(rawResponse[0]); i++ {
			toAppend := fmt.Sprintf("%02x", rawResponse[0][i])
			pairHexStr += toAppend
		}

		// set staking info
		_, err = pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "set-fee-on",
			"tx": {
				"from": "''%s",
				"to": "''%s",
				"value": "0",
				"function": "setFeeOn",
				"arguments": [
					"%s",
					"''%s",
					"str:%s"
				],
				"gasLimit": "10,000,000",
				"gasPrice": "0"
			},
			"expect": {
				"out": ["*"],
				"status": "",
				"logs": [],
				"gas": "*",
				"refund": "*"
			}
		}`,
			string(pfe.ownerAddress),
			string(pfe.routerAddress),
			pairHexStr,
			string(pfe.stakingAddress),
			pfe.wegldTokenId,
		))
		if err != nil {
			return err
		}

		rawOutput, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "",
			"tx": {
				"from": "''%s",
				"to": "%s",
				"value": "0",
				"function": "getFeeState",
				"arguments": [],
				"gasLimit": "10,000,000",
				"gasPrice": "0"
			},
			"expect": {
				"out": ["1"],
				"status": "",
				"logs": [],
				"gas": "*",
				"refund": "*"
			}
		}`,
			string(pfe.ownerAddress),
			pairHexStr,
		))
		if err != nil {
			return err
		}

		Use(rawOutput)
	}

	return nil
}

func (pfe *fuzzDexExecutor) increaseBlockEpoch(epochDelta int) error {
	currentBlockEpoch := uint32(0)
	if pfe.world.CurrentBlockInfo != nil {
		currentBlockEpoch = pfe.world.CurrentBlockInfo.BlockEpoch
	}

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "%d - increase block epoch",
		"currentBlockInfo": {
			"blockEpoch": "%d"
		}
	}`,
		pfe.nextTxIndex(),
		currentBlockEpoch+uint32(epochDelta),
	))
	if err != nil {
		return err
	}

	pfe.log("block epoch: %d ---> %d", currentBlockEpoch, currentBlockEpoch+uint32(epochDelta))
	return nil
}

func (pfe *fuzzDexExecutor) nextTxIndex() int {
	pfe.txIndex++
	return pfe.txIndex
}

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func (pfe *fuzzDexExecutor) doHackishSteps() error {
	lpTokenIndex := 1
	for i := 1; i < pfe.numTokens; i++ {
		for j := i + 1; j <= pfe.numTokens; j++ {
			err := pfe.doHackishStep(pfe.tokenTicker(i), pfe.tokenTicker(j), lpTokenIndex)
			lpTokenIndex += 1
			if err != nil {
				return err
			}
		}
	}
	for i := 1; i <= pfe.numTokens; i++ {
		err := pfe.doHackishStep(pfe.wegldTokenId, pfe.tokenTicker(i), lpTokenIndex)
		lpTokenIndex += 1
		if err != nil {
			return err
		}
	}

	err := pfe.doHachishStepStaking()
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzDexExecutor) doHackishStep(tokenA string, tokenB string, index int) error {
	lpTokenName := pfe.lpTokenTicker(index)

	rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
		"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
	if err != nil {
		return err
	}

	pairHexStr := "0x"
	for i := 0; i < len(rawResponse[0]); i++ {
		toAppend := fmt.Sprintf("%02x", rawResponse[0][i])
		pairHexStr += toAppend
	}

	err = pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "test",
		"accounts": {
			"%s": {
				"nonce": "0",
				"balance": "0",
				"esdtRoles": {
					"str:%s": [
						"ESDTRoleLocalMint",
						"ESDTRoleLocalBurn"
					]
				},
				"storage": {
					"str:first_token_id": "str:%s",
					"str:second_token_id": "str:%s",
					"str:state": "1",
					"str:lpTokenIdentifier": "str:%s",
					"str:router_address": "''%s",
					"str:fee_state": "1",
					"str:fee_address": "''%s",
					"str:fee_token_identifier": "str:%s",
					"str:total_fee_precent": "300",
					"str:special_fee_precent": "100",
					"str:router_owner_address": "address:owner"
				},
				"code": "file:../../../test/dex/v0_1/output/elrond_dex_pair.wasm"
			}
		}
	}`,
		pairHexStr,
		lpTokenName,
		tokenA,
		tokenB,
		lpTokenName,
		string(pfe.routerAddress),
		string(pfe.stakingAddress),
		pfe.wegldTokenId,
	))
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzDexExecutor) doHachishStepStaking() error {
	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "test",
		"accounts": {
			"''%s": {
				"nonce": "0",
				"balance": "0",
				"esdtRoles": {
					"str:%s": [
						"ESDTRoleNFTCreate",
						"ESDTRoleNFTAddQuantity",
						"ESDTRoleNFTBurn"
					],
					"str:%s": [
						"ESDTRoleNFTCreate",
						"ESDTRoleNFTAddQuantity",
						"ESDTRoleNFTBurn"
					]
				},
				"storage": {
					"str:wegld_token_id": "str:%s",
					"str:stake_token_id": "str:%s",
					"str:unstake_token_id": "str:%s",
					"str:router_address": "''%s",
					"str:virtual_token_id": "str:%s",
					"str:state": "1"
				},
				"code": "file:../../../test/dex/v0_1/output/elrond_dex_staking.wasm"
			}
		}
	}`,
		string(pfe.stakingAddress),
		"STAKING-abcdef",
		"UNSTAK-abcdef",
		pfe.wegldTokenId,
		"STAKING-abcdef",
		"UNSTAK-abcdef",
		string(pfe.routerAddress),
		pfe.wegldTokenId,
	))
	if err != nil {
		return err
	}

	return nil
}

// This function allows equality with a += 1
func equalMatrix(left [][]byte, right [][]byte) bool {
	if len(left) != len(right) {
		return false
	}

	for i := 0; i < len(left); i++ {
		if !bytes.Equal(left[i], right[i]) {
			if i == len(left) - 1 {
				leftIncreased := make([]byte, len(left[i]))
				copy(leftIncreased, left[i])
				if len(leftIncreased) > 0 {
					leftIncreased[len(leftIncreased)-1] += 1
				}

				rightIncreased := make([]byte, len(right[i]))
				copy(rightIncreased, right[i])
				if len(rightIncreased) > 0 {
					rightIncreased[len(rightIncreased)-1] += 1
				}

				return bytes.Equal(leftIncreased, right[i]) || bytes.Equal(left[i], rightIncreased)
			}
		}
	}

	return true
}

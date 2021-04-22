package dex

import (
	"bytes"
	"encoding/hex"
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
	mexTokenId					string
	numUsers					int
	numTokens					int
	numEvents					int
	removeLiquidityProb			float32
	addLiquidityProb			float32
	swapProb					float32
	queryPairsProb				float32
	enterFarmProb				float32
	exitFarmProb				float32
	unbondProb					float32
	increaseEpochProb			float32
	removeLiquidityMaxValue		int
	addLiquidityMaxValue 		int
	swapMaxValue 				int
	enterFarmMaxValue			int
	exitFarmMaxValue			int
	blockEpochIncrease			int
	tokensCheckFrequency		int
}

type PairMetadata struct {
	tokenA						string
	tokenB						string
	addr 						string
}

type FarmerInfo struct {
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
	mexTokenId					string
	ownerAddress				[]byte
	routerAddress				[]byte
	wegldFarmingAddress 		[]byte
	mexFarmingAddress			[]byte
	numUsers					int
	numTokens					int
	numEvents					int
	removeLiquidityProb			float32
	addLiquidityProb			float32
	swapProb					float32
	queryPairsProb				float32
	enterFarmProb				float32
	exitFarmProb				float32
	unbondProb					float32
	increaseEpochProb			float32
	removeLiquidityMaxValue		int
	addLiquidityMaxValue 		int
	swapMaxValue 				int
	enterFarmMaxValue			int
	exitFarmMaxValue			int
	unbondMaxValue				int
	blockEpochIncrease			int
	tokensCheckFrequency		int
	currentFarmTokenNonce		int
	farmers						map[int]FarmerInfo
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

	enterFarmHits				int
	enterFarmMisses				int

	exitFarmHits				int
	exitFarmMisses				int
	exitFarmWithRewards			int
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
	esdtString += fmt.Sprintf(`
						"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, pfe.mexTokenId)
	for i := 1; i <= ((pfe.numTokens + 1) * (pfe.numTokens + 2) / 2); i++ {
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
		err := pfe.createPair(pfe.wegldTokenId, pfe.tokenTicker(i))
		if err != nil {
			return err
		}
	}
	for i := 1; i <= pfe.numTokens; i++ {
		err := pfe.createPair(pfe.mexTokenId, pfe.tokenTicker(i))
		if err != nil {
			return err
		}
	}
	err := pfe.createPair(pfe.mexTokenId, pfe.wegldTokenId)
	if err != nil {
		return err
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

	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return err
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

	err, _, _ = pfe.getLpTokenIdentifier(pairHexStr)
	if err != nil {
		return err
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
	var stakeTokens [2]string
	stakeTokens[0] = pfe.wegldTokenId
	stakeTokens[1] = pfe.mexTokenId

	pairs := make([]PairMetadata, 0)

	for j := 0; j < len(stakeTokens); j++ {
		for i := 1; i <= pfe.numTokens; i++ {
			tokenA := stakeTokens[j]
			tokenB := pfe.tokenTicker(i)

			pairHexStr, err := pfe.setFeeOnPair(tokenA, tokenB)
			if err != nil {
				return err
			}
			elem := PairMetadata{
				tokenA: tokenA,
				tokenB: tokenB,
				addr:   pairHexStr,
			}
			pairs = append(pairs, elem)
		}
	}

	pairHex, err := pfe.setFeeOnPair(pfe.wegldTokenId, pfe.mexTokenId)
	if err != nil {
		return err
	}
	elem := PairMetadata{
		tokenA: pfe.wegldTokenId,
		tokenB: pfe.mexTokenId,
		addr:   pairHex,
	}
	pairs = append(pairs, elem)

	for i := 0; i < len(pairs); i++ {
		for j := 0; j < len(pairs); j++ {

			_, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "whitelist",
			"tx": {
				"from": "''%s",
				"to": "%s",
				"value": "0",
				"function": "whitelist",
				"arguments": [
					"%s"
				],
				"gasLimit": "10,000,000",
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
				pairs[i].addr,
				pairs[j].addr,
			))
			if err != nil {
				return err
			}


			_, err = pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "whitelist",
			"tx": {
				"from": "''%s",
				"to": "%s",
				"value": "0",
				"function": "addTrustedSwapPair",
				"arguments": [
					"%s",
					"str:%s",
					"str:%s"
				],
				"gasLimit": "10,000,000",
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
				pairs[i].addr,
				pairs[j].addr,
				pairs[j].tokenA,
				pairs[j].tokenB,
			))
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (pfe *fuzzDexExecutor) setFeeOnPair(tokenA string, tokenB string) (string, error) {
	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return "", err
	}

	err, _, lpTokenHexStr := pfe.getLpTokenIdentifier(pairHexStr)
	if err != nil {
		return "", err
	}

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
				"out": [],
				"status": "",
				"logs": [],
				"gas": "*",
				"refund": "*"
			}
		}`,
		string(pfe.ownerAddress),
		string(pfe.routerAddress),
		pairHexStr,
		string(pfe.wegldFarmingAddress),
		pfe.wegldTokenId,
	))
	if err != nil {
		return "", err
	}

	if tokenA == pfe.wegldTokenId || tokenB == pfe.wegldTokenId {
		_, err = pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "add-accepted-pair-address-and-lp-token",
			"tx": {
				"from": "''%s",
				"to": "''%s",
				"value": "0",
				"function": "addAcceptedPairAddressAndLpToken",
				"arguments": [
					"%s",
					"%s"
				],
				"gasLimit": "10,000,000",
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
			string(pfe.wegldFarmingAddress),
			pairHexStr,
			lpTokenHexStr,
		))
		if err != nil {
			return "", err
		}
	}

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
				"out": [],
				"status": "",
				"logs": [],
				"gas": "*",
				"refund": "*"
			}
		}`,
		string(pfe.ownerAddress),
		string(pfe.routerAddress),
		pairHexStr,
		string(pfe.mexFarmingAddress),
		pfe.mexTokenId,
	))
	if err != nil {
		return "", err
	}

	if tokenA == pfe.mexTokenId || tokenB == pfe.mexTokenId {
		_, err = pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "add-accepted-pair-address-and-lp-token",
			"tx": {
				"from": "''%s",
				"to": "''%s",
				"value": "0",
				"function": "addAcceptedPairAddressAndLpToken",
				"arguments": [
					"%s",
					"%s"
				],
				"gasLimit": "10,000,000",
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
			string(pfe.mexFarmingAddress),
			pairHexStr,
			lpTokenHexStr,
		))
		if err != nil {
			return "", err
		}
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
		return "", err
	}

	Use(rawOutput)
	return pairHexStr, nil
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
	for i := 1; i <= pfe.numTokens; i++ {
		err := pfe.doHackishStep(pfe.mexTokenId, pfe.tokenTicker(i), lpTokenIndex)
		lpTokenIndex += 1
		if err != nil {
			return err
		}
	}
	err := pfe.doHackishStep(pfe.mexTokenId, pfe.wegldTokenId, lpTokenIndex)
	lpTokenIndex += 1
	if err != nil {
		return err
	}

	err = pfe.doHachishStepStaking()
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzDexExecutor) doHackishStep(tokenA string, tokenB string, index int) error {
	lpTokenName := pfe.lpTokenTicker(index)

	err, _, pairHexStr := pfe.getPair(tokenA, tokenB)
	if err != nil {
		return err
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
					"str:total_fee_precent": "300",
					"str:special_fee_precent": "100",
					"str:router_owner_address": "''%s"
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
		string(pfe.ownerAddress),
	))
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzDexExecutor) doHachishStepStaking() error {
	esdt_role_string := ""
	for i := 1; i <= ((pfe.numTokens + 1) * (pfe.numTokens + 2) / 2); i++ {
		esdt_role_string += fmt.Sprintf(`
					"str:%s": [
						"ESDTRoleLocalBurn"
					],`, pfe.lpTokenTicker(i))
	}
	esdt_role_string += fmt.Sprintf(`
					"str:%s": [
						"ESDTRoleLocalBurn"
					],`, pfe.wegldTokenId)
	esdt_role_string += fmt.Sprintf(`
					"str:%s": [
						"ESDTRoleLocalBurn"
					]`, pfe.mexTokenId)

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
					],%s
				},
				"storage": {
					"str:farming_pool_token_id": "str:%s",
					"str:farm_token_id": "str:%s",
					"str:router_address": "''%s",
					"str:state": "1",
					"str:owner": "''%s",
					"str:farm_with_lp_tokens": "1"
				},
				"code": "file:../../../test/dex/v0_1/output/elrond_dex_farm.wasm"
			}
		}
	}`,
		string(pfe.wegldFarmingAddress),
		"FARM-abcdef",
		esdt_role_string,
		pfe.wegldTokenId,
		"FARM-abcdef",
		string(pfe.routerAddress),
		string(pfe.ownerAddress),
	))
	if err != nil {
		return err
	}

	err = pfe.executeStep(fmt.Sprintf(`
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
					],%s
				},
				"storage": {
					"str:farming_pool_token_id": "str:%s",
					"str:farm_token_id": "str:%s",
					"str:router_address": "''%s",
					"str:state": "1",
					"str:owner": "''%s",
					"str:farm_with_lp_tokens": "1"
				},
				"code": "file:../../../test/dex/v0_1/output/elrond_dex_farm.wasm"
			}
		}
	}`,
		string(pfe.mexFarmingAddress),
		"FARM-abcdef",
		esdt_role_string,
		pfe.mexTokenId,
		"FARM-abcdef",
		string(pfe.routerAddress),
		string(pfe.ownerAddress),
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

func (pfe *fuzzDexExecutor) getPair(tokenA string, tokenB string) (error, string, string) {
	rawResponse, err := pfe.querySingleResult(pfe.ownerAddress, pfe.routerAddress,
		"getPair", fmt.Sprintf("\"str:%s\", \"str:%s\"", tokenA, tokenB))
	if err != nil {
		return err, "", ""
	}

	pairHexStr := "0x"
	for i := 0; i < len(rawResponse[0]); i++ {
		toAppend := fmt.Sprintf("%02x", rawResponse[0][i])
		pairHexStr += toAppend
	}

	if (pairHexStr == "0x0000000000000000000000000000000000000000000000000000000000000000") && (tokenA != tokenB) {
		return errors.New("NULL pair for different tokens"), "", ""
	}

	return nil, string(rawResponse[0]), pairHexStr
}

func (pfe *fuzzDexExecutor) getLpTokenIdentifier(pairHexStr string) (error, string, string) {
	rawLpToken, errLpToken := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getLpTokenIdentifier", "")
	if errLpToken != nil {
		return errLpToken, "", ""
	}
	lpTokenHex := ""
	for i := 0; i < len(rawLpToken[0]); i++ {
		toAppend := fmt.Sprintf("%02x", rawLpToken[0][i])
		lpTokenHex += toAppend
	}
	lpToken, err := hex.DecodeString(lpTokenHex)
	if err != nil {
		return err, "", ""
	}
	lpTokenHex = "0x" + lpTokenHex
	return nil, string(lpToken), lpTokenHex
}
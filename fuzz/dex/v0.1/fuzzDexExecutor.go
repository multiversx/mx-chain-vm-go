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
	generatedScenario           *mj.Scenario
}

type eventsStatistics struct {
	swapFixedInputHits			int
	swapFixedInputMisses		int

	swapFixedOutputHits			int
	swapFixedOutputMisses		int

	addLiquidityHits			int
	addLiquidityMisses			int

	removeLiquidityHits			int
	removeLiquidityMisses		int

	queryPairsHits				int
	queryPairsMisses			int
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
			"function": "get_lp_token_identifier",
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
					"%s"
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

func (pfe *fuzzDexExecutor) increaseBlockNonce(nonceDelta int) error {
	currentBlockNonce := uint64(0)
	if pfe.world.CurrentBlockInfo != nil {
		currentBlockNonce = pfe.world.CurrentBlockInfo.BlockNonce
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
		currentBlockNonce+uint64(nonceDelta),
	))
	if err != nil {
		return err
	}

	pfe.log("block nonce: %d ---> %d", currentBlockNonce, currentBlockNonce+uint64(nonceDelta))
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

func (pfe *fuzzDexExecutor) swapFixedInput(user string, tokenA string, amountA int, tokenB string,
	amountB int, statistics *eventsStatistics) error {

	if tokenA == tokenB {
		return nil
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

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "swap-fixed-input",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "swapTokensFixedInput",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
				"str:%s",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pairHexStr,
		tokenA,
		amountA,
		tokenB,
		amountB,
	))
	if output != nil {
		if output.ReturnMessage != "" {
			statistics.swapFixedInputMisses += 1
			pfe.log("swapFixedInput %s -> %s", tokenA, tokenB)
			pfe.log("could not swap because %s", output.ReturnMessage)

			if output.ReturnMessage == "K invariant failed" {
				return errors.New(output.ReturnMessage)
			}
		} else {
			statistics.swapFixedInputHits += 1
		}
	} else {
		return errors.New("NULL output")
	}

	return nil
}

func (pfe *fuzzDexExecutor) swapFixedOutput(user string, tokenA string, amountA int, tokenB string,
	amountB int, statistics *eventsStatistics) error {

	if tokenA == tokenB {
		return nil
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

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "swap-fixed-input",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "swapTokensFixedOutput",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
				"str:%s",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pairHexStr,
		tokenA,
		amountA,
		tokenB,
		amountB,
	))
	if output != nil {
		if output.ReturnMessage != "" {
			statistics.swapFixedOutputMisses += 1
			pfe.log("swapFixedOutput %s -> %s", tokenA, tokenB)
			pfe.log("could not swap because %s", output.ReturnMessage)

			if output.ReturnMessage == "K invariant failed" {
				return errors.New(output.ReturnMessage)
			}
		} else {
			statistics.swapFixedOutputHits += 1
		}
	} else {
		return errors.New("NULL output")
	}

	return nil
}

func (pfe *fuzzDexExecutor) addLiquidity(user string, tokenA string, tokenB string, amountA int,
	amountB int , amountAmin int, amountBmin int, statistics *eventsStatistics) error {

	if tokenA == tokenB {
		return nil
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

	rawEquivalent, errEquivalent := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "accept-esdt-payment",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "acceptEsdtPayment",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "0",
			"message": "",
			"gas": "*",
			"refund": "*"
		}
	}`,
		user,
		pairHexStr,
		tokenA,
		amountA,
	))
	if err != nil {
		return err
	}

	output, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "accept-esdt-payment",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "acceptEsdtPayment",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "0",
			"message": "",
			"gas": "*",
			"refund": "*"
		}
	}`,
		user,
		pairHexStr,
		tokenB,
		amountB,
	))
	if err != nil {
		return err
	}

	output, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "add-liquidity",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "addLiquidity",
			"arguments": [
				"%d",
				"%d",
				"%d",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pairHexStr,
		amountA,
		amountB,
		amountAmin,
		amountBmin,
	))
	if output == nil {
		return errors.New("NULL output")
	}

	if output.ReturnMessage != "" {
		statistics.addLiquidityMisses += 1
		pfe.log("add liquidity %s -> %s", tokenA, tokenB)
		pfe.log("could not add because %s", output.ReturnMessage)

		//In case we get these errors but values are !=0, its an error
		if (output.ReturnMessage == "PAIR: INSSUFICIENT TOKEN A FUNDS SENT" ||
			output.ReturnMessage == "PAIR: INSSUFICIENT TOKEN B FUNDS SENT" ||
			output.ReturnMessage == "PAIR: NO AVAILABLE TOKEN A FUNDS" ||
			output.ReturnMessage == "PAIR: NO AVAILABLE TOKEN B FUNDS") &&
				(amountA > 0 && amountB > 0) {
			return errors.New(output.ReturnMessage)
		}

		if output.ReturnMessage == "Pair: FIRST TOKENS NEEDS TO BE GRATER THAN MINIMUM LIQUIDITY: 1000 * 1000e-18" &&
			amountA > 1000 && amountB > 1000 {
			return errors.New(output.ReturnMessage)
		}

		//No way we should receive this
		if output.ReturnMessage == "K invariant failed" {
			return errors.New(output.ReturnMessage)
		}

		// Other errors are fine
	} else {
		// Add liquidity is good
		statistics.addLiquidityHits += 1

		// Get New price
		rawResponse, err = pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))
		if err != nil {
			return err
		}

		// New and old prices should be the same
		if errEquivalent == nil && !equalMatrix(rawResponse, rawEquivalent) {
			return errors.New("PRICE CHANGED after add liquidity")
		}
	}

	output, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "reclaim-temporary-funds",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "reclaimTemporaryFunds",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "0",
			"message": "",
			"gas": "*",
			"refund": "*"
		}
	}`,
		user,
		pairHexStr,
	))
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzDexExecutor) removeLiquidity(user string, tokenA string, tokenB string, amount int, amountAmin int,
	amountBmin int, statistics *eventsStatistics) error {

	if tokenA == tokenB {
		return nil
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

	rawResponse, err = pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"get_lp_token_identifier", "")

	rawEquivalent, errEquivalent := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	lpToken := rawResponse[0]
	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "remove-liq",
		"tx": {
			"from": "''%s",
			"to": "%s",
			"value": "0",
			"function": "removeLiquidity",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d"
			},
			"arguments": [
				"%d",
				"%d"
			],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		pairHexStr,
		string(lpToken),
		amount,
		amountAmin,
		amountBmin,
	))
	if output == nil {
		return errors.New("NULL Output")
	}

	if output.ReturnMessage != "" {
		pfe.log("remove liquidity %s -> %s", tokenA, tokenB)
		pfe.log("could not remove because %s", output.ReturnMessage)
		statistics.removeLiquidityMisses += 1
	} else {
		statistics.removeLiquidityHits += 1

		rawOutput, err := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
			"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

		if errEquivalent != nil && err != nil && !equalMatrix(rawEquivalent, rawOutput) {
			return errors.New("PRICE CHANGED after success remove")
		}
	}

	return nil
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

func (pfe *fuzzDexExecutor) checkPairViews(user string, tokenA string, tokenB string, stats *eventsStatistics) error {
	if tokenA == tokenB {
		return nil
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

	outputAmountInA, errAmountInA := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	outputAmountOutA, errAmountOutA := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	outputEquivalentOutA, errEquivalentA := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenA, 1000))

	outputAmountInB, errAmountInB := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountIn", fmt.Sprintf("\"str:%s\", \"%d\"", tokenB, 1000))

	outputAmountOutB, errAmountOutB := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getAmountOut", fmt.Sprintf("\"str:%s\", \"%d\"", tokenB, 1000))

	outputEquivalentOutB, errEquivalentB := pfe.querySingleResultStringAddr(pfe.ownerAddress, pairHexStr,
		"getEquivalent", fmt.Sprintf("\"str:%s\", \"%d\"", tokenB, 1000))

	if errAmountInA != nil || errAmountInB != nil || errAmountOutA != nil || errAmountOutB != nil ||
		errEquivalentA != nil || errEquivalentB != nil {
		pfe.log("some query returned errors")
		stats.queryPairsMisses += 1
	} else {
		stats.queryPairsHits += 1
	}

	Use(outputAmountInA, outputAmountInB, outputAmountOutA, outputAmountOutB, outputEquivalentOutA, outputEquivalentOutB)

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
					"str:token_a_name": "str:%s",
					"str:token_b_name": "str:%s",
					"str:state": "1",
					"str:lpTokenIdentifier": "str:%s",
					"str:router_address": "''%s",
					"str:fee_state": "1",
					"str:fee_address": "''%s",
					"str:fee_token_identifier": "str:%s"
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
					]
				},
				"storage": {
					"str:wegld_token_identifier": "str:%s",
					"str:sft_staking_token_identifier": "str:%s",
					"str:router_address": "''%s",
					"str:state": "1"
				},
				"code": "file:../../../test/dex/v0_1/output/elrond_dex_staking.wasm"
			}
		}
	}`,
		string(pfe.stakingAddress),
		"STAKING-abcdef",
		pfe.wegldTokenId,
		"STAKING-abcdef",
		string(pfe.routerAddress),
	))
	if err != nil {
		return err
	}

	return nil
}

func equalMatrix(left [][]byte, right [][]byte) bool {
	if len(left) != len(right) {
		return false
	}

	for i := 0; i < len(left); i++ {
		if !bytes.Equal(left[i], right[i]) {
			return false
		}
	}

	return true
}

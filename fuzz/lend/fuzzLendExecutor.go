package lendFuzz

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwenmandos"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/fileresolver"
	mandosjson "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/json/model"
	parser "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/json/parse"
	jsonWrite "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/json/write"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

const (
	wegld  = "WEGLD-abcdef"
	lwegld = "LWEGLD-abcdef"
	bwegld = "BWEGLD-abcdef"

	busd  = "BUSD-abcdef"
	lbusd = "LBUSD-abcdef"
	bbusd = "BBUSD-abcdef"
)

type fuzzLendExecutorArgs struct {
	wegldTokenID  string
	lwegldTokenID string
	bwegldTokenID string
	busdTokenID   string
	lbusdTokenID  string
	bbusdTokenID  string

	wegldLPAddress string
	busdLPAddress  string

	lendPoolAddress string

	numUsers  int
	numTokens int
	numEvents int
}

type fuzzLendExecutor struct {
	vm                vmcommon.VMExecutionHandler
	arwenTestExecutor *arwenmandos.ArwenTestExecutor
	world             *worldmock.MockWorld
	mandosParser      parser.Parser

	wegldTokenID  string
	lwegldTokenID string
	bwegldTokenID string
	busdTokenID   string
	lbusdTokenID  string
	bbusdTokenID  string

	txIndex int

	ownerAddress string

	wegldLPAddress  string
	busdLPAddress   string
	lendPoolAddress string

	numUsers  int
	numTokens int
	numEvents int

	generatedScenario *mandosjson.Scenario
}

type statistics struct {
	depositHits   int
	depositMisses int

	borrowHits   int
	borrowMisses int

	withdrawHits   int
	withdrawMisses int

	repayHits   int
	repayMisses int
}

func newFuzzLendExecutor(fileResolver fr.FileResolver) (*fuzzLendExecutor, error) {
	arwenTestExecutor, err := arwenmandos.NewArwenTestExecutor()
	if err != nil {
		return nil, err
	}

	mandosParser := parser.NewParser(fileResolver)

	return &fuzzLendExecutor{
		arwenTestExecutor: arwenTestExecutor,
		world:             arwenTestExecutor.World,
		vm:                arwenTestExecutor.GetVM(),
		mandosParser:      mandosParser,
		txIndex:           0,
		generatedScenario: &mandosjson.Scenario{
			Name: "lend fuzz generated",
		},
	}, nil
}

func (e *fuzzLendExecutor) executeStep(stepSnippet string) error {
	step, err := e.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return err
	}

	e.addStep(step)
	return e.arwenTestExecutor.ExecuteStep(step)
}

func (e *fuzzLendExecutor) executeTxStep(stepSnippet string) (*vmcommon.VMOutput, error) {
	step, err := e.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		return nil, err
	}
	e.addStep(step)

	txStep, ok := step.(*mandosjson.TxStep)
	if !ok {
		return nil, errors.New("tx step expected")
	}

	return e.arwenTestExecutor.ExecuteTxStep(txStep)
}

func (e *fuzzLendExecutor) increaseBlockNonce(epochDelta int) error {
	var currBlockNonce uint64
	if e.world.CurrentBlockInfo != nil {
		currBlockNonce = e.world.CurrentBlockInfo.BlockNonce
	}

	err := e.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"comment": "%d - increase block nonce",
		"currentBlockInfo": {
			"blockNonce": "%d"
		}
	}`,
		e.nextTxIndex(),
		currBlockNonce+uint64(epochDelta),
	))
	if err != nil {
		return err
	}

	return nil
}

func (e *fuzzLendExecutor) saveGeneratedScenario() {
	serialized := jsonWrite.ScenarioToJSONString(e.generatedScenario)

	err := ioutil.WriteFile("lend_fuzz.gen.scen.json", []byte(serialized), 0644)
	if err != nil {
		log.Panicln("could not save scenario", "err", err.Error())
	}
}

func (e *fuzzLendExecutor) nextTxIndex() int {
	e.txIndex++
	return e.txIndex
}

func (e *fuzzLendExecutor) addStep(step mandosjson.Step) {
	e.generatedScenario.Steps = append(e.generatedScenario.Steps, step)
}

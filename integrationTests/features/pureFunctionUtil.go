package featuresintegrationtest

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	arwen "github.com/ElrondNetwork/wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/wasm-vm/arwen/host"
	"github.com/ElrondNetwork/wasm-vm/arwen/mock"
	"github.com/ElrondNetwork/wasm-vm/config"
	er "github.com/ElrondNetwork/wasm-vm/mandos-go/expression/reconstructor"
	worldhook "github.com/ElrondNetwork/wasm-vm/mock/world"
	"github.com/stretchr/testify/require"
)

type pureFunctionIO struct {
	functionName    string
	arguments       [][]byte
	expectedStatus  vmi.ReturnCode
	expectedMessage string
	expectedResults [][]byte
}

var testVMType = []byte{0, 0}

type resultInterpreter func([]byte) *big.Int
type logProgress func(testCaseIndex, testCaseCount int)

type pureFunctionExecutor struct {
	world           *worldhook.MockWorld
	vm              vmi.VMExecutionHandler
	contractAddress []byte
	userAddress     []byte
}

func newPureFunctionExecutor() (*pureFunctionExecutor, error) {
	world := worldhook.NewMockWorld()

	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMapForTests()
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldhook.WorldMarshalizer)
	vm, err := arwenHost.NewArwenVM(world, &arwen.VMHostParameters{
		VMType:                   testVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              gasSchedule,
		BuiltInFuncContainer:     builtInFunctions.NewBuiltInFunctionContainer(),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
		ESDTTransferParser:       esdtTransferParser,
		EpochNotifier:            &mock.EpochNotifierStub{},
		EnableEpochsHandler:      &worldmock.EnableEpochsHandlerStub{},
		WasmerSIGSEGVPassthrough: false,
	})
	if err != nil {
		return nil, err
	}
	return &pureFunctionExecutor{
		world: world,
		vm:    vm,
	}, nil
}

func (pfe *pureFunctionExecutor) initAccounts(contractPath string) {
	pfe.contractAddress = []byte("contract_addr_________________s1")
	pfe.userAddress = []byte("user_addr_____________________s1")

	scCode, err := ioutil.ReadFile(contractPath)
	if err != nil {
		panic(err)
	}

	pfe.world.AcctMap.PutAccount(&worldhook.Account{
		Address: pfe.contractAddress,
		Nonce:   0,
		Balance: big.NewInt(0),
		Storage: make(map[string][]byte),
		Code:    scCode,
	})

	pfe.world.AcctMap.PutAccount(&worldhook.Account{
		Address: pfe.userAddress,
		Nonce:   0,
		Balance: big.NewInt(0x100000000),
		Storage: make(map[string][]byte),
		Code:    []byte{},
	})
}

func (pfe *pureFunctionExecutor) scCall(testCase *pureFunctionIO) (*vmi.VMOutput, error) {
	input := &vmi.ContractCallInput{
		RecipientAddr: pfe.contractAddress,
		Function:      testCase.functionName,
		VMInput: vmi.VMInput{
			CallerAddr:  pfe.userAddress,
			Arguments:   testCase.arguments,
			CallValue:   big.NewInt(0),
			GasPrice:    1,
			GasProvided: 100000000,
		},
	}

	return pfe.vm.RunSmartContractCall(input)
}

func (pfe *pureFunctionExecutor) checkTxResults(
	testCase *pureFunctionIO,
	output *vmi.VMOutput,
	resultInterpreter resultInterpreter) error {

	if output.ReturnCode != testCase.expectedStatus {
		return fmt.Errorf("result code mismatch. Want: %d. Have: %d (%s). Message: %s",
			int(testCase.expectedStatus), int(output.ReturnCode), output.ReturnCode.String(), output.ReturnMessage)
	}

	if output.ReturnMessage != testCase.expectedMessage {
		return fmt.Errorf("result message mismatch. Want: %s. Have: %s",
			testCase.expectedMessage, output.ReturnMessage)
	}

	// check result
	if len(output.ReturnData) != len(testCase.expectedResults) {
		rec := er.ExprReconstructor{}
		return fmt.Errorf("result length mismatch. Want: %s. Have: %s",
			rec.ReconstructList(testCase.expectedResults, er.NoHint),
			rec.ReconstructList(output.ReturnData, er.NoHint))
	}
	for i, expected := range testCase.expectedResults {
		wantNum := resultInterpreter(expected)
		haveNum := resultInterpreter(output.ReturnData[i])
		if wantNum.Cmp(haveNum) != 0 {
			var argStr []string
			for _, arg := range testCase.arguments {
				argNum := resultInterpreter(arg)
				argStr = append(argStr, fmt.Sprintf("%d", argNum))
			}
			return fmt.Errorf("result mismatch. Want: %d. Have: %d. Call: %s(%s)",
				wantNum, haveNum, testCase.functionName, strings.Join(argStr, ", "))
		}
	}

	return nil
}

func (pfe *pureFunctionExecutor) executePureFunctionTests(t *testing.T,
	testCases []*pureFunctionIO,
	resultInterpreter resultInterpreter,
	logProgress logProgress) {

	// RUN!
	for testCaseIndex, testCase := range testCases {
		if logProgress != nil {
			logProgress(testCaseIndex, len(testCases))
		}

		output, err := pfe.scCall(testCase)
		require.Nil(t, err)

		err = pfe.checkTxResults(testCase, output, resultInterpreter)
		require.Nil(t, err)

		vmHost := pfe.vm.(arwen.VMHost)
		vmHost.Reset()
	}
}

package featuresintegrationtest

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	er "github.com/multiversx/mx-chain-scenario-go/scenario/expression/reconstructor"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

var defaultHasher = blake2b.NewBlake2b()

type pureFunctionIO struct {
	functionName    string
	arguments       [][]byte
	expectedStatus  vmcommon.ReturnCode
	expectedMessage string
	expectedResults [][]byte
}

var testVMType = []byte{0, 0}

type resultInterpreter func([]byte) *big.Int
type logProgress func(testCaseIndex, testCaseCount int)

type pureFunctionExecutor struct {
	world           *worldmock.MockWorld
	vm              vmcommon.VMExecutionHandler
	contractAddress []byte
	userAddress     []byte
}

func newPureFunctionExecutor() (*pureFunctionExecutor, error) {
	world := worldmock.NewMockWorld()

	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMapForTests()
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	vm, err := hostCore.NewVMHost(
		world,
		&vmhost.VMHostParameters{
			VMType:                   testVMType,
			OverrideVMExecutor:       nil,
			BlockGasLimit:            blockGasLimit,
			GasSchedule:              gasSchedule,
			BuiltInFuncContainer:     builtInFunctions.NewBuiltInFunctionContainer(),
			ProtectedKeyPrefix:       []byte("E" + "L" + "R" + "O" + "N" + "D"),
			ESDTTransferParser:       esdtTransferParser,
			EpochNotifier:            &mock.EpochNotifierStub{},
			EnableEpochsHandler:      worldmock.EnableEpochsHandlerStubNoFlags(),
			WasmerSIGSEGVPassthrough: false,
			Hasher:                   defaultHasher,
		})
	if err != nil {
		return nil, err
	}
	return &pureFunctionExecutor{
		world: world,
		vm:    vm,
	}, nil
}

func (pfe *pureFunctionExecutor) initAccounts(contractCode []byte) {
	pfe.contractAddress = []byte("contract_addr_________________s1")
	pfe.userAddress = []byte("user_addr_____________________s1")

	pfe.world.AcctMap.PutAccount(&worldmock.Account{
		Address: pfe.contractAddress,
		Nonce:   0,
		Balance: big.NewInt(0),
		Storage: make(map[string][]byte),
		Code:    contractCode,
	})

	pfe.world.AcctMap.PutAccount(&worldmock.Account{
		Address: pfe.userAddress,
		Nonce:   0,
		Balance: big.NewInt(0x100000000),
		Storage: make(map[string][]byte),
		Code:    []byte{},
	})
}

func (pfe *pureFunctionExecutor) scCall(testCase *pureFunctionIO) (*vmcommon.VMOutput, error) {
	input := &vmcommon.ContractCallInput{
		RecipientAddr: pfe.contractAddress,
		Function:      testCase.functionName,
		VMInput: vmcommon.VMInput{
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
	output *vmcommon.VMOutput,
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

	defer func() {
		vmHost := pfe.vm.(vmhost.VMHost)
		vmHost.Reset()
	}()

	// RUN!
	for testCaseIndex, testCase := range testCases {
		if logProgress != nil {
			logProgress(testCaseIndex, len(testCases))
		}

		output, err := pfe.scCall(testCase)
		require.Nil(t, err)

		err = pfe.checkTxResults(testCase, output, resultInterpreter)
		require.Nil(t, err)
	}
}

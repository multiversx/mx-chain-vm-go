package hosttest

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
)

func TestGasUsed_OneAsyncCall_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCall()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCalls_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCalls()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsync_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsAsync()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_DifferentTypeOfCallsToSameFunction_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestDifferentTypeOfCallsToSameFunction()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestCallbackCallsAsync()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SimpleSyncAndAsync1_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSimpleSyncAndAsync1()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SimpleSyncAndAsync2_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSimpleSyncAndAsync2()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_GraphTest1_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTest1()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_GraphTest2_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTest2()
	runGraphCallTestTemplate(t, callGraph)
}

type usedGasPerContract struct {
	contractAddress []byte
	gasUsed         uint64
}

func runGraphCallTestTemplate(t *testing.T, callGraph *test.TestCallGraph) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = callGraph.StartNode.GasLimit
	testConfig.GasLockCost = test.DefaultCallGraphLockedGas

	// compute execution order (return data) assertions
	expectedReturnData := make([][]byte, 0)
	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	startNode := executionGraph.GetStartNode()

	// executionOrder := test.CreateRunOrderFromExecutionGraph(executionGraph)
	// for _, testCall := range executionOrder {
	// 	expectedReturnData = append(expectedReturnData, []byte(string(testCall.ContractAddress)+"_"+testCall.FunctionName+test.TestReturnDataSuffix))
	// }

	// compute gas assertions
	gasGraph := executionGraph.CreateGasGraphFromExecutionGraph()
	gasGraph.ComputeRemainingGasBeforeCallbacks()
	gasGraph.ComputeGasStepByStep(func(graph *test.TestCallGraph, step int) {})

	totalGasUsed := uint64(0)
	expectedGasUsagePerContract := make(map[string]*usedGasPerContract)
	gasGraph.DfsGraph(func(path []*test.TestCallNode, parent *test.TestCallNode, node *test.TestCallNode, incomingEdge *test.TestCallEdge) *test.TestCallNode {
		if !node.IsLeaf() {
			return node
		}
		gasPerContract := expectedGasUsagePerContract[string(parent.Call.ContractAddress)]
		if gasPerContract == nil {
			gasPerContract = &usedGasPerContract{
				contractAddress: parent.Call.ContractAddress,
				gasUsed:         0,
			}
			expectedGasUsagePerContract[string(parent.Call.ContractAddress)] = gasPerContract
		}
		gasPerContract.gasUsed += node.GasUsed
		totalGasUsed += node.GasUsed

		expectedNodeRetData := txDataBuilder.NewBuilder()
		expectedNodeRetData.Func(parent.Call.FunctionName)
		expectedNodeRetData.Str(string(parent.Call.ContractAddress) + "_" + parent.Call.FunctionName + test.TestReturnDataSuffix)
		// fmt.Println("exp\t" + string(parent.Call.ContractAddress) + "_" + parent.Call.FunctionName + test.TestReturnDataSuffix)
		expectedNodeRetData.Int64(int64(parent.GasLimit))
		// fmt.Println(fmt.Sprintf("exp\t%d", parent.GasLimit))
		expectedNodeRetData.Int64(int64(parent.GasRemaining))
		// fmt.Println(fmt.Sprintf("exp\t%d", parent.GasRemaining))
		// fmt.Println()
		expectedReturnData = append(expectedReturnData, expectedNodeRetData.ToBytes())

		return node
	})

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContractsFromAsyncTestCallGraph(callGraph, testConfig)...,
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr([]byte(startNode.Call.ContractAddress)).
			WithGasProvided(testConfig.GasProvided).
			WithFunction(startNode.Call.FunctionName).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verifier := verify.
				Ok().
				ReturnDataForGraphTesting(expectedReturnData...).
				GasRemaining(callGraph.StartNode.GasLimit - totalGasUsed)
			for _, gasPerContract := range expectedGasUsagePerContract {
				verifier.GasUsed(gasPerContract.contractAddress, gasPerContract.gasUsed)
			}
		})
}

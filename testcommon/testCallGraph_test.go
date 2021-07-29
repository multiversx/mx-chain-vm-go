package testcommon

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallGraph_Dfs(t *testing.T) {
	callGraph := CreateGraphTest1()

	traversalOrder := make([]TestCall, 0)
	callGraph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		traversalOrder = append(traversalOrder, TestCall{
			ContractAddress: node.Call.ContractAddress,
			FunctionName:    node.Call.FunctionName,
		})
		return node
	})

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc2", "f6"),
		*buildTestCall("sc3", "f7"),
		*buildTestCall("sc1", "cb2"),
		*buildTestCall("sc4", "f5"),
		*buildTestCall("sc2", "cb3"),
		*buildTestCall("sc1", "cb4"),
		*buildTestCall("sc1", "cbg1"),
		*buildTestCall("sc1", "ctxcb"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, traversalOrder))
}

func TestExecutionGraph_Execution_OneAsyncCall(t *testing.T) {
	callGraph := CreateGraphTestOneAsyncCall()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_OneAsyncCallWithGroupCallback(t *testing.T) {
	callGraph := CreateGraphTestOneAsyncCallWithGroupCallback()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
		*buildTestCall("sc1", "cbg1"),
		*buildTestCall("sc1", "ctxcb"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_TwoAsyncCalls(t *testing.T) {
	callGraph := CreateGraphTestTwoAsyncCalls()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc1", "cb1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_AsyncCallsAsync(t *testing.T) {
	callGraph := CreateGraphTestAsyncCallsAsync()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc2", "cb2"),
		*buildTestCall("sc1", "cb1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_AsyncCallsAsync2(t *testing.T) {
	callGraph := CreateGraphTestAsyncCallsAsync2()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc2", "f4"),
		*buildTestCall("sc2", "cb3"),
		*buildTestCall("sc2", "cb2"),
		*buildTestCall("sc1", "cb1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_GroupCallbacks(t *testing.T) {
	callGraph := CreateGraphTestGroupCallbacks()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc2", "cb2"),
		*buildTestCall("sc1", "cb1"),
		*buildTestCall("sc2", "cbg2"),
		*buildTestCall("sc1", "cbg1"),
		*buildTestCall("sc1", "ctxcb"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_SimpleSyncAndAsync1(t *testing.T) {
	callGraph := CreateGraphTestSimpleSyncAndAsync1()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "cb1"),
		*buildTestCall("sc2", "cb1"),
		*buildTestCall("sc4", "f4"),
		*buildTestCall("sc1", "f1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_SimpleSyncAndAsync2(t *testing.T) {
	callGraph := CreateGraphTestSimpleSyncAndAsync2()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc4", "f4"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc1", "cb1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	// gasGraph := computeFinalGasGraph(executionGraph)
	// expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	// require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_GraphTest2(t *testing.T) {
	callGraph := CreateGraphTest2()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc4", "f4"),
		*buildTestCall("sc1", "cb1"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb2"),
		*buildTestCall("sc5", "f5"),
		*buildTestCall("sc4", "f4"),
		*buildTestCall("sc1", "cb1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func TestExecutionGraph_Execution_GraphTest1(t *testing.T) {
	callGraph := CreateGraphTest1()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "cb3"),
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc4", "f5"),
		*buildTestCall("sc1", "cb2"),
		*buildTestCall("sc2", "f6"),
		*buildTestCall("sc1", "cb4"),
		*buildTestCall("sc3", "f7"),
		*buildTestCall("sc1", "cb4"),
		*buildTestCall("sc1", "cbg1"),
		*buildTestCall("sc1", "ctxcb"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	//gasGraph := computeFinalGasGraph(executionGraph)
	//expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	//require.Equal(t, expectedRemainingGas, gasGraph.StartNode.GasRemaining)
}

func computeExpectedRemainingGas(gasGraph *TestCallGraph) uint64 {
	gasInLeafs := 0
	for _, node := range gasGraph.Nodes {
		if node.IsLeaf() {
			gasInLeafs += int(node.GasUsed)
		}
	}
	expectedRemainingGas := gasGraph.StartNode.GasLimit - uint64(gasInLeafs)
	return expectedRemainingGas
}

func computeFinalGasGraph(executionGraph *TestCallGraph) *TestCallGraph {
	gasGraph := executionGraph.CreateGasGraphFromExecutionGraph()
	gasGraph.ComputeRemainingGasBeforeCallbacks()
	// gasGraph.ComputeRemainingGasAfterCallbacks()
	// gasGraph.ComputeGasAccumulation()
	// gasGraph.ComputeRemainingGasAfterGroupCallbacks()
	return gasGraph
}

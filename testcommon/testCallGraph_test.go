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
			CallID:          node.Call.CallID,
		})
		return node
	}, true)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc1", "cb2"),
		*buildTestCall("sc4", "f5"),
		*buildTestCall("sc2", "cb3"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, traversalOrder))
}

func TestExecutionGraph_OneAsyncCall(t *testing.T) {
	callGraph := CreateGraphTestOneAsyncCall()
	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
	}
	runAsserts(callGraph, t, expectedOrder)
}

func TestExecutionGraph_TwoAsyncCalls(t *testing.T) {
	callGraph := CreateGraphTestTwoAsyncCalls()
	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc1", "cb2"),
	}
	runAsserts(callGraph, t, expectedOrder)
}

func TestExecutionGraph_AsyncCallsAsync(t *testing.T) {
	callGraph := CreateGraphTestAsyncCallsAsync()
	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "cb2"),
		*buildTestCall("sc1", "cb1"),
	}
	runAsserts(callGraph, t, expectedOrder)
}

func TestExecutionGraph_SimpleSyncAndAsync1(t *testing.T) {
	callGraph := CreateGraphTestSyncAndAsync1()
	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
	}
	runAsserts(callGraph, t, expectedOrder)
}

func TestExecutionGraph_SimpleSyncAndAsync2(t *testing.T) {
	callGraph := CreateGraphTestSyncAndAsync2()
	expectedOrder := []TestCall{
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc5", "f5"),
		*buildTestCall("sc2", "cb1"),
		*buildTestCall("sc5", "f5"),
		*buildTestCall("sc2", "cb1"),
		*buildTestCall("sc4", "f4"),
		*buildTestCall("sc1", "f1"),
	}
	runAsserts(callGraph, t, expectedOrder)
}

func TestExecutionGraph_GraphTest2(t *testing.T) {
	callGraph := CreateGraphTest2()
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
	runAsserts(callGraph, t, expectedOrder)
}

func TestExecutionGraph_GraphTest1(t *testing.T) {
	callGraph := CreateGraphTest1()
	expectedOrder := []TestCall{
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "cb3"),
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc4", "f5"),
		*buildTestCall("sc1", "cb2"),
	}
	runAsserts(callGraph, t, expectedOrder)
}

func TestExecutionGraph_DifferentTypeOfCallsToSameFunction(t *testing.T) {
	callGraph := CreateGraphTestDifferentTypeOfCallsToSameFunction()
	expectedOrder := []TestCall{
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb2"),
	}
	runAsserts(callGraph, t, expectedOrder)
}

func runAsserts(callGraph *TestCallGraph, t *testing.T, expectedOrder []TestCall) {
	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunOrderFromExecutionGraph(executionGraph)

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))

	gasGraph := computeFinalGasGraph(executionGraph)
	expectedRemainingGas := computeExpectedRemainingGas(gasGraph)
	var actualGasRemaining uint64
	if gasGraph.StartNode.GasAccumulatedAfterCallback != 0 {
		actualGasRemaining = gasGraph.StartNode.GasAccumulatedAfterCallback
	} else {
		actualGasRemaining = gasGraph.StartNode.GasRemaining
	}
	require.Equal(t, expectedRemainingGas, actualGasRemaining)
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
	gasGraph := executionGraph.ComputeGasGraphFromExecutionGraph()
	gasGraph.ComputeRemainingGasBeforeCallbacks()
	gasGraph.ComputeGasStepByStep(func(graph *TestCallGraph, step int) {})
	return gasGraph
}

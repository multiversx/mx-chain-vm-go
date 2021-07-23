package testcommon

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallGraph_Dfs(t *testing.T) {
	callGraph := CreateGraphTest1()

	traversalOrder := make([]TestCall, 0)
	callGraph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		//fmt.Println(string(node.call.ContractAddress) + " " + node.call.FunctionName)
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

func TestExecutionGraph_Execution_GraphTest1(t *testing.T) {
	callGraph := CreateGraphTest1()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunExpectationOrder(executionGraph)

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
}

func TestExecutionGraph_Execution_GraphTestSimple1(t *testing.T) {
	callGraph := CreateGraphTestSimple1()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunExpectationOrder(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc1", "cb1"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc1", "cb1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))
}

func TestExecutionGraph_Execution_GraphTestSimple2(t *testing.T) {
	callGraph := CreateGraphTestSimple2()

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	executionOrder := CreateRunExpectationOrder(executionGraph)

	expectedOrder := []TestCall{
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f3"),
		*buildTestCall("sc2", "cb1"),
		*buildTestCall("sc2", "cb1"),
		*buildTestCall("sc4", "f4"),
		*buildTestCall("sc1", "f1"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))
}

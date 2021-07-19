package testcommon

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallGraph_Dfs(t *testing.T) {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddNode("sc1", "f1")

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddEdge(sc1f1, sc2f2)

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc2f3, "cb2", "gr")

	sc3f4 := callGraph.AddNode("sc3", "f4")
	callGraph.AddEdge(sc2f3, sc3f4)

	callGraph.AddAsyncEdge(sc2f2, sc3f4, "cb3", "gr")

	sc1cb1 := callGraph.AddNode("sc1", "cb2")
	sc4f5 := callGraph.AddNode("sc4", "f5")
	callGraph.AddEdge(sc1cb1, sc4f5)

	sc2cb3 := callGraph.AddNode("sc2", "cb3")
	callGraph.AddEdge(sc2cb3, sc3f4)

	traversalOrder := make([]TestCall, 0)
	callGraph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		//fmt.Println(string(node.asyncCall.ContractAddress) + " " + node.asyncCall.FunctionName)
		traversalOrder = append(traversalOrder, TestCall{
			ContractAddress: node.call.ContractAddress,
			FunctionName:    node.call.FunctionName,
		})
		return node
	})

	expectedOrder := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc1", "cb2"),
		*buildTestCall("sc4", "f5"),
		*buildTestCall("sc2", "cb3"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, traversalOrder))
}

func TestExecutionGraph_Creation(t *testing.T) {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddNode("sc1", "f1")

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddEdge(sc1f1, sc2f2)

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc2f3, "cb2", "gr")

	sc3f4 := callGraph.AddNode("sc3", "f4")
	callGraph.AddEdge(sc2f3, sc3f4)

	callGraph.AddAsyncEdge(sc2f2, sc3f4, "cb3", "gr")

	sc1cb1 := callGraph.AddNode("sc1", "cb2")
	sc4f5 := callGraph.AddNode("sc4", "f5")
	callGraph.AddEdge(sc1cb1, sc4f5)

	sc2cb3 := callGraph.AddNode("sc2", "cb3")
	callGraph.AddEdge(sc2cb3, sc3f4)

	executionOrder := make([]TestCall, 0)

	executionGraph := callGraph.createExecutionGraphFromCallGraph()

	startNode := executionGraph.FindNode(MakeTestSCAddress("sc1"), "f1")
	executionGraph.DfsFromNode(nil, startNode, make([]*TestCallNode, 0), func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		if !node.HasAdjacentNodes() {
			fmt.Println(node.GetCall().FunctionName)
			executionOrder = append(executionOrder, TestCall{
				ContractAddress: node.call.ContractAddress,
				FunctionName:    node.call.FunctionName,
			})
		}
		return node
	}, false)

	expectedOrder := []TestCall{
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "cb3"),
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc4", "f5"),
		*buildTestCall("sc1", "cb2"),
	}

	require.True(t, reflect.DeepEqual(expectedOrder, executionOrder))
}

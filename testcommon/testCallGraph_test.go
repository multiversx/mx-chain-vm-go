package testcommon

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsyncTestCallGraph(t *testing.T) {
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

	result := make([]TestCall, 0)
	callGraph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		//fmt.Println(string(node.asyncCall.ContractAddress) + " " + node.asyncCall.FunctionName)
		result = append(result, TestCall{
			ContractAddress: node.call.ContractAddress,
			FunctionName:    node.call.FunctionName,
		})
		return node
	})

	expected := []TestCall{
		*buildTestCall("sc1", "f1"),
		*buildTestCall("sc2", "f2"),
		*buildTestCall("sc3", "f4"),
		*buildTestCall("sc2", "f3"),
		*buildTestCall("sc1", "cb2"),
		*buildTestCall("sc4", "f5"),
		*buildTestCall("sc2", "cb3"),
	}

	require.True(t, reflect.DeepEqual(expected, result))
}

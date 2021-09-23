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

package testcommon

import (
	"fmt"
	"testing"
)

// func TestAsyncTestCallTree(t *testing.T) {
// 	callTree := CreateCallTree(BuildOldAsyncTestCall("sc1", "g", "f1", ""))
// 	root := callTree.root

// 	ch1 := root.AddChild(BuildOldAsyncTestCall("sc1_2", "g", "f2", "cb2"))
// 	ch2 := root.AddChild(BuildOldAsyncTestCall("sc1_3", "g", "f3", "cb3"))

// 	ch1.AddChild(BuildOldAsyncTestCall("sc2_1", "g", "f4", "cb4"))
// 	ch1.AddChild(BuildOldAsyncTestCall("sc2_2", "g", "f5", "cb5"))

// 	ch2.AddChild(BuildOldAsyncTestCall("sc3_1", "g", "f6", "cb6"))

// 	CreateMockContractsFromAsyncTestCallTree(callTree,
// 		&TestConfig{
// 			ParentBalance: 1000,
// 		})
// }

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

	callGraph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		fmt.Println(string(node.asyncCall.ContractAddress) + " " + node.asyncCall.FunctionName)
		return node
	})

	// testConfig := &TestConfig{
	// 	GasProvided:           2000,
	// 	GasProvidedToChild:    300,
	// 	GasProvidedToCallback: 50,
	// 	GasUsedByParent:       400,
	// 	GasUsedByChild:        200,
	// 	GasUsedByCallback:     100,
	// 	GasLockCost:           150,

	// 	TransferFromParentToChild: 7,

	// 	ParentBalance:        1000,
	// 	ChildBalance:         1000,
	// 	TransferToThirdParty: 3,
	// 	TransferToVault:      4,
	// 	ESDTTokensToTransfer: 0,
	// }

	// testConfig.GasProvided = 10_000
	// testConfig.GasLockCost = 10

	// CreateMockContractsFromAsyncTestCallGraph(callGraph, testConfig)
}

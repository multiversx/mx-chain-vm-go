package testcommon

import (
	"testing"
)

func TestAsyncTestCallTree(t *testing.T) {
	callTree := CreateCallTree(BuildAsyncTestCall("sc1", "g", "f1", ""))
	root := callTree.root

	ch1 := root.AddChild(BuildAsyncTestCall("sc1_2", "g", "f2", "cb2"))
	ch2 := root.AddChild(BuildAsyncTestCall("sc1_3", "g", "f3", "cb3"))

	ch1.AddChild(BuildAsyncTestCall("sc2_1", "g", "f4", "cb4"))
	ch1.AddChild(BuildAsyncTestCall("sc2_2", "g", "f5", "cb5"))

	ch2.AddChild(BuildAsyncTestCall("sc3_1", "g", "f6", "cb6"))

	CreateMockContractsFromAsyncTestCallTree(callTree,
		&TestConfig{
			ParentBalance: 1000,
		})
}

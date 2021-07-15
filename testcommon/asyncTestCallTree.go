package testcommon

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
)

type AsyncTestCall struct {
	ContractAddress []byte
	GroupName       string
	FunctionName    string
	CallbackName    string
}

func (call *AsyncTestCall) ToString() string {
	return "contract=" + string(call.ContractAddress) + " group=" + call.GroupName + " function=" + call.FunctionName + " callback=" + call.CallbackName
}

func BuildAsyncTestCall(contractId string, groupName string, functionName string, callbackName string) *AsyncTestCall {
	return &AsyncTestCall{
		ContractAddress: MakeTestSCAddress(contractId),
		GroupName:       groupName,
		FunctionName:    functionName,
		CallbackName:    callbackName,
	}
}

type AsyncTestCallNode struct {
	asyncCall *AsyncTestCall
	children  []*AsyncTestCallNode
}

type AsyncTestCallTree struct {
	root *AsyncTestCallNode
}

var emptyPath []*AsyncTestCallNode

func CreateCallTree(asyncCall *AsyncTestCall) *AsyncTestCallTree {
	return &AsyncTestCallTree{
		root: &AsyncTestCallNode{
			asyncCall: asyncCall,
		},
	}
}

func (tree *AsyncTestCallTree) GetRoot() *AsyncTestCallNode {
	return tree.root
}

func (node *AsyncTestCallNode) GetAsyncTestCall() *AsyncTestCall {
	return node.asyncCall
}

func (node *AsyncTestCallNode) AddChild(asyncCall *AsyncTestCall) *AsyncTestCallNode {
	newChildNode := &AsyncTestCallNode{asyncCall: asyncCall}
	node.children = append(node.children, newChildNode)
	return newChildNode
}

func (node *AsyncTestCallNode) isLeaf() bool {
	return node.children == nil
}

func copyTree(tree *AsyncTestCallTree) *AsyncTestCallTree {
	newRoot := dfs(nil, tree.root, emptyPath, func(path []*AsyncTestCallNode, parent *AsyncTestCallNode, node *AsyncTestCallNode) *AsyncTestCallNode {
		newNode := &AsyncTestCallNode{asyncCall: node.asyncCall}
		if parent != nil {
			parent.children = append(parent.children, newNode)
		}
		return newNode
	})
	return &AsyncTestCallTree{root: newRoot}
}

func printTree(tree *AsyncTestCallTree) {
	dfsTree(tree, printNode)
}

func dfsTree(tree *AsyncTestCallTree, processNode func([]*AsyncTestCallNode, *AsyncTestCallNode, *AsyncTestCallNode) *AsyncTestCallNode) {
	dfs(nil, tree.root, emptyPath, processNode)
}

func dfs(parent *AsyncTestCallNode, node *AsyncTestCallNode, path []*AsyncTestCallNode, processNode func([]*AsyncTestCallNode, *AsyncTestCallNode, *AsyncTestCallNode) *AsyncTestCallNode) *AsyncTestCallNode {
	path = append(path, node)
	processedParent := processNode(path, parent, node)
	for _, child := range node.children {
		dfs(processedParent, child, path, processNode)
	}
	return processedParent
}

func printNodeSimple(node *AsyncTestCallNode) *AsyncTestCallNode {
	fmt.Println(node.asyncCall.ToString())
	return node
}

func printNode(path []*AsyncTestCallNode, parent *AsyncTestCallNode, node *AsyncTestCallNode) *AsyncTestCallNode {
	level := len(path) - 1
	for t := 0; t < level; t++ {
		fmt.Print("\t")
	}

	fmt.Printf("%s ", node.asyncCall.ToString())
	// if path != nil {
	// 	fmt.Print("[ ")
	// 	for _, nodeInPath := range path {
	// 		fmt.Printf("%s ", node.asyncCall.ToString())
	// 	}
	// 	fmt.Print("]")
	// }
	fmt.Println()

	return node
}

func encodeTree(tree *AsyncTestCallTree) []byte {
	return encodePostorder(tree.root)
}

func encodePostorder(node *AsyncTestCallNode) []byte {
	encodedChildren := make([][]byte, 0)
	for _, child := range node.children {
		encodedChildren = append(encodedChildren, encodePostorder(child))
	}

	nodeCallData := txDataBuilder.NewBuilder()
	nodeCallData.Func(node.asyncCall.FunctionName)
	nodeCallData.Bytes(node.asyncCall.ContractAddress)
	nodeCallData.Str(node.asyncCall.GroupName)
	nodeCallData.Str(node.asyncCall.CallbackName)

	callData := txDataBuilder.NewBuilder()
	callData.Func(".")
	callData.Bytes(nodeCallData.ToBytes())
	for _, encodedChild := range encodedChildren {
		callData.Bytes(encodedChild)
	}

	fmt.Println("encoded " + node.asyncCall.ToString() + " as " + callData.ToString())
	return callData.ToBytes()
}

func decodeRoot(encodedRoot []byte) *AsyncTestCall {
	callArgsParser := parsers.NewCallArgsParser()
	functionName, args, err := callArgsParser.ParseData(string(encodedRoot))
	if err != nil {
		panic(err)
	}
	// fmt.Println(functionName)
	// fmt.Println(string(args[0]))
	// fmt.Println(string(args[1]))
	// fmt.Println(string(args[2]))
	// fmt.Println("-----------")
	return &AsyncTestCall{
		FunctionName:    functionName,
		ContractAddress: args[0],
		GroupName:       string(args[1]),
		CallbackName:    string(args[2]),
	}
}

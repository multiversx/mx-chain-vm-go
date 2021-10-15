package testcommon

import (
	"fmt"
	"strconv"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/crypto/factory"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

// DefaultCallGraphLockedGas is the default gas locked value
const DefaultCallGraphLockedGas = 150

// FakeCallbackName - used by test framework to reprezent visually a callback that is not present
const FakeCallbackName = "<>"

var logTestGraph = logger.GetOrCreate("arwen/testgraph")

// TestCall represents the payload of a node in the call graph
type TestCall struct {
	ContractAddress    []byte
	FunctionName       string
	CallID             []byte
	OriginalContractID string
}

// ToString - string representatin of a TestCall
func (call *TestCall) ToString() string {
	return "contract=" + string(call.ContractAddress) + " function=" + call.FunctionName
}

func (call *TestCall) copy() *TestCall {
	return &TestCall{
		ContractAddress:    call.ContractAddress,
		FunctionName:       call.FunctionName,
		CallID:             call.CallID,
		OriginalContractID: call.OriginalContractID,
	}
}

func buildTestCall(contractID string, functionName string) *TestCall {
	return &TestCall{
		ContractAddress:    MakeTestSCAddress(contractID),
		FunctionName:       functionName,
		CallID:             []byte{1},  // initial callID, should be updated when an edge is added
		OriginalContractID: contractID, // used for fake cross shard callbacks scenarios
	}
}

// TestCallNode is a node in the call graph
type TestCallNode struct {
	ID uint

	// entry point in call graph
	IsStartNode bool

	// node payload
	Call *TestCall
	// connected nodes
	AdjacentEdges []*TestCallEdge

	NonGasEdgeCounter int64

	// group callbacks
	// GroupCallbacks map[string]*TestCallNode
	// context callback
	// ContextCallback *TestCallNode

	// will be reseted after each dfs traversal
	// Visited bool

	// labels used only for visualization & debugging
	VisualLabel string
	// needs to be unique (will be in te form of contract_function_index)
	Label string

	// back pointer / "edge" to parent for trees (not part of the actual graph, not traversed)
	Parent       *TestCallNode
	IncomingEdge *TestCallEdge

	// info used for gas assertions
	// set from an incoming edge edge
	GasLimit  uint64
	GasUsed   uint64
	GasLocked uint64

	// computed info
	ExecutionRound           int
	MaxSubtreeExecutionRound int
	GasRemaining             uint64
	GasAccumulated           uint64

	// set automaticaly when the test is run
	CrtTxHash []byte

	ShardID uint32
	Fail    bool
}

// LeafLabel - special node label for leafs
const LeafLabel = "*"

// GetEdges gets the outgoing edges of the node
func (node *TestCallNode) GetEdges() []*TestCallEdge {
	return node.AdjacentEdges
}

// IsLeaf returns true if the node as any adjacent nodes
func (node *TestCallNode) IsLeaf() bool {
	return len(node.AdjacentEdges) == 0
}

// IsGasLeaf returns true if the node as any adjacent nodes and is "*" node
func (node *TestCallNode) IsGasLeaf() bool {
	return node.IsLeaf() && node.Call.FunctionName == LeafLabel
}

// GetIncomingEdgeType returns the type of the incoming edge (for a tree)
func (node *TestCallNode) GetIncomingEdgeType() TestCallEdgeType {
	if node.IncomingEdge == nil {
		return TestCallEdgeType(0)
	}
	return node.IncomingEdge.Type
}

// IsSync -
func (node *TestCallNode) IsSync() bool {
	return node.GetIncomingEdgeType() == Sync
}

// IsAsync -
func (node *TestCallNode) IsAsync() bool {
	incEdgeType := node.GetIncomingEdgeType()
	return incEdgeType == Async || incEdgeType == AsyncCrossShard
}

// IsCallback -
func (node *TestCallNode) IsCallback() bool {
	incEdgeType := node.GetIncomingEdgeType()
	return incEdgeType == Callback || incEdgeType == CallbackCrossShard
}

// HasCallback -
func (node *TestCallNode) HasCallback() bool {
	for _, edge := range node.AdjacentEdges {
		if edge.Type == Callback || edge.Type == CallbackCrossShard {
			return true
		}
	}
	return false
}

// Copy copyies the node call info into a new node
func (node *TestCallNode) copy() *TestCallNode {
	return &TestCallNode{
		Call:          node.Call.copy(),
		AdjacentEdges: make([]*TestCallEdge, 0),
		// GroupCallbacks:   make(map[string]*TestCallNode, 0),
		// ContextCallback:  nil,
		// Visited:          false,
		IsStartNode:  node.IsStartNode,
		Label:        node.Label,
		GasLimit:     node.GasLimit,
		GasRemaining: node.GasRemaining,
		GasUsed:      node.GasUsed,
		GasLocked:    node.GasLocked,
		// IncomingEdge: node.IncomingEdge,
		Fail: node.Fail,
	}
}

// IsIncomingEdgeFail -
func (node *TestCallNode) IsIncomingEdgeFail() bool {
	if node.IncomingEdge != nil && node.IncomingEdge.Fail {
		return true
	}
	return false
}

// WillNotExecute returns if node will execute based on execution round
func (node *TestCallNode) WillNotExecute() bool {
	return node.ExecutionRound == -1
}

// TestCallEdgeType the type of TestCallEdges
type TestCallEdgeType int

// types of TestCallEdges
const (
	Sync = iota
	Async
	Callback
	AsyncCrossShard
	CallbackCrossShard
	// GroupCallback
	// ContextCallback
)

// TestCallEdge an edge between two nodes of the call graph
type TestCallEdge struct {
	ID uint

	Type TestCallEdgeType

	// callback function name
	Callback string
	// Group    string

	// outgoing node
	To *TestCallNode

	// gas config for the outgoing node (represented by To)
	GasLimit          uint64
	GasUsed           uint64
	GasUsedByCallback uint64
	GasLocked         uint64

	// used only for visualization & debugging
	Label string

	Fail         bool
	CallbackFail bool
}

func (edge *TestCallEdge) copy() *TestCallEdge {
	return &TestCallEdge{
		ID:       edge.ID,
		Type:     edge.Type,
		Callback: edge.Callback,
		// Group:             edge.Group,
		To:                edge.To,
		GasLimit:          edge.GasLimit,
		GasUsed:           edge.GasUsed,
		GasUsedByCallback: edge.GasUsedByCallback,
		GasLocked:         edge.GasLocked,
		Label:             edge.Label,
		Fail:              edge.Fail,
	}
}

// SetGasLimit - builder style setter
func (edge *TestCallEdge) SetGasLimit(gasLimit uint64) *TestCallEdge {
	edge.GasLimit = gasLimit
	return edge
}

// SetGasUsedByCallback - builder style setter
func (edge *TestCallEdge) SetGasUsedByCallback(gasUsedByCallback uint64) *TestCallEdge {
	if edge.Type != Async && edge.Type != AsyncCrossShard {
		panic("Callbacks are only for async edges")
	}
	if edge.Callback == FakeCallbackName && gasUsedByCallback != 0 {
		panic("Callbacks not present, can't use gas")
	}
	edge.GasUsedByCallback = gasUsedByCallback
	return edge
}

// SetGasUsed - builder style setter
func (edge *TestCallEdge) SetGasUsed(gasUsed uint64) *TestCallEdge {
	edge.GasUsed = gasUsed
	return edge
}

// SetFail - builder style setter
func (edge *TestCallEdge) SetFail() *TestCallEdge {
	edge.Fail = true
	edge.To.Fail = true
	return edge
}

// SetCallbackFail - builder style setter
func (edge *TestCallEdge) SetCallbackFail() *TestCallEdge {
	edge.CallbackFail = true
	return edge
}

// SetGasLocked - builder style setter
func (edge *TestCallEdge) SetGasLocked(gasLocked uint64) *TestCallEdge {
	if edge.Type != Async && edge.Type != AsyncCrossShard {
		panic("Gas locked is only for async edges")
	}
	edge.GasLocked = gasLocked
	return edge
}

func (edge *TestCallEdge) copyAttributesFrom(sourceEdge *TestCallEdge) {
	edge.ID = sourceEdge.ID
	edge.Type = sourceEdge.Type
	edge.Callback = sourceEdge.Callback
	// edge.Group = sourceEdge.Group
	edge.GasLimit = sourceEdge.GasLimit
	edge.GasLocked = sourceEdge.GasLocked
	edge.GasUsed = sourceEdge.GasUsed
	edge.GasUsedByCallback = sourceEdge.GasUsedByCallback
	edge.Label = sourceEdge.Label
	edge.Fail = sourceEdge.Fail
}

// TestCallGraph is the call graph
type TestCallGraph struct {
	Nodes     []*TestCallNode
	StartNode *TestCallNode

	Crypto       crypto.VMCrypto
	sequence     uint
	edgeSequence uint
}

// CreateTestCallGraph is the initial build metohd for the call graph
func CreateTestCallGraph() *TestCallGraph {
	return &TestCallGraph{
		Nodes:  make([]*TestCallNode, 0),
		Crypto: factory.NewVMCrypto(),
	}
}

// AddStartNode adds the start node of the call graph
func (graph *TestCallGraph) AddStartNode(contractID string, functionName string, gasLimit uint64, gasUsed uint64) *TestCallNode {
	node := graph.AddNode(contractID, functionName)
	graph.StartNode = node
	node.IsStartNode = true
	node.GasLimit = gasLimit
	node.GasRemaining = 0
	node.GasUsed = gasUsed
	return node
}

// AddNodeCopy adds a copy of a node to the node list
func (graph *TestCallGraph) AddNodeCopy(node *TestCallNode) *TestCallNode {
	nodeCopy := node.copy()
	graph.sequence++
	nodeCopy.ID = graph.sequence
	graph.Nodes = append(graph.Nodes, nodeCopy)
	return nodeCopy
}

// AddNode adds a node to the call graph
func (graph *TestCallGraph) AddNode(contractID string, functionName string) *TestCallNode {
	graph.sequence++
	testCall := buildTestCall(contractID, functionName)
	testNode := &TestCallNode{
		ID:            graph.sequence,
		Call:          testCall,
		AdjacentEdges: make([]*TestCallEdge, 0),
		// GroupCallbacks:  make(map[string]*TestCallNode, 0),
		// ContextCallback: nil,
		// Visited:     false,
		IsStartNode: false,
		Label:       strconv.Quote(contractID + "." + functionName),
	}
	graph.Nodes = append(graph.Nodes, testNode)
	return testNode
}

// AddSyncEdge adds a labeled sync call edge between two nodes of the call graph
func (graph *TestCallGraph) AddSyncEdge(from *TestCallNode, to *TestCallNode) *TestCallEdge {
	edge := graph.addEdge(from, to, true)
	edge.Label = "Sync"
	return edge
}

// addEdge adds an edge between two nodes of the call graph
func (graph *TestCallGraph) addEdge(from *TestCallNode, to *TestCallNode, fillEdgeID bool) *TestCallEdge {
	edge := graph.buildEdge(to, fillEdgeID)
	to.Parent = from
	from.AdjacentEdges = append(from.AdjacentEdges, edge)
	return edge
}

func (graph *TestCallGraph) addEdgeAtStart(from *TestCallNode, to *TestCallNode, fillEdgeID bool) *TestCallEdge {
	edge := graph.buildEdge(to, fillEdgeID)
	to.Parent = from
	from.AdjacentEdges = append([]*TestCallEdge{edge}, from.AdjacentEdges...)
	return edge
}

func (graph *TestCallGraph) buildEdge(to *TestCallNode, fillEdgeID bool) *TestCallEdge {
	ID := uint(0)
	if fillEdgeID {
		graph.edgeSequence++
		ID = graph.edgeSequence
	}
	edge := &TestCallEdge{
		ID:       ID,
		Type:     Sync,
		Callback: "",
		// Group:    "",
		To: to,
	}
	return edge
}

// AddAsyncEdge adds a local async call edge between two nodes of the call graph
func (graph *TestCallGraph) AddAsyncEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) *TestCallEdge {
	return graph.addAsyncEdgeWithType(Async, from, to, callBack, group)
}

// AddAsyncCrossShardEdge adds a local async call edge between two nodes of the call graph
func (graph *TestCallGraph) AddAsyncCrossShardEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) *TestCallEdge {
	return graph.addAsyncEdgeWithType(AsyncCrossShard, from, to, callBack, group)
}

func (graph *TestCallGraph) addAsyncEdgeWithType(edgeType TestCallEdgeType, from *TestCallNode, to *TestCallNode, callBack string, group string) *TestCallEdge {
	gasLocked := uint64(0)
	if callBack != "" {
		gasLocked = DefaultCallGraphLockedGas
	} else if edgeType == AsyncCrossShard {
		callBack = FakeCallbackName
		graph.AddNode(from.Call.OriginalContractID, callBack)
	}

	graph.edgeSequence++
	edge := &TestCallEdge{
		ID:       graph.edgeSequence,
		Type:     edgeType,
		Callback: callBack,
		// Group:     group,
		To:        to,
		GasLocked: gasLocked,
	}
	edge.setAsyncEdgeAttributes(group, callBack)
	from.AdjacentEdges = append(from.AdjacentEdges, edge)
	return edge
}

func (edge *TestCallEdge) setAsyncEdgeAttributes(group string, callBack string) {
	edge.Label = "Async"
	// if group != "" {
	// 	edge.Label += "[" + group + "]"
	// }
	edge.Label += "\n"
	if callBack != "" {
		edge.Label += callBack
	}
}

// GetStartNode - start node getter
func (graph *TestCallGraph) GetStartNode() *TestCallNode {
	return graph.StartNode
}

// SetGroupCallback sets the callback for the specified group id
// func (graph *TestCallGraph) SetGroupCallback(node *TestCallNode, groupID string, groupCallbackNode *TestCallNode,
// 	gasLocked uint64, gasUsed uint64) {
// 	groupCallbackNode.GasLocked = gasLocked
// 	groupCallbackNode.GasUsed = gasUsed
// 	// node.GroupCallbacks[groupID] = groupCallbackNode
// }

// SetContextCallback sets the callback for the async context
// func (graph *TestCallGraph) SetContextCallback(node *TestCallNode, contextCallbackNode *TestCallNode,
// 	gasLocked uint64, gasUsed uint64) {
// 	node.ContextCallback = contextCallbackNode
// 	node.ContextCallback.GasLocked = gasLocked
// 	node.ContextCallback.GasUsed = gasUsed
// }

// FindNode finds the corresponding node in the call graph
func (graph *TestCallGraph) FindNode(contractAddress []byte, functionName string) *TestCallNode {
	// in the future we can use an index / map if this proves to be a performance problem
	// but for test call graphs we are ok
	for _, node := range graph.Nodes {
		if string(node.Call.ContractAddress) == string(contractAddress) &&
			node.Call.FunctionName == functionName {
			return node
		}
	}
	return nil
}

type processNodeFunc func([]*TestCallNode, *TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode

func isVisited(node *TestCallNode, visits map[uint]bool) bool {
	value, exists := visits[node.ID]
	if !exists {
		return false
	} else {
		return value
	}
}

func setVisited(node *TestCallNode, visits map[uint]bool) {
	visits[node.ID] = true
}

// DfsGraph a standard DFS traversal for the call graph
func (graph *TestCallGraph) DfsGraph(
	processNode processNodeFunc,
	followCrossShardEdges bool) {
	visits := make(map[uint]bool)
	for _, node := range graph.Nodes {
		if isVisited(node, visits) {
			continue
		}
		graph.dfsFromNode(nil, node, nil, make([]*TestCallNode, 0), processNode, visits, followCrossShardEdges)
	}
}

// DfsGraphFromNode standard DFS starting from a node
func (graph *TestCallGraph) DfsGraphFromNode(startNode *TestCallNode, processNode processNodeFunc,
	visits map[uint]bool,
	followCrossShardEdges bool) {
	graph.dfsFromNode(startNode.Parent, startNode, nil, make([]*TestCallNode, 0), processNode, visits, followCrossShardEdges)
}

func (graph *TestCallGraph) dfsFromNode(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge, path []*TestCallNode,
	processNode processNodeFunc,
	visits map[uint]bool,
	followCrossShardEdges bool) *TestCallNode {
	if isVisited(node, visits) {
		return node
	}

	path = append(path, node)
	processedParent := processNode(path, parent, node, incomingEdge)
	// a signal to stop DFS for this branch
	if processedParent == nil {
		return node
	}
	setVisited(node, visits)

	for _, edge := range node.AdjacentEdges {
		if !followCrossShardEdges && (edge.Type == AsyncCrossShard || edge.Type == CallbackCrossShard) {
			continue
		}
		graph.dfsFromNode(processedParent, edge.To, edge, path, processNode, visits, followCrossShardEdges)
	}
	return processedParent
}

func (graph *TestCallGraph) dfsFromNodeRunningOrder(
	parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge, path []*TestCallNode,
	processNode processNodeFunc,
	postProcessNode processNodeFunc,
	visits map[uint]bool) *TestCallNode {

	if isVisited(node, visits) {
		return node
	}

	path = append(path, node)
	processNode(path, parent, node, incomingEdge)
	// nodes configured as fail will stop DFS
	if node.Fail {
		// even if failed, async nodes need to traverse the callback branch
		if node.IsAsync() {
			callbackEdge := node.AdjacentEdges[len(node.AdjacentEdges)-1]
			if callbackEdge.To.IsCallback() {
				processedNode := graph.dfsFromNodeRunningOrder(node, callbackEdge.To, callbackEdge, path, processNode, postProcessNode, visits)
				// post proces callback
				postProcessNode(path, node, processedNode, incomingEdge)
				// post process failed async call
				postProcessNode(path, parent, node, parent.IncomingEdge)
			}
		}
		// stop DFS
		return nil
	}

	setVisited(node, visits)

	for _, edge := range node.AdjacentEdges {
		processedNode := graph.dfsFromNodeRunningOrder(node, edge.To, edge, path, processNode, postProcessNode, visits)
		// failed non-async branches will stop the DFS edge processing for current node
		if processedNode == nil && !edge.To.IsAsync() {
			// if the current node is an async node,
			// it will also immediately fail the current node itself
			if !node.IsAsync() {
				return nil
			}
			// for async nodes with failed branches, we call the callback and don't fail the node
			callbackEdge := node.AdjacentEdges[len(node.AdjacentEdges)-1]
			graph.dfsFromNodeRunningOrder(node, callbackEdge.To, callbackEdge, path, processNode, postProcessNode, visits)
			break
		}
		postProcessNode(path, node, processedNode, edge)
	}

	return node
}

// DfsFromNodeUntilFailures will stop DFS going deeper at first encountered fail
func (graph *TestCallGraph) DfsFromNodeUntilFailures(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge, path []*TestCallNode,
	processNode processNodeFunc,
	visits map[uint]bool) *TestCallNode {

	if isVisited(node, visits) {
		return node
	}

	path = append(path, node)
	processNode(path, parent, node, incomingEdge)
	// any failed node stops DFS (configured or not - due failure upstream propagation)
	if incomingEdge != nil && incomingEdge.Fail {
		// evan if failed, async nodes need to traverse the callback branch
		if node.IsAsync() {
			callbackEdge := node.AdjacentEdges[len(node.AdjacentEdges)-1]
			if callbackEdge.To.IsCallback() {
				graph.DfsFromNodeUntilFailures(node, callbackEdge.To, callbackEdge, path, processNode, visits)
			}
		}
		return node
	}

	setVisited(node, visits)

	for _, edge := range node.AdjacentEdges {
		graph.DfsFromNodeUntilFailures(node, edge.To, edge, path, processNode, visits)
	}

	return node
}

// DfsGraphFromNodePostOrder - standard post order DFS
func (graph *TestCallGraph) DfsGraphFromNodePostOrder(startNode *TestCallNode, processNode func(*TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) {
	visits := make(map[uint]bool)
	graph.dfsFromNodePostOrder(nil, startNode, nil, processNode, visits)
}

func (graph *TestCallGraph) dfsFromNodePostOrder(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge, processNode func(*TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode, visits map[uint]bool) *TestCallNode {
	for _, edge := range node.AdjacentEdges {
		graph.dfsFromNodePostOrder(node, edge.To, edge, processNode, visits)
	}

	if isVisited(node, visits) {
		return node
	}

	processedParent := processNode(parent, node, incomingEdge)
	setVisited(node, visits)

	return processedParent
}

func (graph *TestCallGraph) newGraphUsingNodes() *TestCallGraph {
	graphCopy := CreateTestCallGraph()

	for _, node := range graph.Nodes {
		graphCopy.AddNodeCopy(node)
	}

	// for _, nodeCopy := range graphCopy.Nodes {
	// 	node := graph.FindNode(nodeCopy.Call.ContractAddress, nodeCopy.Call.FunctionName)
	// 	for group, callBackNode := range node.GroupCallbacks {
	// 		nodeCopy.GroupCallbacks[group] = graphCopy.FindNode(callBackNode.Call.ContractAddress, callBackNode.Call.FunctionName)
	// 	}
	// }

	// for _, node := range graph.Nodes {
	// 	if node.ContextCallback != nil {
	// 		executionNode := graphCopy.FindNode(node.Call.ContractAddress, node.Call.FunctionName)
	// 		executionNode.ContextCallback = graphCopy.FindNode(node.ContextCallback.Call.ContractAddress, node.ContextCallback.Call.FunctionName)
	// 	}
	// }

	return graphCopy
}

// CreateExecutionGraphFromCallGraph - creates an execution graph from the call graph
func (graph *TestCallGraph) CreateExecutionGraphFromCallGraph() *TestCallGraph {

	executionGraph := graph.newGraphUsingNodes()
	executionGraph.StartNode = executionGraph.FindNode(
		graph.StartNode.Call.ContractAddress,
		graph.StartNode.Call.FunctionName)

	graph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		newSource := executionGraph.FindNode(node.Call.ContractAddress, node.Call.FunctionName)
		if node.IsLeaf() {
			addFinishNodeAsFirstEdge(executionGraph, newSource)
			return node
		}

		addSyncEdgesToExecutionGraph(node, executionGraph, newSource)

		addFinishNode(executionGraph, newSource)

		addAsyncEdgesToExecutionGraph(node, executionGraph, newSource)

		// callbacks were added by async source node processing and must be moved to the end of the node
		// after all other node activity (sync & async calls)
		moveCallbacksToTheEndOfEdges(newSource)

		// add group callbacks calls if any
		// for _, group := range groups {
		// 	groupCallbackNode := newSource.GroupCallbacks[group]
		// 	if groupCallbackNode != nil {
		// 		execEdge := executionGraph.addEdge(newSource, groupCallbackNode)
		// 		execEdge.Type = GroupCallback
		// 		execEdge.Label = "Callback[" + group + "]"
		// 		execEdge.GasUsed = groupCallbackNode.GasUsed
		// 		execEdge.GasLocked = groupCallbackNode.GasLocked
		// 	}
		// }

		// // is start node add context callback
		// if newSource.ContextCallback != nil {
		// 	execEdge := executionGraph.addEdge(newSource, newSource.ContextCallback)
		// 	execEdge.Type = ContextCallback
		// 	execEdge.Label = "Callback\nContext"
		// 	execEdge.GasUsed = newSource.ContextCallback.GasUsed
		// 	execEdge.GasLocked = newSource.ContextCallback.GasLocked
		// }

		return node
	}, true)
	return executionGraph
}

func moveCallbacksToTheEndOfEdges(newSource *TestCallNode) {
	nonCallBackEdges := make([]*TestCallEdge, 0)
	callBackEdges := make([]*TestCallEdge, 0)
	for _, newEdge := range newSource.AdjacentEdges {
		if newEdge.Type == Callback || newEdge.Type == CallbackCrossShard {
			callBackEdges = append(callBackEdges, newEdge)
		} else {
			nonCallBackEdges = append(nonCallBackEdges, newEdge)
		}
	}
	newSource.AdjacentEdges = append(nonCallBackEdges, callBackEdges...)
}

func addSyncEdgesToExecutionGraph(node *TestCallNode, executionGraph *TestCallGraph, newSource *TestCallNode) {
	for _, edge := range node.AdjacentEdges {
		if edge.Type != Sync {
			continue
		}

		originalDestination := edge.To.Call
		newDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)

		execEdge := executionGraph.addEdge(newSource, newDestination, false)
		execEdge.copyAttributesFrom(edge)
	}
}

func addAsyncEdgesToExecutionGraph(node *TestCallNode, executionGraph *TestCallGraph, newSource *TestCallNode) {
	// groups := make([]string, 0)
	for _, edge := range node.AdjacentEdges {
		if edge.Type != Async && edge.Type != AsyncCrossShard {
			continue
		}

		// crtGroup := edge.Group
		// if !isGroupPresent(crtGroup, groups) {
		// 	groups = append(groups, crtGroup)
		// }

		originalDestination := edge.To.Call
		newAsyncDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)

		execEdge := executionGraph.addEdge(newSource, newAsyncDestination, false)
		execEdge.copyAttributesFrom(edge)

		if edge.Callback != "" {
			callbackDestination := executionGraph.FindNode(node.Call.ContractAddress, edge.Callback)
			if callbackDestination == nil {
				panic(fmt.Sprintf("Cant find node %s %s", node.Call.ContractAddress, edge.Callback))
			}
			execEdge := executionGraph.addEdge(newAsyncDestination, callbackDestination, false)
			if edge.CallbackFail {
				execEdge.Fail = true
				execEdge.To.Fail = true
			}
			if edge.Type == Async {
				execEdge.Type = Callback
			} else {
				execEdge.Type = CallbackCrossShard
			}
			execEdge.GasUsedByCallback = edge.GasUsedByCallback
			execEdge.Label = "Callback"
		}
	}
}

// add a new 'finish' edge to a special end of sync execution node
func addFinishNode(graph *TestCallGraph, sourceNode *TestCallNode) {
	addFinishNodeWithEdgeFunc(graph, sourceNode, graph.addEdge)
}

func addFinishNodeAsFirstEdge(graph *TestCallGraph, sourceNode *TestCallNode) {
	addFinishNodeWithEdgeFunc(graph, sourceNode, graph.addEdgeAtStart)
}

func addFinishNodeWithEdgeFunc(graph *TestCallGraph, sourceNode *TestCallNode, addEdge func(*TestCallNode, *TestCallNode, bool) *TestCallEdge) {
	finishNode := buildFinishNode(graph, sourceNode)
	addEdge(sourceNode, finishNode, false)
}

func buildFinishNode(graph *TestCallGraph, sourceNode *TestCallNode) *TestCallNode {
	finishNode := graph.AddNode("", LeafLabel)
	finishNode.Label = LeafLabel
	return finishNode
}

// TestCallPath a path in a tree, len(edges) = len(nodes) - 1
type TestCallPath struct {
	nodes []*TestCallNode
	edges []*TestCallEdge
}

func addToPath(path *TestCallPath, edge *TestCallEdge) *TestCallPath {
	return &TestCallPath{
		nodes: append(path.nodes, edge.To),
		edges: append(path.edges, edge),
	}
}

// deep copy of a path
func (path *TestCallPath) copy() *TestCallPath {
	return &TestCallPath{
		nodes: copyNodesList(path.nodes),
		edges: copyEdgeList(path.edges),
	}
}

// deep copy of a node list
func copyNodesList(source []*TestCallNode) []*TestCallNode {
	dest := make([]*TestCallNode, len(source))
	for idxNode, node := range source {
		dest[idxNode] = node.copy()
	}
	return dest
}

// deep copy of an edge list
func copyEdgeList(source []*TestCallEdge) []*TestCallEdge {
	dest := make([]*TestCallEdge, len(source))
	for idxEdge, edge := range source {
		dest[idxEdge] = edge.copy()
	}
	return dest
}

// gets all the paths (as a list) from a DAG
func (graph *TestCallGraph) getPaths() []*TestCallPath {
	path := &TestCallPath{
		nodes: []*TestCallNode{graph.GetStartNode()},
		edges: make([]*TestCallEdge, 0),
	}
	paths := make([]*TestCallPath, 0)
	graph.getPathsRecursive(path, func(newPath *TestCallPath) {
		newPath.print()
		paths = append(paths, newPath.copy())
	})
	return paths
}

// follow the paths in DAG, bun only the allowed paths
// (we can have multiple INs and OUTs for a node, but not all pairs will be paths)
func (graph *TestCallGraph) getPathsRecursive(path *TestCallPath, addPathToResult func(*TestCallPath)) {
	lastNodeInPath := path.nodes[len(path.nodes)-1]

	if lastNodeInPath.IsLeaf() {
		lastNodeInPath.GasUsed = path.nodes[len(path.nodes)-2].GasUsed
		lastNodeInPath.GasLimit = lastNodeInPath.GasUsed
		addPathToResult(path)
		LogGraph.Trace("end of path")
		return
	}

	var lastEdgeInPath *TestCallEdge
	if len(path.nodes) > 1 {
		lastEdgeInPath = path.edges[len(path.nodes)-2]
	}

	// for each outgoing edge from the last node in path, if it's allowed to continue on that edge from
	// the current path, add the next node to the current path and recurse
	for _, edge := range lastNodeInPath.AdjacentEdges {
		if
		// don't follow a path in the form of sync -> callback
		(lastEdgeInPath != nil && lastEdgeInPath.Type == Sync && (edge.Type == Callback || edge.Type == CallbackCrossShard)) ||
			// don't follow mixed local / cross-shard paths
			(lastEdgeInPath != nil && lastEdgeInPath.Type == Async && edge.Type == CallbackCrossShard) ||
			(lastEdgeInPath != nil && lastEdgeInPath.Type == AsyncCrossShard && edge.Type == Callback) ||
			// don't follow a path from async -> callback that is not your own
			(lastEdgeInPath != nil &&
				((lastEdgeInPath.Type == Async && edge.Type == Callback) ||
					(lastEdgeInPath.Type == AsyncCrossShard && edge.Type == CallbackCrossShard)) &&
				(lastEdgeInPath.Callback != edge.To.Call.FunctionName)) {
			continue
		}

		edge.To.GasLimit = edge.GasLimit
		edge.To.GasRemaining = 0
		edge.To.GasLocked = edge.GasLocked
		if lastEdgeInPath != nil &&
			((lastEdgeInPath.Type == Async && edge.Type == Callback) ||
				(lastEdgeInPath.Type == AsyncCrossShard && edge.Type == CallbackCrossShard)) {
			edge.To.GasUsed = lastEdgeInPath.GasUsedByCallback
		} else {
			edge.To.GasUsed = edge.GasUsed
		}
		edge.To.IncomingEdge = edge

		LogGraph.Trace("add [" + edge.Label + "] " + edge.To.Label)
		newPath := addToPath(path, edge)
		graph.getPathsRecursive(newPath, addPathToResult)
	}
	LogGraph.Trace("end of edges for " + lastNodeInPath.Label)
}

func (path *TestCallPath) print() {
	LogGraph.Trace("path = ")
	for pathIdx, node := range path.nodes {
		if pathIdx > 0 {
			LogGraph.Trace(" [" + eliminateNewLines(path.edges[pathIdx-1].Label) + "#" + strconv.Itoa(int(path.edges[pathIdx-1].ID)) + "] ")
			LogGraph.Trace(" / ")
		}
		LogGraph.Trace(node.Label + " ")
		LogGraph.Trace(" / ")
	}
}

func eliminateNewLines(input string) string {
	result := ""
	for _, c := range input {
		if c != '\n' {
			result += string(c)
		} else {
			result += "|"
		}
	}
	return result
}

// ComputeGasGraphFromExecutionGraph - creates a gas graph from an execution graph
func (graph *TestCallGraph) ComputeGasGraphFromExecutionGraph() *TestCallGraph {
	return pathsTreeFromDag(graph)
}

// This will create a list of paths from the DAG and merge them into a call tree
// The merging is done by following the paths in the building tree and complete with nodes from the paths when necessary
func pathsTreeFromDag(graph *TestCallGraph) *TestCallGraph {
	newGraph := CreateTestCallGraph()

	paths := graph.getPaths()

	LogGraph.Trace("process path")
	var crtNode *TestCallNode
	for _, path := range paths {
	nextNode:
		for pathIdx, node := range path.nodes {
			if pathIdx == 0 {
				crtNode = newGraph.FindNode(node.Call.ContractAddress, node.Call.FunctionName)
				if crtNode == nil {
					crtNode = newGraph.AddNodeCopy(node)
					newGraph.StartNode = crtNode
				}
				continue
			}
			for _, edge := range crtNode.AdjacentEdges {
				crtChild := edge.To
				LogGraph.Trace(edge.Label + "==" + path.edges[pathIdx-1].Label + "\n=>" + strconv.FormatBool(edge.Label == path.edges[pathIdx-1].Label))
				if string(crtChild.Call.ContractAddress) == string(node.Call.ContractAddress) &&
					crtChild.Call.FunctionName == node.Call.FunctionName &&
					edge.Label == path.edges[pathIdx-1].Label &&
					edge.ID == path.edges[pathIdx-1].ID {
					crtNode = crtChild
					continue nextNode
				}
			}
			parent := crtNode
			crtNode = newGraph.AddNodeCopy(node)

			LogGraph.Trace("add edge " + parent.Label + " -> " + crtNode.Label)
			pathEdge := path.edges[pathIdx-1]
			newEdge := newGraph.addEdge(parent, crtNode, false)
			newEdge.To.IncomingEdge = newEdge // these are edges in a tree
			newEdge.copyAttributesFrom(pathEdge)
		}
	}

	return newGraph
}

func isGroupPresent(group string, groups []string) bool {
	for _, crtGroup := range groups {
		if group == crtGroup {
			return true
		}
	}
	return false
}

// PropagateSyncFailures -
func (graph *TestCallGraph) PropagateSyncFailures() {
	graph.DfsGraphFromNodePostOrder(graph.StartNode, func(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		if node.IsLeaf() || node.IsCallback() || node.WillNotExecute() {
			return node
		}

		if node.IsIncomingEdgeFail() {
			// propagate failure to parent until we reach an async node
			if parent != nil && incomingEdge.Type != Async && incomingEdge.Type != AsyncCrossShard {
				parent.IncomingEdge.Fail = true
			}
		}

		return node
	})
}

// AssignExecutionRounds -
func (graph *TestCallGraph) AssignExecutionRounds() {
	visits := make(map[uint]bool)

	// init execution rounds for graph, all -1 except root and it's execution leaf
	for _, node := range graph.Nodes {
		node.ExecutionRound = -1
	}
	graph.StartNode.ExecutionRound = 0
	for _, edge := range graph.StartNode.AdjacentEdges {
		if edge.To.IsLeaf() {
			edge.To.ExecutionRound = 0
			break
		}
	}

	graph.dfsFromNodeRunningOrder(graph.StartNode.Parent, graph.StartNode, nil, make([]*TestCallNode, 0),
		func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
			if incomingEdge == nil || node.IsGasLeaf() {
				return node
			}

			switch incomingEdge.Type {
			case Sync:
				node.ExecutionRound = parent.ExecutionRound
				break
			case Async:
				node.ExecutionRound = parent.ExecutionRound
				break
			case AsyncCrossShard:
				node.ExecutionRound = parent.MaxSubtreeExecutionRound + 1
				break
			case Callback:
				node.ExecutionRound = parent.MaxSubtreeExecutionRound
				break
			case CallbackCrossShard:
				node.ExecutionRound = parent.MaxSubtreeExecutionRound + 1
				break
			}

			node.MaxSubtreeExecutionRound = node.ExecutionRound
			getGasLeaf(node).ExecutionRound = node.ExecutionRound

			return node
		},
		func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
			if node == nil {
				return node
			}
			if parent.MaxSubtreeExecutionRound < node.MaxSubtreeExecutionRound {
				parent.MaxSubtreeExecutionRound = node.MaxSubtreeExecutionRound
			}
			return node
		}, visits)
}

func getGasLeaf(node *TestCallNode) *TestCallNode {
	for _, edge := range node.AdjacentEdges {
		if edge.To.IsGasLeaf() {
			return edge.To
		}
	}
	return nil
}

// ComputeRemainingGasBeforeCallbacks - adjusts the gas graph / tree remaining gas info using the gas provided to children
// this will not take into consideration callback nodes that don't have provided gas info computed yet (see ComputeGasStepByStep)
func (graph *TestCallGraph) ComputeRemainingGasBeforeCallbacks() {
	graph.DfsGraphFromNodePostOrder(graph.StartNode, func(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		if node.IsLeaf() || node.IsCallback() || node.WillNotExecute() {
			return node
		}

		if node.IsIncomingEdgeFail() {
			node.GasRemaining = 0
		} else {
			nodeGasRemaining := int64(node.GasLimit)
			for _, edge := range node.AdjacentEdges {
				if edge.Type != Callback && edge.Type != CallbackCrossShard {
					nodeGasRemaining -= int64(edge.To.GasLimit + edge.To.GasLocked)
					if nodeGasRemaining < 0 {
						var incomingEdgeLabel string
						if incomingEdge != nil {
							incomingEdgeLabel = incomingEdge.Label
						}
						panic(fmt.Sprintf("Bad test gas configuration %s incoming edge '%s'", node.Label, incomingEdgeLabel))
					}
				}
			}

			node.GasRemaining += uint64(nodeGasRemaining)
			if parent != nil && incomingEdge.Type != Async && incomingEdge.Type != AsyncCrossShard {
				parent.GasRemaining += uint64(node.GasRemaining)
			}
		}

		if node.IsAsync() && !node.HasCallback() {
			node.GasAccumulated += node.GasRemaining
		}

		for _, edge := range node.AdjacentEdges {
			if edge.Type == Callback || edge.Type == CallbackCrossShard {
				edge.To.GasLimit = node.GasRemaining + node.GasLocked
			}
		}
		return node
	})
}

// ComputeRemainingGasAfterCallbacks -
func (graph *TestCallGraph) ComputeRemainingGasAfterCallbacks() {
	graph.DfsGraphFromNodePostOrder(graph.StartNode, func(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		// process only executing callbacks
		if !node.IsCallback() || node.WillNotExecute() {
			return node
		}

		nodeGasRemaining := int64(node.GasLimit)
		for _, edge := range node.AdjacentEdges {
			nodeGasRemaining -= int64(edge.To.GasLimit + edge.To.GasLocked)
			if nodeGasRemaining < 0 {
				var incomingEdgeLabel string
				if incomingEdge != nil {
					incomingEdgeLabel = incomingEdge.Label
				}
				panic(fmt.Sprintf("Bad test gas configuration %s incoming edge '%s'", node.Label, incomingEdgeLabel))
			}
		}
		node.GasRemaining += uint64(nodeGasRemaining)

		// propagate remaining gas
		asyncInitiator := node.Parent.Parent
		if asyncInitiator.ExecutionRound != node.ExecutionRound {
			// CallbackCrossShard
			if !node.IsIncomingEdgeFail() {
				node.GasAccumulated += node.GasRemaining
			} else {
				node.GasRemaining = 0
				node.GasAccumulated = 0
			}
		} else {
			// Callback
			if !node.IsIncomingEdgeFail() {
				crtEdgeType := asyncInitiator.GetIncomingEdgeType()
				for crtParent := asyncInitiator; crtParent != nil; crtParent = crtParent.Parent {
					if crtParent == asyncInitiator || crtEdgeType != Sync {
						crtParent.GasAccumulated += node.GasRemaining
					} else {
						crtParent.GasRemaining += node.GasRemaining
					}

					crtEdgeType = crtParent.GetIncomingEdgeType()
					if crtEdgeType == Async || crtEdgeType == AsyncCrossShard {
						for _, edge := range crtParent.AdjacentEdges {
							if edge.Type == Callback || edge.Type == CallbackCrossShard {
								edge.To.GasLimit += node.GasRemaining
							}
						}
					}
					if crtEdgeType == Async || crtEdgeType == AsyncCrossShard ||
						crtEdgeType == Callback || crtEdgeType == CallbackCrossShard {
						break
					}
				}
			} else {
				node.GasRemaining = 0
			}
		}

		return node
	})
}

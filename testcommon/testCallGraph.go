package testcommon

import (
	"fmt"
	"strconv"
)

// DefaultCallGraphLockedGas is the default gas locked value
const DefaultCallGraphLockedGas = 150

// TestCall represents the payload of a node in the call graph
type TestCall struct {
	ContractAddress []byte
	FunctionName    string
}

// ToString - string representatin of a TestCall
func (call *TestCall) ToString() string {
	return "contract=" + string(call.ContractAddress) + " function=" + call.FunctionName
}

func (call *TestCall) copy() *TestCall {
	return &TestCall{
		ContractAddress: call.ContractAddress,
		FunctionName:    call.FunctionName,
	}
}

func buildTestCall(contractID string, functionName string) *TestCall {
	return &TestCall{
		ContractAddress: MakeTestSCAddress(contractID),
		FunctionName:    functionName,
	}
}

// TestCallNode is a node in the call graph
type TestCallNode struct {
	// node payload
	Call *TestCall
	// connected nodes
	AdjacentEdges []*TestCallEdge
	// group callbacks
	GroupCallbacks map[string]*TestCallNode
	// context callback
	ContextCallback *TestCallNode
	// used by execution graphs - by default these nodes are ignored by FindNode() calls
	// TODO matei-p is this necessary? all these are leafs !!!
	// we could rename it, to IsLeaf because in certain cases we don't have yet the edges
	IsEndOfSyncExecutionNode bool
	IsStartNode              bool
	// will be reseted after each dfs traversal
	Visited bool
	// used only for visualization
	VisualLabel string
	Label       string
	// back pointer / "edge" to parent for trees (not part of the graph)
	Parent           *TestCallNode
	IncomingEdgeType TestCallEdgeType
	// info used for gas assertions
	// set from an incoming edge edge
	GasLimit  uint64
	GasUsed   uint64
	GasLocked uint64
	// computed info
	GasRemaining              uint64
	GasRemainingAfterCallback uint64
}

// SpecialLabel - special node label for IsEndOfSyncExecutionNode
const SpecialLabel = "*"

// GetEdges gets the outgoing edges of the node
func (node *TestCallNode) GetEdges() []*TestCallEdge {
	return node.AdjacentEdges
}

// IsLeaf returns true if the node as any adjacent nodes
func (node *TestCallNode) IsLeaf() bool {
	return len(node.AdjacentEdges) == 0
}

// Copy copyies the node call info into a new node
func (node *TestCallNode) copy() *TestCallNode {
	return &TestCallNode{
		Call:                     node.Call.copy(),
		AdjacentEdges:            make([]*TestCallEdge, 0),
		GroupCallbacks:           make(map[string]*TestCallNode, 0),
		ContextCallback:          nil,
		IsEndOfSyncExecutionNode: node.IsEndOfSyncExecutionNode,
		Visited:                  false,
		IsStartNode:              node.IsStartNode,
		Label:                    node.Label,
		GasLimit:                 node.GasLimit,
		GasRemaining:             node.GasRemaining,
		GasUsed:                  node.GasUsed,
		GasLocked:                node.GasLocked,
		IncomingEdgeType:         node.IncomingEdgeType,
	}
}

// TestCallEdgeType the type of TestCallEdges
type TestCallEdgeType int

// types of TestCallEdges
const (
	Sync = iota
	Async
	Callback
	GroupCallback
	ContextCallback
)

// TestCallEdge an edge between two nodes of the call graph
type TestCallEdge struct {
	Type     TestCallEdgeType
	Callback string
	Group    string
	To       *TestCallNode
	// gas config for the outgoing node (To)
	GasLimit          uint64
	GasUsed           uint64
	GasUsedByCallback uint64
	GasLocked         uint64
	// used only for visualization
	Label string
	Color string
}

func (edge *TestCallEdge) copy() *TestCallEdge {
	return &TestCallEdge{
		Type:              edge.Type,
		Callback:          edge.Callback,
		Group:             edge.Group,
		To:                edge.To,
		GasLimit:          edge.GasLimit,
		GasUsed:           edge.GasUsed,
		GasUsedByCallback: edge.GasUsedByCallback,
		GasLocked:         edge.GasLocked,
		Label:             edge.Label,
		Color:             edge.Color,
	}
}

// SetGasLimit - builder style setter
func (edge *TestCallEdge) SetGasLimit(gasLimit uint64) *TestCallEdge {
	edge.GasLimit = gasLimit
	return edge
}

// SetGasUsedByCallback - builder style setter
func (edge *TestCallEdge) SetGasUsedByCallback(gasUsedByCallback uint64) *TestCallEdge {
	if edge.Type != Async {
		panic("Callbacks are only for async edges")
	}
	edge.GasUsedByCallback = gasUsedByCallback
	return edge
}

// SetGasUsed - builder style setter
func (edge *TestCallEdge) SetGasUsed(gasUsed uint64) *TestCallEdge {
	edge.GasUsed = gasUsed
	return edge
}

// SetGasLocked - builder style setter
func (edge *TestCallEdge) SetGasLocked(gasLocked uint64) *TestCallEdge {
	if edge.Type != Async {
		panic("Gas locked is only for async edges")
	}
	edge.GasLocked = gasLocked
	return edge
}

func (edge *TestCallEdge) copyAttributesFrom(sourceEdge *TestCallEdge) {
	edge.Type = sourceEdge.Type
	edge.Callback = sourceEdge.Callback
	edge.Group = sourceEdge.Group
	edge.GasLimit = sourceEdge.GasLimit
	edge.GasLocked = sourceEdge.GasLocked
	edge.GasUsed = sourceEdge.GasUsed
	edge.GasUsedByCallback = sourceEdge.GasUsedByCallback
	edge.Color = sourceEdge.Color
	edge.Label = sourceEdge.Label
}

// TestCallGraph is the call graph
type TestCallGraph struct {
	Nodes     []*TestCallNode
	StartNode *TestCallNode
}

// CreateTestCallGraph is the initial build metohd for the call graph
func CreateTestCallGraph() *TestCallGraph {
	return &TestCallGraph{
		Nodes: make([]*TestCallNode, 0),
	}
}

// AddStartNode adds the start node of the call graph
func (graph *TestCallGraph) AddStartNode(contractID string, functionName string, gasLimit uint64, gasUsed uint64) *TestCallNode {
	node := graph.AddNode(contractID, functionName)
	graph.StartNode = node
	node.IsStartNode = true
	node.GasLimit = gasLimit
	node.GasRemaining = gasLimit
	node.GasUsed = gasUsed
	return node
}

// AddNodeCopy adds a copy of a node to the node list
func (graph *TestCallGraph) AddNodeCopy(node *TestCallNode) *TestCallNode {
	nodeCopy := node.copy()
	graph.Nodes = append(graph.Nodes, nodeCopy)
	return nodeCopy
}

// AddNode adds a node to the call graph
func (graph *TestCallGraph) AddNode(contractID string, functionName string) *TestCallNode {
	testCall := buildTestCall(contractID, functionName)
	testNode := &TestCallNode{
		Call:                     testCall,
		AdjacentEdges:            make([]*TestCallEdge, 0),
		GroupCallbacks:           make(map[string]*TestCallNode, 0),
		ContextCallback:          nil,
		Visited:                  false,
		IsEndOfSyncExecutionNode: false,
		IsStartNode:              false,
		Label:                    strconv.Quote(contractID + "." + functionName),
	}
	graph.Nodes = append(graph.Nodes, testNode)
	return testNode
}

// AddSyncEdge adds a labeled sync call edge between two nodes of the call graph
func (graph *TestCallGraph) AddSyncEdge(from *TestCallNode, to *TestCallNode) *TestCallEdge {
	edge := graph.addEdge(from, to)
	edge.Label = "Sync"
	edge.Color = "blue"
	return edge
}

// addEdge adds an edge between two nodes of the call graph
func (graph *TestCallGraph) addEdge(from *TestCallNode, to *TestCallNode) *TestCallEdge {
	edge := buildEdge(to)
	to.Parent = from
	from.AdjacentEdges = append(from.AdjacentEdges, edge)
	return edge
}

func (graph *TestCallGraph) addEdgeAtStart(from *TestCallNode, to *TestCallNode) *TestCallEdge {
	edge := buildEdge(to)
	to.Parent = from
	from.AdjacentEdges = append([]*TestCallEdge{edge}, from.AdjacentEdges...)
	return edge
}

func buildEdge(to *TestCallNode) *TestCallEdge {
	edge := &TestCallEdge{
		Type:     Sync,
		Callback: "",
		Group:    "",
		To:       to,
	}
	return edge
}

// AddAsyncEdge adds an async call edge between two nodes of the call graph
func (graph *TestCallGraph) AddAsyncEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) *TestCallEdge {
	edge := &TestCallEdge{
		Type:      Async,
		Callback:  callBack,
		Group:     group,
		To:        to,
		GasLocked: DefaultCallGraphLockedGas,
	}
	edge.setAsyncEdgeAttributes(group, callBack)
	from.AdjacentEdges = append(from.AdjacentEdges, edge)
	return edge
}

func (edge *TestCallEdge) setAsyncEdgeAttributes(group string, callBack string) {
	edge.Label = "Async"
	edge.Color = "red"
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
func (graph *TestCallGraph) SetGroupCallback(node *TestCallNode, groupID string, groupCallbackNode *TestCallNode,
	gasLocked uint64, gasUsed uint64) {
	groupCallbackNode.GasLocked = gasLocked
	groupCallbackNode.GasUsed = gasUsed
	node.GroupCallbacks[groupID] = groupCallbackNode
}

// SetContextCallback sets the callback for the async context
func (graph *TestCallGraph) SetContextCallback(node *TestCallNode, contextCallbackNode *TestCallNode,
	gasLocked uint64, gasUsed uint64) {
	node.ContextCallback = contextCallbackNode
	node.ContextCallback.GasLocked = gasLocked
	node.ContextCallback.GasUsed = gasUsed
}

// FindNode finds the corresponding node in the call graph
func (graph *TestCallGraph) FindNode(contractAddress []byte, functionName string) *TestCallNode {
	// in the future we can use an index / map if this proves to be a performance problem
	// but for test call graphs we are ok
	for _, node := range graph.Nodes {
		if string(node.Call.ContractAddress) == string(contractAddress) &&
			node.Call.FunctionName == functionName &&
			!node.IsEndOfSyncExecutionNode {
			return node
		}
	}
	return nil
}

// DfsGraph a standard DFS traversal for the call graph
func (graph *TestCallGraph) DfsGraph(processNode func([]*TestCallNode, *TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) {
	for _, node := range graph.Nodes {
		if node.Visited {
			continue
		}
		graph.dfsFromNode(nil, node, nil, make([]*TestCallNode, 0), processNode)
	}
	graph.clearVisitedNodesFlag()
}

// DfsGraphFromNode standard DFS starting from a node
// stopAtVisited set to false, enables traversal of callexecution graphs (with the risk of infinite cycles)
func (graph *TestCallGraph) DfsGraphFromNode(startNode *TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) {
	graph.dfsFromNode(nil, startNode, nil, make([]*TestCallNode, 0), processNode)
	graph.clearVisitedNodesFlag()
}

func (graph *TestCallGraph) clearVisitedNodesFlag() {
	for _, node := range graph.Nodes {
		node.Visited = false
	}
}

func (graph *TestCallGraph) dfsFromNode(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge, path []*TestCallNode,
	processNode func([]*TestCallNode, *TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) *TestCallNode {
	if node.Visited {
		return node
	}

	path = append(path, node)
	processedParent := processNode(path, parent, node, incomingEdge)
	node.Visited = true

	for _, edge := range node.AdjacentEdges {
		graph.dfsFromNode(processedParent, edge.To, edge, path, processNode)
	}
	return processedParent
}

// DfsGraphFromNodePostOrder - standard post order DFS
func (graph *TestCallGraph) DfsGraphFromNodePostOrder(startNode *TestCallNode, processNode func(*TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) {
	graph.dfsFromNodePostOrder(nil, startNode, nil, processNode)
	graph.clearVisitedNodesFlag()
}

func (graph *TestCallGraph) dfsFromNodePostOrder(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge, processNode func(*TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) *TestCallNode {
	for _, edge := range node.AdjacentEdges {
		graph.dfsFromNodePostOrder(node, edge.To, edge, processNode)
	}

	if node.Visited {
		return node
	}

	processedParent := processNode(parent, node, incomingEdge)
	node.Visited = true

	return processedParent
}

// OneStepDfsFromNodePostOrder -
func (graph *TestCallGraph) OneStepDfsFromNodePostOrder(processNode func(*TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) bool {
	if graph.StartNode.Visited {
		return true
	}
	graph.oneStepDfsFromNodePostOrder(nil, graph.StartNode, nil, processNode)
	return false
}

func (graph *TestCallGraph) oneStepDfsFromNodePostOrder(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge, processNode func(*TestCallNode, *TestCallNode, *TestCallEdge) *TestCallNode) bool {
	for _, edge := range node.AdjacentEdges {
		if graph.oneStepDfsFromNodePostOrder(node, edge.To, edge, processNode) {
			return true
		}
	}

	if node.Visited {
		return false
	}

	processNode(parent, node, incomingEdge)
	node.Visited = true

	return true
}

// TestCallPath a path in a tree
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

func (path *TestCallPath) copy() *TestCallPath {
	return &TestCallPath{
		nodes: copyNodesList(path.nodes),
		edges: copyEdgeList(path.edges),
	}
}

func copyNodesList(source []*TestCallNode) []*TestCallNode {
	dest := make([]*TestCallNode, len(source))
	for idxNode, node := range source {
		dest[idxNode] = node.copy()
	}
	return dest
}

func copyEdgeList(source []*TestCallEdge) []*TestCallEdge {
	dest := make([]*TestCallEdge, len(source))
	for idxEdge, edge := range source {
		dest[idxEdge] = edge.copy()
	}
	return dest
}

func (graph *TestCallGraph) getPaths() []*TestCallPath {
	path := &TestCallPath{
		nodes: []*TestCallNode{graph.GetStartNode()},
		edges: make([]*TestCallEdge, 0),
	}
	paths := make([]*TestCallPath, 0)
	graph.getPathsRecursive(path, func(newPath *TestCallPath) {
		paths = append(paths, newPath.copy())
	})
	return paths
}

func (graph *TestCallGraph) getPathsRecursive(path *TestCallPath, addPathToResult func(*TestCallPath)) {
	lastNodeInPath := path.nodes[len(path.nodes)-1]
	if lastNodeInPath.IsLeaf() {
		lastNodeInPath.GasUsed = path.nodes[len(path.nodes)-2].GasUsed
		// path.nodes[len(path.nodes)-2].GasUsed = 0 // TODO matei-p
		lastNodeInPath.GasLimit = lastNodeInPath.GasUsed
		addPathToResult(path)
		// path.print()
		return
	}

	var lastEdgeInPath *TestCallEdge
	if len(path.nodes) > 1 {
		lastEdgeInPath = path.edges[len(path.nodes)-2]
	}

	// if lastEdgeInPath != nil {
	// 	fmt.Println("-> [" + lastEdgeInPath.Label + "] " + lastEdgeInPath.To.Label)
	// }

	// for each outgoing edge from the last node in path, if it's allowed to continue on that edge from
	// the current path, add the next node to the current path and recurse
	for _, edge := range lastNodeInPath.AdjacentEdges {
		if
		// don't follow a path in the form of sync -> callback
		(lastEdgeInPath != nil && lastEdgeInPath.Type == Sync && edge.Type == Callback) ||
			// don't follow a path from async -> callback that is not your own
			(lastEdgeInPath != nil && lastEdgeInPath.Type == Async &&
				edge.Type == Callback && lastEdgeInPath.Callback != edge.To.Call.FunctionName) {
			continue
		}

		edge.To.GasLimit = edge.GasLimit
		edge.To.GasRemaining = edge.GasLimit
		edge.To.GasLocked = edge.GasLocked
		if lastEdgeInPath != nil && lastEdgeInPath.Type == Async && edge.Type == Callback {
			edge.To.GasUsed = lastEdgeInPath.GasUsedByCallback
		} else {
			edge.To.GasUsed = edge.GasUsed
		}
		edge.To.IncomingEdgeType = edge.Type

		// fmt.Println("add [" + edge.Label + "] " + edge.To.Label)
		newPath := addToPath(path, edge)
		graph.getPathsRecursive(newPath, addPathToResult)
	}
	// fmt.Println("end of edges for " + lastNodeInPath.Label)
}

func (path *TestCallPath) print() {
	fmt.Println()
	fmt.Print("path = ")
	for pathIdx, node := range path.nodes {
		if pathIdx > 0 {
			fmt.Print(" [" + path.edges[pathIdx-1].Label + "] ")
		}
		fmt.Print(node.Label + " ")
	}
	fmt.Println()
}

func (graph *TestCallGraph) newGraphUsingNodes() *TestCallGraph {
	graphCopy := CreateTestCallGraph()

	for _, node := range graph.Nodes {
		graphCopy.AddNodeCopy(node)
	}

	for _, nodeCopy := range graphCopy.Nodes {
		node := graph.FindNode(nodeCopy.Call.ContractAddress, nodeCopy.Call.FunctionName)
		for group, callBackNode := range node.GroupCallbacks {
			nodeCopy.GroupCallbacks[group] = graphCopy.FindNode(callBackNode.Call.ContractAddress, callBackNode.Call.FunctionName)
		}
	}

	for _, node := range graph.Nodes {
		if node.ContextCallback != nil {
			executionNode := graphCopy.FindNode(node.Call.ContractAddress, node.Call.FunctionName)
			executionNode.ContextCallback = graphCopy.FindNode(node.ContextCallback.Call.ContractAddress, node.ContextCallback.Call.FunctionName)
		}
	}

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

		// process sync edges
		for _, edge := range node.AdjacentEdges {
			if edge.Type != Sync {
				continue
			}

			originalDestination := edge.To.Call
			newDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)

			execEdge := executionGraph.addEdge(newSource, newDestination)
			execEdge.copyAttributesFrom(edge)
		}

		addFinishNode(executionGraph, newSource)

		groups := make([]string, 0)

		// add async and callback edges
		for _, edge := range node.AdjacentEdges {
			if edge.Type != Async {
				continue
			}

			crtGroup := edge.Group
			if !isGroupPresent(crtGroup, groups) {
				groups = append(groups, crtGroup)
			}

			originalDestination := edge.To.Call
			newAsyncDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)

			execEdge := executionGraph.addEdge(newSource, newAsyncDestination)
			execEdge.copyAttributesFrom(edge)

			if edge.Callback != "" {
				callbackDestination := executionGraph.FindNode(node.Call.ContractAddress, edge.Callback)
				execEdge := executionGraph.addEdge(newAsyncDestination, callbackDestination)
				execEdge.Type = Callback
				execEdge.GasUsedByCallback = edge.GasUsedByCallback
				execEdge.Label = "Callback"
				execEdge.Color = "grey"
			}
		}

		// add group callbacks calls if any
		for _, group := range groups {
			groupCallbackNode := newSource.GroupCallbacks[group]
			if groupCallbackNode != nil {
				execEdge := executionGraph.addEdge(newSource, groupCallbackNode)
				execEdge.Type = GroupCallback
				execEdge.Label = "Callback[" + group + "]"
				execEdge.Color = "gray"
				execEdge.GasUsed = groupCallbackNode.GasUsed
				execEdge.GasLocked = groupCallbackNode.GasLocked
			}
		}

		// callbacks were added by async source node processing and must be moved to the end of the node
		// after all other node activity (sync & async calls)
		nonCallBackEdges := make([]*TestCallEdge, 0)
		callBackEdges := make([]*TestCallEdge, 0)
		for _, newEdge := range newSource.AdjacentEdges {
			if newEdge.Type == Callback {
				callBackEdges = append(callBackEdges, newEdge)
			} else {
				nonCallBackEdges = append(nonCallBackEdges, newEdge)
			}
		}
		newSource.AdjacentEdges = append(nonCallBackEdges, callBackEdges...)

		// is start node add context callback
		if newSource.ContextCallback != nil {
			execEdge := executionGraph.addEdge(newSource, newSource.ContextCallback)
			execEdge.Type = ContextCallback
			execEdge.Label = "Callback\nContext"
			execEdge.Color = "gray"
			execEdge.GasUsed = newSource.ContextCallback.GasUsed
			execEdge.GasLocked = newSource.ContextCallback.GasLocked
		}

		return node
	})
	return executionGraph
}

// add a new 'finish' edge to a special end of sync execution node
func addFinishNode(graph *TestCallGraph, sourceNode *TestCallNode) {
	addFinishNodeWithEdgeFunc(graph, sourceNode, graph.addEdge)
}

func addFinishNodeAsFirstEdge(graph *TestCallGraph, sourceNode *TestCallNode) {
	addFinishNodeWithEdgeFunc(graph, sourceNode, graph.addEdgeAtStart)
}

func addFinishNodeWithEdgeFunc(graph *TestCallGraph, sourceNode *TestCallNode, addEdge func(*TestCallNode, *TestCallNode) *TestCallEdge) {
	finishNode := buildFinishNode(graph, sourceNode)
	addEdge(sourceNode, finishNode)
}

func buildFinishNode(graph *TestCallGraph, sourceNode *TestCallNode) *TestCallNode {
	finishNode := graph.AddNode("", SpecialLabel)
	finishNode.Label = SpecialLabel
	finishNode.IsEndOfSyncExecutionNode = true
	return finishNode
}

// CreateGasGraphFromExecutionGraph - creates a gas graph from an execution graph
func (graph *TestCallGraph) CreateGasGraphFromExecutionGraph() *TestCallGraph {
	return pathsTreeFromDag(graph)
}

func pathsTreeFromDag(graph *TestCallGraph) *TestCallGraph {
	newGraph := CreateTestCallGraph()

	paths := graph.getPaths()
	// fmt.Println()

	var crtNode *TestCallNode
	for _, path := range paths {
		// fmt.Println("process path")
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
				//fmt.Println(edge.Label + "==" + path.edges[pathIdx-1].Label + "\n=>" + strconv.FormatBool(edge.Label == path.edges[pathIdx-1].Label))
				if string(crtChild.Call.ContractAddress) == string(node.Call.ContractAddress) &&
					crtChild.Call.FunctionName == node.Call.FunctionName &&
					edge.Label == path.edges[pathIdx-1].Label {
					crtNode = crtChild
					continue nextNode
				}
			}
			parent := crtNode
			crtNode = newGraph.AddNodeCopy(node)

			// fmt.Println("add edge " + parent.Label + " -> " + crtNode.Label)
			pathEdge := path.edges[pathIdx-1]
			newEdge := newGraph.addEdge(parent, crtNode)
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

// ComputeRemainingGasBeforeCallbacks -
func (graph *TestCallGraph) ComputeRemainingGasBeforeCallbacks() {
	graph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		if node.IsLeaf() ||
			(!node.IsStartNode && (incomingEdge.Type == Callback || incomingEdge.Type == GroupCallback || incomingEdge.Type == ContextCallback)) {
			return node
		}
		nodeGasRemaining := int64(node.GasLimit)
		for _, edge := range node.AdjacentEdges {
			nodeGasRemaining -= int64(edge.To.GasLimit + edge.To.GasLocked)
			if nodeGasRemaining < 0 {
				panic(fmt.Sprintf("Bad test gas configuration %s incoming edge %s", node.Label, incomingEdge.Label))
			}
		}
		node.GasRemaining = uint64(nodeGasRemaining)
		return node
	})
}

// ComputeGasStepByStep -
func (graph *TestCallGraph) ComputeGasStepByStep(executeAfterEachStep func(graph *TestCallGraph, step int)) {
	step := 1
	finishedOneStepDfs := false
	for ; !finishedOneStepDfs; step++ {
		finishedOneStepDfs = graph.OneStepDfsFromNodePostOrder(func(parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
			if parent != nil {
				if node.IsLeaf() && node.Parent.IncomingEdgeType == Callback {
					callBackNode := parent
					asyncNode := callBackNode.Parent
					asyncInitiator := asyncNode.Parent

					var asyncInitiatorGasRemaining uint64
					if asyncInitiator.GasRemainingAfterCallback == 0 {
						asyncInitiatorGasRemaining = asyncInitiator.GasRemaining
					} else {
						asyncInitiatorGasRemaining = asyncInitiator.GasRemainingAfterCallback
					}
					var asyncNodeGasRemaining uint64
					if asyncNode.GasRemainingAfterCallback == 0 {
						asyncNodeGasRemaining = asyncNode.GasRemaining
					} else {
						asyncNodeGasRemaining = asyncNode.GasRemainingAfterCallback
					}
					callBackNode.GasLimit = asyncInitiatorGasRemaining + asyncNodeGasRemaining + asyncNode.GasLocked

					callBackNodeGasRemaining := int64(callBackNode.GasLimit)
					for _, edge := range callBackNode.GetEdges() {
						if edge.Type != Async {
							callBackNodeGasRemaining -= int64(edge.To.GasLimit - edge.To.GasRemaining)
						} else {
							callBackNodeGasRemaining -= int64(edge.To.GasLimit + edge.To.GasLocked)
						}
						if callBackNodeGasRemaining < 0 {
							panic(fmt.Sprintf("Bad test callback gas configuration %s (%s)", node.Label, callBackNode.Label))
						}
						callBackNode.GasRemaining = uint64(callBackNodeGasRemaining)
					}
				} else if !node.IsLeaf() && node.IncomingEdgeType == Callback {
					var callBackNodeGasRemaining uint64
					if node.GasRemainingAfterCallback == 0 {
						callBackNodeGasRemaining = node.GasRemaining
					} else {
						callBackNodeGasRemaining = node.GasRemainingAfterCallback
					}
					node.Parent.Parent.GasRemainingAfterCallback = callBackNodeGasRemaining
				} else if !node.IsLeaf() && node.IncomingEdgeType == Sync {
					if node.GasRemainingAfterCallback == 0 {
						parent.GasRemaining += node.GasRemaining
					} else {
						parent.GasRemaining += node.GasRemainingAfterCallback
					}
				}
			}
			return node
		})
		executeAfterEachStep(graph, step)
	}
	graph.clearVisitedNodesFlag()
}

package testcommon

import "fmt"

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
	Call          *TestCall
	AdjacentEdges []*TestCallEdge
	// group callbacks
	GroupCallbacks map[string]*TestCallNode
	// context callback
	ContextCallback *TestCallNode
	// used by execution graphs - by default these nodes are ignored by FindNode() calls
	IsEndOfSyncExecutionNode bool
	// will be reseted after each dfs traversal
	// will be ignored if stopAtVisited flag is set to false (for execution graph traversal)
	Visited     bool
	IsStartNode bool
	// used only for visualization
	Label string
}

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
	}
}

// TestCallEdgeType the type of TestCallEdges
type TestCallEdgeType int

// types of TestCallEdges
const (
	Sync = iota
	Async
	Callback
)

// TestCallEdge an edge between two nodes of the call graph
type TestCallEdge struct {
	Type           TestCallEdgeType
	CallbackToCall string
	Group          string
	To             *TestCallNode
	// used only for visualization
	Label string
	Color string
}

func (edge *TestCallEdge) copy() *TestCallEdge {
	return &TestCallEdge{
		Type:           edge.Type,
		CallbackToCall: edge.CallbackToCall,
		Group:          edge.Group,
		To:             edge.To,
		Label:          edge.Label,
		Color:          edge.Color,
	}
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
func (graph *TestCallGraph) AddStartNode(contractID string, functionName string) *TestCallNode {
	node := graph.AddNode(contractID, functionName)
	graph.StartNode = node
	node.IsStartNode = true
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
		Label:                    contractID + "_" + functionName,
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
	from.AdjacentEdges = append(from.AdjacentEdges, edge)
	return edge
}

func (graph *TestCallGraph) addEdgeAtStart(from *TestCallNode, to *TestCallNode) *TestCallEdge {
	edge := buildEdge(to)
	from.AdjacentEdges = append([]*TestCallEdge{edge}, from.AdjacentEdges...)
	return edge
}

func buildEdge(to *TestCallNode) *TestCallEdge {
	edge := &TestCallEdge{
		Type:           Sync,
		CallbackToCall: "",
		Group:          "",
		To:             to,
	}
	return edge
}

// AddAsyncEdge adds an async call edge between two nodes of the call graph
func (graph *TestCallGraph) AddAsyncEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) {
	edge := &TestCallEdge{
		Type:           Async,
		CallbackToCall: callBack,
		Group:          group,
		To:             to,
	}
	edge.Label = "Async"
	edge.Color = "red"
	if group != "" {
		edge.Label += "_" + group
	}
	if callBack != "" {
		edge.Label += "_" + callBack
	}
	from.AdjacentEdges = append(from.AdjacentEdges, edge)
}

// GetStartNode - start node getter
func (graph *TestCallGraph) GetStartNode() *TestCallNode {
	return graph.StartNode
}

// SetGroupCallback sets the callback for the specified group id
func (graph *TestCallGraph) SetGroupCallback(node *TestCallNode, groupID string, groupCallbackNode *TestCallNode) {
	node.GroupCallbacks[groupID] = groupCallbackNode
}

// SetContextCallback sets the callback for the async context
func (graph *TestCallGraph) SetContextCallback(node *TestCallNode, contextCallbackNode *TestCallNode) {
	node.ContextCallback = contextCallbackNode
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
func (graph *TestCallGraph) DfsGraph(processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) {
	for _, node := range graph.Nodes {
		if node.Visited {
			continue
		}
		graph.dfsFromNode(nil, node, make([]*TestCallNode, 0), processNode)
	}
	graph.clearVisitedNodesFlag()
}

// DfsGraphFromNode standard DFS starting from a node
// stopAtVisited set to false, enables traversal of callexecution graphs (with the risk of infinite cycles)
func (graph *TestCallGraph) DfsGraphFromNode(startNode *TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) {
	graph.dfsFromNode(nil, startNode, make([]*TestCallNode, 0), processNode)
	graph.clearVisitedNodesFlag()
}

func (graph *TestCallGraph) clearVisitedNodesFlag() {
	for _, node := range graph.Nodes {
		node.Visited = false
	}
}

func (graph *TestCallGraph) dfsFromNode(parent *TestCallNode, node *TestCallNode, path []*TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) *TestCallNode {
	if node.Visited {
		return node
	}

	path = append(path, node)
	processedParent := processNode(path, parent, node)
	node.Visited = true

	for _, edge := range node.AdjacentEdges {
		graph.dfsFromNode(processedParent, edge.To, path, processNode)
	}
	return processedParent
}

// DfsGraphFromNodePostOrder - standard post order DFS
func (graph *TestCallGraph) DfsGraphFromNodePostOrder(startNode *TestCallNode, processNode func(*TestCallNode, *TestCallNode) *TestCallNode) {
	graph.dfsFromNodePostOrder(nil, startNode, processNode)
	graph.clearVisitedNodesFlag()
}

func (graph *TestCallGraph) dfsFromNodePostOrder(parent *TestCallNode, node *TestCallNode, processNode func(*TestCallNode, *TestCallNode) *TestCallNode) *TestCallNode {
	for _, edge := range node.AdjacentEdges {
		graph.dfsFromNodePostOrder(node, edge.To, processNode)
	}

	if node.Visited {
		return node
	}

	processedParent := processNode(parent, node)
	node.Visited = true

	return processedParent
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
		dest[idxNode] = node
	}
	return dest
}

func copyEdgeList(source []*TestCallEdge) []*TestCallEdge {
	dest := make([]*TestCallEdge, len(source))
	for idxEdge, edge := range source {
		dest[idxEdge] = edge
	}
	return dest
}

func (graph *TestCallGraph) getPaths() []*TestCallPath {
	path := &TestCallPath{
		nodes: []*TestCallNode{graph.GetStartNode()},
		edges: make([]*TestCallEdge, 0),
	}
	paths := make([]*TestCallPath, 0)
	graph.getPathsInternal(path, func(newPath *TestCallPath) {
		paths = append(paths, newPath.copy())
	})
	return paths
}

func (graph *TestCallGraph) getPathsInternal(path *TestCallPath, addPathToResult func(*TestCallPath)) {
	lastNodeInPath := path.nodes[len(path.nodes)-1]

	if lastNodeInPath.IsLeaf() {
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

	for _, edge := range lastNodeInPath.AdjacentEdges {
		if
		// don't follow a path in the form of sync -> callback
		(lastEdgeInPath != nil && lastEdgeInPath.Type == Sync && edge.Type == Callback) ||
			// don't follow a path from async -> callback that is not your own
			(lastEdgeInPath != nil && lastEdgeInPath.Type == Async &&
				edge.Type == Callback && lastEdgeInPath.CallbackToCall != edge.To.Call.FunctionName) {
			continue
		}
		// fmt.Println("add [" + edge.Label + "] " + edge.To.Label)
		newPath := addToPath(path, edge)
		graph.getPathsInternal(newPath, addPathToResult)
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
	executionGraph := CreateTestCallGraph()

	for _, node := range graph.Nodes {
		executionGraph.AddNodeCopy(node)
	}

	for _, executionNode := range executionGraph.Nodes {
		node := graph.FindNode(executionNode.Call.ContractAddress, executionNode.Call.FunctionName)
		for group, callBackNode := range node.GroupCallbacks {
			executionNode.GroupCallbacks[group] = executionGraph.FindNode(callBackNode.Call.ContractAddress, callBackNode.Call.FunctionName)
		}
	}

	for _, node := range graph.Nodes {
		if node.ContextCallback != nil {
			executionNode := executionGraph.FindNode(node.Call.ContractAddress, node.Call.FunctionName)
			executionNode.ContextCallback = executionGraph.FindNode(node.ContextCallback.Call.ContractAddress, node.ContextCallback.Call.FunctionName)
		}
	}

	return executionGraph
}

// CreateExecutionGraphFromCallGraph - creates an execution graph from the call graph
func (graph *TestCallGraph) CreateExecutionGraphFromCallGraph() *TestCallGraph {

	executionGraph := graph.newGraphUsingNodes()
	executionGraph.StartNode = executionGraph.FindNode(
		graph.StartNode.Call.ContractAddress,
		graph.StartNode.Call.FunctionName)

	graph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {

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
			execEdge.Type = Sync
			execEdge.Color = "blue"
			execEdge.Label = "Sync"
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

			// for execution tree, this will be a regular edge
			execEdge := executionGraph.addEdge(newSource, newAsyncDestination)
			execEdge.Type = Async
			execEdge.CallbackToCall = edge.CallbackToCall
			execEdge.Label = "Async"
			execEdge.Color = "red"

			if edge.Group != "" {
				execEdge.Label += "_" + edge.Group
			}
			if edge.CallbackToCall != "" {
				execEdge.Label += "_" + edge.CallbackToCall
				callbackDestination := executionGraph.FindNode(node.Call.ContractAddress, edge.CallbackToCall)
				execEdge := executionGraph.addEdge(newAsyncDestination, callbackDestination)
				execEdge.Type = Callback
				execEdge.Label = "Callback"
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

		// add group callbacks calls if any
		for _, group := range groups {
			groupCallbackNode := newSource.GroupCallbacks[group]
			if groupCallbackNode != nil {
				execEdge := executionGraph.addEdge(newSource, groupCallbackNode)
				execEdge.Label = "Group_Callback_" + group
			}
		}

		// is start node add context callback
		if newSource.ContextCallback != nil {
			execEdge := executionGraph.addEdge(newSource, newSource.ContextCallback)
			execEdge.Label = "Context_Callback"
		}

		return node
	})
	return executionGraph
}

// // add a new 'finish' edge to a special end of sync execution node
func addFinishNode(graph *TestCallGraph, sourceNode *TestCallNode) {
	finishNode := buildFinishNode(graph)
	graph.addEdge(sourceNode, finishNode)
}

func buildFinishNode(graph *TestCallGraph) *TestCallNode {
	finishNode := graph.AddNode("", "X")
	finishNode.Label = "X"
	finishNode.IsEndOfSyncExecutionNode = true
	return finishNode
}

func addFinishNodeAsFirstEdge(graph *TestCallGraph, sourceNode *TestCallNode) {
	finishNode := buildFinishNode(graph)
	graph.addEdgeAtStart(sourceNode, finishNode)
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
			newEdge := newGraph.addEdge(parent, crtNode)
			pathEdge := path.edges[pathIdx-1]
			newEdge.Label = pathEdge.Label
			newEdge.Color = pathEdge.Color
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

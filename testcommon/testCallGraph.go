package testcommon

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
	Call               *TestCall
	OriginalContractID string
	adjacentEdges      []*TestCallEdge
	// group callbacks
	groupCallbacks map[string]*TestCallNode
	// context callback
	contextCallback *TestCallNode
	// used by execution graphs - by default these nodes are ignored by FindNode() calls
	isEndOfSyncExecutionNode bool
	// will be reseted after each dfs traversal
	// will be ignored if stopAtVisited flag is set to false (for execution graph traversal)
	visited     bool
	IsStartNode bool
}

// GetEdges gets the outgoing edges of the node
func (node *TestCallNode) GetEdges() []*TestCallEdge {
	return node.adjacentEdges
}

// HasAdjacentNodes returns true if the node as any adjacent nodes
func (node *TestCallNode) HasAdjacentNodes() bool {
	return len(node.adjacentEdges) != 0
}

// TestCallEdge an edge between two nodes of the call graph
type TestCallEdge struct {
	Async    bool
	callBack string
	group    string
	To       *TestCallNode
	// optional, used only for visualization
	Label string
	Color string
}

// TestCallGraph is the call graph
type TestCallGraph struct {
	Nodes     []*TestCallNode
	startNode *TestCallNode
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
	graph.startNode = node
	node.IsStartNode = true
	return node
}

// AddNode adds a node to the call graph
func (graph *TestCallGraph) AddNode(contractID string, functionName string) *TestCallNode {
	testCall := buildTestCall(contractID, functionName)
	testNode := &TestCallNode{
		Call:                     testCall,
		OriginalContractID:       contractID,
		adjacentEdges:            make([]*TestCallEdge, 0),
		groupCallbacks:           make(map[string]*TestCallNode, 0),
		contextCallback:          nil,
		visited:                  false,
		isEndOfSyncExecutionNode: false,
		IsStartNode:              false,
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
	edge := &TestCallEdge{
		Async:    false,
		callBack: "",
		group:    "",
		To:       to,
	}
	from.adjacentEdges = append(from.adjacentEdges, edge)
	return edge
}

// AddAsyncEdge adds an async call edge between two nodes of the call graph
func (graph *TestCallGraph) AddAsyncEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) {
	edge := &TestCallEdge{
		Async:    true,
		callBack: callBack,
		group:    group,
		To:       to,
	}
	edge.Label = "Async"
	if group != "" {
		edge.Label += "_" + group
		edge.Color = "red"
	}
	if callBack != "" {
		edge.Label += "_" + callBack
	}
	from.adjacentEdges = append(from.adjacentEdges, edge)
}

// GetStartNode - start node getter
func (graph *TestCallGraph) GetStartNode() *TestCallNode {
	return graph.startNode
}

// SetGroupCallback sets the callback for the specified group id
func (graph *TestCallGraph) SetGroupCallback(node *TestCallNode, groupID string, groupCallbackNode *TestCallNode) {
	node.groupCallbacks[groupID] = groupCallbackNode
}

// SetContextCallback sets the callback for the async context
func (graph *TestCallGraph) SetContextCallback(node *TestCallNode, contextCallbackNode *TestCallNode) {
	node.contextCallback = contextCallbackNode
}

// FindNode finds the corresponding node in the call graph
func (graph *TestCallGraph) FindNode(contractAddress []byte, functionName string) *TestCallNode {
	// in the future we can use an index / map if this proves to be a performance problem
	// but for test call graphs we are ok
	for _, node := range graph.Nodes {
		if string(node.Call.ContractAddress) == string(contractAddress) &&
			node.Call.FunctionName == functionName &&
			!node.isEndOfSyncExecutionNode {
			return node
		}
	}
	return nil
}

// DfsGraph a standard DFS traversal for the call graph
func (graph *TestCallGraph) DfsGraph(processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) {
	for _, node := range graph.Nodes {
		if node.visited {
			continue
		}
		graph.dfsFromNode(nil, node, make([]*TestCallNode, 0), processNode, true)
	}
	graph.clearVisitedNodesFlag()
}

// DfsGraphFromNode standard dfs starting from a node
// stopAtVisited set to false, enables traversal of callexecution graphs (with the risk of infinite cycles)
func (graph *TestCallGraph) DfsGraphFromNode(startNode *TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode, stopAtVisited bool) {
	graph.dfsFromNode(nil, startNode, make([]*TestCallNode, 0), processNode, stopAtVisited)
	graph.clearVisitedNodesFlag()
}

func (graph *TestCallGraph) clearVisitedNodesFlag() {
	for _, node := range graph.Nodes {
		node.visited = false
	}
}

func (graph *TestCallGraph) dfsFromNode(parent *TestCallNode, node *TestCallNode, path []*TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode, stopAtVisited bool) *TestCallNode {
	if stopAtVisited && node.visited {
		return node
	}

	path = append(path, node)
	processedParent := processNode(path, parent, node)
	node.visited = true

	for _, edge := range node.adjacentEdges {
		graph.dfsFromNode(processedParent, edge.To, path, processNode, stopAtVisited)
	}
	return processedParent
}

func (graph *TestCallGraph) newGraphUsingNodes() *TestCallGraph {
	executionGraph := CreateTestCallGraph()

	for _, node := range graph.Nodes {
		executionGraph.Nodes = append(executionGraph.Nodes, &TestCallNode{
			Call:                     node.Call.copy(),
			OriginalContractID:       node.OriginalContractID,
			adjacentEdges:            make([]*TestCallEdge, 0),
			groupCallbacks:           make(map[string]*TestCallNode, 0),
			contextCallback:          nil,
			isEndOfSyncExecutionNode: false,
			visited:                  false,
			IsStartNode:              node.IsStartNode,
		})
	}

	for _, executionNode := range executionGraph.Nodes {
		node := graph.FindNode(executionNode.Call.ContractAddress, executionNode.Call.FunctionName)
		for group, callBackNode := range node.groupCallbacks {
			executionNode.groupCallbacks[group] = executionGraph.FindNode(callBackNode.Call.ContractAddress, callBackNode.Call.FunctionName)
		}
	}

	for _, node := range graph.Nodes {
		if node.contextCallback != nil {
			executionNode := executionGraph.FindNode(node.Call.ContractAddress, node.Call.FunctionName)
			executionNode.contextCallback = executionGraph.FindNode(node.contextCallback.Call.ContractAddress, node.contextCallback.Call.FunctionName)
		}
	}

	return executionGraph
}

// CreateExecutionGraphFromCallGraph - creates an execution graph from the call graph
func (graph *TestCallGraph) CreateExecutionGraphFromCallGraph() *TestCallGraph {
	executionGraph := graph.newGraphUsingNodes()

	executionGraph.startNode = executionGraph.FindNode(
		graph.startNode.Call.ContractAddress,
		graph.startNode.Call.FunctionName)

	graph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		if !node.HasAdjacentNodes() {
			return node
		}

		newSource := executionGraph.FindNode(node.Call.ContractAddress, node.Call.FunctionName)

		// process sync edges
		for _, edge := range node.adjacentEdges {
			if edge.Async {
				continue
			}
			originalDestination := edge.To.Call
			newDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)
			execEdge := executionGraph.addEdge(newSource, newDestination)
			execEdge.Color = "blue"
			execEdge.Label = "Sync"
		}

		// add a new 'finish' edge to a special end of sync execution node
		finishNode := executionGraph.AddNode("", newSource.Call.FunctionName)
		finishNode.Call.ContractAddress = newSource.Call.ContractAddress
		finishNode.isEndOfSyncExecutionNode = true
		executionGraph.addEdge(newSource, finishNode)

		groups := make([]string, 0)

		// add async and callback edges
		for _, edge := range node.adjacentEdges {
			if !edge.Async {
				continue
			}

			crtGroup := edge.group
			if !isGroupPresent(crtGroup, groups) {
				groups = append(groups, crtGroup)
			}

			originalDestination := edge.To.Call
			newDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)

			// for execution tree, this will be a regular edge
			execEdge := executionGraph.addEdge(newSource, newDestination)
			execEdge.Label = "Async"
			execEdge.Color = "red"
			if edge.group != "" {
				execEdge.Label += "_" + edge.group
			}
			if edge.callBack != "" {
				execEdge.Label += "_" + edge.callBack
			}

			if edge.callBack != "" {
				callbackDestination := executionGraph.FindNode(node.Call.ContractAddress, edge.callBack)
				execEdge := executionGraph.addEdge(newSource, callbackDestination)
				execEdge.Label = "Callback" + "_" + edge.To.OriginalContractID + "_" + edge.To.Call.FunctionName
			}
		}

		// add group callbacks calls if any
		for _, group := range groups {
			groupCallbackNode := newSource.groupCallbacks[group]
			if groupCallbackNode != nil {
				execEdge := executionGraph.addEdge(newSource, groupCallbackNode)
				execEdge.Label = "GroupCallBack" + "_" + group
			}
		}

		// is start node add context callback
		if newSource.contextCallback != nil {
			execEdge := executionGraph.addEdge(newSource, newSource.contextCallback)
			execEdge.Label = "ContextCallBack"
		}

		return node
	})
	return executionGraph
}

func isGroupPresent(group string, groups []string) bool {
	for _, crtGroup := range groups {
		if group == crtGroup {
			return true
		}
	}
	return false
}

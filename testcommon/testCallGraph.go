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
	call          *TestCall
	adjacentEdges []*TestCallEdge
	// used by execution graphs - by default these nodes are ignored by FindNode() calls
	isEndOfSyncExecutionNode bool
	// will be reseted after each dfs traversal
	// will be ignored if stopAtVisited flag is set to false (for execution graph traversal)
	visited bool
}

// GetCall gets the payload of a node in the call graph
func (node *TestCallNode) GetCall() *TestCall {
	return node.call
}

// HasAdjacentNodes returns true if the node as any adjacent nodes
func (node *TestCallNode) HasAdjacentNodes() bool {
	return len(node.adjacentEdges) != 0
}

// TestCallEdge an edge between two nodes of the call graph
type TestCallEdge struct {
	async    bool
	callBack string
	group    string
	to       *TestCallNode
}

// TestCallGraph is the call graph
type TestCallGraph struct {
	nodes []*TestCallNode
}

// CreateTestCallGraph is the initial build metohd for the call graph
func CreateTestCallGraph() *TestCallGraph {
	return &TestCallGraph{
		nodes: make([]*TestCallNode, 0),
	}
}

// AddNode adds a node to the call graph
func (graph *TestCallGraph) AddNode(contractID string, functionName string) *TestCallNode {
	testCall := buildTestCall(contractID, functionName)
	testNode := &TestCallNode{
		call:                     testCall,
		adjacentEdges:            make([]*TestCallEdge, 0),
		visited:                  false,
		isEndOfSyncExecutionNode: false,
	}
	graph.nodes = append(graph.nodes, testNode)
	return testNode
}

// AddEdge adds a sync call edge between two nodes of the call graph
func (graph *TestCallGraph) AddEdge(from *TestCallNode, to *TestCallNode) {
	from.adjacentEdges = append(from.adjacentEdges, &TestCallEdge{
		async:    false,
		callBack: "",
		group:    "",
		to:       to,
	})
}

// AddAsyncEdge adds an async call edge between two nodes of the call graph
func (graph *TestCallGraph) AddAsyncEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) {
	from.adjacentEdges = append(from.adjacentEdges, &TestCallEdge{
		async:    true,
		callBack: callBack,
		group:    group,
		to:       to,
	})
}

// FindNode finds the corresponding node in the call graph
func (graph *TestCallGraph) FindNode(contractAddress []byte, functionName string) *TestCallNode {
	// in the future we can use an index / map if this proves to be a performance problem
	// but for test call graphs we are ok
	for _, node := range graph.nodes {
		if string(node.call.ContractAddress) == string(contractAddress) &&
			node.call.FunctionName == functionName &&
			!node.isEndOfSyncExecutionNode {
			return node
		}
	}
	return nil
}

// DfsGraph a standard DFS traversal for the call graph
func (graph *TestCallGraph) DfsGraph(processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) {
	for _, node := range graph.nodes {
		if node.visited {
			continue
		}
		graph.DfsFromNode(nil, node, make([]*TestCallNode, 0), processNode, true)
	}
	for _, node := range graph.nodes {
		node.visited = false
	}
}

// DfsFromNode standard dfs starting from a node
// stopAtVisited set to false, enables traversal of callexecution graphs (with the risk of infinite cycles)
func (graph *TestCallGraph) DfsFromNode(parent *TestCallNode, node *TestCallNode, path []*TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode, stopAtVisited bool) *TestCallNode {
	if stopAtVisited && node.visited {
		return node
	}
	node.visited = true

	path = append(path, node)
	processedParent := processNode(path, parent, node)
	node.visited = true

	for _, edge := range node.adjacentEdges {
		graph.DfsFromNode(processedParent, edge.to, path, processNode, stopAtVisited)
	}
	return processedParent
}

func (graph *TestCallGraph) newGraphUsingNodes() *TestCallGraph {
	result := CreateTestCallGraph()
	for _, node := range graph.nodes {
		result.nodes = append(result.nodes, &TestCallNode{
			call:          node.call.copy(),
			adjacentEdges: make([]*TestCallEdge, 0),
			visited:       false,
		})
	}
	return result
}

func (graph *TestCallGraph) createExecutionGraphFromCallGraph() *TestCallGraph {
	executionGraph := graph.newGraphUsingNodes()
	graph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {

		if !node.HasAdjacentNodes() {
			return node
		}

		newSource := executionGraph.FindNode(node.call.ContractAddress, node.call.FunctionName)

		// process sync edges
		for _, edge := range node.adjacentEdges {
			if edge.async {
				continue
			}
			originalDestination := edge.to.call
			newDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)
			executionGraph.AddEdge(newSource, newDestination)
		}

		// add a new 'finish' edge to a special end of sync execution node
		finishNode := executionGraph.AddNode("", newSource.call.FunctionName)
		finishNode.call.ContractAddress = newSource.call.ContractAddress
		finishNode.isEndOfSyncExecutionNode = true
		executionGraph.AddEdge(newSource, finishNode)

		// add async and callback edges
		for _, edge := range node.adjacentEdges {
			if !edge.async {
				continue
			}
			originalDestination := edge.to.call
			newDestination := executionGraph.FindNode(originalDestination.ContractAddress, originalDestination.FunctionName)
			// for execution tree, this will be a regular edge
			executionGraph.AddEdge(newSource, newDestination)

			callbackDestination := executionGraph.FindNode(node.call.ContractAddress, edge.callBack)
			executionGraph.AddEdge(newSource, callbackDestination)
		}

		return node
	})
	return executionGraph
}

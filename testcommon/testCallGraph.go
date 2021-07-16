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

func buildTestCall(contractID string, functionName string) *TestCall {
	return &TestCall{
		ContractAddress: MakeTestSCAddress(contractID),
		FunctionName:    functionName,
	}
}

// TestCallNode is a node in the call graph
type TestCallNode struct {
	call          *TestCall
	adjacentNodes []*TestCallEdge
	// needs to be reseted after each traversal!
	visited bool
}

// GetCall gets the payload of a node in the call graph
func (node *TestCallNode) GetCall() *TestCall {
	return node.call
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
		call:          testCall,
		adjacentNodes: make([]*TestCallEdge, 0),
		visited:       false,
	}
	graph.nodes = append(graph.nodes, testNode)
	return testNode
}

// AddEdge adds a sync call edge between two nodes of the call graph
func (graph *TestCallGraph) AddEdge(from *TestCallNode, to *TestCallNode) {
	from.adjacentNodes = append(from.adjacentNodes, &TestCallEdge{
		async:    false,
		callBack: "",
		group:    "",
		to:       to,
	})
}

// AddAsyncEdge adds an async call edge between two nodes of the call graph
func (graph *TestCallGraph) AddAsyncEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) {
	from.adjacentNodes = append(from.adjacentNodes, &TestCallEdge{
		async:    true,
		callBack: callBack,
		group:    group,
		to:       to,
	})
}

// FindNode finds the corresponding node in the call graph
func (graph *TestCallGraph) FindNode(contractID string, functionName string) *TestCallNode {
	for _, node := range graph.nodes {
		if string(node.call.ContractAddress) == contractID && node.call.FunctionName == functionName {
			return node
		}
	}
	return nil
}

// DfsGraph a standard DFS traversal for the call graph
func (graph *TestCallGraph) DfsGraph(processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) {
	foundUnvisitedNodes := true
	for foundUnvisitedNodes {
		foundUnvisitedNodes = false
		for _, node := range graph.nodes {
			if node.visited {
				continue
			}
			foundUnvisitedNodes = true
			graph.dfs(nil, node, make([]*TestCallNode, 0), processNode)
		}
	}
	for _, node := range graph.nodes {
		node.visited = false
	}
}

func (graph *TestCallGraph) dfs(parent *TestCallNode, node *TestCallNode, path []*TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) *TestCallNode {
	if node.visited {
		return node
	}
	node.visited = true

	path = append(path, node)
	processedParent := processNode(path, parent, node)
	node.visited = true

	for _, edge := range node.adjacentNodes {
		graph.dfs(processedParent, edge.to, path, processNode)
	}
	return processedParent
}

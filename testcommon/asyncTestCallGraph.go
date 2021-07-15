package testcommon

// TestCall -
type TestCall struct {
	ContractAddress []byte
	FunctionName    string
}

// ToString -
func (call *TestCall) ToString() string {
	return "contract=" + string(call.ContractAddress) + " function=" + call.FunctionName
}

// buildTestCall -
func buildTestCall(contractID string, functionName string) *TestCall {
	return &TestCall{
		ContractAddress: MakeTestSCAddress(contractID),
		FunctionName:    functionName,
	}
}

// TestCallNode -
type TestCallNode struct {
	asyncCall     *TestCall
	adjacentNodes []*TestCallEdge
	// needs to be reseted after each traversal!
	visited bool
}

// GetAsyncCall -
func (node *TestCallNode) GetAsyncCall() *TestCall {
	return node.asyncCall
}

// TestCallEdge -
type TestCallEdge struct {
	async    bool
	callBack string
	group    string
	to       *TestCallNode
}

// TestCallGraph -
type TestCallGraph struct {
	nodes []*TestCallNode
}

// CreateTestCallGraph -
func CreateTestCallGraph() *TestCallGraph {
	return &TestCallGraph{
		nodes: make([]*TestCallNode, 0),
	}
}

// AddNode -
func (graph *TestCallGraph) AddNode(contractID string, functionName string) *TestCallNode {
	testCall := buildTestCall(contractID, functionName)
	testNode := &TestCallNode{
		asyncCall:     testCall,
		adjacentNodes: make([]*TestCallEdge, 0),
		visited:       false,
	}
	graph.nodes = append(graph.nodes, testNode)
	return testNode
}

// AddEdge -
func (graph *TestCallGraph) AddEdge(from *TestCallNode, to *TestCallNode) {
	from.adjacentNodes = append(from.adjacentNodes, &TestCallEdge{
		async:    false,
		callBack: "",
		group:    "",
		to:       to,
	})
}

// AddAsyncEdge -
func (graph *TestCallGraph) AddAsyncEdge(from *TestCallNode, to *TestCallNode, callBack string, group string) {
	from.adjacentNodes = append(from.adjacentNodes, &TestCallEdge{
		async:    true,
		callBack: callBack,
		group:    group,
		to:       to,
	})
}

// FindNode -
func (graph *TestCallGraph) FindNode(contractID string, functionName string) *TestCallNode {
	for _, node := range graph.nodes {
		if string(node.asyncCall.ContractAddress) == contractID && node.asyncCall.FunctionName == functionName {
			return node
		}
	}
	return nil
}

// DfsGraph -
func (graph *TestCallGraph) DfsGraph(processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) {
	foundUnvisitedNodes := true
	for foundUnvisitedNodes {
		foundUnvisitedNodes = false
		for _, node := range graph.nodes {
			if node.visited {
				continue
			}
			foundUnvisitedNodes = true
			graph.Dfs(nil, node, make([]*TestCallNode, 0), processNode)
		}
	}
	for _, node := range graph.nodes {
		node.visited = false
	}
}

// Dfs -
func (graph *TestCallGraph) Dfs(parent *TestCallNode, node *TestCallNode, path []*TestCallNode, processNode func([]*TestCallNode, *TestCallNode, *TestCallNode) *TestCallNode) *TestCallNode {
	if node.visited {
		return node
	}
	node.visited = true

	path = append(path, node)
	processedParent := processNode(path, parent, node)
	node.visited = true

	for _, edge := range node.adjacentNodes {
		graph.Dfs(processedParent, edge.to, path, processNode)
	}
	return processedParent
}

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/awalterschulze/gographviz"
)

func main() {
	callGraph := test.CreateGraphTest1()
	// callGraph := test.CreateGraphTestSimple1()
	// callGraph := test.CreateGraphTestSimple2()

	///////////////////

	graphviz := toGraphviz(callGraph)
	createSvg("call-graph", graphviz)

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	graphviz = toGraphviz(executionGraph)
	createSvg("execution-graph", graphviz)

	gasGraph := executionGraph.CreateGasGraphFromExecutionGraph()
	graphviz = toGraphviz(gasGraph)
	createSvg("gas-graph", graphviz)
}

func createSvg(file string, graphviz *gographviz.Graph) {
	location := os.Args[1]

	destDot := location + file + ".dot"

	output := graphviz.String()
	err := ioutil.WriteFile(destDot, []byte(output), 0644)
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("dot", "-Tsvg", destDot).Output()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(location+file+".svg", out, 0644)
	if err != nil {
		panic(err)
	}
}

func toGraphviz(graph *test.TestCallGraph) *gographviz.Graph {
	graphviz := gographviz.NewGraph()
	graphviz.Directed = true
	graphName := "G"

	nodeCounters := make(map[string]int)
	for _, node := range graph.Nodes {
		nodeLabel := getGraphvizNodeLabel(node, nodeCounters)
		node.Label = nodeLabel
	}

	for _, node := range graph.Nodes {
		attrs := make(map[string]string)
		if node.IsStartNode {
			attrs["shape"] = "box"
		}
		if !node.IsEndOfSyncExecutionNode {
			attrs["bgcolor"] = "grey"
			attrs["style"] = "filled"
		}
		from := getGraphvizNodeLabel(node, nil)
		graphviz.AddNode(graphName, from, attrs)
		for _, edge := range node.GetEdges() {
			to := getGraphvizNodeLabel(edge.To, nil)
			attrs := make(map[string]string)
			if edge.To.IsEndOfSyncExecutionNode {
				attrs["style"] = "dotted"
			}
			if edge.Label != "" {
				attrs["label"] = edge.Label
			}
			if edge.Color != "" {
				attrs["color"] = edge.Color
			} else {
				attrs["color"] = "black"
			}
			graphviz.AddEdge(from, to, true, attrs)
		}
	}

	return graphviz
}

func getGraphvizNodeLabel(node *test.TestCallNode, nodeCounters map[string]int) string {
	if nodeCounters == nil {
		return node.Label
	}

	var prefix string
	if node.Call.FunctionName == "X" {
		prefix = "X"
	} else {
		prefix = node.Label
	}

	counter, present := nodeCounters[prefix]
	if !present {
		counter = 0
	}
	counter++
	nodeCounters[prefix] = counter

	suffix := ""
	if counter > 1 {
		suffix = "_" + strconv.Itoa(counter)
	}
	return prefix + suffix
}

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/awalterschulze/gographviz"
)

func main() {
	callGraph := test.CreateGraphTest1()
	graphviz := toGraphviz(callGraph)
	createSvg("call-graph", graphviz)

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	graphviz = toGraphviz(executionGraph)
	createSvg("execution-graph", graphviz)
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

	for _, node := range graph.Nodes {
		attrs := make(map[string]string)
		if node.IsStartNode {
			attrs["shape"] = "box"
		}
		attrs["bgcolor"] = "grey"
		attrs["style"] = "filled"
		graphviz.AddNode(graphName, getGraphvizNodeLabel(node), attrs)
		for _, edge := range node.GetEdges() {
			from := getGraphvizNodeLabel(node)
			to := getGraphvizNodeLabel(edge.To)
			if edge.To.OriginalContractID != "" {
				attrs := make(map[string]string)
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
	}

	toRemove := make([]string, 0)
	for _, gvNode := range graphviz.Nodes.Nodes {
		if strings.HasPrefix(gvNode.Name, "_") {
			toRemove = append(toRemove, gvNode.Name)
		}
	}

	for _, nodeToRemove := range toRemove {
		graphviz.RemoveNode(graphName, nodeToRemove)
	}

	return graphviz
}

func getGraphvizNodeLabel(node *test.TestCallNode) string {
	return node.OriginalContractID + "_" + node.Call.FunctionName
}

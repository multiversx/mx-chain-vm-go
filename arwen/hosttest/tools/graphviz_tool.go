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
	// callGraph := test.CreateGraphTestSimple3()

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
		node.Label, node.VisualLabel = getGraphvizNodeLabel(node, nodeCounters)
	}

	for _, node := range graph.Nodes {
		attrs := make(map[string]string)
		if node.IsStartNode {
			attrs["shape"] = "box"
		}
		setGasLabel(node, attrs)
		if !node.IsEndOfSyncExecutionNode {
			attrs["bgcolor"] = "grey"
			attrs["style"] = "filled"
		}
		from := node.Label
		attrs["label"] = node.VisualLabel
		graphviz.AddNode(graphName, from, attrs)
		for _, edge := range node.GetEdges() {
			to := edge.To.Label
			attrs := make(map[string]string)
			// if edge.To.IsEndOfSyncExecutionNode {
			// 	attrs["style"] = "dotted"
			// }
			if edge.Label != "" {
				attrs["label"] = strconv.Quote(edge.Label)
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

func getGraphvizNodeLabel(node *test.TestCallNode, nodeCounters map[string]int) (string, string) {
	if nodeCounters == nil {
		return node.Label, node.Label
	}

	var prefix string
	if node.Call.FunctionName == test.SpecialLabel {
		prefix = test.SpecialLabel
	} else {
		prefix, _ = strconv.Unquote(node.Label)
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
	return strconv.Quote(prefix + suffix), strconv.Quote(prefix)
}

func setGasLabel(node *test.TestCallNode, attrs map[string]string) {
	if node.GasLimit == 0 {
		return
	}
	gasLimit := strconv.Itoa(int(node.GasLimit))
	gasRemaining := strconv.Itoa(int(node.GasRemaining))
	var prefix string
	if node.IsEndOfSyncExecutionNode {
		prefix = "U"
	} else {
		prefix = "L"
	}
	attrs["xlabel"] = "<<font color='green'>" + prefix + gasLimit + "<br/>R" + gasRemaining + "</font>>"
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/awalterschulze/gographviz"
)

func main() {
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShard5()
	// callGraph := test.CreateGraphTestCallbackCallsAsyncLocalLocal()
	// callGraph := test.CreateGraphTestCallbackCallsAsyncLocalCross()
	callGraph := test.CreateGraphTestCallbackCallsAsyncCrossLocal()
	// callGraph := test.CreateGraphTestCallbackCallsAsyncCrossCross()
	// callGraph := test.CreateGraphTestAsyncCallsAsync()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncCrossLocal()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncLocalCross()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncCrossShard()
	// callGraph := test.CreateGraphTestTwoAsyncCalls()
	// callGraph := test.CreateGraphTestTwoAsyncCallsLocalCross()
	// callGraph := test.CreateGraphTestTwoAsyncCallsCrossLocal()
	// callGraph := test.CreateGraphTestTwoAsyncCallsCrossShard()
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShard4()

	/////////////////////////////////////////////////////////////////////////////////////////

	// callGraph := test.CreateGraphTestSyncCalls()
	// callGraph := test.CreateGraphTestOneAsyncCall()
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShard()
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShard2()
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShard3()
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShardComplex()

	// callGraph := test.CreateGraphTestTwoAsyncCalls()
	// callGraph := test.CreateGraphTestAsyncCallsAsync2() // not allowed to run!
	// callGraph := test.CreateGraphTestDifferentTypeOfCallsToSameFunction()

	// callGraph := test.CreateGraphTestCallbackCallsSync()
	// callGraph := test.CreateGraphTestSimpleSyncAndAsync1()
	// callGraph := test.CreateGraphTestSimpleSyncAndAsync2()
	// callGraph := test.CreateGraphTest1()
	// callGraph := test.CreateGraphTest2()

	///////////////////

	graphviz := toGraphviz(callGraph, true)
	createSvg("1 call-graph", graphviz)

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	graphviz = toGraphviz(executionGraph, true)
	createSvg("2 execution-graph", graphviz)

	gasGraph := executionGraph.ComputeGasGraphFromExecutionGraph()
	graphviz = toGraphviz(gasGraph, false)
	createSvg("3 tree-call-graph", graphviz)

	gasGraph.ComputeRemainingGasBeforeCallbacks()
	graphviz = toGraphviz(gasGraph, false)
	createSvg("4 gas-graph-gasbeforecallbacks", graphviz)

	gasGraph.ComputeGasStepByStep(func(graph *test.TestCallGraph, step int) {
		graphviz = toGraphviz(gasGraph, false)
		createSvg(fmt.Sprintf("step %d", step), graphviz)
	})
}

func createSvg(file string, graphviz *gographviz.Graph) {
	location := os.Args[1]

	destDot := location + file + ".dot"

	output := graphviz.String()
	err := ioutil.WriteFile(destDot, []byte(output), 0644)
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("dot" /*"-extent 800x1500",*/, "-Tsvg", destDot).Output()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(location+file+".svg", out, 0644)
	if err != nil {
		panic(err)
	}
}

func toGraphviz(graph *test.TestCallGraph, showGasEdgeLabels bool) *gographviz.Graph {
	graphviz := gographviz.NewGraph()
	graphviz.Directed = true
	graphName := "G"
	graphviz.Attrs["nodesep"] = "1.5"

	nodeCounters := make(map[string]int)
	for _, node := range graph.Nodes {
		node.Label, node.VisualLabel = computeUniqueGraphvizNodeLabel(node, nodeCounters)
	}

	for _, node := range graph.Nodes {
		nodeAttrs := make(map[string]string)
		setNodeAttributes(node, nodeAttrs)
		from := node.Label
		graphviz.AddNode(graphName, from, nodeAttrs)
		for _, edge := range node.GetEdges() {
			to := edge.To.Label
			edgeAttrs := make(map[string]string)
			if edge.Label != "" {
				setEdgeLabel(edgeAttrs, edge, showGasEdgeLabels)
			}
			setEdgeAttributes(edge, edgeAttrs)
			graphviz.AddEdge(from, to, true, edgeAttrs)
		}
	}

	return graphviz
}

func setNodeAttributes(node *test.TestCallNode, attrs map[string]string) {
	if node.IsStartNode {
		attrs["shape"] = "box"
	}
	if node.Visited {
		attrs["penwidth"] = "4"
	}
	setGasLabelForNode(node, attrs)
	if !node.IsGasLeaf() {
		attrs["bgcolor"] = "grey"
		attrs["style"] = "filled"
		attrs["label"] = node.VisualLabel
	}
}

func setEdgeLabel(attrs map[string]string, edge *test.TestCallEdge, showGasEdgeLabels bool) {
	attrs["label"] = edge.Label
	if showGasEdgeLabels && edge.Type != test.Callback && edge.Type != test.CallbackCrossShard {
		attrs["label"] += "\n" +
			"P" + strconv.Itoa(int(edge.GasLimit)) +
			"/U" + strconv.Itoa(int(edge.GasUsed))
		if edge.Type == test.Async || edge.Type == test.AsyncCrossShard {
			attrs["label"] += "/CU" + strconv.Itoa(int(edge.GasUsedByCallback))
		}
	}
	attrs["label"] = strconv.Quote(attrs["label"])
}

func setEdgeAttributes(edge *test.TestCallEdge, attrs map[string]string) {
	if edge.To.IsGasLeaf() {
		attrs["color"] = "black"
		return
	}
	switch edge.Type {
	case test.Sync:
		attrs["color"] = "blue"
	case test.Async:
		attrs["color"] = "red"
	case test.AsyncCrossShard:
		attrs["color"] = "red"
		attrs["style"] = "dashed"
	case test.Callback:
		attrs["color"] = "grey"
	case test.CallbackCrossShard:
		attrs["color"] = "grey"
		attrs["style"] = "dashed"
	default:
		attrs["color"] = "black"
	}
}

// generates unique graphviz node label using "_number" suffix if necessary
func computeUniqueGraphvizNodeLabel(node *test.TestCallNode, nodeCounters map[string]int) (string, string) {
	if nodeCounters == nil {
		return node.Label, node.Label
	}
	if node.VisualLabel != "" {
		return node.Label, node.VisualLabel
	}

	var prefix string
	if node.Call.FunctionName == test.LeafLabel {
		prefix = test.LeafLabel
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

const gasFontStart = "<<font color='green'>"
const gasFontEnd = "</font>>"

func setGasLabelForNode(node *test.TestCallNode, attrs map[string]string) {
	if node.GasLimit == 0 && node.GasUsed == 0 {
		// special label for end nodes without gas info
		if node.IsGasLeaf() {
			attrs["label"] = strconv.Quote("*")
		}
		return
	}

	gasLimit := strconv.Itoa(int(node.GasLimit))
	gasUsed := strconv.Itoa(int(node.GasUsed))
	gasRemaining := strconv.Itoa(int(node.GasRemaining))
	gasRemainingAfterCallback := strconv.Itoa(int(node.GasAccumulatedAfterCallback))
	gasLocked := strconv.Itoa(int(node.GasLocked))
	var xlabel string
	if node.IsGasLeaf() {
		attrs["label"] = gasFontStart + gasUsed + gasFontEnd
	} else {
		// display only gas locked for uncomputed gas values (for group callbacks and context callbacks)
		if node.GasLimit == 0 {
			// xlabel += gasFontStart + "L" + gasLocked + gasFontEnd
			// attrs["xlabel"] = xlabel
			return
		}
		xlabel = gasFontStart
		xlabel += "P" + gasLimit
		if node.GasLocked != 0 {
			xlabel += "/L" + gasLocked
		}

		// only for debug
		// xlabel += "/U" + strconv.Itoa(int(node.GasUsed))

		xlabel += "<br/>R" + gasRemaining
		if node.GasAccumulatedAfterCallback != 0 {
			xlabel += "<br/>A" + gasRemainingAfterCallback
		}
		xlabel += gasFontEnd
		attrs["xlabel"] = xlabel
	}
}

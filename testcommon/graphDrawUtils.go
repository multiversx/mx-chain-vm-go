package testcommon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

// GenerateSVGforGraph -
func GenerateSVGforGraph(callGraph *TestCallGraph, folder string, name string) {
	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	gasGraph := executionGraph.ComputeGasGraphFromExecutionGraph()
	gasGraph.PropagateSyncFailures()
	gasGraph.AssignExecutionRounds(nil)
	gasGraph.ComputeRemainingGasBeforeCallbacks(nil)
	gasGraph.ComputeRemainingGasAfterCallbacks()
	graphviz := ToGraphviz(gasGraph, false)
	CreateSvgWithLocation(folder, name, graphviz)
}

// CreateSvg -
func CreateSvg(file string, graphviz *gographviz.Graph) {
	location := os.Args[1]
	CreateSvgWithLocation(location, file, graphviz)
}

// CreateSvgWithLocation -
func CreateSvgWithLocation(folder string, file string, graphviz *gographviz.Graph) {
	destDot := folder + file + ".dot"

	output := graphviz.String()
	err := os.WriteFile(destDot, []byte(output), 0644)
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("dot" /*"-extent 800x1500",*/, "-Tsvg", destDot).Output()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(folder+file+".svg", out, 0644)
	if err != nil {
		panic(err)
	}

	_ = os.Remove(destDot)
}

// ToGraphviz -
func ToGraphviz(graph *TestCallGraph, showGasEdgeLabels bool) *gographviz.Graph {
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
		_ = graphviz.AddNode(graphName, from, nodeAttrs)
		for _, edge := range node.GetEdges() {
			to := edge.To.Label
			edgeAttrs := make(map[string]string)
			if edge.Label != "" {
				setEdgeLabel(edgeAttrs, edge, showGasEdgeLabels)
			}
			setEdgeAttributes(edge, edgeAttrs)
			_ = graphviz.AddEdge(from, to, true, edgeAttrs)
		}
	}

	return graphviz
}

func setNodeAttributes(node *TestCallNode, attrs map[string]string) {
	if node.IsStartNode {
		attrs["shape"] = "box"
	}
	// if node.Visited {
	// 	attrs["penwidth"] = "4"
	// }
	setGasLabelForNode(node, attrs)
	if !node.IsGasLeaf() {
		if node.Fail || node.IsIncomingEdgeFail() || node.HasFailSyncEdge() {
			attrs["fillcolor"] = "hotpink"
		} else {
			attrs["fillcolor"] = "lightgrey"
		}
		attrs["style"] = "filled"
		attrs["label"] = node.VisualLabel
	}
}

func setEdgeLabel(attrs map[string]string, edge *TestCallEdge, showGasEdgeLabels bool) {
	attrs["label"] = edge.Label
	if showGasEdgeLabels && edge.Type != Callback && edge.Type != CallbackCrossShard {
		attrs["label"] += "\n" +
			"P" + strconv.Itoa(int(edge.GasLimit)) +
			"/U" + strconv.Itoa(int(edge.GasUsed))
		if edge.Type == Async || edge.Type == AsyncCrossShard {
			attrs["label"] += "/CU" + strconv.Itoa(int(edge.GasUsedByCallback))
		}
	}
	attrs["label"] = strconv.Quote(attrs["label"])
}

func setEdgeAttributes(edge *TestCallEdge, attrs map[string]string) {
	if edge.To.IsGasLeaf() {
		attrs["color"] = "black"
		return
	}
	switch edge.Type {
	case Sync:
		attrs["color"] = "blue"
	case Async:
		attrs["color"] = "red"
	case AsyncCrossShard:
		attrs["color"] = "red"
		attrs["style"] = "dashed"
	case Callback:
		attrs["color"] = "grey"
	case CallbackCrossShard:
		attrs["color"] = "grey"
		attrs["style"] = "dashed"
	default:
		attrs["color"] = "black"
	}
}

// generates unique graphviz node label using "_number" suffix if necessary
func computeUniqueGraphvizNodeLabel(node *TestCallNode, nodeCounters map[string]int) (string, string) {
	if nodeCounters == nil {
		return node.Label, node.Label
	}
	if node.VisualLabel != "" {
		return node.Label, node.VisualLabel
	}

	var prefix string
	if node.Call.FunctionName == LeafLabel {
		prefix = LeafLabel
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

	var visualLabel string
	if node.WillNotExecute() {
		visualLabel = strconv.Quote(prefix)
	} else {
		visualLabel = strconv.Quote(fmt.Sprintf("%s [%d]", prefix, node.ExecutionRound))
	}

	return strconv.Quote(prefix + suffix), visualLabel
}

const gasFontStart = "<<font color='green'>"
const gasFontEnd = "</font>>"

func setGasLabelForNode(node *TestCallNode, attrs map[string]string) {
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
	gasAccumulated := strconv.Itoa(int(node.GasAccumulated))
	gasLocked := strconv.Itoa(int(node.GasLocked))
	var xlabel string
	if node.IsGasLeaf() {
		if node.WillNotExecute() {
			attrs["label"] = strconv.Quote(LeafLabel)
		} else {
			parent := node.Parent
			if node.IsGasLeaf() && parent != nil && parent.IsIncomingEdgeFail() {
				attrs["label"] = strconv.Quote(LeafLabel)
			} else {
				attrs["label"] = gasFontStart + gasUsed + gasFontEnd
			}
		}
	} else {
		// display only gas locked for uncomputed gas values (for group callbacks and context callbacks)
		if node.GasLimit == 0 || node.WillNotExecute() {
			return
		}
		xlabel = gasFontStart
		xlabel += "P" + gasLimit
		if node.GasLocked != 0 {
			xlabel += "/L" + gasLocked
		}

		xlabel += "<br/>R" + gasRemaining
		if node.GasAccumulated != 0 {
			xlabel += "<br/>A" + gasAccumulated
		}
		xlabel += gasFontEnd
		attrs["xlabel"] = xlabel
	}
}

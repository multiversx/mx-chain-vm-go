package main

import (
	"github.com/ElrondNetwork/wasm-vm/testcommon"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
)

func main() {
	callGraph := test.CreateGraphTestOneAsyncCallNoCallback()

	graphviz := testcommon.ToGraphviz(callGraph, true)
	testcommon.CreateSvg("1 call-graph", graphviz)

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	graphviz = testcommon.ToGraphviz(executionGraph, true)
	testcommon.CreateSvg("2 execution-graph", graphviz)

	gasGraph := executionGraph.ComputeGasGraphFromExecutionGraph()
	gasGraph.PropagateSyncFailures()
	gasGraph.AssignExecutionRounds(nil)

	graphviz = testcommon.ToGraphviz(gasGraph, false)
	testcommon.CreateSvg("3 initial-gas-graph", graphviz)

	gasGraph.ComputeRemainingGasBeforeCallbacks(nil)
	graphviz = testcommon.ToGraphviz(gasGraph, false)
	testcommon.CreateSvg("4 gas-graph-gasbeforecallbacks", graphviz)

	gasGraph.ComputeRemainingGasAfterCallbacks()
	graphviz = testcommon.ToGraphviz(gasGraph, false)
	testcommon.CreateSvg("5 gas-graph-gasaftercallbacks-norestore", graphviz)
}

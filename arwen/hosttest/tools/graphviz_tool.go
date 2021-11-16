package main

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
)

func main() {

	/*
		1 lvl of async calls
	*/
	// callGraph := test.CreateGraphTestOneAsyncCall()
	// callGraph := test.CreateGraphTestOneAsyncCallNoCallback()
	// callGraph := test.CreateGraphTestOneAsyncCallFail()
	// callGraph := test.CreateGraphTestOneAsyncCallNoCallbackFail()
	// callGraph := test.CreateGraphTestAsyncCallIndirectFail()
	// callGraph := test.CreateGraphTestOneAsyncCallbackFail()
	// callGraph := test.CreateGraphTestAsyncCallbackIndirectFail()
	// callGraph := test.CreateGraphTestAsyncCallIndirectFailCrossShard()
	// callGraph := test.CreateGraphTestOneAsyncCallFailCrossShard()
	// callGraph := test.CreateGraphTestOneAsyncCallbackFailCrossShard()
	// callGraph := test.CreateGraphTestTwoAsyncCallsSecondCallbackFailLocalCross()
	// callGraph := test.CreateGraphTestAsyncCallbackIndirectFailCrossShard()
	// callGraph := test.CreateGraphTestSyncCalls()
	// callGraph := test.CreateGraphTestSyncCalls2()
	// callGraph := test.CreateGraphTestOneAsyncCall()
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShard()
	// callGraph := test.CreateGraphTestOneAsyncCallCrossShard2() //!
	// callGraph := test.CreateGraphTestOneAsyncCallNoCallbackCrossShard()
	// callGraph := test.CreateGraphTestOneAsyncCallFailNoCallbackCrossShard()
	// callGraph := test.CreateGraphTestTwoAsyncCalls()
	// callGraph := test.CreateGraphTestTwoAsyncCallsOneFail()
	// callGraph := test.CreateGraphTestTwoAsyncCallsLocalCross()
	// callGraph := test.CreateGraphTestTwoAsyncCallsCrossLocal()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncSecondFail()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncLocalCross()
	// callGraph := test.CreateGraphTestCallbackCallsSync()
	// callGraph := test.CreateGraphTestSyncAndAsync1()
	// callGraph := test.CreateGraphTestSyncAndAsync2()
	// callGraph := test.CreateGraphTestSyncAndAsync3()
	// callGraph := test.CreateGraphTestSyncAndAsync6()
	// callGraph := test.CreateGraphTestSyncAndAsync7()
	// callGraph := test.CreateGraphTestSyncAndAsync8()
	// callGraph := test.CreateGraphTestTwoAsyncCallsCrossShard()
	// callGraph := test.CreateGraphTestTwoAsyncCallsFirstCallbackFailCrossShard()
	// callGraph := test.CreateGraphTestSyncCallsFailPropagation()
	// callGraph := test.CreateGraphTestTwoAsyncCallsFirstFail()
	// callGraph := test.CreateGraphTestTwoAsyncCallsFirstFailLocalCross()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncFirstNoCallbackLocalCross()
	// callGraph := test.CreateGraphTestOneAsyncCallCustomGasLocked()

	/*
		multi lvl of async calls
	*/
	// callGraph := test.CreateGraphTestAsyncCallsAsync()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncLocalCross()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncCrossShard()
	// callGraph := test.CreateGraphTestAsyncsOnMultiLevelFail1()
	// callGraph := test.CreateGraphTestCallbackCallsAsyncCrossCross()
	// callGraph := test.CreateGraphTestAsyncCallsCrossShard6()
	// callGraph := test.CreateGraphTestAsyncCallsCrossShard7()
	callGraph := test.CreateGraphTestSyncAndAsync5()
	// callGraph := test.CreateGraphTestDifferentTypeOfCallsToSameFunction()
	// callGraph := test.CreateGraphTestCallbackCallsAsyncLocalLocal()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncSecondFail()
	// callGraph := test.CreateGraphTestCallbackCallsAsyncCrossLocal()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncSecondCallbackFailCrossShard()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncBothCallbacksFailLocalCross()
	// callGraph := test.CreateGraphTestAsyncCallsAsyncSecondCallbackFailLocalCross()

	///////////////////

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

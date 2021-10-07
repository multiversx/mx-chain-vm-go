package hosttest

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/contexts"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func TestGasUsed_SyncCalls_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncCalls()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncCalls2_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncCalls2()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_OneAsyncCall_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCall()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_OneAsyncCallFail_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCallFail()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallIndirectFail_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallIndirectFail()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_OneAsyncCallbackFail_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCallbackFail()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallbackIndirectFail_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallbackIndirectFail()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_OneAsyncCallCrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCallCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_OneAsyncCallCrossShard2_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCallCrossShard2()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_OneAsyncCallFailCrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCallFailCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallIndirectFailCrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallIndirectFailCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_OneAsyncCallbackFailCrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCallbackFailCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallbackIndirectFailCrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallbackIndirectFailCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCalls_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCalls()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCallsOneFail_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCallsOneFail()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCalls_LocalCross_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCallsLocalCross()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCalls_CrossLocal_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCallsCrossLocal()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCallsCrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCallsCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsync_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsAsync()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsync_CrossLocal_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsAsyncCrossLocal()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsync_LocalCross_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsAsyncLocalCross()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsyncCrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsAsyncCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_DifferentTypeOfCallsToSameFunction_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestDifferentTypeOfCallsToSameFunction()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsSync_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestCallbackCallsSync()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_LocalLocal_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestCallbackCallsAsyncLocalLocal()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_LocalCross_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestCallbackCallsAsyncLocalCross()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_CrossLocal_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestCallbackCallsAsyncCrossLocal()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_CrossCross_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestCallbackCallsAsyncCrossCross()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync1_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncAndAsync1()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync2_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncAndAsync2()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync3_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncAndAsync3()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync4_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestSyncAndAsync4()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync5_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestSyncAndAsync5()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync6_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncAndAsync6()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync7_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncAndAsync7()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync8_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncAndAsync8()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync9_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestSyncAndAsync9()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall2_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard2()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall3_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard3()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall4_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard4()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall5_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard5()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall6_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard6()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall7_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard7()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall8_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard8()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall9_CrossShard_CallGraph(t *testing.T) {
	// t.Skip()
	callGraph := test.CreateGraphTestAsyncCallsCrossShard9()
	runGraphCallTestTemplate(t, callGraph)
}

func runGraphCallTestTemplate(t *testing.T, callGraph *test.TestCallGraph) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = callGraph.StartNode.GasLimit
	testConfig.GasLockCost = test.DefaultCallGraphLockedGas

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()

	gasGraph := executionGraph.ComputeGasGraphFromExecutionGraph()
	gasGraph.PropagateSyncFailures()
	gasGraph.AssignExecutionRounds()
	gasGraph.ComputeRemainingGasBeforeCallbacks()
	gasGraph.ComputeRemainingGasAfterCallbacks()

	startNode := gasGraph.GetStartNode()
	crossShardCallsQueue := test.NewCrossShardCallQueue()
	crossShardCallsQueue.Enqueue(test.UserAddress, startNode, vm.DirectCall, []byte{})

	computeCallIDs(gasGraph)

	// compute execution order (return data) assertions and compute gas assertions
	expectedCallFinishData := computeExpectedValues(gasGraph)
	totalGasUsed, totalGasRemaining := computeExpectedTotalGasValues(gasGraph)

	// graph gas sanity check
	require.Equal(t, int(gasGraph.StartNode.GasLimit), int(totalGasUsed+totalGasRemaining), "Expected Gas Sanity Check")

	crtTxNumber := 0

	var currentVMOutput *vmcommon.VMOutput
	// var lastErr error

	runtimeConfigsForCalls := make(map[string]*test.RuntimeConfigOfCall)
	callsFinishData := &test.CallsFinishData{
		Data: make([]*test.CallFinishDataItem, 0),
	}

	world := worldmock.NewMockWorld()

	// create contracts
	mockInstancesTestTemplate := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContractsFromAsyncTestCallGraph(callGraph, callsFinishData, runtimeConfigsForCalls, testConfig)...,
		)
	contractsInitialized := false

	var crossShardCall *test.CrossShardCall
	for !crossShardCallsQueue.IsEmpty() {
		crossShardCall = crossShardCallsQueue.Dequeue()
		startNode = crossShardCall.StartNode

		crtTxNumber++
		crtTxHash := big.NewInt(int64(crtTxNumber)).Bytes()
		crossShardCall.StartNode.CrtTxHash = crtTxHash

		crossShardEdges := getCrossShardEdgesFromSubtree(gasGraph, startNode, crossShardCallsQueue)

		arguments := [][]byte{}
		if len(crossShardCall.Data) != 0 {
			_, parsedArguments, err := parsers.NewCallArgsParser().ParseData(string(crossShardCall.Data))
			if err != nil {
				panic(err)
			}
			arguments = parsedArguments
		}

		currentVMOutput, _ /*lastErr*/ = mockInstancesTestTemplate.
			WithInput(test.CreateTestContractCallInputBuilder().
				WithCallerAddr(crossShardCall.CallerAddress).
				WithRecipientAddr([]byte(startNode.Call.ContractAddress)).
				WithFunction(startNode.Call.FunctionName).
				WithGasProvided(startNode.GasLimit).
				WithGasLocked(startNode.GasLocked).
				WithCallType(crossShardCall.CallType).
				WithArguments(arguments...).
				WithPrevTxHash(big.NewInt(int64(crtTxNumber-1)).Bytes()).
				WithCurrentTxHash(crtTxHash).
				Build()).
			WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
				world.SelfShardID = world.GetShardOfAddress(startNode.Call.ContractAddress)
				setZeroCodeCosts(host)
				setAsyncCosts(host, testConfig.GasLockCost)
			}).
			AndAssertResultsWithWorld(world, !contractsInitialized, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				// TODO matei-p adapt depending on run config
				// verify.Ok()
				// verify.ReturnCode(vmcommon.ExecutionFailed)
			})
		contractsInitialized = true

		extractAndPersistStores(t, world, currentVMOutput)

		extractOuptutTransferCalls(currentVMOutput, crossShardEdges, crossShardCallsQueue)
	}

	checkThatStoreIsEmpty(t, world)

	checkReturnDataWithGasValuesForGraphTesting(t, expectedCallFinishData, callsFinishData.Data)

	// TODO matei-p adapt depending on run config
	// test.NewVMOutputVerifier(t, currentVMOutput, lastErr).
	// 	Ok().
	// ReturnCode(vmcommon.ExecutionFailed)
	// GasRemaining(callGraph.StartNode.GasLimit - totalGasUsed)
}

func checkThatStoreIsEmpty(t testing.TB, world *worldmock.MockWorld) {
	for address, account := range world.AcctMap {
		for key, value := range account.Storage {
			require.Equal(t, []byte{}, value, fmt.Sprintf("Value present in storage for address '%s' key '%s'", address, key))
		}
	}
}

func getCrossShardEdgesFromSubtree(gasGraph *test.TestCallGraph, startNode *test.TestCallNode, crossShardCallsQueue *test.CrossShardCallsQueue) []*test.TestCallEdge {
	crossShardEdges := make([]*test.TestCallEdge, 0)
	visits := make(map[uint]bool)
	gasGraph.DfsGraphFromNode(startNode, func(path []*test.TestCallNode, parent *test.TestCallNode, node *test.TestCallNode, incomingEdge *test.TestCallEdge) *test.TestCallNode {
		for _, edge := range node.AdjacentEdges {
			if edge.Type == test.AsyncCrossShard || edge.Type == test.CallbackCrossShard {
				crossShardEdges = append(crossShardEdges, edge)
			}
		}
		return node
	}, visits, false /* don't followCrossShardEdges */)

	// if a parent context async exists, add it's callback edge also
	incomingEdgeType := startNode.GetIncomingEdgeType()
	if incomingEdgeType == test.CallbackCrossShard &&
		startNode.Parent != nil && startNode.Parent.Parent != nil {
		prevPrevNode := startNode.Parent.Parent
		if prevPrevNode.IsAsync() {
			for _, edge := range prevPrevNode.AdjacentEdges {
				if edge.Type == test.Callback || edge.Type == test.CallbackCrossShard {
					crossShardEdges = append(crossShardEdges, edge)
				}
			}
		}
	}

	return crossShardEdges
}

func executionOrderTraversal(gasGraph *test.TestCallGraph, nodeProcessing func(node *test.TestCallNode)) {
	sortedNodes := make(NodesList, 0)
	for _, node := range gasGraph.Nodes {
		if node.WillNotExecute() {
			continue
		}
		sortedNodes = append(sortedNodes, node)
	}
	sort.Stable(sortedNodes)

	for _, node := range sortedNodes {
		nodeProcessing(node)
	}
}

type NodesList []*test.TestCallNode

func (nodes NodesList) Len() int           { return len(nodes) }
func (nodes NodesList) Less(i, j int) bool { return nodes[i].ExecutionRound < nodes[j].ExecutionRound }
func (nodes NodesList) Swap(i, j int)      { nodes[i], nodes[j] = nodes[j], nodes[i] }

func computeCallIDs(gasGraph *test.TestCallGraph) {
	executionOrderTraversal(gasGraph, func(node *test.TestCallNode) {
		if node.IsLeaf() || node.Parent == nil {
			return
		}

		var parent *test.TestCallNode
		if node.GetIncomingEdgeType() == test.Callback {
			parent = node.Parent.Parent
		} else {
			parent = node.Parent
		}

		if parent != nil {
			parent.NonGasEdgeCounter++
			newCallID := append(parent.Call.CallID, big.NewInt(parent.NonGasEdgeCounter).Bytes()...)
			newCallID, _ = gasGraph.Crypto.Sha256(newCallID)
			node.Call.CallID = newCallID
		}
	})

}

func computeExpectedValues(gasGraph *test.TestCallGraph) []*test.CallFinishDataItem {
	expectedCallsFinishData := make([]*test.CallFinishDataItem, 0)

	executionOrderTraversal(gasGraph, func(node *test.TestCallNode) {
		parent := node.Parent
		if !node.IsLeaf() ||
			(parent != nil && parent.IncomingEdge != nil && parent.IncomingEdge.Fail) {
			return
		}

		expectedCallFinishData := &test.CallFinishDataItem{
			ContractAndFunction: string(parent.Call.ContractAddress) + "_" + parent.Call.FunctionName + test.TestReturnDataSuffix,
			GasProvided:         parent.GasLimit,
			GasRemaining:        parent.GasRemaining,
		}

		expectedCallsFinishData = append(expectedCallsFinishData, expectedCallFinishData)
		return
	})

	return expectedCallsFinishData
}

func computeExpectedTotalGasValues(graph *test.TestCallGraph) (uint64, uint64) {
	visits := make(map[uint]bool)
	totalGasUsed := uint64(0)
	totalGasRemaining := uint64(0)

	graph.DfsFromNodeUntilFailures(graph.StartNode.Parent, graph.StartNode, nil, make([]*test.TestCallNode, 0),
		func(path []*test.TestCallNode, parent *test.TestCallNode, node *test.TestCallNode, incomingEdge *test.TestCallEdge) *test.TestCallNode {
			if node.IsLeaf() {
				// fmt.Printf("node %s used %d\n", node.VisualLabel, node.GasUsed)
				totalGasUsed += node.GasUsed
				return node
			}

			// all gas is used for failed callss
			if node.IncomingEdge != nil && node.IncomingEdge.Fail {
				// fmt.Printf("failed %s used %d\n", node.VisualLabel, node.GasLimit)
				totalGasUsed += node.GasLimit
				return node
			}

			if parent == nil {
				totalGasRemaining += node.GasRemaining + node.GasAccumulated
			} else if node.IsCallback() {
				totalGasRemaining += node.GasAccumulated
			}

			return node
		}, visits)

	return totalGasUsed, totalGasRemaining
}

func extractOuptutTransferCalls(vmOutput *vmcommon.VMOutput, crossShardEdges []*test.TestCallEdge, crossShardCallsQueue *test.CrossShardCallsQueue) {
	for _, crossShardEdge := range crossShardEdges {
		edgeToAddress := string(crossShardEdge.To.Call.ContractAddress)
		for _, outputAccount := range vmOutput.OutputAccounts {
			transferDestinationAddress := string(outputAccount.Address)
			if edgeToAddress != transferDestinationAddress {
				continue
			}
			for _, outputTransfer := range outputAccount.OutputTransfers {
				callType := outputTransfer.CallType

				argParser := parsers.NewCallArgsParser()
				function, parsedArgs, _ := argParser.ParseData(string(outputTransfer.Data))

				callID := parsedArgs[0]

				var encodedArgs []byte
				if bytes.Equal(callID, crossShardEdge.To.Call.CallID) {
					fmt.Println(
						"Found transfer from sender", string(outputTransfer.SenderAddress),
						"to", string(outputAccount.Address),
						"gas limit", outputTransfer.GasLimit,
						"callType", callType,
						"data", contexts.DebugCallIDAsString(outputTransfer.Data))
					if callType == vm.AsynchronousCall {
						encodedArgs = outputTransfer.Data
					} else if callType == vm.AsynchronousCallBack {
						callData := txDataBuilder.NewBuilder()
						callData.Func(function)
						for _, arg := range parsedArgs {
							callData.Bytes(arg)
						}
						encodedArgs = callData.ToBytes()
					}
					crossShardCallsQueue.Enqueue(outputTransfer.SenderAddress, crossShardEdge.To, callType, encodedArgs)
				}
			}
		}
	}
}

func extractAndPersistStores(t testing.TB, world *worldmock.MockWorld, vmOutput *vmcommon.VMOutput) {
	// check if accounts with storage from OutputAccounts have the same shardID as world mock
	for _, outputAccount := range vmOutput.OutputAccounts {
		if len(outputAccount.StorageUpdates) != 0 {
			require.Equal(t, world.SelfShardID, world.GetShardOfAddress(outputAccount.Address), fmt.Sprintf("Incorrect shard for account with address '%s'", string(outputAccount.Address)))
		}
	}

	world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func extractGasUsedPerContract(vmOutput *vmcommon.VMOutput, gasUsedPerContract map[string]uint64) {
	for _, outputAccount := range vmOutput.OutputAccounts {
		if _, ok := gasUsedPerContract[string(outputAccount.Address)]; !ok {
			gasUsedPerContract[string(outputAccount.Address)] = 0
		}
		gasUsedPerContract[string(outputAccount.Address)] += outputAccount.GasUsed
	}
}

func checkReturnDataWithGasValuesForGraphTesting(t testing.TB, expectedCallsFinishData []*test.CallFinishDataItem, callsFinishData []*test.CallFinishDataItem) {
	require.Equal(t, len(expectedCallsFinishData), len(callsFinishData), "CallFinishData length")
	for idx := range expectedCallsFinishData {
		expectedCallFinishData := expectedCallsFinishData[idx]
		actualCallFinishData := callsFinishData[idx]
		require.Equal(t, expectedCallFinishData.ContractAndFunction, actualCallFinishData.ContractAndFunction, "CallFinishData - Call")
		require.Equal(t, int(expectedCallFinishData.GasProvided), int(actualCallFinishData.GasProvided), fmt.Sprintf("CallFinishData - Gas Limit for '%s'", actualCallFinishData.ContractAndFunction))
		require.Equal(t, int(expectedCallFinishData.GasRemaining), int(actualCallFinishData.GasRemaining), fmt.Sprintf("CallFinishData - Gas Remaining for '%s'", actualCallFinishData.ContractAndFunction))
	}
}

func CheckUsedGasPerContract(t testing.TB, expectedGasUsagePerContract map[string]uint64, gasUsedPerContract map[string]uint64) {
	for expectedContract, expectedGas := range expectedGasUsagePerContract {
		require.Equal(t, int(expectedGas), int(gasUsedPerContract[expectedContract]), fmt.Sprintf("Used gas for contract %s", expectedContract))
	}
}

package hosttest

import (
	"bytes"
	"fmt"
	"math/big"
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

// TODO matei-p add error test cases

func TestGasUsed_SyncCalls_CallGraph(t *testing.T) {
	// arwen.SetLoggingForTests()
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

func TestGasUsed_OneAsyncCallCrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCallCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCalls_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCalls()
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
	callGraph := test.CreateGraphTestAsyncCallsAsync()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsync_CrossLocal_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsAsyncCrossLocal()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsync_LocalCross_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsAsyncLocalCross()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsyncCrossShard_CallGraph(t *testing.T) {
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
	callGraph := test.CreateGraphTestCallbackCallsAsyncLocalLocal()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_LocalCross_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestCallbackCallsAsyncLocalCross()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_CrossLocal_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestCallbackCallsAsyncCrossLocal()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_CrossCross_CallGraph(t *testing.T) {
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
	callGraph := test.CreateGraphTestSyncAndAsync4()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SyncAndAsync5_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSyncAndAsync5()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_GraphTest1_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTest1()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_GraphTest2_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTest2()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall2_CrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsCrossShard2()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall3_CrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsCrossShard3()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall4_CrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsCrossShard4()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall5_CrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsCrossShard5()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCall6_CrossShard_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsCrossShard6()
	runGraphCallTestTemplate(t, callGraph)
}

func runGraphCallTestTemplate(t *testing.T, callGraph *test.TestCallGraph) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = callGraph.StartNode.GasLimit
	testConfig.GasLockCost = test.DefaultCallGraphLockedGas

	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()

	gasGraph := executionGraph.ComputeGasGraphFromExecutionGraph()
	gasGraph.ComputeRemainingGasBeforeCallbacks()
	gasGraph.ComputeGasStepByStep(func(graph *test.TestCallGraph, step int) {})

	startNode := gasGraph.GetStartNode()
	crossShardCallsQueue := test.NewCrossShardCallQueue()
	crossShardCallsQueue.Enqueue(test.UserAddress, startNode, vm.DirectCall, []byte{})

	// compute execution order (return data) assertions and compute gas assertions
	// totalGasUsed, expectedGasUsagePerContract, expectedReturnData := computeExpectedValues(gasGraph)
	_, _, expectedReturnData := computeExpectedValues(gasGraph)
	computeCallIDs(gasGraph)

	// account -> (key -> value)
	storage := make(map[string]map[string][]byte)
	// gasUsedPerContract := make(map[string]uint64)

	globalReturnData := make([][]byte, 0)
	crtTxNumber := 0

	var currentVMOutput *vmcommon.VMOutput
	var lastErr error

	var crossShardCall *test.CrossShardCall
	for !crossShardCallsQueue.IsEmpty() {
		crossShardCall = crossShardCallsQueue.Dequeue()
		startNode = crossShardCall.StartNode

		crtTxNumber++
		crtTxHash := big.NewInt(int64(crtTxNumber)).Bytes()
		//fmt.Println("set tx hash for " + crossShardCall.StartNode.Label + " to " + fmt.Sprintf("%d", crtTxNumber))
		crossShardCall.StartNode.CrtTxHash = crtTxHash

		crossShardEdges := preprocessLocalCallSubtree(gasGraph, startNode, crossShardCallsQueue)

		arguments := [][]byte{}
		if len(crossShardCall.Data) != 0 {
			_, parsedArguments, err := parsers.NewCallArgsParser().ParseData(string(crossShardCall.Data))
			if err != nil {
				panic(err)
			}
			arguments = parsedArguments
		}

		currentVMOutput, lastErr = test.BuildMockInstanceCallTest(t).
			WithContracts(
				test.CreateMockContractsFromAsyncTestCallGraph(callGraph, testConfig)...,
			).
			WithInput(test.CreateTestContractCallInputBuilder().
				WithCallerAddr(crossShardCall.CallerAddress).
				WithRecipientAddr([]byte(startNode.Call.ContractAddress)).
				WithFunction(startNode.Call.FunctionName).
				WithGasProvided(startNode.GasLimit).
				WithGasLocked(startNode.GasLocked).
				WithCallType(crossShardCall.CallType).
				WithArguments(arguments...).
				WithPrevTxHash(big.NewInt(int64(crtTxNumber - 1)).Bytes()).
				WithCurrentTxHash(crtTxHash).
				Build()).
			WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
				world.SelfShardID = world.GetShardOfAddress(startNode.Call.ContractAddress)
				persistStorageUpdatesToWorld(storage, world)
				setZeroCodeCosts(host)
				setAsyncCosts(host, testConfig.GasLockCost)
			}).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				verify.Ok()
			})

		extractStores(currentVMOutput, storage)
		// extractGasUsedPerContract(currentVMOutput, gasUsedPerContract)
		globalReturnData = append(globalReturnData, currentVMOutput.ReturnData...)

		extractOuptutTransferCalls(currentVMOutput, crossShardEdges, crossShardCallsQueue)
	}

	CheckReturnDataWithGasValuesForGraphTesting(t, expectedReturnData, globalReturnData)
	test.NewVMOutputVerifier(t, currentVMOutput, lastErr).
		Ok() //.
	// GasRemaining(callGraph.StartNode.GasLimit - totalGasUsed)
	// CheckUsedGasPerContract(t, expectedGasUsagePerContract, gasUsedPerContract)
}

func persistStorageUpdatesToWorld(storage map[string]map[string][]byte, world *worldmock.MockWorld) {
	for address, store := range storage {
		account := world.AcctMap.GetAccount([]byte(address))
		if account == nil {
			continue
		}
		for key, value := range store {
			account.SaveKeyValue([]byte(key), value)
		}
	}
}

func preprocessLocalCallSubtree(gasGraph *test.TestCallGraph, startNode *test.TestCallNode, crossShardCallsQueue *test.CrossShardCallsQueue) []*test.TestCallEdge {
	crossShardEdges := make([]*test.TestCallEdge, 0)
	gasGraph.DfsGraphFromNode(startNode, func(path []*test.TestCallNode, parent *test.TestCallNode, node *test.TestCallNode, incomingEdge *test.TestCallEdge) *test.TestCallNode {
		for _, edge := range node.AdjacentEdges {
			if edge.Type == test.AsyncCrossShard || edge.Type == test.CallbackCrossShard {
				crossShardEdges = append(crossShardEdges, edge)
			}
		}
		return node
	}, false /* don't followCrossShardEdges */)

	// if a parent context async exists, add it's callback edge also
	if startNode.IncomingEdgeType == test.CallbackCrossShard &&
		startNode.Parent != nil && startNode.Parent.Parent != nil {
		prevPrevNode := startNode.Parent.Parent
		if prevPrevNode.IncomingEdgeType == test.Async ||
			prevPrevNode.IncomingEdgeType == test.AsyncCrossShard {
			for _, edge := range prevPrevNode.AdjacentEdges {
				if edge.Type == test.Callback || edge.Type == test.CallbackCrossShard {
					crossShardEdges = append(crossShardEdges, edge)
				}
			}
		}
	}

	return crossShardEdges
}

func executionOrderTraversal(gasGraph *test.TestCallGraph, nodeProcessing func(parent *test.TestCallNode, node *test.TestCallNode) *test.TestCallNode) {
	crossShardCallsQueue := test.NewCrossShardCallQueue()
	crossShardCallsQueue.Enqueue(test.UserAddress, gasGraph.StartNode, vm.DirectCall, []byte{})
	var crossShardCall *test.CrossShardCall
	for !crossShardCallsQueue.IsEmpty() {
		crossShardCall = crossShardCallsQueue.Dequeue()
		startNode := crossShardCall.StartNode

		if (startNode.IncomingEdgeType == test.Callback || startNode.IncomingEdgeType == test.CallbackCrossShard) &&
			!crossShardCallsQueue.CanExecuteLocalCallback(startNode) {
			crossShardCallsQueue.Requeue(crossShardCall)
			continue
		}

		gasGraph.DfsGraphFromNode(startNode, func(path []*test.TestCallNode, parent *test.TestCallNode, node *test.TestCallNode, incomingEdge *test.TestCallEdge) *test.TestCallNode {
			for _, edge := range node.AdjacentEdges {
				if edge.Type == test.AsyncCrossShard || edge.Type == test.CallbackCrossShard {
					destinationNode := edge.To
					var callType vm.CallType
					switch edge.Type {
					case test.AsyncCrossShard:
						callType = vm.AsynchronousCall
					case test.CallbackCrossShard:
						callType = vm.AsynchronousCallBack
					}
					crossShardCallsQueue.Enqueue(node.Call.ContractAddress, destinationNode, callType, nil)
				}
			}

			if incomingEdge != nil && incomingEdge.Type == test.Callback {
				if !crossShardCallsQueue.CanExecuteLocalCallback(node) {
					crossShardCallsQueue.Enqueue(node.Call.ContractAddress, incomingEdge.To, vm.AsynchronousCallBack, nil)
					// stop DFS for this branch
					return nil
				}
			}

			return nodeProcessing(parent, node)
		}, false)
	}
}

func computeCallIDs(gasGraph *test.TestCallGraph) {
	executionOrderTraversal(gasGraph, func(parent *test.TestCallNode, node *test.TestCallNode) *test.TestCallNode {
		if node.IsLeaf() || node.Parent == nil {
			return node
		}

		if node.IncomingEdgeType == test.Callback {
			parent = node.Parent.Parent
		} else {
			parent = node.Parent
		}
		// fmt.Println("compute callID for", string(node.Call.ContractAddress)+"_"+node.Call.FunctionName, "gasLimit", parent.GasLimit)

		if parent != nil {
			parent.NonGasEdgeCounter++
			newCallID := append(parent.Call.CallID, big.NewInt(parent.NonGasEdgeCounter).Bytes()...)
			newCallID, _ = gasGraph.Crypto.Sha256(newCallID)
			node.Call.CallID = newCallID
			// fmt.Println("label", node.Label, "CallID", newCallID)
		}

		return node
	})
}

func computeExpectedValues(gasGraph *test.TestCallGraph) (uint64, map[string]uint64, [][]byte) {
	totalGasUsed := uint64(0)
	expectedGasUsagePerContract := make(map[string]uint64)
	expectedReturnData := make([][]byte, 0)

	executionOrderTraversal(gasGraph, func(parent *test.TestCallNode, node *test.TestCallNode) *test.TestCallNode {
		if !node.IsLeaf() {
			return node
		}

		contractAddr := string(parent.Call.ContractAddress)
		if _, ok := expectedGasUsagePerContract[contractAddr]; !ok {
			expectedGasUsagePerContract[contractAddr] = 0
		}
		expectedGasUsagePerContract[contractAddr] += node.GasUsed
		totalGasUsed += node.GasUsed

		expectedNodeRetData := txDataBuilder.NewBuilder()
		expectedNodeRetData.Func(parent.Call.FunctionName)
		expectedNodeRetData.Str(string(parent.Call.ContractAddress) + "_" + parent.Call.FunctionName + test.TestReturnDataSuffix)
		expectedNodeRetData.Int64(int64(parent.GasLimit))
		expectedNodeRetData.Int64(int64(parent.GasRemaining))
		expectedReturnData = append(expectedReturnData, expectedNodeRetData.ToBytes())
		// fmt.Println("add expected call to ", string(parent.Call.ContractAddress)+"_"+parent.Call.FunctionName, "gasLimit", parent.GasLimit)

		return node
	})

	return totalGasUsed, expectedGasUsagePerContract, expectedReturnData
}

func extractOuptutTransferCalls(vmOutput *vmcommon.VMOutput, crossShardEdges []*test.TestCallEdge, crossShardCallsQueue *test.CrossShardCallsQueue) {
	for _, crossShardEdge := range crossShardEdges {
		edgeFromAddress := string(crossShardEdge.To.Parent.Call.ContractAddress)
		edgeToAddress := string(crossShardEdge.To.Call.ContractAddress)
		for _, outputAccount := range vmOutput.OutputAccounts {
			transferDestinationAddress := string(outputAccount.Address)
			if edgeToAddress != transferDestinationAddress {
				continue
			}
			for _, outputTransfer := range outputAccount.OutputTransfers {
				callType := outputTransfer.CallType
				transferSenderAddress := string(outputTransfer.SenderAddress)

				argParser := parsers.NewCallArgsParser()
				function, parsedArgs, _ := argParser.ParseData(string(outputTransfer.Data))

				callID := parsedArgs[0]

				var encodedArgs []byte
				if edgeFromAddress == transferSenderAddress &&
					bytes.Equal(callID, crossShardEdge.To.Call.CallID) &&
					(callType == vm.AsynchronousCall || callType == vm.AsynchronousCallBack) {
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
						// will be read and removed by arwen
						callData.Bytes(parsedArgs[0])
						callData.Bytes(parsedArgs[1])
						callData.Bytes(parsedArgs[2])
						callData.Bytes(parsedArgs[3])
						// return code
						callData.Bytes(parsedArgs[3])
						// testing framework info
						callData.Bytes(big.NewInt(int64(crossShardEdge.Type)).Bytes())
						callData.Int64(int64(crossShardEdge.GasUsedByCallback))

						encodedArgs = callData.ToBytes()
					}
					crossShardCallsQueue.Enqueue(outputTransfer.SenderAddress, crossShardEdge.To, callType, encodedArgs)
					// "delete" info so we use the output transfer only once
					outputTransfer.SenderAddress = nil
				}
			}
		}
	}
}

func extractStores(vmOutput *vmcommon.VMOutput, storage map[string]map[string][]byte) {
	for _, outputAccount := range vmOutput.OutputAccounts {
		for _, storageUpdate := range outputAccount.StorageUpdates {
			accountStorage := storage[string(outputAccount.Address)]
			if accountStorage == nil {
				accountStorage = make(map[string][]byte)
				storage[string(outputAccount.Address)] = accountStorage
			}
			storage[string(outputAccount.Address)][string(storageUpdate.Offset)] = storageUpdate.Data
		}
	}
}

func extractGasUsedPerContract(vmOutput *vmcommon.VMOutput, gasUsedPerContract map[string]uint64) {
	for _, outputAccount := range vmOutput.OutputAccounts {
		if _, ok := gasUsedPerContract[string(outputAccount.Address)]; !ok {
			gasUsedPerContract[string(outputAccount.Address)] = 0
		}
		gasUsedPerContract[string(outputAccount.Address)] += outputAccount.GasUsed
	}
}

// CheckReturnDataWithGasValuesForGraphTesting verifies if ReturnData is the same as the provided one
func CheckReturnDataWithGasValuesForGraphTesting(t testing.TB, expectedReturnData [][]byte, returnData [][]byte) {
	processedReturnData := make([][]byte, 0)
	argParser := parsers.NewCallArgsParser()

	// eliminte from the final return data the gas used for callback arguments
	// in order to be able to compare them with the provided return data
	for i := 0; i < len(returnData); i++ {
		retDataItem := returnData[i]
		if len(retDataItem) == 1 && retDataItem[0] == test.Callback {
			i++ // jump over next item
			continue
		}
		processedReturnData = append(processedReturnData, retDataItem)
	}
	require.Equal(t, len(expectedReturnData), len(processedReturnData), "ReturnData length")
	for idx := range expectedReturnData {
		_, expRetData, _ := argParser.ParseData(string(expectedReturnData[idx]))
		_, actualRetData, _ := argParser.ParseData(string(processedReturnData[idx]))
		expectedContractAndFunction := string(expRetData[0])
		actualContractAndFunction := string(actualRetData[0])
		require.Equal(t, expectedContractAndFunction, actualContractAndFunction, "ReturnData - Call")
		// TODO matei-p re-enable gas assertions
		// expectedGasLimitForCall := big.NewInt(0).SetBytes(expRetData[1])
		// actualGasLimitForCall := big.NewInt(0).SetBytes(actualRetData[1])
		// require.Equal(t, expectedGasLimitForCall, actualGasLimitForCall, fmt.Sprintf("ReturnData - Gas Limit for '%s'", expRetData[0]))
		// expectedGasRemainingForCall := big.NewInt(0).SetBytes(expRetData[2])
		// actualGasRemainingForCall := big.NewInt(0).SetBytes(actualRetData[2])
		// require.Equal(t, expectedGasRemainingForCall, actualGasRemainingForCall, fmt.Sprintf("ReturnData - Gas Remaining for '%s'", expRetData[0]))
	}
}

func CheckUsedGasPerContract(t testing.TB, expectedGasUsagePerContract map[string]uint64, gasUsedPerContract map[string]uint64) {
	for expectedContract, expectedGas := range expectedGasUsagePerContract {
		require.Equal(t, int(expectedGas), int(gasUsedPerContract[expectedContract]), fmt.Sprintf("Used gas for contract %s", expectedContract))
	}
}

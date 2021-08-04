package hosttest

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
)

func TestGasUsed_OneAsyncCall_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestOneAsyncCall()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_TwoAsyncCalls_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestTwoAsyncCalls()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_AsyncCallsAsync_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestAsyncCallsAsync()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_DifferentTypeOfCallsToSameFunction_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestDifferentTypeOfCallsToSameFunction()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_CallbackCallsAsync_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestCallbackCallsAsync()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SimpleSyncAndAsync1_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSimpleSyncAndAsync1()
	runGraphCallTestTemplate(t, callGraph)
}

func TestGasUsed_SimpleSyncAndAsync2_CallGraph(t *testing.T) {
	callGraph := test.CreateGraphTestSimpleSyncAndAsync2()
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

func TestGasUsed_OneAsyncCall_CrossShard_CallGraph(t *testing.T) {
	arwen.SetLoggingForTests()
	callGraph := test.CreateGraphTestOneAsyncCallCrossShard()
	runGraphCallTestTemplate(t, callGraph)
}

type usedGasPerContract struct {
	contractAddress []byte
	gasUsed         uint64
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
	crossShardCallsQueue.Enqueue(test.UserAddress, startNode, vm.DirectCall, [][]byte{})

	// compute execution order (return data) assertions and compute gas assertions
	// totalGasUsed, expectedGasUsagePerContract, expectedReturnData := computeExpectedValues(gasGraph)

	// account -> (key -> value)
	storage := make(map[string]map[string][]byte)

	var crossShardCall *test.CrossShardCall
	for !crossShardCallsQueue.IsEmpty() {
		crossShardCall = crossShardCallsQueue.Dequeue()
		startNode = crossShardCall.StartNode

		crossShardNodes := make([]*test.TestCallNode, 0)
		gasGraph.DfsGraphFromNode(startNode, func(path []*test.TestCallNode, parent *test.TestCallNode, node *test.TestCallNode, incomingEdge *test.TestCallEdge) *test.TestCallNode {
			for _, edge := range node.AdjacentEdges {
				if edge.Type == test.AsyncCrossShard || edge.Type == test.CallbackCrossShard {
					destinationNode := edge.To
					crossShardNodes = append(crossShardNodes, destinationNode)

					var args [][]byte
					var callType vm.CallType
					switch edge.Type {
					case test.AsyncCrossShard:
						callType = vm.AsynchronousCall
						args = [][]byte{big.NewInt(int64(test.Async)).Bytes(), big.NewInt(int64(edge.GasUsed)).Bytes(), big.NewInt(int64(edge.GasUsedByCallback)).Bytes()}
					case test.CallbackCrossShard:
						callType = vm.AsynchronousCallBack
						args = [][]byte{{0}, big.NewInt(int64(test.Callback)).Bytes(), big.NewInt(int64(edge.GasUsedByCallback)).Bytes()}
					}

					crossShardCallsQueue.Enqueue(node.Call.ContractAddress, destinationNode, callType, args)
				}
			}
			return node
		}, false)

		vmOutput := test.BuildMockInstanceCallTest(t).
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
				WithArguments(crossShardCall.Arguments...).
				Build()).
			WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
				for _, crossShardAccount := range crossShardNodes {
					world.AcctMap.DeleteAccount(crossShardAccount.Call.ContractAddress)
				}
				for address, store := range storage {
					account := world.AcctMap.GetAccount([]byte(address))
					if account == nil {
						continue
					}
					for key, value := range store {
						account.SaveKeyValue([]byte(key), value)
					}
				}
				setZeroCodeCosts(host)
				setAsyncCosts(host, testConfig.GasLockCost)
			}).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				/*verifier := */ verify.
					Ok()
				// ReturnDataForGraphTesting(expectedReturnData...)
				// GasRemaining(callGraph.StartNode.GasLimit - totalGasUsed)
				// for _, gasPerContract := range expectedGasUsagePerContract {
				// 	verifier.GasUsed(gasPerContract.contractAddress, gasPerContract.gasUsed)
				// }
			})

		extractStores(vmOutput, storage)
		// extractOuptutTransferCalls(vmOutput, crossShardCallsQueue)
	}
}

func computeExpectedValues(gasGraph *test.TestCallGraph) (uint64, map[string]*usedGasPerContract, [][]byte) {
	totalGasUsed := uint64(0)
	expectedGasUsagePerContract := make(map[string]*usedGasPerContract)
	expectedReturnData := make([][]byte, 0)
	gasGraph.DfsGraphFromNode(gasGraph.StartNode, func(path []*test.TestCallNode, parent *test.TestCallNode, node *test.TestCallNode, incomingEdge *test.TestCallEdge) *test.TestCallNode {
		if !node.IsLeaf() {
			return node
		}
		gasPerContract := expectedGasUsagePerContract[string(parent.Call.ContractAddress)]
		if gasPerContract == nil {
			gasPerContract = &usedGasPerContract{
				contractAddress: parent.Call.ContractAddress,
				gasUsed:         0,
			}
			expectedGasUsagePerContract[string(parent.Call.ContractAddress)] = gasPerContract
		}
		gasPerContract.gasUsed += node.GasUsed
		totalGasUsed += node.GasUsed

		expectedNodeRetData := txDataBuilder.NewBuilder()
		expectedNodeRetData.Func(parent.Call.FunctionName)
		expectedNodeRetData.Str(string(parent.Call.ContractAddress) + "_" + parent.Call.FunctionName + test.TestReturnDataSuffix)
		expectedNodeRetData.Int64(int64(parent.GasLimit))
		expectedNodeRetData.Int64(int64(parent.GasRemaining))
		expectedReturnData = append(expectedReturnData, expectedNodeRetData.ToBytes())

		return node
	}, true)
	return totalGasUsed, expectedGasUsagePerContract, expectedReturnData
}

func extractOuptutTransferCalls(vmOutput *vmcommon.VMOutput, crossShardCallsQueue *test.CrossShardCallsQueue) {
	for _, outputAccount := range vmOutput.OutputAccounts {
		for _, outputTransfer := range outputAccount.OutputTransfers {
			fmt.Println(outputTransfer)
			// crossShardCallsQueue.Enqueue()
		}
	}
	// TODO matei-p
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

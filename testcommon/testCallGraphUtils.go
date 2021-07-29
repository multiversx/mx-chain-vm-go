package testcommon

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/elrondapi"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

var logGenContr = logger.GetOrCreate("arwen/graph")

// TestReturnDataSuffix -
var TestReturnDataSuffix = "_returnData"

// TestCallbackPrefix -
var TestCallbackPrefix = "callback_"

// TestContextCallbackFunction -
var TestContextCallbackFunction = "contextCallback"

// CreateMockContractsFromAsyncTestCallGraph creates the contracts
// with functions that reflect the behavior specified by the call graph
func CreateMockContractsFromAsyncTestCallGraph(callGraph *TestCallGraph, testConfig *TestConfig) []MockTestSmartContract {
	contracts := make(map[string]*MockTestSmartContract)
	callGraph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		contractAddressAsString := string(node.Call.ContractAddress)
		if contracts[contractAddressAsString] == nil {
			newContract := CreateMockContract(node.Call.ContractAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(func(instanceMock *mock.InstanceMock, testConfig *TestConfig) {
					for functionName := range contracts[contractAddressAsString].tempFunctionsList {
						instanceMock.AddMockMethod(functionName, func() *mock.InstanceMock {
							host := instanceMock.Host
							instance := mock.GetMockInstance(host)
							t := instance.T

							async := host.Async()
							crtFunctionCalled := host.Runtime().Function()

							logGenContr.Trace("Executing graph node", "sc", string(host.Runtime().GetSCAddress()), "func", crtFunctionCalled)
							// fmt.Println("Executing graph node", "sc", string(host.Runtime().GetSCAddress()), "func", crtFunctionCalled)

							crtNode := callGraph.FindNode(host.Runtime().GetSCAddress(), crtFunctionCalled)

							var gasUsed int64

							if crtNode.IsStartNode {
								// for start node we get no arguments to read gas used from
								gasUsed = int64(crtNode.GasUsed)
							} else {
								arguments := host.Runtime().Arguments()
								if len(arguments) > 0 {
									edgeTypeArgIndex := 0
									gasUsedArgIndex := 1
									if host.Runtime().GetVMInput().CallType == vm.AsynchronousCallBack {
										// for callbacks, arguments[0] is the return code of the async call
										edgeTypeArgIndex = 1
										gasUsedArgIndex = 2
									}
									edgeType := big.NewInt(0).SetBytes(arguments[edgeTypeArgIndex]).Int64()
									if edgeType == Async {
										host.Output().Finish(big.NewInt(int64(Callback)).Bytes())
										host.Output().Finish(arguments[2]) // gas used by callback
									}

									gasUsed = big.NewInt(0).SetBytes(arguments[gasUsedArgIndex]).Int64()
								}
							}

							// burn gas for function
							logGenContr.Trace("Burning", gasUsed, "gas for", crtFunctionCalled)
							// fmt.Println("Burning", uint64(gasUsed), "gas for", crtFunctionCalled)
							host.Metering().UseGasBounded(uint64(gasUsed))

							if crtNode.ContextCallback != nil {
								err := async.SetContextCallback(crtNode.ContextCallback.Call.FunctionName, []byte{}, 0)
								require.Nil(t, err)
							}

							value := big.NewInt(testConfig.TransferFromParentToChild)

							for _, edge := range crtNode.AdjacentEdges {
								destFunctionName := edge.To.Call.FunctionName
								destAddress := edge.To.Call.ContractAddress
								if edge.Type == Sync {
									logGenContr.Trace("Sync call to ", string(destAddress), " func ", destFunctionName, " gas ", edge.GasLimit)
									// fmt.Println("Sync call to " + string(destAddress) + " func " + destFunctionName + " gas " + strconv.Itoa(int(edge.GasLimit)))
									elrondapi.ExecuteOnDestContextWithTypedArgs(
										host,
										int64(edge.GasLimit),
										value,
										[]byte(destFunctionName),
										destAddress,
										[][]byte{
											big.NewInt(int64(Sync)).Bytes(),
											big.NewInt(int64(edge.GasUsed)).Bytes()}) // args
								} else {
									logGenContr.Trace("Register async call", "to", string(destAddress), "func", destFunctionName, "gas", edge.GasLimit)
									// fmt.Println("Register async call", "to", string(destAddress), "func", destFunctionName, "gas", strconv.Itoa(int(edge.GasLimit)))

									callData := txDataBuilder.NewBuilder()
									callData.Func(destFunctionName)
									callData.Bytes(big.NewInt(int64(Async)).Bytes())
									callData.Int64(int64(edge.GasUsed))
									callData.Int64(int64(edge.GasUsedByCallback))

									err := async.RegisterAsyncCall(edge.Group, &arwen.AsyncCall{
										Status:          arwen.AsyncCallPending,
										Destination:     destAddress,
										Data:            callData.ToBytes(),
										ValueBytes:      value.Bytes(),
										GasLimit:        edge.GasLimit,
										SuccessCallback: edge.Callback,
										ErrorCallback:   edge.Callback,
									})
									require.Nil(t, err)
								}
							}

							for group, groupCallbackNode := range crtNode.GroupCallbacks {
								err := async.SetGroupCallback(group, groupCallbackNode.Call.FunctionName, []byte{}, 0)
								require.Nil(t, err)
							}

							host.Output().Finish([]byte(string(host.Runtime().GetSCAddress()) + "_" + crtFunctionCalled + TestReturnDataSuffix))
							logGenContr.Trace("End of ", crtFunctionCalled, " on ", string(host.Runtime().GetSCAddress()))
							// fmt.Println("End of " + crtFunctionCalled + " on " + string(host.Runtime().GetSCAddress()))

							return instance
						})
					}
				})
			contracts[contractAddressAsString] = &newContract
		}
		functionName := node.Call.FunctionName
		contract := contracts[contractAddressAsString]
		//fmt.Println("Add " + functionName + " to " + contractAddressAsString)
		addFunctionToTempList(contract, functionName, true)
		return node
	})
	contractsList := make([]MockTestSmartContract, 0)
	for _, contract := range contracts {
		contractsList = append(contractsList, *contract)
	}
	return contractsList
}

func addFunctionToTempList(contract *MockTestSmartContract, functionName string, isCallBack bool) {
	_, functionPresent := contract.tempFunctionsList[functionName]
	if !functionPresent {
		contract.tempFunctionsList[functionName] = isCallBack
	}
}

// CreateRunOrderFromExecutionGraph returns an exepected execution order starting from an execution graph
func CreateRunOrderFromExecutionGraph(executionGraph *TestCallGraph) []TestCall {
	executionOrder := make([]TestCall, 0)
	pathsTree := pathsTreeFromDag(executionGraph)
	pathsTree.DfsGraphFromNode(pathsTree.StartNode, func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		if node.IsEndOfSyncExecutionNode {
			//fmt.Println("end exec " + parent.Label)
			executionOrder = append(executionOrder, TestCall{
				ContractAddress: parent.Call.ContractAddress,
				FunctionName:    parent.Call.FunctionName,
			})
		}
		return node
	})
	return executionOrder
}

// CreateGraphTest1 -
func CreateGraphTest1() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 0, 0)

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc2f3, "cb2", "gr1")

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2)

	sc2f6 := callGraph.AddNode("sc2", "f6")
	callGraph.AddAsyncEdge(sc1f1, sc2f6, "cb4", "gr1")

	sc3f7 := callGraph.AddNode("sc3", "f7")
	callGraph.AddAsyncEdge(sc1f1, sc3f7, "cb4", "gr2")

	sc3f4 := callGraph.AddNode("sc3", "f4")
	callGraph.AddSyncEdge(sc2f3, sc3f4)

	callGraph.AddAsyncEdge(sc2f2, sc3f4, "cb3", "gr3")

	sc1cb1 := callGraph.AddNode("sc1", "cb2")
	sc4f5 := callGraph.AddNode("sc4", "f5")
	callGraph.AddSyncEdge(sc1cb1, sc4f5)

	callGraph.AddNode("sc2", "cb3")
	callGraph.AddNode("sc1", "cb4")

	sc1cbg1 := callGraph.AddNode("sc1", "cbg1")
	callGraph.SetGroupCallback(sc1f1, "gr1", sc1cbg1, 0, 0)

	ctxcb := callGraph.AddNode("sc1", "ctxcb")
	callGraph.SetContextCallback(sc1f1, ctxcb, 0, 0)
	return callGraph
}

// CreateGraphTestAsyncCallsAsync -
func CreateGraphTestAsyncCallsAsync() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(100).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetGasLocked(10)

	callGraph.AddNode("sc1", "cb1")

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc2f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(12)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsync2 -
func CreateGraphTestAsyncCallsAsync2() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(100).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetGasLocked(10)

	callGraph.AddNode("sc1", "cb1")

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc2f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(12)

	callGraph.AddNode("sc2", "cb2")

	sc2f4 := callGraph.AddNode("sc2", "f4")
	callGraph.AddAsyncEdge(sc2f3, sc2f4, "cb3", "gr3").
		SetGasLimit(10).
		SetGasUsed(5).
		SetGasUsedByCallback(2).
		SetGasLocked(3)

	callGraph.AddNode("sc2", "cb3")

	return callGraph
}

// CreateGraphTestGroupCallbacks -
func CreateGraphTestGroupCallbacks() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(100).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetGasLocked(10)

	callGraph.AddNode("sc1", "cb1")

	sc1cbg1 := callGraph.AddNode("sc1", "cbg1")
	callGraph.SetGroupCallback(sc1f1, "gr1", sc1cbg1, 20, 10)

	ctxcb := callGraph.AddNode("sc1", "ctxcb")
	callGraph.SetContextCallback(sc1f1, ctxcb, 15, 10)

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc2f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(12)

	callGraph.AddNode("sc2", "cb2")

	sc2cbg2 := callGraph.AddNode("sc2", "cbg2")
	callGraph.SetGroupCallback(sc2f2, "gr2", sc2cbg2, 9, 20)

	return callGraph
}

// CreateGraphTestOneAsyncCall -
func CreateGraphTestOneAsyncCall() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6).
		SetGasLocked(9)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestOneAsyncCallWithGroupCallback -
func CreateGraphTestOneAsyncCallWithGroupCallback() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(50).
		SetGasUsed(20).
		SetGasUsedByCallback(10).
		SetGasLocked(20)

	callGraph.AddNode("sc1", "cb1")

	sc1cbg1 := callGraph.AddNode("sc1", "cbg1")
	callGraph.SetGroupCallback(sc1f1, "gr1", sc1cbg1, 10, 30)

	ctxcb := callGraph.AddNode("sc1", "ctxcb")
	callGraph.SetContextCallback(sc1f1, ctxcb, 20, 10)

	return callGraph
}

// CreateGraphTestTwoAsyncCalls -
func CreateGraphTestTwoAsyncCalls() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(50).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetGasLocked(0)

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc2f3, "cb1", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(0)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestSimpleSyncAndAsync1 -
func CreateGraphTestSimpleSyncAndAsync1() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(100).
		SetGasUsed(7)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb1", "").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(12)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1f1, sc4f4).
		SetGasLimit(20).
		SetGasUsed(1)

	sc2cb1 := callGraph.AddNode("sc2", "cb1")
	callGraph.AddSyncEdge(sc4f4, sc2cb1).
		SetGasLimit(10).
		SetGasUsed(5)

	return callGraph
}

// CreateGraphTestDifferentTypeOfCallsToSameFunction -
func CreateGraphTestDifferentTypeOfCallsToSameFunction() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(20).
		SetGasUsed(5)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetGasLocked(10)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb2", "").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(12)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(6).
		SetGasUsed(6)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTest2 -
func CreateGraphTest2() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 200, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(20).
		SetGasUsed(5)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetGasLocked(10)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb2", "").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(12)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(6).
		SetGasUsed(6)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	sc1cb1 := callGraph.AddNode("sc1", "cb1")
	callGraph.AddSyncEdge(sc1cb1, sc4f4).
		SetGasLimit(10).
		SetGasUsed(5)

	sc1cb2 := callGraph.AddNode("sc1", "cb2")
	sc5f5 := callGraph.AddNode("sc5", "f5")

	callGraph.AddAsyncEdge(sc1cb2, sc5f5, "cb1", "").
		SetGasLimit(20).
		SetGasUsed(4).
		SetGasUsedByCallback(3).
		SetGasLocked(6)

	return callGraph
}

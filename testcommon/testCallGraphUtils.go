package testcommon

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/elrondapi"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

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
				WithMethods(func(instanceMock *mock.InstanceMock, config interface{}) {
					for functionName := range contracts[contractAddressAsString].tempFunctionsList {
						instanceMock.AddMockMethod(functionName, func() *mock.InstanceMock {
							host := instanceMock.Host
							instance := mock.GetMockInstance(host)
							t := instance.T

							async := host.Async()
							crtFunctionCalled := host.Runtime().Function()
							gasProvided := host.Runtime().GetVMInput().GasProvided
							gasForChildren := gasProvided / 4

							crtNode := callGraph.FindNode(host.Runtime().GetSCAddress(), crtFunctionCalled)
							if crtNode.ContextCallback != nil {
								err := async.SetContextCallback(crtNode.ContextCallback.Call.FunctionName, []byte{}, 0)
								require.Nil(t, err)
							}
							log.Trace("Executing graph node", "sc", string(host.Runtime().GetSCAddress()), "func", crtFunctionCalled)

							value := big.NewInt(testConfig.TransferFromParentToChild)

							for _, edge := range crtNode.AdjacentEdges {
								destFunctionName := edge.To.Call.FunctionName
								destAddress := edge.To.Call.ContractAddress
								if edge.Type == Sync {
									log.Trace("Sync call", "to", string(destAddress), "func", destFunctionName, "gas", gasForChildren)
									elrondapi.ExecuteOnDestContextWithTypedArgs(
										host,
										int64(gasForChildren),
										value,
										[]byte(destFunctionName),
										destAddress,
										make([][]byte, 0)) // args
								} else {
									log.Trace("Async call", "to", string(destAddress), "func", destFunctionName, "gas", gasForChildren)

									callData := txDataBuilder.NewBuilder()
									callData.Func(destFunctionName)

									err := async.RegisterAsyncCall(edge.Group, &arwen.AsyncCall{
										Status:          arwen.AsyncCallPending,
										Destination:     destAddress,
										Data:            callData.ToBytes(),
										ValueBytes:      value.Bytes(),
										GasLimit:        gasForChildren,
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
							log.Trace("End of function", "sc", string(host.Runtime().GetSCAddress()), "func", crtFunctionCalled)

							return instance
						})
					}
				})
			contracts[contractAddressAsString] = &newContract
		}
		functionName := node.Call.FunctionName
		contract := contracts[contractAddressAsString]
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

// CreateRunExpectationOrder returns an exepected execution order starting from an execution graph
func CreateRunExpectationOrder(executionGraph *TestCallGraph) []TestCall {
	executionOrder := make([]TestCall, 0)
	pathsTree := pathsTreeFromDag(executionGraph)
	pathsTree.DfsGraphFromNode(pathsTree.StartNode, func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode, incomingEdge *TestCallEdge) *TestCallNode {
		if node.IsEndOfSyncExecutionNode {
			logTestGraph.Trace("end execution", "label", parent.Label)
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
	callGraph.SetGroupCallback(sc1f1, "gr1", sc1cbg1)

	ctxcb := callGraph.AddNode("sc1", "ctxcb")
	callGraph.SetContextCallback(sc1f1, ctxcb)
	return callGraph
}

// CreateGraphTestSimple1 -
func CreateGraphTestSimple1() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 100, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetGasLocked(10)

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc2f3, "cb1", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3).
		SetGasLocked(12)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestSimple2 -
func CreateGraphTestSimple2() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 0, 0)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb1", "")

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1f1, sc4f4)

	sc2cb1 := callGraph.AddNode("sc2", "cb1")
	callGraph.AddSyncEdge(sc4f4, sc2cb1)
	// callGraph.AddSyncEdge(sc2cb1, sc3f3)

	return callGraph
}

// CreateGraphTestSimple3 -
func CreateGraphTestSimple3() *TestCallGraph {
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

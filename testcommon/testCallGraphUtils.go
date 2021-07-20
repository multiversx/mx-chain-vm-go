package testcommon

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen/elrondapi"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
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
	callGraph.DfsGraph(func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		contractAddressAsString := string(node.call.ContractAddress)
		if contracts[contractAddressAsString] == nil {
			newContract := CreateMockContract(node.call.ContractAddress).
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

							crtNode := callGraph.FindNode(host.Runtime().GetSCAddress(), crtFunctionCalled)
							if crtNode.contextCallback != nil {
								err := async.SetContextCallback(crtNode.contextCallback.call.FunctionName, []byte{}, 0)
								require.Nil(t, err)
							}
							fmt.Println("Executing " + crtFunctionCalled + " on " + string(host.Runtime().GetSCAddress()))

							value := big.NewInt(testConfig.TransferFromParentToChild)

							for _, edge := range crtNode.adjacentEdges {
								destFunctionName := edge.to.call.FunctionName
								destAddress := edge.to.call.ContractAddress
								if !edge.async {
									fmt.Println("Sync call to " + destFunctionName + " on " + string(destAddress))
									elrondapi.ExecuteOnDestContextWithTypedArgs(
										host,
										int64(testConfig.GasProvidedToChild),
										value,
										[]byte(destFunctionName),
										destAddress,
										make([][]byte, 0)) // args
								} else {
									fmt.Println("Async call to " + destFunctionName + " on " + string(destAddress))
									callData := txDataBuilder.NewBuilder()
									callData.Func(destFunctionName)

									err := async.RegisterAsyncCall(edge.group, &arwen.AsyncCall{
										Status:          arwen.AsyncCallPending,
										Destination:     destAddress,
										Data:            callData.ToBytes(),
										ValueBytes:      value.Bytes(),
										GasLimit:        testConfig.GasProvidedToChild,
										SuccessCallback: edge.callBack,
										ErrorCallback:   edge.callBack,
									})
									require.Nil(t, err)
								}
							}

							for group, groupCallbackNode := range crtNode.groupCallbacks {
								err := async.SetGroupCallback(group, groupCallbackNode.call.FunctionName, []byte{}, 0)
								require.Nil(t, err)
							}

							host.Output().Finish([]byte(string(host.Runtime().GetSCAddress()) + "_" + crtFunctionCalled + TestReturnDataSuffix))
							fmt.Println("End of " + crtFunctionCalled + " on " + string(host.Runtime().GetSCAddress()))

							return instance
						})
					}
				})
			contracts[contractAddressAsString] = &newContract
		}
		functionName := node.call.FunctionName
		contract := contracts[contractAddressAsString]
		fmt.Println("Add " + functionName + " to " + contractAddressAsString)
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
	executionGraph.DfsGraphFromNode(executionGraph.startNode, func(path []*TestCallNode, parent *TestCallNode, node *TestCallNode) *TestCallNode {
		if !node.HasAdjacentNodes() {
			fmt.Println("leaf " + node.GetCall().FunctionName)
			executionOrder = append(executionOrder, TestCall{
				ContractAddress: node.call.ContractAddress,
				FunctionName:    node.call.FunctionName,
			})
		}
		return node
	}, false)
	return executionOrder
}

// CreateGraphTest1 -
func CreateGraphTest1() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1")

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc2f3, "cb2", "gr1")

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddEdge(sc1f1, sc2f2)

	sc2f6 := callGraph.AddNode("sc2", "f6")
	callGraph.AddAsyncEdge(sc1f1, sc2f6, "cb4", "gr1")

	sc3f7 := callGraph.AddNode("sc3", "f7")
	callGraph.AddAsyncEdge(sc1f1, sc3f7, "cb4", "gr2")

	sc3f4 := callGraph.AddNode("sc3", "f4")
	callGraph.AddEdge(sc2f3, sc3f4)

	callGraph.AddAsyncEdge(sc2f2, sc3f4, "cb3", "gr3")

	sc1cb1 := callGraph.AddNode("sc1", "cb2")
	sc4f5 := callGraph.AddNode("sc4", "f5")
	callGraph.AddEdge(sc1cb1, sc4f5)

	sc2cb3 := callGraph.AddNode("sc2", "cb3")
	callGraph.AddEdge(sc2cb3, sc3f4)

	callGraph.AddNode("sc1", "cb4")

	sc1cbg1 := callGraph.AddNode("sc1", "cbg1")
	callGraph.SetGroupCallback(sc1f1, "gr1", sc1cbg1)

	ctxcb := callGraph.AddNode("sc1", "ctxcb")
	callGraph.SetContextCallback(sc1f1, ctxcb)
	return callGraph
}

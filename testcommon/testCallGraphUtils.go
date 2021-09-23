package testcommon

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/elrondapi"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

// LogGraph -
var LogGraph = logger.GetOrCreate("arwen/graph")

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
			var shardID uint32
			if node.ShardID == 0 {
				if incomingEdge == nil {
					shardID = 1
				} else if incomingEdge.Type == Sync || incomingEdge.Type == Async {
					shardID = parent.ShardID
				} else if incomingEdge.Type == AsyncCrossShard {
					shardID = contracts[string(parent.Call.ContractAddress)].shardID + 1
				}
				node.ShardID = shardID
			} else {
				shardID = node.ShardID
			}
			newContract := CreateMockContract(node.Call.ContractAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithShardID(shardID).
				WithMethods(func(instanceMock *mock.InstanceMock, testConfig *TestConfig) {
					for functionName := range contracts[contractAddressAsString].tempFunctionsList {
						instanceMock.AddMockMethod(functionName, func() *mock.InstanceMock {
							host := instanceMock.Host
							instance := mock.GetMockInstance(host)
							t := instance.T

							crtFunctionCalled := host.Runtime().Function()
							LogGraph.Trace("Executing graph node", "sc", string(host.Runtime().GetSCAddress()), "func", crtFunctionCalled, "txHash", host.Runtime().GetCurrentTxHash())

							crtNode := callGraph.FindNode(host.Runtime().GetSCAddress(), crtFunctionCalled)
							gasUsed := readGasUsedFromArguments(crtNode, host)

							// burn gas for function
							fmt.Println("Burning", "gas", gasUsed, "function", crtFunctionCalled)
							// LogGraph.Trace("Burning", "gas", gasUsed, "function", crtFunctionCalled)
							host.Metering().UseGasBounded(uint64(gasUsed))

							for _, edge := range crtNode.AdjacentEdges {
								if edge.Type == Sync {
									makeSyncCallFromEdge(host, edge, testConfig)
								} else {
									err := makeAsyncCallFromEdge(host, edge, testConfig)
									require.Nil(t, err)
								}
							}

							computeReturnData(crtFunctionCalled, host)

							return instance
						})
					}
				})
			contracts[contractAddressAsString] = &newContract
		}
		functionName := node.Call.FunctionName
		contract := contracts[contractAddressAsString]
		node.ShardID = contract.shardID
		addFunctionToTempList(contract, functionName, true)
		return node
	}, true)
	contractsList := make([]MockTestSmartContract, 0)
	for _, contract := range contracts {
		contractsList = append(contractsList, *contract)
	}
	return contractsList
}

func makeSyncCallFromEdge(host arwen.VMHost, edge *TestCallEdge, testConfig *TestConfig) {
	value := big.NewInt(testConfig.TransferFromParentToChild)
	destFunctionName := edge.To.Call.FunctionName
	destAddress := edge.To.Call.ContractAddress
	arguments := [][]byte{
		big.NewInt(int64(Sync)).Bytes(),
		big.NewInt(int64(edge.GasUsed)).Bytes()}

	LogGraph.Trace("Sync call to ", string(destAddress), " func ", destFunctionName, " gas ", edge.GasLimit)
	elrondapi.ExecuteOnDestContextWithTypedArgs(
		host,
		int64(edge.GasLimit),
		value,
		[]byte(destFunctionName),
		destAddress,
		arguments)
}

func makeAsyncCallFromEdge(host arwen.VMHost, edge *TestCallEdge, testConfig *TestConfig) error {
	async := host.Async()
	destFunctionName := edge.To.Call.FunctionName
	destAddress := edge.To.Call.ContractAddress
	value := big.NewInt(testConfig.TransferFromParentToChild)

	LogGraph.Trace("Register async call", "to", string(destAddress), "func", destFunctionName, "gas", edge.GasLimit)

	callData := txDataBuilder.NewBuilder()
	callData.Func(destFunctionName)
	callData.Bytes(big.NewInt(int64(edge.Type)).Bytes())
	callData.Int64(int64(edge.GasUsed))
	callData.Int64(int64(edge.GasUsedByCallback))

	err := async.RegisterAsyncCall("", &arwen.AsyncCall{
		Status:          arwen.AsyncCallPending,
		Destination:     destAddress,
		Data:            callData.ToBytes(),
		ValueBytes:      value.Bytes(),
		GasLimit:        edge.GasLimit,
		SuccessCallback: edge.Callback,
		ErrorCallback:   edge.Callback,
	})
	return err
}

// return data is encoded using standard txDataBuilder
// format is function@nodeLabel@providedGas@remainingGas
func computeReturnData(crtFunctionCalled string, host arwen.VMHost) {
	runtime := host.Runtime()
	metering := host.Metering()
	async := host.Async()

	returnData := txDataBuilder.NewBuilder()
	returnData.Func(crtFunctionCalled)
	returnData.Str(string(runtime.GetSCAddress()) + "_" + crtFunctionCalled + TestReturnDataSuffix)
	returnData.Int64(int64(runtime.GetVMInput().GasProvided))
	returnData.Int64(int64(metering.GasLeft()))
	returnData.Bytes(async.GetCallID())
	returnData.Bytes(async.GetCallbackAsyncInitiatorCallID())
	returnData.Bool(async.IsCrossShard())
	host.Output().Finish(returnData.ToBytes())
	LogGraph.Trace("End of ", crtFunctionCalled, " on ", string(host.Runtime().GetSCAddress()))
	// TODO matei-p remove for logging
	fmt.Println(
		"Return Data -> callID", async.GetCallID(),
		"CallbackAsyncInitiatorCallID", async.GetCallbackAsyncInitiatorCallID(),
		"IsCrossShard", async.IsCrossShard(),
		"For contract ", string(runtime.GetSCAddress()), "/ "+crtFunctionCalled+"\t",
		"Gas provided", fmt.Sprintf("%d\t", runtime.GetVMInput().GasProvided),
		"Gas remaining", fmt.Sprintf("%d\t", metering.GasLeft()))
}

func readGasUsedFromArguments(crtNode *TestCallNode, host arwen.VMHost) int64 {
	var gasUsed int64
	if crtNode.IsStartNode {
		// for start node we get no arguments to read gas used from
		gasUsed = int64(crtNode.GasUsed)
	} else {
		arguments := host.Runtime().Arguments()
		callType := host.Runtime().GetVMInput().CallType
		if len(arguments) > 0 {
			if len(arguments) == 1 && callType == vm.AsynchronousCallBack {
				// group callback, we are limited to one argument ...
				gasUsed = big.NewInt(0).SetBytes(arguments[0]).Int64()
			} else {
				edgeTypeArgIndex := 0
				gasUsedArgIndex := 1
				if host.Runtime().GetVMInput().CallType == vm.AsynchronousCallBack {
					// for callbacks, first argument is return code
					edgeTypeArgIndex = 1
					gasUsedArgIndex = 2
				}
				edgeType := big.NewInt(0).SetBytes(arguments[edgeTypeArgIndex]).Int64()
				if edgeType == Async {
					host.Output().Finish(big.NewInt(int64(Callback)).Bytes()) // edge type
					host.Output().Finish(arguments[2])                        // gas used by callback
				}

				gasUsed = big.NewInt(0).SetBytes(arguments[gasUsedArgIndex]).Int64() // gas used for call
			}
		}
	}
	return gasUsed
}

func addFunctionToTempList(contract *MockTestSmartContract, functionName string, isCallBack bool) {
	_, functionPresent := contract.tempFunctionsList[functionName]
	if !functionPresent {
		contract.tempFunctionsList[functionName] = isCallBack
	}
}

// CreateGraphTest1 -
func CreateGraphTest1() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 5000, 10)

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc2f3, "cb2", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(500).
		SetGasUsed(7)

	// sc2f6 := callGraph.AddNode("sc2", "f6")
	// callGraph.AddAsyncEdge(sc1f1, sc2f6, "cb4", "gr1").
	// 	SetGasLimit(400).
	// 	SetGasUsed(7).
	// 	SetGasUsedByCallback(5)

	// sc3f7 := callGraph.AddNode("sc3", "f7")
	// callGraph.AddAsyncEdge(sc1f1, sc3f7, "cb4", "gr2").
	// 	SetGasLimit(30).
	// 	SetGasUsed(5).
	// 	SetGasUsedByCallback(5)

	// callGraph.AddNode("sc1", "cb4")

	sc3f4 := callGraph.AddNode("sc3", "f4")

	// callGraph.AddSyncEdge(sc2f3, sc3f4).
	// 	SetGasLimit(20).
	// 	SetGasUsed(12)

	callGraph.AddAsyncEdge(sc2f2, sc3f4, "cb3", "gr3").
		SetGasLimit(100).
		SetGasUsed(2).
		SetGasUsedByCallback(3)

	sc1cb2 := callGraph.AddNode("sc1", "cb2")
	sc4f5 := callGraph.AddNode("sc4", "f5")
	callGraph.AddSyncEdge(sc1cb2, sc4f5).
		SetGasLimit(4).
		SetGasUsed(2)

	callGraph.AddNode("sc2", "cb3")

	return callGraph
}

// CreateGraphTestAsyncCallsAsync -
func CreateGraphTestAsyncCallsAsync() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncCrossLocal -
func CreateGraphTestAsyncCallsAsyncCrossLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncLocalCross -
func CreateGraphTestAsyncCallsAsyncLocalCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncCrossShard -
func CreateGraphTestAsyncCallsAsyncCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestCallbackCallsSync -
func CreateGraphTestCallbackCallsSync() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(1000).
		SetGasUsed(70).
		SetGasUsedByCallback(500)

	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddSyncEdge(sc1cb1, sc2f3).
		SetGasLimit(200).
		SetGasUsed(60)

	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestCallbackCallsAsyncLocalLocal -
func CreateGraphTestCallbackCallsAsyncLocalLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(1000).
		SetGasUsed(70).
		SetGasUsedByCallback(500)

	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc1cb1, sc3f3, "cb2", "gr2").
		SetGasLimit(200).
		SetGasUsed(60).
		SetGasUsedByCallback(30)

	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestCallbackCallsAsyncLocalCross -
func CreateGraphTestCallbackCallsAsyncLocalCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(1000).
		SetGasUsed(70).
		SetGasUsedByCallback(500)

	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc1cb1, sc3f3, "cb2", "gr2").
		SetGasLimit(200).
		SetGasUsed(60).
		SetGasUsedByCallback(30)

	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestCallbackCallsAsyncCrossLocal -
func CreateGraphTestCallbackCallsAsyncCrossLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(1000).
		SetGasUsed(70).
		SetGasUsedByCallback(500)

	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc1cb1, sc3f3, "cb2", "gr2").
		SetGasLimit(200).
		SetGasUsed(60).
		SetGasUsedByCallback(30)

	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestCallbackCallsAsyncCrossCross -
func CreateGraphTestCallbackCallsAsyncCrossCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(1000).
		SetGasUsed(70).
		SetGasUsedByCallback(500)

	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncCrossShardEdge(sc1cb1, sc2f3, "cb2", "gr2").
		SetGasLimit(200).
		SetGasUsed(60).
		SetGasUsedByCallback(30)

	callGraph.AddNode("sc1", "cb2")

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
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc2f3 := callGraph.AddNode("sc2", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc2f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	sc2f4 := callGraph.AddNode("sc2", "f4")
	callGraph.AddAsyncEdge(sc2f3, sc2f4, "cb3", "gr3").
		SetGasLimit(10).
		SetGasUsed(5).
		SetGasUsedByCallback(2)

	callGraph.AddNode("sc2", "cb3")

	return callGraph
}

// CreateGraphTestSyncCalls -
func CreateGraphTestSyncCalls() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(100).
		SetGasUsed(7)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc1f1, sc3f3).
		SetGasLimit(100).
		SetGasUsed(7)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(35).
		SetGasUsed(7)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc3f3, sc5f5).
		SetGasLimit(35).
		SetGasUsed(7)

	return callGraph
}

// CreateGraphTestSyncCalls2 -
func CreateGraphTestSyncCalls2() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(100).
		SetGasUsed(7)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(50).
		SetGasUsed(7)

	return callGraph
}

// CreateGraphTestOneAsyncCall -
func CreateGraphTestOneAsyncCall() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestTwoAsyncCalls -
func CreateGraphTestTwoAsyncCalls() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(20).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc3f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestTwoAsyncCallsLocalCross -
func CreateGraphTestTwoAsyncCallsLocalCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(20).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	s3f3 := callGraph.AddNode("s3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc1f1, s3f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestTwoAsyncCallsCrossLocal -
func CreateGraphTestTwoAsyncCallsCrossLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(20).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	s3f3 := callGraph.AddNode("s3", "f3")
	callGraph.AddAsyncEdge(sc1f1, s3f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestTwoAsyncCallsCrossShard -
func CreateGraphTestTwoAsyncCallsCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(20).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc3f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc1", "cb2")

	return callGraph
}

// CreateGraphTestSyncAndAsync1 -
func CreateGraphTestSyncAndAsync1() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(400).
		SetGasUsed(60).
		SetGasUsedByCallback(100)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(200).
		SetGasUsed(10)

	return callGraph
}

// CreateGraphTestSyncAndAsync2 -
func CreateGraphTestSyncAndAsync2() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(700).
		SetGasUsed(70)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb1", "").
		SetGasLimit(400).
		SetGasUsed(60).
		SetGasUsedByCallback(100)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1f1, sc4f4).
		SetGasLimit(200).
		SetGasUsed(10)

	sc2cb1 := callGraph.AddNode("sc2", "cb1")
	callGraph.AddSyncEdge(sc4f4, sc2cb1).
		SetGasLimit(80).
		SetGasUsed(50)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc2cb1, sc5f5).
		SetGasLimit(20).
		SetGasUsed(10)

	return callGraph
}

// CreateGraphTestSyncAndAsync3 -
func CreateGraphTestSyncAndAsync3() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(55).
		SetGasUsed(7).
		SetGasUsedByCallback(6)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc2f2, sc5f5).
		SetGasLimit(10).
		SetGasUsed(4)

	sc3f3 := callGraph.AddNode("sc3", "f3")

	sc1cb1 := callGraph.AddNode("sc1", "cb1")
	callGraph.AddSyncEdge(sc1cb1, sc3f3).
		SetGasLimit(20).
		SetGasUsed(3)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(10).
		SetGasUsed(5)

	return callGraph
}

// CreateGraphTestSyncAndAsync4 -
func CreateGraphTestSyncAndAsync4() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1f1, sc4f4).
		SetGasLimit(400).
		SetGasUsed(40)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddAsyncEdge(sc4f4, sc5f5, "cb2", "").
		SetGasLimit(200).
		SetGasUsed(20).
		SetGasUsedByCallback(100)

	callGraph.AddNode("sc4", "cb2")

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(400).
		SetGasUsed(60).
		SetGasUsedByCallback(100)

	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc1cb1, sc3f3).
		SetGasLimit(10).
		SetGasUsed(5)

	return callGraph
}

// CreateGraphTestSyncAndAsync5 -
func CreateGraphTestSyncAndAsync5() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(800).
		SetGasUsed(50).
		SetGasUsedByCallback(20)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(500).
		SetGasUsed(20)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(300).
		SetGasUsed(15)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncEdge(sc4f4, sc5f5, "cb4", "").
		SetGasLimit(100).
		SetGasUsed(50).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc4", "cb4")

	return callGraph
}

// CreateGraphTestSyncAndAsync6 -
func CreateGraphTestSyncAndAsync6() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc0f0 := callGraph.AddStartNode("sc0", "f0", 3000, 10)

	sc1f1 := callGraph.AddNode("sc1", "f1")

	callGraph.AddAsyncEdge(sc0f0, sc1f1, "cb0", "").
		SetGasLimit(1300).
		SetGasUsed(60).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc0", "cb0")

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddSyncEdge(sc1f1, sc6f6).
		SetGasLimit(1000).
		SetGasUsed(5)

	return callGraph
}

// CreateGraphTestSyncAndAsync7 -
func CreateGraphTestSyncAndAsync7() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc0f0 := callGraph.AddStartNode("sc0", "f0", 3000, 10)

	sc1f1 := callGraph.AddNode("sc1", "f1")
	callGraph.AddSyncEdge(sc0f0, sc1f1).
		SetGasLimit(1000).
		SetGasUsed(100)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(500).
		SetGasUsed(40)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc1f1, sc3f3).
		SetGasLimit(100).
		SetGasUsed(10)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddAsyncEdge(sc2f2, sc4f4, "cb2", "").
		SetGasLimit(300).
		SetGasUsed(60).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestDifferentTypeOfCallsToSameFunction -
func CreateGraphTestDifferentTypeOfCallsToSameFunction() *TestCallGraph {
	callGraph := CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 2000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(20).
		SetGasUsed(5)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb2", "").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

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
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 5000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(20).
		SetGasUsed(5)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb2", "").
		SetGasLimit(600).
		SetGasUsed(6).
		SetGasUsedByCallback(400)

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
		SetGasUsedByCallback(3)

	return callGraph
}

// CreateGraphTestOneAsyncCallCrossShard -
func CreateGraphTestOneAsyncCallCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")

	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard2 -
func CreateGraphTestAsyncCallsCrossShard2() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 800, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(220).
		SetGasUsed(27).
		SetGasUsedByCallback(23)

	sc3f6 := callGraph.AddNode("sc3", "f6")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f6, "cb2", "").
		SetGasLimit(30).
		SetGasUsed(10).
		SetGasUsedByCallback(6)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard3 -
func CreateGraphTestAsyncCallsCrossShard3() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 800, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(220).
		SetGasUsed(27).
		SetGasUsedByCallback(23)

	sc3f6 := callGraph.AddNode("sc3", "f6")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f6, "cb2", "").
		SetGasLimit(30).
		SetGasUsed(10).
		SetGasUsedByCallback(6)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard4 -
func CreateGraphTestAsyncCallsCrossShard4() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	//callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "").
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(800).
		SetGasUsed(50).
		SetGasUsedByCallback(20)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(500).
		SetGasUsed(20)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(300).
		SetGasUsed(15)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncCrossShardEdge(sc4f4, sc5f5, "cb4", "").
		SetGasLimit(100).
		SetGasUsed(50).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc4", "cb4")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard5 -
func CreateGraphTestAsyncCallsCrossShard5() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 3000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(2000).
		SetGasUsed(50).
		SetGasUsedByCallback(20)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(500).
		SetGasUsed(20)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(300).
		SetGasUsed(15)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncCrossShardEdge(sc4f4, sc5f5, "cb4", "").
		SetGasLimit(100).
		SetGasUsed(50).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc1", "cb1")
	callGraph.AddNode("sc4", "cb4")

	sc3f6 := callGraph.AddNode("sc3", "f6")
	callGraph.AddSyncEdge(sc2f2, sc3f6).
		SetGasLimit(500).
		SetGasUsed(20)

	sc4f7 := callGraph.AddNode("sc4", "f7")
	callGraph.AddSyncEdge(sc3f6, sc4f7).
		SetGasLimit(300).
		SetGasUsed(15)

	sc5f8 := callGraph.AddNode("sc5", "f8")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncCrossShardEdge(sc4f7, sc5f8, "cb5", "").
		SetGasLimit(100).
		SetGasUsed(50).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc4", "cb5")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard6 -
func CreateGraphTestAsyncCallsCrossShard6() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc0f0 := callGraph.AddStartNode("sc0", "f0", 3000, 10)

	sc1f1 := callGraph.AddNode("sc1", "f1")

	callGraph.AddAsyncEdge(sc0f0, sc1f1, "cb0", "").
		SetGasLimit(1300).
		SetGasUsed(60).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc0", "cb0")

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddSyncEdge(sc1f1, sc6f6).
		SetGasLimit(1000).
		SetGasUsed(5)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc6f6, sc2f2, "cb1", "").
		SetGasLimit(800).
		SetGasUsed(50).
		SetGasUsedByCallback(20)

	callGraph.AddNode("sc6", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(500).
		SetGasUsed(20)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(300).
		SetGasUsed(15)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncCrossShardEdge(sc4f4, sc5f5, "cb4", "").
		SetGasLimit(100).
		SetGasUsed(50).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc4", "cb4")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard7 -
func CreateGraphTestAsyncCallsCrossShard7() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc0f0 := callGraph.AddStartNode("sc0", "f0", 3000, 10)

	sc1f1 := callGraph.AddNode("sc1", "f1")

	callGraph.AddAsyncEdge(sc0f0, sc1f1, "cb0", "").
		SetGasLimit(1300).
		SetGasUsed(60).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc0", "cb0")

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddSyncEdge(sc1f1, sc6f6).
		SetGasLimit(1000).
		SetGasUsed(5)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc6f6, sc2f2, "cb1", "").
		SetGasLimit(800).
		SetGasUsed(50).
		SetGasUsedByCallback(20)

	callGraph.AddNode("sc6", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(500).
		SetGasUsed(20)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(300).
		SetGasUsed(15)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddNode("sc2", "cb2")
	callGraph.AddAsyncEdge(sc4f4, sc5f5, "cb4", "").
		SetGasLimit(100).
		SetGasUsed(50).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc4", "cb4")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard8 -
func CreateGraphTestAsyncCallsCrossShard8() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc0f0 := callGraph.AddStartNode("sc0", "f0", 3000, 10)

	sc1f1 := callGraph.AddNode("sc1", "f1")

	callGraph.AddAsyncEdge(sc0f0, sc1f1, "cb0", "").
		SetGasLimit(1300).
		SetGasUsed(60).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc0", "cb0")

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddSyncEdge(sc1f1, sc6f6).
		SetGasLimit(1000).
		SetGasUsed(5)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc6f6, sc2f2, "cb1", "").
		SetGasLimit(800).
		SetGasUsed(50).
		SetGasUsedByCallback(20)

	callGraph.AddNode("sc6", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(600).
		SetGasUsed(20)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc3f3, sc4f4).
		SetGasLimit(500).
		SetGasUsed(15)

	callGraph.AddNode("sc2", "cb2")

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddAsyncCrossShardEdge(sc4f4, sc5f5, "cb4", "").
		SetGasLimit(70).
		SetGasUsed(50).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc4", "cb4")

	sc7f7 := callGraph.AddNode("sc7", "f7")
	callGraph.AddAsyncCrossShardEdge(sc4f4, sc7f7, "cb7", "").
		SetGasLimit(30).
		SetGasUsed(20).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc4", "cb7")

	return callGraph
}

// CreateGraphTestAsyncCallsCrossShard9 -
func CreateGraphTestAsyncCallsCrossShard9() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc0f0 := callGraph.AddStartNode("sc0", "f0", 5500, 10)

	sc1f1 := callGraph.AddNode("sc1", "f1")
	callGraph.AddAsyncEdge(sc0f0, sc1f1, "cb0", "").
		SetGasLimit(5000).
		SetGasUsed(60).
		SetGasUsedByCallback(10)

	callGraph.AddNode("sc0", "cb0")

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1f1, sc4f4).
		SetGasLimit(3000).
		SetGasUsed(25)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddAsyncCrossShardEdge(sc4f4, sc5f5, "cb4", "").
		SetGasLimit(2800).
		SetGasUsed(50).
		SetGasUsedByCallback(35)

	sc4cb4 := callGraph.AddNode("sc4", "cb4")

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddAsyncCrossShardEdge(sc5f5, sc6f6, "cb5", "").
		SetGasLimit(100).
		SetGasUsed(10).
		SetGasUsedByCallback(2)

	callGraph.AddNode("sc5", "cb5")

	sc7f7 := callGraph.AddNode("sc7", "f7")
	callGraph.AddAsyncEdge(sc4cb4, sc7f7, "cb44", "").
		SetGasLimit(1500).
		SetGasUsed(45).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc4", "cb44")

	sc8f8 := callGraph.AddNode("sc8", "f8")
	callGraph.AddAsyncCrossShardEdge(sc7f7, sc8f8, "cb71", "").
		SetGasLimit(500).
		SetGasUsed(20).
		SetGasUsedByCallback(45)

	callGraph.AddNode("sc7", "cb71")

	sc9f9 := callGraph.AddNode("sc9", "f9")
	callGraph.AddAsyncCrossShardEdge(sc7f7, sc9f9, "cb72", "").
		SetGasLimit(600).
		SetGasUsed(10).
		SetGasUsedByCallback(55)

	callGraph.AddNode("sc7", "cb72")

	// ############

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(1000).
		SetGasUsed(15)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "").
		SetGasLimit(700).
		SetGasUsed(40).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc2", "cb2")

	sc3f33 := callGraph.AddNode("sc3", "f33")
	callGraph.AddSyncEdge(sc3f3, sc3f33).
		SetGasLimit(300).
		SetGasUsed(100)

	return callGraph
}

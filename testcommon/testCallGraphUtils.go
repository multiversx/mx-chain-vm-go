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

type argIndexesForGraphCall struct {
	edgeTypeIdx          int
	failureIdx           int
	gasUsedIdx           int
	gasUsedByCallbackIdx int
	callbackFailureIdx   int
}

var syncCallArgIndexes = argIndexesForGraphCall{
	edgeTypeIdx: 0,
	failureIdx:  1,
	gasUsedIdx:  2,
}

var asyncCallArgIndexes = argIndexesForGraphCall{
	edgeTypeIdx:          0,
	failureIdx:           1,
	gasUsedIdx:           2,
	gasUsedByCallbackIdx: 3,
	callbackFailureIdx:   4,
}

var callbackCallArgIndexes = argIndexesForGraphCall{
	edgeTypeIdx: 1,
	failureIdx:  2,
	gasUsedIdx:  3,
}

// RuntimeConfigOfCall -
type RuntimeConfigOfCall struct {
	gasUsed           uint64
	gasUsedByCallback uint64
	edgeType          TestCallEdgeType
	willFail          bool
	willCallbackFail  bool
}

// CallsFinishData -
type CallsFinishData struct {
	Data []*CallFinishDataItem
}

// CreateMockContractsFromAsyncTestCallGraph creates the contracts
// with functions that reflect the behavior specified by the call graph
func CreateMockContractsFromAsyncTestCallGraph(callGraph *TestCallGraph, callsFinishData *CallsFinishData, runtimeConfigsForCalls map[string]*RuntimeConfigOfCall, testConfig *TestConfig) []MockTestSmartContract {
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
				// fmt.Println("-> shard of ", string(node.Call.ContractAddress), "is", node.ShardID)
			} else {
				shardID = node.ShardID
			}

			newContract := CreateMockContract(node.Call.ContractAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithShardID(shardID).
				WithMethods(func(instanceMock *mock.InstanceMock, testConfig *TestConfig) {
					for functionName := range contracts[contractAddressAsString].tempFunctionsList {
						if functionName == FakeCallbackName {
							continue
						}
						instanceMock.AddMockMethod(functionName, func() *mock.InstanceMock {
							host := instanceMock.Host
							instance := mock.GetMockInstance(host)
							t := instance.T

							crtFunctionCalled := host.Runtime().Function()
							LogGraph.Trace("Executing graph node", "sc", string(host.Runtime().GetSCAddress()), "func", crtFunctionCalled)

							crtNode := callGraph.FindNode(host.Runtime().GetSCAddress(), crtFunctionCalled)
							var runtimeConfig *RuntimeConfigOfCall
							if crtNode.IsStartNode {
								runtimeConfig = &RuntimeConfigOfCall{
									gasUsed: crtNode.GasUsed,
								}
							} else {
								runtimeConfig = readGasUsedFromArguments(host, runtimeConfigsForCalls)
							}
							runtimeConfigsForCalls[string(host.Async().GetCallID())] = runtimeConfig

							// prepare arguments for callback
							if runtimeConfig.edgeType == Async || runtimeConfig.edgeType == AsyncCrossShard {
								arguments := argumentsFromRuntimeConfig(runtimeConfig)
								if !runtimeConfig.willFail {
									createFinishDataFromArguments(host.Output(), arguments)
								} else {
									host.Runtime().FailExecution(fmt.Errorf("async call fail"))
									return instance
								}
							} else if runtimeConfig.edgeType == Sync && runtimeConfig.willFail {
								host.Runtime().FailExecution(fmt.Errorf("sync call fail"))
								return instance
							} else if runtimeConfig.willFail {
								host.Runtime().FailExecution(fmt.Errorf("callback fail"))
								return instance
							}

							// TODO matei-p change to debug logging
							// burn gas for function
							fmt.Println("Burning", "gas", runtimeConfig.gasUsed, "function", crtFunctionCalled)
							host.Metering().UseGasBounded(uint64(runtimeConfig.gasUsed))

							for _, edge := range crtNode.AdjacentEdges {
								if edge.Type == Sync {
									breakPointValue := makeSyncCallFromEdge(host, edge, testConfig)
									if breakPointValue == arwen.BreakpointExecutionFailed {
										host.Runtime().FailExecution(fmt.Errorf("breakpoint sync execution failed"))
										return instance
									}
								} else {
									err := makeAsyncCallFromEdge(host, edge, testConfig)
									require.Nil(t, err)
								}
							}

							callFinishData := computeReturnDataForTestFramework(crtFunctionCalled, host)
							callsFinishData.Data = append(callsFinishData.Data, callFinishData)

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

func argumentsFromRuntimeConfig(runtimeConfig *RuntimeConfigOfCall) [][]byte {
	arguments := make([][]byte, 3)
	var callbackEdgeType TestCallEdgeType
	if runtimeConfig.edgeType == Async {
		callbackEdgeType = Callback
	} else {
		callbackEdgeType = CallbackCrossShard
	}
	setGraphCallArg(arguments, syncCallArgIndexes.edgeTypeIdx, int(callbackEdgeType))
	setGraphCallArg(arguments, syncCallArgIndexes.failureIdx, failAsInt(runtimeConfig.willCallbackFail))
	setGraphCallArg(arguments, syncCallArgIndexes.gasUsedIdx, int(runtimeConfig.gasUsedByCallback))
	return arguments
}

func makeSyncCallFromEdge(host arwen.VMHost, edge *TestCallEdge, testConfig *TestConfig) arwen.BreakpointValue {
	value := big.NewInt(testConfig.TransferFromParentToChild)
	destFunctionName := edge.To.Call.FunctionName
	destAddress := edge.To.Call.ContractAddress

	arguments := make([][]byte, 3)
	setGraphCallArg(arguments, syncCallArgIndexes.edgeTypeIdx, Sync)
	setGraphCallArg(arguments, syncCallArgIndexes.failureIdx, failAsInt(edge.Fail))
	setGraphCallArg(arguments, syncCallArgIndexes.gasUsedIdx, int(edge.GasUsed))

	LogGraph.Trace("Sync call to ", string(destAddress), " func ", destFunctionName, " gas ", edge.GasLimit)
	elrondapi.ExecuteOnDestContextWithTypedArgs(
		host,
		int64(edge.GasLimit),
		value,
		[]byte(destFunctionName),
		destAddress,
		arguments)

	return host.Runtime().GetRuntimeBreakpointValue()
}

func failAsInt(fail bool) int {
	failAsInt := 0
	if fail {
		failAsInt = 1
	}
	return failAsInt
}

func setGraphCallArg(arguments [][]byte, index int, value int) {
	arguments[index] = big.NewInt(int64(value)).Bytes()
}

func makeAsyncCallFromEdge(host arwen.VMHost, edge *TestCallEdge, testConfig *TestConfig) error {
	async := host.Async()
	destFunctionName := edge.To.Call.FunctionName
	destAddress := edge.To.Call.ContractAddress
	value := big.NewInt(testConfig.TransferFromParentToChild)

	LogGraph.Trace("Register async call", "to", string(destAddress), "func", destFunctionName, "gas", edge.GasLimit)

	arguments := make([][]byte, 5)
	setGraphCallArg(arguments, asyncCallArgIndexes.edgeTypeIdx, int(edge.Type))
	setGraphCallArg(arguments, asyncCallArgIndexes.failureIdx, failAsInt(edge.Fail))
	setGraphCallArg(arguments, asyncCallArgIndexes.gasUsedIdx, int(edge.GasUsed))
	setGraphCallArg(arguments, asyncCallArgIndexes.gasUsedByCallbackIdx, int(edge.GasUsedByCallback))
	setGraphCallArg(arguments, asyncCallArgIndexes.callbackFailureIdx, failAsInt(edge.CallbackFail))

	callDataAsBytes, _ := createEncodedDataFromArguments(destFunctionName, arguments)

	callbackName := edge.Callback
	if callbackName == FakeCallbackName {
		callbackName = ""
	}

	err := async.RegisterAsyncCall("", &arwen.AsyncCall{
		Status:          arwen.AsyncCallPending,
		Destination:     destAddress,
		Data:            callDataAsBytes,
		ValueBytes:      value.Bytes(),
		GasLimit:        edge.GasLimit,
		SuccessCallback: callbackName,
		ErrorCallback:   callbackName,
	})
	return err
}

func createEncodedDataFromArguments(destFunctionName string, arguments [][]byte) ([]byte, string) {
	callData := txDataBuilder.NewBuilder()
	callData.Func(destFunctionName)
	for _, arg := range arguments {
		callData.Bytes(arg)
	}
	return callData.ToBytes(), callData.ToString()
}

func createFinishDataFromArguments(output arwen.OutputContext, arguments [][]byte) {
	for _, arg := range arguments {
		output.Finish(arg)
	}
}

// CallFinishDataItem -
type CallFinishDataItem struct {
	ContractAndFunction          string
	GasProvided                  uint64
	GasRemaining                 uint64
	CallID                       []byte
	CallbackAsyncInitiatorCallID []byte
	IsCrossShard                 bool
}

// return data is encoded using standard txDataBuilder
// format is function@nodeLabel@providedGas@remainingGas
func computeReturnDataForTestFramework(crtFunctionCalled string, host arwen.VMHost) *CallFinishDataItem {
	runtime := host.Runtime()
	metering := host.Metering()
	async := host.Async()
	LogGraph.Trace("End of ", crtFunctionCalled, " on ", string(host.Runtime().GetSCAddress()))

	/// TODO matei-p change to debug logging
	fmt.Println(
		"CallFinishDataItem -> callID", async.GetCallID(),
		"CallbackAsyncInitiatorCallID", async.GetCallbackAsyncInitiatorCallID(),
		"IsCrossShard", async.IsCrossShard(),
		"For contract ", string(runtime.GetSCAddress()), "/ "+crtFunctionCalled+"\t",
		"Gas provided", fmt.Sprintf("%d\t", runtime.GetVMInput().GasProvided),
		"Gas remaining", fmt.Sprintf("%d\t", metering.GasLeft()))

	return &CallFinishDataItem{
		ContractAndFunction:          string(runtime.GetSCAddress()) + "_" + crtFunctionCalled + TestReturnDataSuffix,
		GasProvided:                  runtime.GetVMInput().GasProvided,
		GasRemaining:                 metering.GasLeft(),
		CallID:                       async.GetCallID(),
		CallbackAsyncInitiatorCallID: async.GetCallbackAsyncInitiatorCallID(),
		IsCrossShard:                 async.IsCrossShard(),
	}
}

func readGasUsedFromArguments(host arwen.VMHost, runtimeConfigsForCalls map[string]*RuntimeConfigOfCall) *RuntimeConfigOfCall {
	runtimeConfig := &RuntimeConfigOfCall{}
	var argIndexes argIndexesForGraphCall

	arguments := host.Runtime().Arguments()
	callType := host.Runtime().GetVMInput().CallType

	if callType == vm.DirectCall {
		argIndexes = syncCallArgIndexes
	} else if callType == vm.AsynchronousCall {
		argIndexes = asyncCallArgIndexes
		runtimeConfig.gasUsedByCallback = big.NewInt(0).SetBytes(arguments[argIndexes.gasUsedByCallbackIdx]).Uint64()
		runtimeConfig.willCallbackFail = (big.NewInt(0).SetBytes(arguments[argIndexes.callbackFailureIdx]).Int64() == 1)
	} else if callType == vm.AsynchronousCallBack {
		// for callbacks, first argument is return code
		returnCode := big.NewInt(0).SetBytes(arguments[0]).Int64()
		argIndexes = callbackCallArgIndexes
		if returnCode != 0 {
			runtimeConfig = runtimeConfigsForCalls[string(host.Async().GetCallerCallID())]
			if runtimeConfig.edgeType == Async {
				runtimeConfig.edgeType = Callback
			} else if runtimeConfig.edgeType == AsyncCrossShard {
				runtimeConfig.edgeType = CallbackCrossShard
			} else {
				panic("this should not happen!")
			}
			runtimeConfig.gasUsed = runtimeConfig.gasUsedByCallback
			runtimeConfig.willFail = runtimeConfig.willCallbackFail
			if runtimeConfig != nil {
				return runtimeConfig
			}
		}
	}

	runtimeConfig.edgeType = TestCallEdgeType(big.NewInt(0).SetBytes(arguments[argIndexes.edgeTypeIdx]).Int64())
	runtimeConfig.gasUsed = big.NewInt(0).SetBytes(arguments[argIndexes.gasUsedIdx]).Uint64()
	runtimeConfig.willFail = (big.NewInt(0).SetBytes(arguments[argIndexes.failureIdx]).Int64() == 1)

	return runtimeConfig
}

func addFunctionToTempList(contract *MockTestSmartContract, functionName string, isCallBack bool) {
	_, functionPresent := contract.tempFunctionsList[functionName]
	if !functionPresent {
		contract.tempFunctionsList[functionName] = isCallBack
	}
}

// CreateGraphTestSyncAndAsync8 -
func CreateGraphTestSyncAndAsync8() *TestCallGraph {
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

	sc3f4 := callGraph.AddNode("sc3", "f4")

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

// CreateGraphTestAsyncCallsAsyncFirstNoCallback -
func CreateGraphTestAsyncCallsAsyncFirstNoCallback() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "", "gr1").
		SetGasLimit(500).
		SetGasUsed(7)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondNoCallback -
func CreateGraphTestAsyncCallsAsyncSecondNoCallback() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "", "gr2").
		SetGasLimit(100).
		SetGasUsed(6)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstFail -
func CreateGraphTestAsyncCallsAsyncFirstFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstFail -
func CreateGraphTestAsyncCallsAsyncFirstCallbackFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetCallbackFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondFail -
func CreateGraphTestAsyncCallsAsyncSecondFail() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetFail()

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondCallbackFail -
func CreateGraphTestAsyncCallsAsyncSecondCallbackFail() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetGasUsedByCallback(5)

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

// CreateGraphTestAsyncCallsAsyncFirstNoCallbackCrossLocal -
func CreateGraphTestAsyncCallsAsyncFirstNoCallbackCrossLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "", "gr1").
		SetGasLimit(500).
		SetGasUsed(7)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondNoCallbackCrossLocal -
func CreateGraphTestAsyncCallsAsyncSecondNoCallbackCrossLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "", "gr2").
		SetGasLimit(100).
		SetGasUsed(6)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstFailCrossLocal -
func CreateGraphTestAsyncCallsAsyncFirstFailCrossLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstCallbackFailCrossLocal -
func CreateGraphTestAsyncCallsAsyncFirstCallbackFailCrossLocal() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetCallbackFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondFailCrossLocal -
func CreateGraphTestAsyncCallsAsyncSecondFailCrossLocal() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetFail()

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondCallbackFailCrossLocal -
func CreateGraphTestAsyncCallsAsyncSecondCallbackFailCrossLocal() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetCallbackFail()

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

// CreateGraphTestAsyncCallsAsyncFirstNoCallbackLocalCross -
func CreateGraphTestAsyncCallsAsyncFirstNoCallbackLocalCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "", "gr1").
		SetGasLimit(500).
		SetGasUsed(7)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondNoCallbackLocalCross -
func CreateGraphTestAsyncCallsAsyncSecondNoCallbackLocalCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "", "gr2").
		SetGasLimit(100).
		SetGasUsed(6)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstFailLocalCross -
func CreateGraphTestAsyncCallsAsyncFirstFailLocalCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondFailLocalCross -
func CreateGraphTestAsyncCallsAsyncSecondFailLocalCross() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetFail()

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstCallbackFailLocalCross -
func CreateGraphTestAsyncCallsAsyncFirstCallbackFailLocalCross() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetCallbackFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondCallbackFailLocalCross -
func CreateGraphTestAsyncCallsAsyncSecondCallbackFailLocalCross() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetCallbackFail()

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

// CreateGraphTestAsyncCallsAsyncFirstNoCallbackCrossShard -
func CreateGraphTestAsyncCallsAsyncFirstNoCallbackCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "", "gr1").
		SetGasLimit(500).
		SetGasUsed(7)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondNoCallbackCrossShard -
func CreateGraphTestAsyncCallsAsyncSecondNoCallbackCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5)

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "", "gr2").
		SetGasLimit(100).
		SetGasUsed(6)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstFailCrossShard -
func CreateGraphTestAsyncCallsAsyncFirstFailCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondFailCrossShard -
func CreateGraphTestAsyncCallsAsyncSecondFailCrossShard() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetFail()

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncFirstCallbackFailCrossShard -
func CreateGraphTestAsyncCallsAsyncFirstCallbackFailCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(500).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetCallbackFail()

	callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb2", "gr2").
		SetGasLimit(100).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	callGraph.AddNode("sc2", "cb2")

	return callGraph
}

// CreateGraphTestAsyncCallsAsyncSecondCallbackFailCrossShard -
func CreateGraphTestAsyncCallsAsyncSecondCallbackFailCrossShard() *TestCallGraph {
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
		SetGasUsedByCallback(3).
		SetCallbackFail()

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

// CreateGraphTestOneAsyncCallNoCallback -
func CreateGraphTestOneAsyncCallNoCallback() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestOneAsyncCallFail -
func CreateGraphTestOneAsyncCallFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6).
		SetFail()
	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(10).
		SetGasUsed(4)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1cb1, sc4f4).
		SetGasLimit(100).
		SetGasUsed(40)

	return callGraph
}

// CreateGraphTestOneAsyncCallNoCallbackFail -
func CreateGraphTestOneAsyncCallNoCallbackFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6).
		SetFail()

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestAsyncCallIndirectFailCrossShard -
func CreateGraphTestAsyncCallIndirectFailCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)
	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddSyncEdge(sc2f2, sc6f6).
		SetGasLimit(8).
		SetGasUsed(3)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(10).
		SetGasUsed(4)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc3f3, sc5f5).
		SetGasLimit(3).
		SetGasUsed(1).
		SetFail()

	sc7f7 := callGraph.AddNode("sc7", "f7")
	callGraph.AddSyncEdge(sc2f2, sc7f7).
		SetGasLimit(5).
		SetGasUsed(2)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1cb1, sc4f4).
		SetGasLimit(100).
		SetGasUsed(40)

	return callGraph
}

// CreateGraphTestAsyncCallIndirectFail -
func CreateGraphTestAsyncCallIndirectFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)
	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddSyncEdge(sc2f2, sc6f6).
		SetGasLimit(8).
		SetGasUsed(3)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(10).
		SetGasUsed(4)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc3f3, sc5f5).
		SetGasLimit(3).
		SetGasUsed(1).
		SetFail()

	sc7f7 := callGraph.AddNode("sc7", "f7")
	callGraph.AddSyncEdge(sc2f2, sc7f7).
		SetGasLimit(5).
		SetGasUsed(2)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1cb1, sc4f4).
		SetGasLimit(100).
		SetGasUsed(40)

	return callGraph
}

// CreateGraphTestOneAsyncCallbackFail -
func CreateGraphTestOneAsyncCallbackFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6).
		SetCallbackFail()
	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(10).
		SetGasUsed(4)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1cb1, sc4f4).
		SetGasLimit(100).
		SetGasUsed(40)

	return callGraph
}

// CreateGraphTestAsyncCallbackIndirectFail -
func CreateGraphTestAsyncCallbackIndirectFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)
	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(10).
		SetGasUsed(4)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1cb1, sc4f4).
		SetGasLimit(100).
		SetGasUsed(40)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc1cb1, sc5f5).
		SetGasLimit(10).
		SetGasUsed(4).
		SetFail()

	return callGraph
}

// CreateGraphTestOneAsyncCallbackFailCrossShard -
func CreateGraphTestOneAsyncCallbackFailCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc1f1, sc5f5).
		SetGasLimit(400).
		SetGasUsed(6)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc5f5, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6).
		SetCallbackFail()
	sc5cb1 := callGraph.AddNode("sc5", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(10).
		SetGasUsed(4)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc5cb1, sc4f4).
		SetGasLimit(100).
		SetGasUsed(40)

	return callGraph
}

// CreateGraphTestAsyncCallbackIndirectFailCrossShard -
func CreateGraphTestAsyncCallbackIndirectFailCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)
	sc1cb1 := callGraph.AddNode("sc1", "cb1")

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddSyncEdge(sc2f2, sc3f3).
		SetGasLimit(10).
		SetGasUsed(4)

	sc6f6 := callGraph.AddNode("sc6", "f6")
	callGraph.AddSyncEdge(sc1cb1, sc6f6).
		SetGasLimit(10).
		SetGasUsed(4)

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc1cb1, sc4f4).
		SetGasLimit(100).
		SetGasUsed(40)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc4f4, sc5f5).
		SetGasLimit(5).
		SetGasUsed(2).
		SetFail()

	sc7f7 := callGraph.AddNode("sc7", "f7")
	callGraph.AddSyncEdge(sc1cb1, sc7f7).
		SetGasLimit(12).
		SetGasUsed(2)

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

// CreateGraphTestTwoAsyncCallsOneFail -
func CreateGraphTestTwoAsyncCallsOneFail() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 1000, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddAsyncEdge(sc1f1, sc2f2, "cb1", "gr1").
		SetGasLimit(20).
		SetGasUsed(7).
		SetGasUsedByCallback(5).
		SetFail()

	sc4f4 := callGraph.AddNode("sc4", "f4")
	callGraph.AddSyncEdge(sc2f2, sc4f4).
		SetGasLimit(5).
		SetGasUsed(2)

	sc3f3 := callGraph.AddNode("sc3", "f3")
	callGraph.AddAsyncEdge(sc1f1, sc3f3, "cb2", "gr2").
		SetGasLimit(30).
		SetGasUsed(6).
		SetGasUsedByCallback(3)

	sc5f5 := callGraph.AddNode("sc5", "f5")
	callGraph.AddSyncEdge(sc3f3, sc5f5).
		SetGasLimit(4).
		SetGasUsed(1)

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

// CreateGraphTestSyncAndAsync9 -
func CreateGraphTestSyncAndAsync9() *TestCallGraph {
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

// CreateGraphTestOneAsyncCallNoCallbackCrossShard -
func CreateGraphTestOneAsyncCallNoCallbackCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")

	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "", "").
		SetGasLimit(35).
		SetGasUsed(7)

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestOneAsyncCallCrossShard2 -
func CreateGraphTestOneAsyncCallCrossShard2() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")
	callGraph.AddSyncEdge(sc1f1, sc2f2).
		SetGasLimit(300).
		SetGasUsed(20)

	sc3f3 := callGraph.AddNode("sc3", "f3")

	callGraph.AddAsyncCrossShardEdge(sc2f2, sc3f3, "cb1", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6)

	callGraph.AddNode("sc2", "cb1")

	return callGraph
}

// CreateGraphTestOneAsyncCallFailCrossShard -
func CreateGraphTestOneAsyncCallFailCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")

	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "cb1", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetGasUsedByCallback(6).
		SetFail()

	callGraph.AddNode("sc1", "cb1")

	return callGraph
}

// CreateGraphTestOneAsyncCallFailNoCallbackCrossShard -
func CreateGraphTestOneAsyncCallFailNoCallbackCrossShard() *TestCallGraph {
	callGraph := CreateTestCallGraph()

	sc1f1 := callGraph.AddStartNode("sc1", "f1", 500, 10)

	sc2f2 := callGraph.AddNode("sc2", "f2")

	callGraph.AddAsyncCrossShardEdge(sc1f1, sc2f2, "", "").
		SetGasLimit(35).
		SetGasUsed(7).
		SetFail()

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

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

							crtFunctionCalled := host.Runtime().Function()

							crtNode := callGraph.FindNode(string(host.Runtime().GetSCAddress()), crtFunctionCalled)
							fmt.Println("Executing " + crtFunctionCalled + " on " + string(host.Runtime().GetSCAddress()))
							//fmt.Println("Node " + string(crtNode.asyncCall.ContractAddress) + " / " + crtNode.asyncCall.FunctionName)

							value := big.NewInt(testConfig.TransferFromParentToChild)

							for _, edge := range crtNode.adjacentNodes {
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
									async := host.Async()
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

							host.Output().Finish([]byte(crtFunctionCalled + TestReturnDataSuffix))
							fmt.Println("End of " + crtFunctionCalled + " on " + string(host.Runtime().GetSCAddress()))

							return instance
						})
					}
				})
			contracts[contractAddressAsString] = &newContract
		}
		functionName := node.call.FunctionName
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

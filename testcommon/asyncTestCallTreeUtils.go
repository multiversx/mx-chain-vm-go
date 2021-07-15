package testcommon

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

var AsyncReturnDataSuffix = "_returnData"
var AsyncCallbackPrefix = "callback_"

var AsyncContextCallbackFunction = "contextCallback"

func CreateMockContractsFromAsyncTestCallTree(callTree *AsyncTestCallTree, testConfig *TestConfig) []MockTestSmartContract {
	contracts := make(map[string]*MockTestSmartContract)
	dfsTree(callTree, func(path []*AsyncTestCallNode, parent *AsyncTestCallNode, node *AsyncTestCallNode) *AsyncTestCallNode {
		contractAddressAsString := string(node.asyncCall.ContractAddress)
		if contracts[contractAddressAsString] == nil {
			newContract := CreateMockContract(node.asyncCall.ContractAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(func(instanceMock *mock.InstanceMock, testConfig *TestConfig) {
					for functionName, isCallBack := range contracts[contractAddressAsString].tempFunctionsList {
						//fmt.Println("Add mock method " + functionName + " to " + contractAddressAsString)
						if isCallBack {
							instanceMock.AddMockMethod(functionName,
								WasteGasWithReturnDataMockMethod(
									instanceMock,
									testConfig.GasUsedByCallback,
									[]byte(AsyncCallbackPrefix+functionName)))
						} else {
							instanceMock.AddMockMethod(functionName, func() *mock.InstanceMock {
								host := instanceMock.Host
								fmt.Println("Executing " + host.Runtime().Function() + " on " + string(host.Runtime().GetSCAddress()))
								instance := mock.GetMockInstance(host)
								t := instance.T

								async := host.Async()
								value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

								for _, child := range node.children {
									callData := txDataBuilder.NewBuilder()
									childFunctionName := child.asyncCall.FunctionName
									callData.Func(childFunctionName)

									err := async.RegisterAsyncCall(child.asyncCall.GroupName, &arwen.AsyncCall{
										Status:          arwen.AsyncCallPending,
										Destination:     child.asyncCall.ContractAddress,
										Data:            callData.ToBytes(),
										ValueBytes:      value,
										GasLimit:        testConfig.GasProvidedToChild,
										SuccessCallback: child.asyncCall.CallbackName,
										ErrorCallback:   child.asyncCall.CallbackName,
									})
									require.Nil(t, err)
								}

								host.Output().Finish([]byte(functionName + AsyncReturnDataSuffix))

								return instance
							})
						}
					}
				})
			contracts[contractAddressAsString] = &newContract
		}
		functionName := node.asyncCall.FunctionName
		contract := contracts[contractAddressAsString]
		addFunctionToTempList(contract, functionName, false)
		callbackName := node.asyncCall.CallbackName
		if callbackName != "" {
			addFunctionToTempList(contracts[string(parent.asyncCall.ContractAddress)], callbackName, true)
		}
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

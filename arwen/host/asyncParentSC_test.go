package host

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func createTestAsyncParentContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock, testConfig *asyncCallTestConfig) {
	parentInstance := imb.CreateAndStoreInstanceMock(t, host, parentAddress, testConfig.parentBalance)
	addAsyncParentMethodsToInstanceMock(parentInstance, testConfig)
}

func addAsyncParentMethodsToInstanceMock(instanceMock *mock.InstanceMock, testConfig *asyncCallTestConfig) {
	input := DefaultTestContractCallInput()
	input.GasProvided = testConfig.gasProvidedToChild

	instanceMock.AddMockMethodWithError("performAsyncCall", func() {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(testConfig.gasUsedByParent)

		host.Storage().SetStorage(parentKeyA, parentDataA)
		host.Storage().SetStorage(parentKeyB, parentDataB)
		host.Output().Finish(parentFinishA)
		host.Output().Finish(parentFinishB)

		err := host.Output().Transfer(thirdPartyAddress, host.Runtime().GetSCAddress(), 0, 0,
			big.NewInt(testConfig.transferToThirdParty), []byte("hello"), 0)
		require.Nil(t, err)

		arguments := host.Runtime().Arguments()

		callData := txDataBuilder.NewBuilder()
		// funcion to be called on child
		callData.Func("transferToThirdParty")
		// value to send to third party
		callData.Int64(testConfig.transferToThirdParty)
		// data for child -> third party tx
		callData.Str(" there")
		// behavior param for child
		callData.Bytes(arguments[0])

		// amount to transfer from parent to child
		value := big.NewInt(testConfig.transferFromParentToChild).Bytes()

		err = host.Runtime().ExecuteAsyncCall(childAddress, callData.ToBytes(), value)
		require.Nil(t, err)
	}, errors.New("breakpoint / failed to call function"))

	instanceMock.AddMockMethod("callBack", func() {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.gasUsedByCallback)

		if len(arguments) < 2 {
			return
		}

		loadedData := host.Storage().GetStorage(parentKeyB)

		status := bytes.Compare(loadedData, parentDataB)
		if status != 0 {
			status = 1
		}

		handleParentBehaviorArgument(host, arguments[1][0])
		err := handleTransferToVault(host, arguments)
		require.Nil(t, err)

		finishResult(host, status)
	})
}

func handleParentBehaviorArgument(host arwen.VMHost, behavior byte) {
	if behavior == 3 {
		host.Runtime().SignalUserError("callBack error")
	}
	if behavior == 4 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}

	host.Output().Finish([]byte{behavior})
}

func mustTransferToVault(arguments [][]byte) bool {
	vault := "vault"
	numArgs := len(arguments)
	if numArgs == 3 {
		if string(arguments[2]) == vault {
			return false
		}
	}

	if numArgs == 4 {
		if string(arguments[3]) == vault {
			return false
		}
	}

	return true
}

func handleTransferToVault(host arwen.VMHost, arguments [][]byte) error {
	err := error(nil)
	if mustTransferToVault(arguments) {
		valueToTransfer := big.NewInt(4)
		err = host.Output().Transfer(vaultAddress, host.Runtime().GetSCAddress(), 0, 0, valueToTransfer, arguments[1], 0)
	}

	return err
}

func finishResult(host arwen.VMHost, result int) {
	outputContext := host.Output()
	if result == 0 {
		outputContext.Finish([]byte("succ"))
	}
	if result == 1 {
		outputContext.Finish([]byte("fail"))
	}
	if result != 0 && result != 1 {
		outputContext.Finish([]byte("unkn"))
	}
}

func argumentsToHexString(functionName string, args ...[]byte) []byte {
	separator := byte('@')
	output := append([]byte(functionName))
	for _, arg := range args {
		output = append(output, separator)
		output = append(output, hex.EncodeToString(arg)...)
	}
	return output
}

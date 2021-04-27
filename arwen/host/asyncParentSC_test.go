package host

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/stretchr/testify/require"
)

func createTestAsyncParentContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock, testConfig *asyncCallTestConfig) {
	parentInstance := imb.CreateAndStoreInstanceMock(t, host, parentAddress, testConfig.parentBalance)
	addAsyncParentMethodsToInstanceMock(parentInstance, testConfig)
}

func addAsyncParentMethodsToInstanceMock(instance *mock.InstanceMock, testConfig *asyncCallTestConfig) {
	input := DefaultTestContractCallInput()
	input.GasProvided = testConfig.gasProvidedToChild

	t := instance.T

	instance.AddMockMethodWithError("performAsyncCall", func() {
		host := instance.Host
		host.Metering().UseGas(testConfig.gasUsedByParent)

		host.Storage().SetStorage(parentKeyA, parentDataA)
		host.Storage().SetStorage(parentKeyB, parentDataB)
		host.Output().Finish(parentFinishA)
		host.Output().Finish(parentFinishB)

		arguments := host.Runtime().Arguments()

		callData := argumentsToHexString(
			// funcion to be called on child
			"transferToThirdParty",
			// value to send to third party
			big.NewInt(testConfig.transferToThirdParty).Bytes(),
			// data for child -> third party tx
			[]byte(" there"),
			// behavior param for child
			[]byte{byte(big.NewInt(0).SetBytes(arguments[0]).Uint64() + '0')})

		// amount to transfer from parent to child
		value := big.NewInt(testConfig.transferFromParentToChild).Bytes()

		err := host.Runtime().ExecuteAsyncCall(childAddress, callData, value)
		require.Nil(t, err)
	}, errors.New("breakpoint / failed to call function"))

	handleBehaviorArgument := func(behavior byte) {
		host := instance.Host

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

	mustTransferToVault := func(arguments [][]byte) bool {
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

	handleTransferToVault := func(host arwen.VMHost, arguments [][]byte) {
		if mustTransferToVault(arguments) {
			valueToTransfer := big.NewInt(4)
			err := host.Output().Transfer(vaultAddress, host.Runtime().GetSCAddress(), 0, 0, valueToTransfer, arguments[1], 0)
			require.Nil(t, err)
		}
	}

	instance.AddMockMethod("callBack", func() {
		host := instance.Host
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.gasUsedByCallback)

		if len(arguments) < 2 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return
		}

		loadedData := host.Storage().GetStorage(parentKeyB)

		status := bytes.Compare(loadedData, parentDataB)
		if status != 0 {
			status = 1
		}

		handleBehaviorArgument(arguments[2][0])
		handleTransferToVault(host, arguments)

		finishResult(status, host)
	})
}

func finishResult(result int, host arwen.VMHost) {
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

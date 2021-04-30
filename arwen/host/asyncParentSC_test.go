package host

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/testscommon/txDataBuilder"
	"github.com/stretchr/testify/require"
)

func addAsyncParentMethodsToInstanceMock(instanceMock *mock.InstanceMock, config interface{}) {

	testConfig := config.(*asyncCallTestConfig)

	instanceMock.AddMockMethod("performAsyncCall", func() *mock.InstanceMock {
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
		callData.Bytes(append(arguments[0]))

		// amount to transfer from parent to child
		value := big.NewInt(testConfig.transferFromParentToChild).Bytes()

		err = host.Runtime().ExecuteAsyncCall(childAddress, callData.ToBytes(), value)
		require.Nil(t, err)

		return instance

	})

	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()

		host.Metering().UseGas(testConfig.gasUsedByCallback)

		if len(arguments) < 2 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		loadedData := host.Storage().GetStorage(parentKeyB)

		status := bytes.Compare(loadedData, parentDataB)
		if status != 0 {
			status = 1
		}

		if len(arguments) >= 4 {
			err := handleParentBehaviorArgument(host, big.NewInt(0).SetBytes(arguments[1]))
			if err != nil {
				return instance
			}
		}
		err := handleTransferToVault(host, arguments)
		require.Nil(t, err)

		finishResult(host, status)

		return instance
	})
}

func handleParentBehaviorArgument(host arwen.VMHost, behavior *big.Int) error {
	if behavior.Cmp(big.NewInt(3)) == 0 {
		host.Runtime().SignalUserError("callBack error")
		return errors.New("behavior / parent error")
	}
	if behavior.Cmp(big.NewInt(4)) == 0 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}

	behaviorBytes := behavior.Bytes()
	if len(behaviorBytes) == 0 {
		behaviorBytes = []byte{0}
	}
	host.Output().Finish(behaviorBytes)

	return nil
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

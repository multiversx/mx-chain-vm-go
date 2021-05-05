package host

import (
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/stretchr/testify/require"
)

func transferToThirdPartyAsyncChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*asyncCallTestConfig)
	instanceMock.AddMockMethod("transferToThirdParty", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T

		host.Metering().UseGas(testConfig.gasUsedByChild)

		arguments := host.Runtime().Arguments()
		outputContext := host.Output()

		if len(arguments) != 3 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		err := handleChildBehaviorArgument(host, big.NewInt(0).SetBytes(arguments[2]))
		if err != nil {
			return instance
		}

		scAddress := host.Runtime().GetSCAddress()
		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])
		err = outputContext.Transfer(
			thirdPartyAddress,
			scAddress,
			0,
			0,
			valueToTransfer,
			arguments[1],
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("thirdparty"))

		valueToTransfer = big.NewInt(testConfig.transferToVault)
		err = outputContext.Transfer(
			vaultAddress,
			scAddress,
			0,
			0,
			valueToTransfer,
			[]byte{},
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("vault"))

		host.Storage().SetStorage(childKey, childData)

		return instance
	})
}

func handleChildBehaviorArgument(host arwen.VMHost, behavior *big.Int) error {
	if behavior.Cmp(big.NewInt(1)) == 0 {
		host.Runtime().SignalUserError("child error")
		return errors.New("behavior / child error")
	}
	if behavior.Cmp(big.NewInt(2)) == 0 {
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

package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/elrondapi"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/stretchr/testify/require"
)

func DeployContractFromSourceMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("deployContractFromSource", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		t := instance.T

		arguments := host.Runtime().Arguments()

		if len(arguments) < 3 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		sourceContractAddress := arguments[0]
		codeMetadata := arguments[1]
		gasForInit := big.NewInt(0).SetBytes(arguments[2])

		newAddress, err :=
			elrondapi.DeployFromSourceContractWithTypedArgs(
				host,
				sourceContractAddress,
				codeMetadata,
				big.NewInt(0),
				[][]byte{},
				gasForInit.Int64(),
			)

		if err != nil {
			host.Runtime().FailExecution(err)
			return instance
		}

		require.NotNil(t, newAddress)

		host.Output().Finish(newAddress)

		return instance
	})
}

// InitMockMethod
func InitMockMethod(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*testcommon.TestConfig)
	instanceMock.AddMockMethod("init", testcommon.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByInit))
}

type CallbackTestConfig interface {
	CallbackFails() bool
}

// CallbackMockMethodThatCouldFail
func CallbackMockMethodThatCouldFail(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*testcommon.TestConfig)
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		if testConfig.CallbackFails {
			host.Runtime().SignalUserError("fail")
			return instance
		}
		return instance
	})
}

func UpdateContractFromSourceMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("updateContractFromSource", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		arguments := host.Runtime().Arguments()

		if len(arguments) < 4 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		sourceContractAddress := arguments[0]
		destinationContractAddress := arguments[1]
		codeMetadata := arguments[2]
		gasForInit := big.NewInt(0).SetBytes(arguments[3])

		elrondapi.UpgradeFromSourceContractWithTypedArgs(
			host,
			sourceContractAddress,
			destinationContractAddress,
			big.NewInt(0).Bytes(),
			[][]byte{},
			gasForInit.Int64(),
			codeMetadata,
		)

		return instance
	})
}
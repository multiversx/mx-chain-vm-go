package worldmock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/integrationTests/mock"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/builtInFunctions"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

type BuiltinFunctionsWrapper struct {
	Container process.BuiltInFunctionContainer
}

func NewBuiltinFunctionsWrapper(
	shardCoordinator sharding.Coordinator,
	accounts state.AccountsAdapter,
	gasMap config.GasScheduleMap,
) (*BuiltinFunctionsWrapper, error) {
	mapDNSAddresses := make(map[string]struct{})
	// if !check.IfNil(tpn.SmartContractParser) {
	// 	mapDNSAddresses, _ = tpn.SmartContractParser.GetDeployedSCAddresses(genesis.DNSType)
	// }

	// defaults.FillGasMapInternal(gasMap, 1)
	// gasSchedule := mock.NewGasScheduleNotifierMock(gasMap)
	gasSchedule := mock.NewGasScheduleNotifierMock(gasMap)
	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasSchedule:      gasSchedule,
		MapDNSAddresses:  mapDNSAddresses,
		Marshalizer:      nil,
		Accounts:         accounts,
		ShardCoordinator: shardCoordinator,
	}
	builtInFuncFactory, err := builtInFunctions.NewBuiltInFunctionsFactory(argsBuiltIn)
	if err != nil {
		return nil, err
	}

	builtInFuncs, err := builtInFuncFactory.CreateBuiltInFunctionContainer()
	if err != nil {
		return nil, err
	}

	builtinFuncsWrapper := &BuiltinFunctionsWrapper{
		Container: builtInFuncs,
	}

	return builtinFuncsWrapper, nil
}

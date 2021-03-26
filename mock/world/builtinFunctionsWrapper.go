package worldmock

import (
	"errors"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data/state"
	integrationTests "github.com/ElrondNetwork/elrond-go/integrationTests/mock"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/builtInFunctions"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

type BuiltinFunctionsWrapper struct {
	Container       process.BuiltInFunctionContainer
	MapDNSAddresses map[string]struct{}
}

func NewBuiltinFunctionsWrapper(
	shardCoordinator sharding.Coordinator,
	accounts state.AccountsAdapter,
	gasMap config.GasScheduleMap,
) (*BuiltinFunctionsWrapper, error) {

	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasSchedule:      integrationTests.NewGasScheduleNotifierMock(gasMap),
		MapDNSAddresses:  make(map[string]struct{}),
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
		Container:       builtInFuncs,
		MapDNSAddresses: argsBuiltIn.MapDNSAddresses,
	}

	return builtinFuncsWrapper, nil
}

func (bf *BuiltinFunctionsWrapper) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	// TODO implement
	return nil, errors.New("not implemented")
}

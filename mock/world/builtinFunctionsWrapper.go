package worldmock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data/state"
	integrationTests "github.com/ElrondNetwork/elrond-go/integrationTests/mock"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/builtInFunctions"
)

type BuiltinFunctionsWrapper struct {
	Container       process.BuiltInFunctionContainer
	MapDNSAddresses map[string]struct{}
	World           *MockWorld
	Marshalizer     marshal.Marshalizer
}

func NewBuiltinFunctionsWrapper(
	world *MockWorld,
	gasMap config.GasScheduleMap,
) (*BuiltinFunctionsWrapper, error) {
	marshalizer := &marshal.GogoProtoMarshalizer{}

	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasSchedule:      integrationTests.NewGasScheduleNotifierMock(gasMap),
		MapDNSAddresses:  make(map[string]struct{}),
		Marshalizer:      marshalizer,
		Accounts:         NewMockAccountsAdapter(world.AcctMap),
		ShardCoordinator: world,
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
		World:           world,
		Marshalizer:     marshalizer,
	}

	return builtinFuncsWrapper, nil
}

func (bf *BuiltinFunctionsWrapper) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	caller := bf.getAccountSharded(input.CallerAddr)
	recipient := bf.getAccountSharded(input.RecipientAddr)

	function, err := bf.Container.Get(input.Function)
	if err != nil {
		return nil, err
	}

	return function.ProcessBuiltinFunction(caller, recipient, input)
}

func (bf *BuiltinFunctionsWrapper) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return bf.Container.Keys()
}

// TODO change AccountMap to support this instead
func (bf *BuiltinFunctionsWrapper) getAccountSharded(address []byte) state.UserAccountHandler {
	accountShard := bf.World.ComputeId(address)
	if accountShard != bf.World.SelfId() {
		return nil
	}
	return bf.World.AcctMap.GetAccount(address)
}

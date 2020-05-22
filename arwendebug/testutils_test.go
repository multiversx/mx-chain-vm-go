package arwendebug

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

type testContext struct {
	t       *testing.T
	worldID string
	facade  *DebugFacade
}

func newTestContext(t *testing.T) *testContext {
	worldID := fmt.Sprintf("%s_%d", time.Now().Format("20060102150405"), rand.Intn(100))

	return &testContext{
		t:       t,
		worldID: worldID,
		facade:  &DebugFacade{},
	}
}

func (context *testContext) createAccount(address string, balance string) {
	request := CreateAccountRequest{
		RequestBase: context.createRequestBase(),
		Address:     address,
		Balance:     balance,
		Nonce:       0,
	}

	response, err := context.facade.CreateAccount(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
}

func (context *testContext) deployContract(codePath string, impersonated string, arguments ...string) *DeployResponse {
	request := DeployRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:  context.createRequestBase(),
			Impersonated: impersonated,
		},
		CodePath:  codePath,
		Arguments: arguments,
	}

	response, err := context.facade.DeploySmartContract(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Nil(t, response.Error)
	require.Equal(t, vmcommon.Ok, response.Output.ReturnCode)

	return response
}

func (context *testContext) runContract(contract string, impersonated string, function string, arguments ...string) {
	request := RunRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:  context.createRequestBase(),
			Impersonated: impersonated,
		},
		ContractAddress: contract,
		Function:        function,
		Arguments:       arguments,
	}

	response, err := context.facade.RunSmartContract(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Nil(t, response.Error)
	require.Equal(t, vmcommon.Ok, response.Output.ReturnCode)
}

func (context *testContext) queryContract(contract string, impersonated string, function string, arguments ...string) *QueryResponse {
	request := QueryRequest{
		RunRequest: RunRequest{
			ContractRequestBase: ContractRequestBase{
				RequestBase:  context.createRequestBase(),
				Impersonated: impersonated,
			},
			ContractAddress: contract,
			Function:        function,
			Arguments:       arguments,
		},
	}

	response, err := context.facade.QuerySmartContract(request)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Nil(t, response.Error)
	require.Equal(t, vmcommon.Ok, response.Output.ReturnCode)

	return response
}

func (response *ContractResponseBase) getFirstResultAsInt64() int64 {
	result, err := response.Output.GetFirstReturnData(vmcommon.AsBigInt)
	if err != nil {
		return 0
	}

	asBigInt := result.(*big.Int)
	return asBigInt.Int64()
}

func (context *testContext) createRequestBase() RequestBase {
	randomOutcome := fmt.Sprintf("%s_%d", time.Now().Format("20060102150405"), rand.Intn(100))

	return RequestBase{
		DatabasePath: databasePath,
		World:        context.worldID,
		Outcome:      randomOutcome,
	}
}

func (context *testContext) loadWorld() *world {
	database := NewDatabase(databasePath)
	world, err := database.loadWorld(context.worldID)
	require.Nil(context.t, err)

	return world
}

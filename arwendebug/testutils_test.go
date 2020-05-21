package arwendebug

import (
	"fmt"
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
	createAccountRequest := CreateAccountRequest{
		RequestBase: context.createRequestBase(),
		Address:     address,
		Balance:     balance,
		Nonce:       0,
	}

	createAccountResponse, err := context.facade.CreateAccount(createAccountRequest)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, createAccountResponse)
}

func (context *testContext) deployContract(codePath string, impersonated string, arguments ...string) {
	deployRequest := DeployRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:  context.createRequestBase(),
			Impersonated: impersonated,
		},
		CodePath:  codePath,
		Arguments: arguments,
	}

	deployResponse, err := context.facade.DeploySmartContract(deployRequest)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, deployResponse)
	require.Nil(t, deployResponse.Error)
	require.Equal(t, vmcommon.Ok, deployResponse.Output.ReturnCode)
}

func (context *testContext) runContract(contract string, impersonated string, function string, arguments ...string) {
	runRequest := RunRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:  context.createRequestBase(),
			Impersonated: impersonated,
		},
		ContractAddress: contract,
		Function:        function,
		Arguments:       arguments,
	}

	runResponse, err := context.facade.RunSmartContract(runRequest)

	t := context.t
	require.Nil(t, err)
	require.NotNil(t, runResponse)
	require.Nil(t, runResponse.Error)
	require.Equal(t, vmcommon.Ok, runResponse.Output.ReturnCode)
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

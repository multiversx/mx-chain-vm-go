package arwendebug

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var databasePath = "./testdata/db"

func init() {
	os.RemoveAll(databasePath)
}

func TestFacade_CreateAccount(t *testing.T) {
	context := newTestContext(t)
	context.createAccount("alice", "42")

	world := context.loadWorld()
	exists, err := world.blockchainHook.AccountExists(fixTestAddress("alice"))
	require.Nil(t, err)
	require.True(t, exists)
}

func TestFacade_RunContract_Counter(t *testing.T) {
	context := newTestContext(t)

	counterKey := string([]byte{'m', 'y', 'c', 'o', 'u', 'n', 't', 'e', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	context.createAccount("alice", "42")
	context.deployContract("../test/contracts/counter/counter.wasm", "alice")
	context.runContract("contract0000000000000000000alice", "alice", "increment")

	world := context.loadWorld()
	exists, err := world.blockchainHook.AccountExists([]byte("contract0000000000000000000alice"))
	require.Nil(t, err)
	require.True(t, exists)

	world = context.loadWorld()
	state, err := world.blockchainHook.GetAllState([]byte("contract0000000000000000000alice"))
	require.Nil(t, err)
	require.NotNil(t, state)
	require.Equal(t, []byte{2}, state[counterKey])
}

func TestFacade_RunContract_ERC20(t *testing.T) {
	context := newTestContext(t)

	context.createAccount("alice", "42")
	context.createAccount("bob", "40")
	context.createAccount("carol", "30")
	context.deployContract("../test/contracts/erc20/erc20.wasm", "alice", "64")
}

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

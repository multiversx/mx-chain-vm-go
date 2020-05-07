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

func init() {
	os.RemoveAll("./testdata/db")
}

func Test_CreateAccount(t *testing.T) {
	facade := &DebugFacade{}
	request := CreateAccountRequest{
		RequestBase: createRequestBase(),
		Address:     "alice",
		Balance:     "42",
		Nonce:       1,
	}

	response, err := facade.CreateAccount(request)
	require.Nil(t, err)
	require.NotNil(t, response)

	database := NewDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	require.Nil(t, err)

	exists, err := world.blockchainHook.AccountExists(fixTestAddress("alice"))
	require.Nil(t, err)
	require.True(t, exists)
}

func Test_RunContract(t *testing.T) {
	facade := &DebugFacade{}
	requestBase := createRequestBase()

	worldID := requestBase.World
	databasePath := requestBase.DatabasePath

	database := NewDatabase(databasePath)
	world, err := database.loadWorld(worldID)
	require.Nil(t, err)

	counterKey := string([]byte{'m', 'y', 'c', 'o', 'u', 'n', 't', 'e', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	createAccountRequest := CreateAccountRequest{
		RequestBase: requestBase,
		Address:     "alice",
		Balance:     "42",
		Nonce:       1,
	}

	deployRequest := DeployRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:  requestBase,
			Impersonated: "alice",
		},
		CodePath: "../test/contracts/counter/counter.wasm",
	}

	runRequest := RunRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:  requestBase,
			Impersonated: "alice",
		},
		ContractAddress: "contract0000000000000000000alice",
		Function:        "increment",
	}

	createAccountResponse, err := facade.CreateAccount(createAccountRequest)
	require.Nil(t, err)
	require.NotNil(t, createAccountResponse)

	deployResponse, err := facade.DeploySmartContract(deployRequest)
	require.Nil(t, err)
	require.NotNil(t, deployResponse)
	require.Nil(t, deployResponse.Error)
	require.Equal(t, vmcommon.Ok, deployResponse.Output.ReturnCode)

	world, _ = database.loadWorld(worldID)
	exists, err := world.blockchainHook.AccountExists([]byte("contract0000000000000000000alice"))
	require.Nil(t, err)
	require.True(t, exists)

	runResponse, err := facade.RunSmartContract(runRequest)
	require.Nil(t, err)
	require.NotNil(t, runResponse)
	require.Nil(t, runResponse.Error)
	require.Equal(t, vmcommon.Ok, runResponse.Output.ReturnCode)

	world, _ = database.loadWorld(worldID)
	state, err := world.blockchainHook.GetAllState([]byte("contract0000000000000000000alice"))
	require.Nil(t, err)
	require.NotNil(t, state)
	require.Equal(t, []byte{2}, state[counterKey])
}

func createRequestBase() RequestBase {
	randomWorld := fmt.Sprintf("%s_%d", time.Now().Format("20060102150405"), rand.Intn(100))
	randomOutcome := fmt.Sprintf("%s_%d", time.Now().Format("20060102150405"), rand.Intn(100))

	return RequestBase{
		DatabasePath: "./testdata/db",
		World:        randomWorld,
		Outcome:      randomOutcome,
	}
}

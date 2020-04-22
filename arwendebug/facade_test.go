package arwendebug

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	os.RemoveAll("./testdata/db")
}

func Test_CreateAccount(t *testing.T) {
	facade := &DebugFacade{}
	request := CreateAccountRequest{
		RequestBase: createRequestBase(),
		Address:     "erd1alice",
		Balance:     "42",
		Nonce:       1,
	}

	response, err := facade.CreateAccount(request)
	require.Nil(t, err)
	require.NotNil(t, response)

	database := NewDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	require.Nil(t, err)

	exists, err := world.blockchainHook.AccountExists([]byte("erd1alice"))
	require.Nil(t, err)
	require.True(t, exists)
}

func Test_DeployContract(t *testing.T) {
	facade := &DebugFacade{}

	createAccountRequest := CreateAccountRequest{
		RequestBase: createRequestBase(),
		Address:     "erd1alice",
		Balance:     "42",
		Nonce:       1,
	}

	deployRequest := DeployRequest{
		ContractRequestBase: ContractRequestBase{
			RequestBase:  createRequestBase(),
			Impersonated: "erd1alice",
		},
		Code: "x",
	}

	createAccountResponse, err := facade.CreateAccount(createAccountRequest)
	require.Nil(t, err)
	require.NotNil(t, createAccountResponse)

	deployResponse, err := facade.DeploySmartContract(deployRequest)
	require.Nil(t, err)
	require.NotNil(t, deployResponse)

	// database := NewDatabase(deployRequest.DatabasePath)
	// world, err := database.loadWorld(deployRequest.World)
	// require.Nil(t, err)

	// exists, err := world.blockchainHook.AccountExists([]byte("erd1alice"))
	// require.Nil(t, err)
	// require.True(t, exists)
}

func createRequestBase() RequestBase {
	return RequestBase{
		DatabasePath: "./testdata/db",
		World:        time.Now().Format("20060102150405"),
	}
}

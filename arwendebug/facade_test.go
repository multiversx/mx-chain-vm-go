package arwendebug

import (
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

func Test_DeployContract(t *testing.T) {
	facade := &DebugFacade{}
	requestBase := createRequestBase()

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
		CodePath: "./testdata/counter.wasm",
	}

	createAccountResponse, err := facade.CreateAccount(createAccountRequest)
	require.Nil(t, err)
	require.NotNil(t, createAccountResponse)

	deployResponse, err := facade.DeploySmartContract(deployRequest)
	require.Nil(t, err)
	require.NotNil(t, deployResponse)
	require.Nil(t, deployResponse.Error)
	require.Equal(t, vmcommon.Ok, deployResponse.Output.ReturnCode)

	database := NewDatabase(deployRequest.DatabasePath)
	world, err := database.loadWorld(deployRequest.World)
	require.Nil(t, err)

	exists, err := world.blockchainHook.AccountExists([]byte("contract0000000000000000000alice"))
	require.Nil(t, err)
	require.True(t, exists)
}

func Test_RunContract(t *testing.T) {
	facade := &DebugFacade{}
	requestBase := createRequestBase()

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
		CodePath: "./testdata/counter.wasm",
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

	runResponse, err := facade.RunSmartContract(runRequest)
	require.Nil(t, err)
	require.NotNil(t, runResponse)
	require.Nil(t, runResponse.Error)
	require.Equal(t, vmcommon.Ok, runResponse.Output.ReturnCode)
}

func createRequestBase() RequestBase {
	return RequestBase{
		DatabasePath: "./testdata/db",
		World:        time.Now().Format("20060102150405"),
	}
}

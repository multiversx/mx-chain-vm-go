package arwendebug

import (
	"os"
	"testing"

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

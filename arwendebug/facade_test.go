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
	deployResponse := context.deployContract("../test/contracts/counter/counter.wasm", "alice")
	contractAddress := deployResponse.ContractAddress
	require.Equal(t, "contract0000000000000000000alice", contractAddress)
	context.runContract(contractAddress, "alice", "increment")

	world := context.loadWorld()
	exists, err := world.blockchainHook.AccountExists([]byte(contractAddress))
	require.Nil(t, err)
	require.True(t, exists)

	world = context.loadWorld()
	state, err := world.blockchainHook.GetAllState([]byte(contractAddress))
	require.Nil(t, err)
	require.NotNil(t, state)
	require.Equal(t, []byte{2}, state[counterKey])
}

func TestFacade_RunContract_ERC20(t *testing.T) {
	context := newTestContext(t)

	aliceHex := "303030303030303030303030303030303030303030303030303030616c696365"
	bobHex := "3030303030303030303030303030303030303030303030303030303030626f62"
	carolHex := "3030303030303030303030303030303030303030303030303030306361726f6c"
	context.createAccount("alice", "42")
	context.createAccount("bob", "40")
	context.createAccount("carol", "30")
	deployResponse := context.deployContract("../test/contracts/erc20/erc20.wasm", "alice", "64")
	contractAddress := deployResponse.ContractAddress
	require.Equal(t, "contract0000000000000000000alice", contractAddress)

	// Initial state
	totalSupply := context.queryContract(contractAddress, "alice", "totalSupply").getFirstResultAsInt64()
	balanceOfAlice := context.queryContract(contractAddress, "alice", "balanceOf", aliceHex).getFirstResultAsInt64()
	require.Equal(t, int64(100), totalSupply)
	require.Equal(t, int64(100), balanceOfAlice)

	// Transfers
	context.runContract(contractAddress, "alice", "transferToken", aliceHex, "0A")
	context.runContract(contractAddress, "alice", "transferToken", bobHex, "0A")
	context.runContract(contractAddress, "alice", "transferToken", carolHex, "0A")
	context.runContract(contractAddress, "bob", "transferToken", carolHex, "05")

	balanceOfAlice = context.queryContract(contractAddress, "alice", "balanceOf", aliceHex).getFirstResultAsInt64()
	balanceOfBob := context.queryContract(contractAddress, "alice", "balanceOf", bobHex).getFirstResultAsInt64()
	balanceOfCarol := context.queryContract(contractAddress, "alice", "balanceOf", carolHex).getFirstResultAsInt64()
	require.Equal(t, int64(80), balanceOfAlice)
	require.Equal(t, int64(5), balanceOfBob)
	require.Equal(t, int64(15), balanceOfCarol)
}

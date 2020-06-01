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
	context.createAccount(newDummyAddress("alice").hex, "42")

	world := context.loadWorld()
	exists, err := world.blockchainHook.AccountExists(newDummyAddress("alice").raw)
	require.Nil(t, err)
	require.True(t, exists)
}

func TestFacade_RunContract_Counter(t *testing.T) {
	context := newTestContext(t)

	counterKey := string([]byte{'m', 'y', 'c', 'o', 'u', 'n', 't', 'e', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	alice := newDummyAddress("alice")
	context.createAccount(alice.hex, "42")
	deployResponse := context.deployContract("../test/contracts/counter/counter.wasm", alice.hex)
	contractAddress := deployResponse.ContractAddress
	require.Equal(t, "contract0000000000000000000alice", contractAddress)
	context.runContract(contractAddress, alice.hex, "increment")

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

	alice := newDummyAddress("alice")
	bob := newDummyAddress("bob")
	carol := newDummyAddress("carol")
	context.createAccount(alice.hex, "42")
	context.createAccount(bob.hex, "40")
	context.createAccount(carol.hex, "30")
	deployResponse := context.deployContract("../test/contracts/erc20/erc20.wasm", alice.hex, "64")
	contractAddress := deployResponse.ContractAddress
	require.Equal(t, "contract0000000000000000000alice", contractAddress)

	// Initial state
	totalSupply := context.queryContract(contractAddress, alice.hex, "totalSupply").getFirstResultAsInt64()
	balanceOfAlice := context.queryContract(contractAddress, alice.hex, "balanceOf", alice.hex).getFirstResultAsInt64()
	balanceOfBob := context.queryContract(contractAddress, alice.hex, "balanceOf", bob.hex).getFirstResultAsInt64()
	balanceOfCarol := context.queryContract(contractAddress, alice.hex, "balanceOf", carol.hex).getFirstResultAsInt64()
	require.Equal(t, int64(100), totalSupply)
	require.Equal(t, int64(100), balanceOfAlice)
	require.Equal(t, int64(0), balanceOfBob)
	require.Equal(t, int64(0), balanceOfCarol)

	// Transfers
	context.runContract(contractAddress, alice.hex, "transferToken", alice.hex, "0A")
	context.runContract(contractAddress, alice.hex, "transferToken", bob.hex, "0A")
	context.runContract(contractAddress, alice.hex, "transferToken", carol.hex, "0A")
	context.runContract(contractAddress, bob.hex, "transferToken", carol.hex, "05")

	balanceOfAlice = context.queryContract(contractAddress, alice.hex, "balanceOf", alice.hex).getFirstResultAsInt64()
	balanceOfBob = context.queryContract(contractAddress, alice.hex, "balanceOf", bob.hex).getFirstResultAsInt64()
	balanceOfCarol = context.queryContract(contractAddress, alice.hex, "balanceOf", carol.hex).getFirstResultAsInt64()
	require.Equal(t, int64(80), balanceOfAlice)
	require.Equal(t, int64(5), balanceOfBob)
	require.Equal(t, int64(15), balanceOfCarol)
}

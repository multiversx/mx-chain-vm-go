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

	context.createAccount("alice", "42")
	context.createAccount("bob", "40")
	context.createAccount("carol", "30")
	deployResponse := context.deployContract("../test/contracts/erc20/erc20.wasm", "alice", "64")
	contractAddress := deployResponse.ContractAddress
	require.Equal(t, "contract0000000000000000000alice", contractAddress)

	totalSupply := context.queryContract(contractAddress, "alice", "totalSupply").getFirstResultAsInt64()
	balanceOfAlice := context.queryContract(contractAddress, "alice", "balanceOf", "alicehex").getFirstResultAsInt64()
	require.Equal(t, int64(100), totalSupply)
	require.Equal(t, int64(100), balanceOfAlice)

}

// err := context.DeploySC("../testdata/erc20-c-03/wrc20_arwen.wasm", "00"+arwen.FormatHexNumber(42000))
// require.Nil(t, err)

// // Assertion
// require.Equal(t, uint64(42000), context.QuerySCInt("totalSupply", [][]byte{}))
// require.Equal(t, uint64(42000), context.QuerySCInt("balanceOf", [][]byte{context.Owner.Address}))

// // Minting
// err = context.ExecuteSC(owner, "transferToken@"+alice.AddressHex()+"@00"+arwen.FormatHexNumber(1000))
// require.Nil(t, err)
// err = context.ExecuteSC(owner, "transferToken@"+bob.AddressHex()+"@00"+arwen.FormatHexNumber(1000))
// require.Nil(t, err)

// // Regular transfers
// err = context.ExecuteSC(alice, "transferToken@"+bob.AddressHex()+"@00"+arwen.FormatHexNumber(200))
// require.Nil(t, err)
// err = context.ExecuteSC(bob, "transferToken@"+alice.AddressHex()+"@00"+arwen.FormatHexNumber(400))
// require.Nil(t, err)

// // Assertion
// require.Equal(t, uint64(1200), context.QuerySCInt("balanceOf", [][]byte{alice.Address}))
// require.Equal(t, uint64(800), context.QuerySCInt("balanceOf", [][]byte{bob.Address}))

// // Approve and transfer
// err = context.ExecuteSC(alice, "approve@"+bob.AddressHex()+"@00"+arwen.FormatHexNumber(500))
// require.Nil(t, err)
// err = context.ExecuteSC(bob, "approve@"+alice.AddressHex()+"@00"+arwen.FormatHexNumber(500))
// require.Nil(t, err)
// err = context.ExecuteSC(alice, "transferFrom@"+bob.AddressHex()+"@"+carol.AddressHex()+"@00"+arwen.FormatHexNumber(25))
// require.Nil(t, err)
// err = context.ExecuteSC(bob, "transferFrom@"+alice.AddressHex()+"@"+carol.AddressHex()+"@00"+arwen.FormatHexNumber(25))
// require.Nil(t, err)

// require.Equal(t, uint64(1175), context.QuerySCInt("balanceOf", [][]byte{alice.Address}))
// require.Equal(t, uint64(775), context.QuerySCInt("balanceOf", [][]byte{bob.Address}))
// require.Equal(t, uint64(50), context.QuerySCInt("balanceOf", [][]byte{carol.Address}))

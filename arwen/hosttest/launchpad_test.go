package hosttest

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/esdtconvert"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var launchpad = test.MakeTestSCAddress("launchpad")
var esdtTokenKeyPrefix = []byte(core.ElrondProtectedKeyPrefix + core.ESDTKeyIdentifier)
var tokenID = []byte("DEBEN-000000")
var tokenIDAlt = []byte("WHAT-999999")
var egldID = []byte("")

var confirmationStartEpoch = 1
var winnerSelectionStartEpoch = 2
var claimStartEpoch = 3
var tokensPerWinningTicket = 10
var ticketPrice = 4
var nrWinningTickets = 30
var requiredTokenAmount = int(tokensPerWinningTicket) * int(nrWinningTickets)

func TestLaunchpadInit(t *testing.T) {
	_, _ = setupLaunchpad(t)
}

func setupLaunchpad(t *testing.T) (arwen.VMHost, *worldmock.MockWorld) {
	launchpadCode := test.GetTestSCCode("launchpad", "../../")
	host, world := test.DefaultTestArwenWithWorldMock(t)
	ownerAccount := world.AcctMap.CreateAccount(owner, world)
	ownerAccount.Balance = big.NewInt(1000)
	ownerESDT := makeDefaultOwnerESDT()
	esdtconvert.WriteMockESDTToStorage(ownerESDT, ownerAccount.Storage)

	deployInput := test.CreateTestContractCreateInputBuilder().
		WithCallerAddr(owner).
		WithContractCode(launchpadCode).
		WithGasProvided(1000000).
		WithArguments(
			tokenID,
			makeBytes(tokensPerWinningTicket),
			egldID,
			makeBytes(ticketPrice),
			makeBytes(nrWinningTickets),
			makeBytes(confirmationStartEpoch),
			makeBytes(winnerSelectionStartEpoch),
			makeBytes(claimStartEpoch),
		).
		Build()

	launchpadAccount := world.AcctMap.CreateSmartContractAccount(owner, launchpad, launchpadCode, world)
	_ = runSCCreate(t, host, world, deployInput, owner, launchpad)

	launchpadESDT := makeDefaultLaunchpadESDT()
	esdtconvert.WriteMockESDTToStorage(launchpadESDT, launchpadAccount.Storage)

	depositInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(owner).
		WithGasProvided(1000000).
		WithFunction("depositLaunchpadTokens").
		WithRecipientAddr(launchpad).
		Build()

	_ = runSCCall(t, host, world, depositInput)

	return host, world
}

func runSCCreate(
	t *testing.T,
	host arwen.VMHost,
	world *worldmock.MockWorld,
	input *vmcommon.ContractCreateInput,
	owner []byte,
	contractAddress []byte,
) *vmcommon.VMOutput {

	addressMock := &worldmock.NewAddressMock{
		CreatorAddress: owner,
		CreatorNonce:   0,
		NewAddress:     launchpad,
	}
	world.NewAddressMocks = append(world.NewAddressMocks, addressMock)

	vmOutput, err := host.RunSmartContractCreate(input)
	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok()

	err = world.UpdateAccounts(vmOutput.OutputAccounts, vmOutput.DeletedAccounts)
	require.Nil(t, err)

	return vmOutput
}

func runSCCall(
	t *testing.T,
	host arwen.VMHost,
	world *worldmock.MockWorld,
	input *vmcommon.ContractCallInput,
) *vmcommon.VMOutput {

	vmOutput, err := host.RunSmartContractCall(input)
	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok()

	err = world.UpdateAccounts(vmOutput.OutputAccounts, vmOutput.DeletedAccounts)
	require.Nil(t, err)

	return vmOutput
}

func makeDefaultOwnerESDT() []*esdtconvert.MockESDTData {
	return []*esdtconvert.MockESDTData{
		{
			TokenIdentifier: tokenID,
			Instances:       []*esdt.ESDigitalToken{makeESDTInstance(0, 100)},
		},
		{
			TokenIdentifier: tokenIDAlt,
			Instances:       []*esdt.ESDigitalToken{makeESDTInstance(0, 100)},
		},
	}
}

func makeDefaultLaunchpadESDT() []*esdtconvert.MockESDTData {
	return []*esdtconvert.MockESDTData{
		{
			TokenIdentifier: tokenID,
			Instances:       []*esdt.ESDigitalToken{makeESDTInstance(0, requiredTokenAmount)},
		},
	}
}

func makeTokenKey(tokenName []byte, nonce uint64) []byte {
	nonceBytes := big.NewInt(0).SetUint64(nonce).Bytes()
	tokenKey := append(esdtTokenKeyPrefix, tokenName...)
	tokenKey = append(tokenKey, nonceBytes...)
	return tokenKey
}

func makeESDTInstance(nonce uint64, value int) *esdt.ESDigitalToken {
	return &esdt.ESDigitalToken{
		Type:  uint32(core.Fungible),
		Value: big.NewInt(int64(value)),
		TokenMetaData: &esdt.MetaData{
			Nonce: nonce,
		},
	}
}

func makeBytes(a int) []byte {
	return []byte{byte(a)}
}

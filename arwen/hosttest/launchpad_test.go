package hosttest

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/esdtconvert"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-go-core/core"
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

func TestLaunchpadInit(t *testing.T) {
	_, _ = setupLaunchpad(t)
}

func setupLaunchpad(t *testing.T) (arwen.VMHost, *worldmock.MockWorld) {
	launchpadCode := test.GetTestSCCode("launchpad", "../../")
	host, world := test.DefaultTestArwenWithWorldMock(t)
	world.NewAddressMocks = []*worldmock.NewAddressMock{
		{
			CreatorAddress: owner,
			CreatorNonce:   0,
			NewAddress:     launchpad,
		},
	}

	ownerAccount := world.AcctMap.CreateAccount(owner, world)
	ownerAccount.Balance = big.NewInt(1000)
	ownerESDT := []*mj.ESDTData{
		{
			TokenIdentifier: mj.JSONBytesFromString{Value: tokenID},
			Instances: []*mj.ESDTInstance{
				{Nonce: mj.JSONUint64{Value: 0}, Balance: mj.JSONBigInt{Value: big.NewInt(100)}},
			},
		},
		{
			TokenIdentifier: mj.JSONBytesFromString{Value: tokenIDAlt},
			Instances: []*mj.ESDTInstance{
				{Nonce: mj.JSONUint64{Value: 0}, Balance: mj.JSONBigInt{Value: big.NewInt(100)}},
			},
		},
	}
	esdtconvert.WriteMandosESDTToStorage(ownerESDT, ownerAccount.Storage)

	input := test.CreateTestContractCreateInputBuilder().
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
	vmOutput, err := host.RunSmartContractCreate(input)
	verify := test.NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok()

	err = world.UpdateAccounts(vmOutput.OutputAccounts, vmOutput.DeletedAccounts)
	require.Nil(t, err)

	requiredTokenAmount := int64(tokensPerWinningTicket) * int64(nrWinningTickets)

	launchpadESDT := []*mj.ESDTData{
		{
			TokenIdentifier: mj.JSONBytesFromString{Value: tokenID},
			Instances: []*mj.ESDTInstance{
				{Nonce: mj.JSONUint64{Value: 0}, Balance: mj.JSONBigInt{Value: big.NewInt(requiredTokenAmount)}},
			},
		},
	}
	esdtconvert.WriteMandosESDTToStorage(launchpadESDT, launchpadAccount.Storage)

	depositInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(owner).
		WithGasProvided(1000000).
		WithFunction("depositLaunchpadTokens").
		WithRecipientAddr(launchpad).
		Build()

	vmOutput, err = host.RunSmartContractCall(depositInput)
	verify = test.NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok()

	return host, world
}

func makeTokenKey(tokenName []byte, nonce uint64) []byte {
	nonceBytes := big.NewInt(0).SetUint64(nonce).Bytes()
	tokenKey := append(esdtTokenKeyPrefix, tokenName...)
	tokenKey = append(tokenKey, nonceBytes...)
	return tokenKey
}

func makeBytes(a int) []byte {
	return []byte{byte(a)}
}

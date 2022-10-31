package hosttest

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
	"github.com/stretchr/testify/require"
)

func Test_RunDEXPairBenchmark(t *testing.T) {
	setupMEXPair(t)
}

func setupMEXPair(t *testing.T) {
	owner := arwen.MakeTestWalletAddress("owner")
	world, ownerAccount, host, err := prepare(t, owner)
	require.Nil(t, err)

	userAddress := arwen.MakeTestWalletAddress("user")
	userAccount := world.AcctMap.CreateAccount(userAddress, world)

	mex := NewMEXSetup(t, host, world, ownerAccount, userAccount)
	mex.Deploy()

	mex.ApplyInitialSetup()

	mex.AddLiquidity(
		userAddress,
		mex.UserWEGLDBalance,
		1,
		mex.UserMEXBalance,
		1)
}

type MEXSetup struct {
	WEGLDToken               []byte
	MEXToken                 []byte
	LPToken                  []byte
	OwnerAccount             *worldmock.Account
	OwnerAddress             Address
	RouterAddress            Address
	PairAddress              Address
	TotalFeePercent          uint64
	SpecialFeePercent        uint64
	MaxObservationsPerRecord int
	Code                     []byte
	UserAccount              *worldmock.Account
	UserWEGLDBalance         uint64
	UserMEXBalance           uint64

	T     *testing.T
	Host  arwen.VMHost
	World *worldmock.MockWorld
}

func NewMEXSetup(
	t *testing.T,
	host arwen.VMHost,
	world *worldmock.MockWorld,
	ownerAccount *worldmock.Account,
	userAccount *worldmock.Account,
) *MEXSetup {
	return &MEXSetup{
		WEGLDToken:               []byte("WEGLD-abcdef"),
		MEXToken:                 []byte("MEX-abcdef"),
		LPToken:                  []byte("LPTOK-abcdef"),
		OwnerAccount:             ownerAccount,
		OwnerAddress:             ownerAccount.Address,
		RouterAddress:            ownerAccount.Address,
		PairAddress:              test.MakeTestSCAddress("pairSC"),
		TotalFeePercent:          300,
		SpecialFeePercent:        50,
		MaxObservationsPerRecord: 10,
		Code:                     test.GetTestSCCode("pair", "../../"),
		UserAccount:              userAccount,
		UserWEGLDBalance:         5_000_000_000,
		UserMEXBalance:           5_000_000_000,

		T:     t,
		Host:  host,
		World: world,
	}
}

func (mex *MEXSetup) Deploy() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCreateInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithContractCode(mex.Code).
		WithArguments(
			mex.WEGLDToken,
			mex.MEXToken,
			mex.OwnerAddress,
			mex.RouterAddress,
			big.NewInt(int64(mex.TotalFeePercent)).Bytes(),
			big.NewInt(int64(mex.SpecialFeePercent)).Bytes(),
		).
		WithGasProvided(0xFFFFFFFFFFFFFFFF).
		Build()

	world.NewAddressMocks = append(world.NewAddressMocks, &worldmock.NewAddressMock{
		CreatorAddress: mex.OwnerAddress,
		CreatorNonce:   mex.OwnerAccount.Nonce,
		NewAddress:     mex.PairAddress,
	})

	mex.OwnerAccount.Nonce++ // nonce increases before deploy
	vmOutput, err := host.RunSmartContractCreate(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) ApplyInitialSetup() {
	mex.setLPToken()
	mex.setActiveState()
	mex.setMaxObservationsPerRecord()
	mex.setRequiredTokenRoles()
	mex.setESDTBalances()
}

func (mex *MEXSetup) setLPToken() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("setLpTokenIdentifier").
		WithArguments(mex.LPToken).
		WithGasProvided(0xFFFFFFFFFFFFFFFF).
		Build()

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) setActiveState() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("resume").
		WithGasProvided(0xFFFFFFFFFFFFFFFF).
		Build()

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) setMaxObservationsPerRecord() {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInput := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.OwnerAddress).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("setMaxObservationsPerRecord").
		WithArguments(big.NewInt(int64(mex.MaxObservationsPerRecord)).Bytes()).
		WithGasProvided(0xFFFFFFFFFFFFFFFF).
		Build()

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) setRequiredTokenRoles() {
	world := mex.World
	pairAccount := world.AcctMap.GetAccount(mex.PairAddress)

	roles := []string{core.ESDTRoleLocalMint, core.ESDTRoleLocalBurn}
	pairAccount.SetTokenRolesAsStrings(mex.LPToken, roles)
}

func (mex *MEXSetup) setESDTBalances() {
	mex.UserAccount.SetTokenBalanceUint64(mex.WEGLDToken, 0, mex.UserWEGLDBalance)
	mex.UserAccount.SetTokenBalanceUint64(mex.MEXToken, 0, mex.UserMEXBalance)
}

func (mex *MEXSetup) AddLiquidity(
	userAddress Address,
	WEGLDAmount uint64,
	minWEGLDAmount uint64,
	MEXAmount uint64,
	minMEXAmount uint64,
) {
	t := mex.T
	host := mex.Host
	world := mex.World

	vmInputBuiler := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.UserAccount.Address).
		WithRecipientAddr(mex.PairAddress).
		WithFunction("addLiquidity").
		WithArguments(
			big.NewInt(int64(minWEGLDAmount)).Bytes(),
			big.NewInt(int64(minMEXAmount)).Bytes(),
		).
		WithGasProvided(0xFFFFFFFFFFFFFFFF)

	vmInputBuiler.
		WithESDTTokenName(mex.WEGLDToken).
		WithESDTValue(big.NewInt(int64(WEGLDAmount))).
		NextESDTTransfer().
		WithESDTTokenName(mex.MEXToken).
		WithESDTValue(big.NewInt(int64(MEXAmount)))

	vmInput := vmInputBuiler.Build()

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

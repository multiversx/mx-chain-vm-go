package hostCoretest

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"

	"github.com/stretchr/testify/require"
)

func Test_RunDEXPairBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}
	_ = logger.SetLogLevel("vm/benchmark:TRACE")
	owner := vmhost.MakeTestSCAddress("owner")
	user := vmhost.MakeTestSCAddress("user")

	_, host, mex := setupMEXPair(t, owner, user)

	mex.AddLiquidity(user, 1_000_000, 1, 1_000_000, 1)

	numBatches := 10
	numSwapsPerBatch := 2000

	wegldToMexESDT, wegldToMexESDTSwap := mex.CreateSwapVMInputs(mex.WEGLDToken, 100, mex.MEXToken, 1)
	mexToWegldESDT, mexToWegldSwap := mex.CreateSwapVMInputs(mex.MEXToken, 100, mex.WEGLDToken, 1)

	for batch := 0; batch < numBatches; batch++ {
		start := time.Now()
		for i := 0; i < numSwapsPerBatch/2; i++ {
			mex.ExecuteSwap(wegldToMexESDT, wegldToMexESDTSwap)
			mex.ExecuteSwap(mexToWegldESDT, mexToWegldSwap)
		}
		elapsedTime := time.Since(start)
		logBenchmark.Trace(
			"swap batch finished",
			"numSwapsPerBatch",
			numSwapsPerBatch,
			"duration",
			elapsedTime,
		)
	}

	defer func() {
		host.Reset()
	}()
}

func setupMEXPair(t *testing.T, owner Address, user Address) (*worldmock.MockWorld, vmhost.VMHost, *MEXSetup) {
	world, ownerAccount, host, err := prepare(t, owner)
	require.Nil(t, err)

	userAccount := world.AcctMap.CreateAccount(user, world)
	userAccount.Balance = big.NewInt(100)
	mex := NewMEXSetup(t, host, world, ownerAccount, userAccount)
	mex.Deploy()

	mex.ApplyInitialSetup()

	return world, host, mex
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
	Host  vmhost.VMHost
	World *worldmock.MockWorld
}

func NewMEXSetup(
	t *testing.T,
	host vmhost.VMHost,
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
		WithContractCodeMetadata([]byte{5, 6}).
		WithArguments(
			mex.WEGLDToken,
			mex.MEXToken,
			mex.OwnerAddress,
			mex.RouterAddress,
			big.NewInt(int64(mex.TotalFeePercent)).Bytes(),
			big.NewInt(int64(mex.SpecialFeePercent)).Bytes(),
		).
		WithGasProvided(math.MaxInt64).
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
		WithGasProvided(math.MaxInt64).
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
		WithGasProvided(math.MaxInt64).
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
		WithGasProvided(math.MaxInt64).
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
	_ = pairAccount.SetTokenRolesAsStrings(mex.LPToken, roles)
}

func (mex *MEXSetup) setESDTBalances() {
	_ = mex.UserAccount.SetTokenBalanceUint64(mex.WEGLDToken, 0, mex.UserWEGLDBalance)
	_ = mex.UserAccount.SetTokenBalanceUint64(mex.MEXToken, 0, mex.UserMEXBalance)
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
		WithGasProvided(math.MaxInt64)

	vmInputBuiler.
		WithESDTTokenName(mex.WEGLDToken).
		WithESDTValue(big.NewInt(int64(WEGLDAmount))).
		NextESDTTransfer().
		WithESDTTokenName(mex.MEXToken).
		WithESDTValue(big.NewInt(int64(MEXAmount)))

	vmInput := vmInputBuiler.Build()

	addLiquidityESDT := mex.createMultiESDTTransferVMInput(
		vmInput.CallerAddr,
		vmInput.RecipientAddr,
		vmInput.ESDTTransfers,
	)

	mex.performMultiESDTTransfer(addLiquidityESDT)

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (mex *MEXSetup) CreateSwapVMInputs(
	leftToken []byte,
	leftAmount uint64,
	rightToken []byte,
	rightAmount uint64,
) (*vmcommon.ContractCallInput, *vmcommon.ContractCallInput) {
	vmInputBuiler := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(mex.UserAccount.Address).
		WithRecipientAddr(mex.PairAddress)

	vmInputBuiler.
		WithESDTTokenName(leftToken).
		WithESDTValue(big.NewInt(int64(leftAmount))).
		WithFunction("swapTokensFixedInput").
		WithArguments(
			rightToken,
			big.NewInt(int64(rightAmount)).Bytes(),
		).
		WithGasProvided(math.MaxInt64)

	vmInput := vmInputBuiler.Build()

	multiTransferInput := mex.createMultiESDTTransferVMInput(
		vmInput.CallerAddr,
		vmInput.RecipientAddr,
		vmInput.ESDTTransfers,
	)

	return multiTransferInput, vmInput

}

func (mex *MEXSetup) ExecuteSwap(
	multiTransferInput *vmcommon.ContractCallInput,
	vmInput *vmcommon.ContractCallInput,
) {
	t := mex.T
	host := mex.Host
	world := mex.World

	mex.performMultiESDTTransfer(multiTransferInput)

	vmOutput, err := host.RunSmartContractCall(vmInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)

	err = world.CommitChanges()
	require.Nil(t, err)
}

func (mex *MEXSetup) createMultiESDTTransferVMInput(
	sender Address,
	receiver Address,
	esdtTransfers []*vmcommon.ESDTTransfer,
) *vmcommon.ContractCallInput {
	nrTransfers := len(esdtTransfers)
	nrTransfersAsBytes := big.NewInt(0).SetUint64(uint64(nrTransfers)).Bytes()

	multiTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasPrice:    1,
			GasProvided: math.MaxInt64,
			GasLocked:   0,
		},
		RecipientAddr:     sender,
		Function:          core.BuiltInFunctionMultiESDTNFTTransfer,
		AllowInitFunction: false,
	}
	multiTransferInput.Arguments = append(multiTransferInput.Arguments, receiver, nrTransfersAsBytes)

	for i := 0; i < nrTransfers; i++ {
		token := esdtTransfers[i].ESDTTokenName
		nonceAsBytes := big.NewInt(0).SetUint64(esdtTransfers[i].ESDTTokenNonce).Bytes()
		value := esdtTransfers[i].ESDTValue.Bytes()

		multiTransferInput.Arguments = append(multiTransferInput.Arguments, token, nonceAsBytes, value)
	}

	return multiTransferInput
}

func (mex *MEXSetup) performMultiESDTTransfer(
	multiTransferInput *vmcommon.ContractCallInput,
) {
	t := mex.T
	world := mex.World

	vmOutput, err := world.BuiltinFuncs.ProcessBuiltInFunction(multiTransferInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, "", vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	_ = world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

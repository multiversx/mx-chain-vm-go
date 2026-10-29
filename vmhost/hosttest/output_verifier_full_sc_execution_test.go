package hostCoretest

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	vm "github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/contexts"
	"github.com/multiversx/mx-chain-vm-go/vmhost/verifier"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOutputVerifier_ComplicatedMultiContractCall_WithDevRewardsAndTransfers verifies that
// a multi-contract execution chain with developer reward claiming, cross-contract calls,
// and external transfers passes exit-boundary verification cleanly.
func TestOutputVerifier_ComplicatedMultiContractCall_WithDevRewardsAndTransfers(t *testing.T) {
	t.Parallel()

	callValue := int64(500000)
	claimedRewards := big.NewInt(50000)
	transferToUser := big.NewInt(300000)

	world := worldmock.NewMockWorld()
	world.AcctMap.CreateAccount(test.UserAddress, world)

	input := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(test.UserAddress).
		WithRecipientAddr(test.ParentAddress).
		WithFunction("executeComplicatedWorkflow").
		WithCallValue(callValue).
		WithGasProvided(500000).
		Build()

	builder := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, _ interface{}) {
					parentInstance.AddMockMethod("executeComplicatedWorkflow", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksLegacy(host)

						// 1. Cross-contract call from Parent to Child
						childAddrBuf := managedTypes.NewManagedBufferFromBytes(test.ChildAddress)
						fnBuf := managedTypes.NewManagedBufferFromBytes([]byte("claimAndTransfer"))
						valBuf := managedTypes.NewBigInt(big.NewInt(0))
						argsVecBuf := managedTypes.NewManagedBuffer()

						res := hooks.ManagedExecuteOnDestContext(
							100000,
							childAddrBuf,
							valBuf,
							fnBuf,
							argsVecBuf,
							managedTypes.NewManagedBuffer(),
						)
						require.Equal(t, int32(0), res)

						// 2. Parent performs external transfer of 300,000 to user
						err := host.Output().Transfer(
							test.UserAddress,
							test.ParentAddress,
							0,
							0,
							transferToUser,
							nil,
							nil,
							vm.DirectCall,
						)
						require.Nil(t, err)

						// 3. Set Parent's BalanceDelta:
						// Inflow = 500,000 (CallValue) + 50,000 (claimed by Child and credited to Parent)
						// Outflow = 300,000 (transferred to user)
						// Net delta remaining on Parent = +250,000
						parentOutputAcc, _ := host.Output().GetOutputAccount(test.ParentAddress)
						parentOutputAcc.BalanceDelta = big.NewInt(250000)

						return parentInstance
					})
				}),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(100000).
				WithMethods(func(childInstance *contextmock.InstanceMock, _ interface{}) {
					childInstance.AddMockMethod("claimAndTransfer", func() *contextmock.InstanceMock {
						host := childInstance.Host

						// Child claims developer rewards by emitting authentic ClaimDeveloperRewards log
						host.Output().WriteLogWithIdentifier(
							test.ParentAddress,
							[][]byte{claimedRewards.Bytes()},
							nil,
							[]byte(core.BuiltInFunctionClaimDeveloperRewards),
						)

						return childInstance
					})
				}),
		).
		WithInput(input)

	_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
}

// TestOutputVerifier_ComplicatedESDTAndNFTLifecycle verifies that complex token lifecycles
// (incoming payment, minting NFT, partial transfer, partial burn, post-execution balance)
// are accurately verified by the exit boundary verifier.
func TestOutputVerifier_ComplicatedESDTAndNFTLifecycle(t *testing.T) {
	t.Parallel()

	tokenID := "ESDT-123456"
	nftID := "NFT-abcdef"
	incomingESDT := big.NewInt(1000)
	mintedNFTAmount := big.NewInt(10)
	transferredNFTAmount := big.NewInt(4)
	transferredESDTAmount := big.NewInt(400)
	burnedESDTAmount := big.NewInt(100)

	world := worldmock.NewMockWorld()
	world.AcctMap.CreateAccount(test.UserAddress, world)

	input := test.CreateTestContractCallInputBuilder().
		WithCallerAddr(test.UserAddress).
		WithRecipientAddr(test.ParentAddress).
		WithFunction("manageTokens").
		WithGasProvided(500000).
		Build()
	input.ESDTTransfers = []*vmcommon.ESDTTransfer{
		{
			ESDTTokenName:  []byte(tokenID),
			ESDTTokenNonce: 0,
			ESDTValue:      incomingESDT,
		},
	}

	builder := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, _ interface{}) {
					parentInstance.AddMockMethod("manageTokens", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						userAcc, _ := host.Output().GetOutputAccount(test.UserAddress)

						// 1. Emit NFT Create event (mint 10 units of NFT-abcdef nonce 1)
						host.Output().WriteLogWithIdentifier(
							test.ParentAddress,
							[][]byte{
								[]byte(nftID),
								big.NewInt(1).Bytes(),
								mintedNFTAmount.Bytes(),
							},
							nil,
							[]byte(core.BuiltInFunctionESDTNFTCreate),
						)

						// 2. Transfer 4 units of NFT to user
						nftTransferData := []byte("ESDTNFTTransfer@" +
							hex.EncodeToString([]byte(nftID)) + "@" +
							hex.EncodeToString(big.NewInt(1).Bytes()) + "@" +
							hex.EncodeToString(transferredNFTAmount.Bytes()) + "@" +
							hex.EncodeToString(test.UserAddress))
						contexts.AppendOutputTransfers(userAcc, userAcc.OutputTransfers, vmcommon.OutputTransfer{
							Data:          nftTransferData,
							SenderAddress: test.ParentAddress,
							Value:         big.NewInt(0),
						})
						host.Output().WriteLogWithIdentifier(
							test.ParentAddress,
							[][]byte{
								[]byte(nftID),
								big.NewInt(1).Bytes(),
								transferredNFTAmount.Bytes(),
								test.UserAddress,
							},
							nil,
							[]byte(core.BuiltInFunctionESDTNFTTransfer),
						)

						// 3. Transfer 400 units of ESDT to user
						esdtTransferData := []byte("ESDTTransfer@" +
							hex.EncodeToString([]byte(tokenID)) + "@" +
							hex.EncodeToString(transferredESDTAmount.Bytes()))
						contexts.AppendOutputTransfers(userAcc, userAcc.OutputTransfers, vmcommon.OutputTransfer{
							Data:          esdtTransferData,
							SenderAddress: test.ParentAddress,
							Value:         big.NewInt(0),
						})
						host.Output().WriteLogWithIdentifier(
							test.ParentAddress,
							[][]byte{
								[]byte(tokenID),
								big.NewInt(0).Bytes(),
								transferredESDTAmount.Bytes(),
								test.UserAddress,
							},
							nil,
							[]byte(core.BuiltInFunctionESDTTransfer),
						)

						// 4. Burn 100 units of ESDT
						host.Output().WriteLogWithIdentifier(
							test.ParentAddress,
							[][]byte{
								[]byte(tokenID),
								big.NewInt(0).Bytes(),
								burnedESDTAmount.Bytes(),
							},
							nil,
							[]byte(core.BuiltInFunctionESDTBurn),
						)

						return parentInstance
					})
				}),
		).
		WithInput(input).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			// Set mock ESDT post-balance in storage to 500 (1000 in - 400 out - 100 burn)
			parentAccount := world.AcctMap.GetAccount(test.ParentAddress)
			require.NotNil(t, parentAccount)
			_ = parentAccount.SetTokenBalanceUint64([]byte(tokenID), 0, 500)
		})

	_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
}

// TestOutputVerifier_Defense_DetectsAndReverts_UnbackedESDTOutflow verifies that when a contract
// attempts to output an unbacked ESDT transfer, the exit-boundary verifier catches the violation
// and reverts execution with ErrTokenConservationViolated.
func TestOutputVerifier_Defense_DetectsAndReverts_UnbackedESDTOutflow(t *testing.T) {
	t.Parallel()

	unbackedTokenID := "UNBACKED-999999"
	unbackedAmount := big.NewInt(1000)

	world := worldmock.NewMockWorld()
	world.AcctMap.CreateAccount(test.UserAddress, world)

	builder := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, _ interface{}) {
					parentInstance.AddMockMethod("attackUnbackedESDT", func() *contextmock.InstanceMock {
						host := parentInstance.Host

						// Contract attempts to send 1,000 tokens out of thin air
						esdtTransferData := []byte("ESDTTransfer@" +
							hex.EncodeToString([]byte(unbackedTokenID)) + "@" +
							hex.EncodeToString(unbackedAmount.Bytes()))
						userAcc, _ := host.Output().GetOutputAccount(test.UserAddress)
						contexts.AppendOutputTransfers(userAcc, userAcc.OutputTransfers, vmcommon.OutputTransfer{
							Data:          esdtTransferData,
							SenderAddress: test.ParentAddress,
							Value:         big.NewInt(0),
						})

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.UserAddress).
			WithRecipientAddr(test.ParentAddress).
			WithFunction("attackUnbackedESDT").
			WithGasProvided(500000).
			Build())

	_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
	// Exit-boundary verifier actively caught the unbacked transfer!
	require.Error(t, err)
	require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
	require.Contains(t, vmOutput.ReturnMessage, "token conservation violated")
}

// TestOutputVerifier_Defense_DetectsAndReverts_NativeBalanceInflation verifies that when a contract
// attempts to inflate its native EGLD balance beyond its deposits, the exit boundary verifier reverts.
func TestOutputVerifier_Defense_DetectsAndReverts_NativeBalanceInflation(t *testing.T) {
	t.Parallel()

	world := worldmock.NewMockWorld()
	world.AcctMap.CreateAccount(test.UserAddress, world)

	builder := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, _ interface{}) {
					parentInstance.AddMockMethod("attackInflateBalance", func() *contextmock.InstanceMock {
						host := parentInstance.Host

						// Artificially inject +500,000 EGLD delta with 0 deposit
						parentOutputAcc, _ := host.Output().GetOutputAccount(test.ParentAddress)
						parentOutputAcc.BalanceDelta = big.NewInt(500000)

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.UserAddress).
			WithRecipientAddr(test.ParentAddress).
			WithFunction("attackInflateBalance").
			WithCallValue(100). // Deposit is only 100
			WithGasProvided(500000).
			Build())

	_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
	require.Error(t, err)
	require.ErrorIs(t, err, verifier.ErrNativeTokenConservationViolated)
	require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
	require.Contains(t, vmOutput.ReturnMessage, "native token conservation violated")
}

// TestOutputVerifier_Defense_DetectsAndReverts_AddressKeyMismatch_UF208 verifies that the verifier
// enforces the UF-208 identity aliasing invariant (map key == acc.Address).
func TestOutputVerifier_Defense_DetectsAndReverts_AddressKeyMismatch_UF208(t *testing.T) {
	t.Parallel()

	world := worldmock.NewMockWorld()
	world.AcctMap.CreateAccount(test.UserAddress, world)

	builder := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, _ interface{}) {
					parentInstance.AddMockMethod("attackKeyMismatch", func() *contextmock.InstanceMock {
						host := parentInstance.Host

						// Corrupt output accounts by inserting account under mismatched key
						attackerKey := string(bytes.Repeat([]byte("A"), 32))
						host.Output().GetOutputAccounts()[attackerKey] = &vmcommon.OutputAccount{
							Address:      test.ChildAddress, // Address != key!
							BalanceDelta: big.NewInt(0),
						}

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.UserAddress).
			WithRecipientAddr(test.ParentAddress).
			WithFunction("attackKeyMismatch").
			WithGasProvided(500000).
			Build())

	_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
	require.Error(t, err)
	require.ErrorIs(t, err, verifier.ErrAddressKeyMismatch)
	require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
	require.Contains(t, vmOutput.ReturnMessage, "map key")
}

// TestOutputVerifier_Defense_DetectsAndReverts_ReadOnlyPurityViolation verifies that when
// a contract call is executed under ReadOnly, any state modifications are rejected.
func TestOutputVerifier_Defense_DetectsAndReverts_ReadOnlyPurityViolation(t *testing.T) {
	t.Parallel()

	world := worldmock.NewMockWorld()
	world.AcctMap.CreateAccount(test.UserAddress, world)

	builder := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, _ interface{}) {
					parentInstance.AddMockMethod("viewMethodWithStorageMutation", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						host.Runtime().SetReadOnly(true)

						// Illicitly inject storage updates into output account during read-only call
						acc, _ := host.Output().GetOutputAccount(test.ParentAddress)
						acc.StorageUpdates = map[string]*vmcommon.StorageUpdate{
							"illicit_key": {
								Offset:  []byte("illicit_key"),
								Data:    []byte("illicit_data"),
								Written: true,
							},
						}

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.UserAddress).
			WithRecipientAddr(test.ParentAddress).
			WithFunction("viewMethodWithStorageMutation").
			WithGasProvided(500000).
			Build()).
		WithSetup(func(host vmhost.VMHost, _ *worldmock.MockWorld) {
			enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
			enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
				return true
			}
		})

	_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
	require.Error(t, err)
	require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange)
	require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
	require.Contains(t, vmOutput.ReturnMessage, "read-only")
}

// TestOutputVerifier_Defense_DetectsAndReverts_NegativeOutputTransferValue verifies that
// output transfers with negative values are rejected by the exit boundary verifier.
func TestOutputVerifier_Defense_DetectsAndReverts_NegativeOutputTransferValue(t *testing.T) {
	t.Parallel()

	world := worldmock.NewMockWorld()
	world.AcctMap.CreateAccount(test.UserAddress, world)

	builder := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, _ interface{}) {
					parentInstance.AddMockMethod("attackNegativeTransfer", func() *contextmock.InstanceMock {
						host := parentInstance.Host

						// Directly add output transfer with negative value
						acc, _ := host.Output().GetOutputAccount(test.ParentAddress)
						acc.OutputTransfers = append(acc.OutputTransfers, vmcommon.OutputTransfer{
							Value:         big.NewInt(-1000),
							SenderAddress: test.ParentAddress,
						})

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.UserAddress).
			WithRecipientAddr(test.ParentAddress).
			WithFunction("attackNegativeTransfer").
			WithGasProvided(500000).
			Build())

	_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
	require.Error(t, err)
	require.ErrorIs(t, err, verifier.ErrNegativeOutputTransferValue)
	require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
	assert.Contains(t, vmOutput.ReturnMessage, "negative value")
}

package hostCoretest

import (
	"math/big"
	"testing"

	vm "github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/require"
)

func TestHost_DestinationAddressMutationResistance(t *testing.T) {
	parentBalance := int64(1000000)
	childBalance := int64(50000)
	gasProvided := uint64(1000000)

	victimAddress := make([]byte, len(test.ChildAddress))
	copy(victimAddress, []byte("victim_address_123456789012345678"))
	attackerRecipient := []byte("attacker_dest_000000000000000000")
	stealAmount := big.NewInt(10000)

	var capturedOutput *vmcommon.VMOutput

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("attack", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksLegacy(host)

						// 1. Create managed buffer containing child contract's address
						childAddrBuf := managedTypes.NewManagedBufferFromBytes(test.ChildAddress)

						// 2. Prepare arguments to call Child.doDebitTransfer
						fnBuf := managedTypes.NewManagedBufferFromBytes([]byte("doDebitTransfer"))
						valBuf := managedTypes.NewBigInt(big.NewInt(0)) // 0 call value
						argsVecBuf := managedTypes.NewManagedBuffer()

						// 3. Call Child via ManagedExecuteOnDestContext
						res := hooks.ManagedExecuteOnDestContext(
							50000,
							childAddrBuf,
							valBuf,
							fnBuf,
							argsVecBuf,
							managedTypes.NewManagedBuffer(),
						)
						require.Equal(t, int32(0), res)

						// 4. Now Child has returned, and Child's output (with sender ChildAddress debited)
						// is merged into Parent. In Parent's outputState, Child's OutputAccount.Address
						// points to the backing slice of childAddrBuf!
						// We now mutate childAddrBuf in place to victimAddress!
						_, err := vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, childAddrBuf, 0, int32(len(victimAddress)), victimAddress)
						handleHookError(host, err)

						// 5. Trigger a second nested execution (dummy call to Child)
						// This pushes Parent's state to stateStack, executes dummy call,
						// and upon completion, PopMergeActiveState merges prevState into active state,
						// rekeying the debited account under victimAddress!
						dummyFnBuf := managedTypes.NewManagedBufferFromBytes([]byte("dummy"))
						dummyAddrBuf := managedTypes.NewManagedBufferFromBytes(test.ChildAddress)
						res = hooks.ManagedExecuteOnDestContext(
							20000,
							dummyAddrBuf,
							valBuf,
							dummyFnBuf,
							managedTypes.NewManagedBuffer(),
							managedTypes.NewManagedBuffer(),
						)
						// Re-create Child in OutputAccounts so meteringContext.setGasUsedToOutputAccounts passes!
						_, _ = host.Output().GetOutputAccount(test.ChildAddress)

						return parentInstance
					})
				}),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(childBalance).
				WithMethods(func(childInstance *contextmock.InstanceMock, config interface{}) {
					childInstance.AddMockMethod("doDebitTransfer", func() *contextmock.InstanceMock {
						host := childInstance.Host

						// Child transfers 10,000 EGLD to attackerRecipient
						// This causes sender (Child, whose Address points to Parent's buffer)
						// to get BalanceDelta = -10,000!
						err := host.Output().Transfer(
							attackerRecipient,
							host.Runtime().GetContextAddress(),
							0,
							0,
							stealAmount,
							nil,
							nil,
							vm.DirectCall,
						)
						require.Nil(t, err)

						return childInstance
					})
					childInstance.AddMockMethod("dummy", func() *contextmock.InstanceMock {
						return childInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("attack").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			capturedOutput = verify.VmOutput
		})

	require.Nil(t, err)
	require.NotNil(t, capturedOutput)

	t.Logf("ReturnCode: %v, ReturnMessage: %s, NumAccounts: %d", capturedOutput.ReturnCode, capturedOutput.ReturnMessage, len(capturedOutput.OutputAccounts))
	for key, acc := range capturedOutput.OutputAccounts {
		t.Logf("  Key: %x (%s), Address: %x (%s), BalanceDelta: %s",
			[]byte(key), key, acc.Address, string(acc.Address), acc.BalanceDelta.String())
	}

	// Verify that victim was not debited
	victimAcc, hasVictim := capturedOutput.OutputAccounts[string(victimAddress)]
	require.False(t, hasVictim, "Victim account must not be present in OutputAccounts")
	require.Nil(t, victimAcc)

	// Verify child account retained the debit instead
	childAcc, hasChild := capturedOutput.OutputAccounts[string(test.ChildAddress)]
	require.True(t, hasChild, "Child account must be present")
	require.Equal(t, big.NewInt(0).Neg(stealAmount), childAcc.BalanceDelta, "Child retains its own debit")
}

func TestHost_StorageOverwriteResistance(t *testing.T) {
	parentBalance := int64(1000000)
	childBalance := int64(50000)
	gasProvided := uint64(1000000)

	victimAddress := make([]byte, len(test.ChildAddress))
	copy(victimAddress, []byte("victim_contract_0000000000000000"))
	storageKey := []byte("target_storage_key")
	storageValue := []byte("maliciously_written_value")

	var capturedOutput *vmcommon.VMOutput

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(parentBalance).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("attackStorage", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksLegacy(host)

						// 1. Create managed buffer containing child contract's address
						childAddrBuf := managedTypes.NewManagedBufferFromBytes(test.ChildAddress)

						// 2. Prepare arguments to call Child.doStorageWrite
						fnBuf := managedTypes.NewManagedBufferFromBytes([]byte("doStorageWrite"))
						valBuf := managedTypes.NewBigInt(big.NewInt(0))
						argsVecBuf := managedTypes.NewManagedBuffer()

						// 3. Call Child via ManagedExecuteOnDestContext
						res := hooks.ManagedExecuteOnDestContext(
							50000,
							childAddrBuf,
							valBuf,
							fnBuf,
							argsVecBuf,
							managedTypes.NewManagedBuffer(),
						)
						require.Equal(t, int32(0), res)

						// 4. Mutate childAddrBuf to victimAddress
						_, err := vmhooks.ManagedBufferSetByteSliceWithTypedArgs(host, childAddrBuf, 0, int32(len(victimAddress)), victimAddress)
						handleHookError(host, err)

						// 5. Trigger a second nested execution to rekey via PopMergeActiveState
						dummyFnBuf := managedTypes.NewManagedBufferFromBytes([]byte("dummy"))
						dummyAddrBuf := managedTypes.NewManagedBufferFromBytes(test.ChildAddress)
						res = hooks.ManagedExecuteOnDestContext(
							20000,
							dummyAddrBuf,
							valBuf,
							dummyFnBuf,
							managedTypes.NewManagedBuffer(),
							managedTypes.NewManagedBuffer(),
						)
						// Re-create Child in OutputAccounts so meteringContext.setGasUsedToOutputAccounts passes
						_, _ = host.Output().GetOutputAccount(test.ChildAddress)

						return parentInstance
					})
				}),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(childBalance).
				WithMethods(func(childInstance *contextmock.InstanceMock, config interface{}) {
					childInstance.AddMockMethod("doStorageWrite", func() *contextmock.InstanceMock {
						host := childInstance.Host

						// Child writes to storage
						_, err := host.Storage().SetStorage(storageKey, storageValue)
						require.Nil(t, err)

						return childInstance
					})
					childInstance.AddMockMethod("dummy", func() *contextmock.InstanceMock {
						return childInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("attackStorage").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			capturedOutput = verify.VmOutput
		})

	require.Nil(t, err)
	require.NotNil(t, capturedOutput)

	t.Logf("ReturnCode: %v, ReturnMessage: %s, NumAccounts: %d", capturedOutput.ReturnCode, capturedOutput.ReturnMessage, len(capturedOutput.OutputAccounts))
	for key, acc := range capturedOutput.OutputAccounts {
		t.Logf("  Key: %s, Address: %s, StorageUpdates: %v", key, string(acc.Address), acc.StorageUpdates)
	}

	// Verify if victim's account was not modified
	victimAcc, hasVictim := capturedOutput.OutputAccounts[string(victimAddress)]
	require.False(t, hasVictim, "Victim account must not be present in OutputAccounts")
	require.Nil(t, victimAcc)

	// Verify that Child account retained its own storage updates
	childAcc, hasChild := capturedOutput.OutputAccounts[string(test.ChildAddress)]
	require.True(t, hasChild, "Child account must be present")
	require.NotNil(t, childAcc.StorageUpdates[string(storageKey)], "Child account retains its own storage key")
	require.Equal(t, storageValue, childAcc.StorageUpdates[string(storageKey)].Data)
}

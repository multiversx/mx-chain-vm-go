package hostCoretest

import (
	"bytes"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/mock/contracts"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHost_ExecuteDelete_ReadOnly_And_Permission(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			// Subtest 1: ReadOnly mode rejection
			t.Run("ReadOnly mode", func(t *testing.T) {
				builder := test.BuildMockInstanceCallTest(t).
					WithContracts(
						test.CreateMockContract(test.ParentAddress).
							WithBalance(100000).
							WithOwnerAddress(test.ParentAddress).
							WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
							WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
								parentInstance.AddMockMethod("testDelete", func() *contextmock.InstanceMock {
									host := parentInstance.Host
									host.Runtime().SetReadOnly(true)

									childInput := test.CreateTestContractCallInputBuilder().
										WithCallerAddr(test.ParentAddress).
										WithRecipientAddr(test.ChildAddress).
										WithFunction("deleteContract").
										WithGasProvided(10000).
										Build()
									childInput.AsyncArguments = &vmcommon.AsyncArguments{
										CallID:       []byte{},
										CallerCallID: []byte{},
									}

									_, _, deleteErr := host.ExecuteOnDestContext(childInput)
									if flagEnabled {
										require.Equal(t, vmhost.ErrInvalidCallOnReadOnlyMode, deleteErr)
									} else {
										require.Nil(t, deleteErr)
									}

									return parentInstance
								})
							}),
						test.CreateMockContract(test.ChildAddress).
							WithBalance(0).
							WithOwnerAddress(test.ParentAddress).
							WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
							WithMethods(),
					).
					WithInput(test.CreateTestContractCallInputBuilder().
						WithRecipientAddr(test.ParentAddress).
						WithGasProvided(100000).
						WithFunction("testDelete").
						Build()).
					WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
						enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
						enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
							if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
								return flagEnabled
							}
							return true
						}
					})

				world := worldmock.NewMockWorld()
				_, _, err := builder.RunTest(world, true, test.RunTest)
				if flagEnabled {
					assert.Nil(t, err)
				} else {
					require.Error(t, err)
					require.Contains(t, err.Error(), "read-only execution violation")
				}
			})

			// Subtest 2: Permission failure propagation (non-upgradeable contract)
			t.Run("Permission failure non-upgradeable", func(t *testing.T) {
				_, err := test.BuildMockInstanceCallTest(t).
					WithContracts(
						test.CreateMockContract(test.ParentAddress).
							WithBalance(100000).
							WithOwnerAddress(test.ParentAddress).
							WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
							WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
								parentInstance.AddMockMethod("testDeleteNotUpgradeable", func() *contextmock.InstanceMock {
									host := parentInstance.Host

									childInput := test.CreateTestContractCallInputBuilder().
										WithCallerAddr(test.ParentAddress).
										WithRecipientAddr(test.ChildAddress).
										WithFunction("deleteContract").
										WithGasProvided(10000).
										Build()
									childInput.AsyncArguments = &vmcommon.AsyncArguments{
										CallID:       []byte{},
										CallerCallID: []byte{},
									}

									_, _, deleteErr := host.ExecuteOnDestContext(childInput)
									if flagEnabled {
										require.Equal(t, vmhost.ErrUpgradeNotAllowed, deleteErr)
									} else {
										require.Nil(t, deleteErr)
									}

									return parentInstance
								})
							}),
						test.CreateMockContract(test.ChildAddress).
							WithBalance(0).
							WithOwnerAddress(test.ParentAddress).
							WithCodeMetadata([]byte{0, 0}). // NOT upgradeable!
							WithMethods(),
					).
					WithInput(test.CreateTestContractCallInputBuilder().
						WithRecipientAddr(test.ParentAddress).
						WithGasProvided(100000).
						WithFunction("testDeleteNotUpgradeable").
						Build()).
					WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
						enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
						enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
							if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
								return flagEnabled
							}
							return true
						}
					}).
					AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
						verify.Ok()
					})
				assert.Nil(t, err)
			})
		})
	}
}

func TestHost_ExecuteUpgrade_ReadOnly(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithOwnerAddress(test.ParentAddress).
						WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testUpgradeReadOnly", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								host.Runtime().SetReadOnly(true)

								childInput := test.CreateTestContractCallInputBuilder().
									WithCallerAddr(test.ParentAddress).
									WithRecipientAddr(test.ChildAddress).
									WithFunction("upgradeContract").
									WithGasProvided(10000).
									Build()
								childInput.AsyncArguments = &vmcommon.AsyncArguments{
									CallID:       []byte{},
									CallerCallID: []byte{},
								}

								_, _, upgradeErr := host.ExecuteOnDestContext(childInput)
								if flagEnabled {
									require.Equal(t, vmhost.ErrInvalidCallOnReadOnlyMode, upgradeErr)
								} else {
									// When flag is disabled, it proceeds past read-only check to checkUpgradePermission / argument extraction
									require.NotEqual(t, vmhost.ErrInvalidCallOnReadOnlyMode, upgradeErr)
								}

								return parentInstance
							})
						}),
					test.CreateMockContract(test.ChildAddress).
						WithBalance(0).
						WithOwnerAddress(test.ParentAddress).
						WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
						WithMethods(),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(100000).
					WithFunction("testUpgradeReadOnly").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
				})
			assert.Nil(t, err)
		})
	}
}

func TestHost_CreateContract_DefensiveCodeClone(t *testing.T) {
	t.Parallel()

	rawCode := test.GetTestSCCode("counter", "../../")
	codeToDeploy := append([]byte(nil), rawCode...)
	firstByteOriginal := codeToDeploy[0]

	world := worldmock.NewMockWorld()
	host := test.NewTestHostBuilder(t).
		WithBlockchainHook(world).
		WithBuiltinFunctions().
		Build()
	defer func() {
		host.Reset()
	}()

	world.NewAddressMocks = append(world.NewAddressMocks, &worldmock.NewAddressMock{
		CreatorAddress: test.UserAddress,
		CreatorNonce:   0,
		NewAddress:     test.ParentAddress,
	})
	world.AcctMap.CreateAccount(test.UserAddress, world)

	deployInput := test.CreateTestContractCreateInputBuilder().
		WithContractCode(codeToDeploy).
		WithCallerAddr(test.UserAddress).
		WithGasProvided(10000000).
		Build()

	vmOutput, err := host.RunSmartContractCreate(deployInput)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	// Mutate the original slice after deployment call
	codeToDeploy[0] = 0xff

	// Find the created account in OutputAccounts
	var deployedCodeInOutput []byte
	for _, acc := range vmOutput.OutputAccounts {
		if bytes.Equal(acc.Address, test.ParentAddress) && len(acc.Code) > 0 {
			deployedCodeInOutput = acc.Code
			break
		}
	}

	require.NotNil(t, deployedCodeInOutput, "Deployed code must be present in output accounts")
	require.Equal(t, firstByteOriginal, deployedCodeInOutput[0], "Deployed code must not be affected by post-deploy mutation of input slice")
}

func TestHost_DestinationAddressLength_FlagGated(t *testing.T) {
	t.Parallel()

	shortDest := []byte("short_dest_not_32_bytes")

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			builder := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(1000000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testShortDest", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								managedTypes := host.ManagedTypes()
								hooks := vmhooks.NewVMHooksImpl(host)

								destBuf := managedTypes.NewManagedBufferFromBytes(shortDest)
								valueBuf := managedTypes.NewBigIntFromInt64(0)
								funcBuf := managedTypes.NewManagedBufferFromBytes([]byte("dummyFunc"))
								argsBuf := managedTypes.NewManagedBufferFromBytes(nil)

								hooks.ManagedAsyncCall(destBuf, valueBuf, funcBuf, argsBuf)

								allErrors := host.Runtime().GetAllErrors()
								if flagEnabled {
									require.ErrorIs(t, allErrors, vmhost.ErrLengthOfBufferNotCorrect)
								} else {
									require.False(t, errors.Is(allErrors, vmhost.ErrLengthOfBufferNotCorrect))
								}

								return parentInstance
							})
						}),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(1000000).
					WithFunction("testShortDest").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				})

			world := worldmock.NewMockWorld()
			_, _, err := builder.RunTest(world, true, test.RunTest)
			if flagEnabled {
				assert.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), "must be exactly 32 bytes")
			}
		})
	}
}

func TestHost_ManagedHooksAddressLength_FlagGated(t *testing.T) {
	t.Parallel()

	shortAddr := []byte("short_address_16")

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(1000000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testAddressHooks", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								managedTypes := host.ManagedTypes()
								hooks := vmhooks.NewVMHooksImpl(host)

								addrHandle := managedTypes.NewManagedBufferFromBytes(shortAddr)
								tokenHandle := managedTypes.NewManagedBufferFromBytes([]byte("TOKEN-123456"))
								valHandle := managedTypes.NewBigIntFromInt64(0)
								respHandle := managedTypes.NewManagedBuffer()

								// 1. ManagedGetESDTBalance
								hooks.ManagedGetESDTBalance(addrHandle, tokenHandle, 0, valHandle)
								allErrors := host.Runtime().GetAllErrors()
								if flagEnabled {
									require.ErrorIs(t, allErrors, vmhost.ErrLengthOfBufferNotCorrect)
								}

								// 2. ManagedGetCodeMetadata
								hooks.ManagedGetCodeMetadata(addrHandle, respHandle)
								allErrors = host.Runtime().GetAllErrors()
								if flagEnabled {
									require.ErrorIs(t, allErrors, vmhost.ErrLengthOfBufferNotCorrect)
								}

								// 3. ManagedGetCodeHash
								hooks.ManagedGetCodeHash(addrHandle, respHandle)
								allErrors = host.Runtime().GetAllErrors()
								if flagEnabled {
									require.ErrorIs(t, allErrors, vmhost.ErrLengthOfBufferNotCorrect)
								}

								// 4. ManagedIsESDTFrozen
								res := hooks.ManagedIsESDTFrozen(addrHandle, tokenHandle, 0)
								if flagEnabled {
									require.Equal(t, int32(-1), res)
									require.ErrorIs(t, host.Runtime().GetAllErrors(), vmhost.ErrLengthOfBufferNotCorrect)
								}

								// 5. MBufferStorageLoadFromAddress
								keyHandle := managedTypes.NewManagedBufferFromBytes([]byte("key"))
								hooks.MBufferStorageLoadFromAddress(addrHandle, keyHandle, respHandle)
								if flagEnabled {
									require.ErrorIs(t, host.Runtime().GetAllErrors(), vmhost.ErrLengthOfBufferNotCorrect)
								}

								return parentInstance
							})
						}),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(1000000).
					WithFunction("testAddressHooks").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

			assert.Nil(t, err)
		})
	}
}

func TestHost_VerifyBLSAggregated_Hardening_BUG524(t *testing.T) {
	t.Parallel()

	valid96Key := make([]byte, 96)
	valid96Key[0] = 0x01
	valid96Key[95] = 0x02

	invalidKey := make([]byte, 95) // malformed length

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			// Subtest 1: Exceeds 400 signers cap
			t.Run("Exceeds max signers 400", func(t *testing.T) {
				_, err := test.BuildMockInstanceCallTest(t).
					WithContracts(
						test.CreateMockContract(test.ParentAddress).
							WithBalance(100000000).
							WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
								parentInstance.AddMockMethod("testExceedsSigners", func() *contextmock.InstanceMock {
									host := parentInstance.Host
									managedTypes := host.ManagedTypes()

									keysList := make([][]byte, 401)
									for i := range keysList {
										keysList[i] = valid96Key
									}
									keysHandle := managedTypes.NewManagedBuffer()
									_ = managedTypes.WriteManagedVecOfManagedBuffers(keysList, keysHandle)

									msgHandle := managedTypes.NewManagedBufferFromBytes([]byte("message"))
									sigHandle := managedTypes.NewManagedBufferFromBytes(make([]byte, 48))

									res := vmhooks.ManagedVerifyBLSWithHost(
										host,
										keysHandle,
										msgHandle,
										sigHandle,
										"verifyBLSAggregatedSignature")

									allErrors := host.Runtime().GetAllErrors()
									if flagEnabled {
										require.Equal(t, int32(1), res)
										require.ErrorIs(t, allErrors, vmhost.ErrArgOutOfRange)
									}

									return parentInstance
								})
							}),
					).
					WithInput(test.CreateTestContractCallInputBuilder().
						WithRecipientAddr(test.ParentAddress).
						WithGasProvided(100000000).
						WithFunction("testExceedsSigners").
						Build()).
					WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
						enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
						enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
							if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
								return flagEnabled
							}
							return true
						}
					}).
					AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

				assert.Nil(t, err)
			})

			// Subtest 2: Malformed key length
			t.Run("Malformed key length", func(t *testing.T) {
				_, err := test.BuildMockInstanceCallTest(t).
					WithContracts(
						test.CreateMockContract(test.ParentAddress).
							WithBalance(100000000).
							WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
								parentInstance.AddMockMethod("testMalformedKey", func() *contextmock.InstanceMock {
									host := parentInstance.Host
									managedTypes := host.ManagedTypes()

									keysList := [][]byte{invalidKey}
									keysHandle := managedTypes.NewManagedBuffer()
									_ = managedTypes.WriteManagedVecOfManagedBuffers(keysList, keysHandle)

									msgHandle := managedTypes.NewManagedBufferFromBytes([]byte("message"))
									sigHandle := managedTypes.NewManagedBufferFromBytes(make([]byte, 48))

									res := vmhooks.ManagedVerifyBLSWithHost(
										host,
										keysHandle,
										msgHandle,
										sigHandle,
										"verifyBLSAggregatedSignature")

									allErrors := host.Runtime().GetAllErrors()
									if flagEnabled {
										require.Equal(t, int32(1), res)
										require.ErrorIs(t, allErrors, vmhost.ErrLengthOfBufferNotCorrect)
									}

									return parentInstance
								})
							}),
					).
					WithInput(test.CreateTestContractCallInputBuilder().
						WithRecipientAddr(test.ParentAddress).
						WithGasProvided(100000000).
						WithFunction("testMalformedKey").
						Build()).
					WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
						enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
						enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
							if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
								return flagEnabled
							}
							return true
						}
					}).
					AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

				assert.Nil(t, err)
			})

			// Subtest 3: Empty signers list
			t.Run("Empty signers list", func(t *testing.T) {
				_, err := test.BuildMockInstanceCallTest(t).
					WithContracts(
						test.CreateMockContract(test.ParentAddress).
							WithBalance(100000000).
							WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
								parentInstance.AddMockMethod("testEmptySigners", func() *contextmock.InstanceMock {
									host := parentInstance.Host
									managedTypes := host.ManagedTypes()

									keysHandle := managedTypes.NewManagedBuffer()
									_ = managedTypes.WriteManagedVecOfManagedBuffers(nil, keysHandle)

									msgHandle := managedTypes.NewManagedBufferFromBytes([]byte("message"))
									sigHandle := managedTypes.NewManagedBufferFromBytes(make([]byte, 48))

									res := vmhooks.ManagedVerifyBLSWithHost(
										host,
										keysHandle,
										msgHandle,
										sigHandle,
										"verifyBLSAggregatedSignature")

									allErrors := host.Runtime().GetAllErrors()
									if flagEnabled {
										require.Equal(t, int32(1), res)
										require.ErrorIs(t, allErrors, vmhost.ErrInvalidArgument)
									}

									return parentInstance
								})
							}),
					).
					WithInput(test.CreateTestContractCallInputBuilder().
						WithRecipientAddr(test.ParentAddress).
						WithGasProvided(100000000).
						WithFunction("testEmptySigners").
						Build()).
					WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
						enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
						enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
							if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
								return flagEnabled
							}
							return true
						}
					}).
					AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

				assert.Nil(t, err)
			})
		})
	}
}

func TestHost_GetCodeHash_UpgradedInTransaction_Overlay(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testGetCodeHash", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								targetSC := test.ChildAddress

								// Simulate an in-transaction code deployment/upgrade for targetSC
								newCode := []byte("new-upgraded-contract-code-bytes-12345")
								host.Output().DeployCode(vmhost.CodeDeployInput{
									ContractAddress:     targetSC,
									ContractCode:        newCode,
									CodeDeployerAddress: test.ParentAddress,
								})

								codeHash := host.Blockchain().GetCodeHash(targetSC)
								expectedNewHash := host.Runtime().ComputeCodeHash(newCode)

								if flagEnabled {
									require.Equal(t, expectedNewHash, codeHash, "under flag, GetCodeHash must return hash of upgraded code from overlay")
								} else {
									require.NotEqual(t, expectedNewHash, codeHash, "without flag, GetCodeHash returns committed state")
								}

								return parentInstance
							})
						}),
					test.CreateMockContract(test.ChildAddress).
						WithBalance(100000).
						WithMethods(),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(100000000).
					WithFunction("testGetCodeHash").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

			assert.Nil(t, err)
		})
	}
}

func TestHost_MBufferNewFromBytes_MaxBufferCap_FlagGated(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testMaxBuffer", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								hooks := vmhooks.NewVMHooksImpl(host)

								// Request buffer creation exceeding maxBufferCopyBytes (10MB + 1)
								oversizedLength := int32(10*1024*1024 + 1)
								handle := hooks.MBufferNewFromBytes(0, oversizedLength)

								allErrors := host.Runtime().GetAllErrors()
								if flagEnabled {
									require.Equal(t, int32(-1), handle)
									require.ErrorIs(t, allErrors, vmhost.ErrArgOutOfRange)
								}

								return parentInstance
							})
						}),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(100000000).
					WithFunction("testMaxBuffer").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

			assert.Nil(t, err)
		})
	}
}

func TestHost_MBufferGetArgument_And_BigIntGetUnsignedArgument_ConsumeGas(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			largeArg := bytes.Repeat([]byte{0x42}, 1000)
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testArgGas", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								hooks := vmhooks.NewVMHooksImpl(host)

								gasBefore := host.Metering().GasLeft()

								// Call MBufferGetArgument for arg 0
								destHandle := host.ManagedTypes().NewManagedBuffer()
								res := hooks.MBufferGetArgument(0, destHandle)
								require.Equal(t, int32(0), res)

								gasAfterMBuffer := host.Metering().GasLeft()
								consumedByMBuffer := gasBefore - gasAfterMBuffer

								// Call BigIntGetUnsignedArgument for arg 0
								bigIntHandle := host.ManagedTypes().NewBigIntFromInt64(0)
								hooks.BigIntGetUnsignedArgument(0, bigIntHandle)

								gasAfterBigInt := host.Metering().GasLeft()
								consumedByBigInt := gasAfterMBuffer - gasAfterBigInt

								copyPerByte := host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte
								expectedByteGas := uint64(len(largeArg)) * copyPerByte

								if flagEnabled {
									require.GreaterOrEqual(t, consumedByMBuffer, expectedByteGas)
									require.GreaterOrEqual(t, consumedByBigInt, expectedByteGas)
								} else {
									require.Less(t, consumedByBigInt, expectedByteGas)
								}

								return parentInstance
							})
						}),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(100000000).
					WithFunction("testArgGas").
					WithArguments(largeArg).
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

			assert.Nil(t, err)
		})
	}
}

func TestHost_GetArgumentsFromMemory_MaxCapped_FlagGated(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testMaxArgs", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								hooks := vmhooks.NewVMHooksImpl(host)

								// Setup memory: dest at 0, balance at 32, func at 64
								_ = parentInstance.MemStore(0, test.ChildAddress)
								_ = parentInstance.MemStore(32, make([]byte, 32))
								_ = parentInstance.MemStore(64, []byte("someFunc"))

								// TransferValueExecute with 1025 arguments (exceeding maxNumArgumentsFromMemoryCapped = 1024)
								res := hooks.TransferValueExecute(
									0,
									32,
									100000,
									64,
									int32(len("someFunc")),
									1025,
									100,
									200,
								)

								allErrors := host.Runtime().GetAllErrors()
								if flagEnabled {
									require.Equal(t, int32(1), res)
									require.NotNil(t, allErrors)
									require.Contains(t, allErrors.Error(), "numArguments out of range (1025)")
								} else {
									if allErrors != nil {
										require.NotContains(t, allErrors.Error(), "numArguments out of range (1025)")
									}
								}

								return parentInstance
							})
						}),
					test.CreateMockContract(test.ChildAddress).
						WithBalance(100000).
						WithMethods(),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(100000000).
					WithFunction("testMaxArgs").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

			assert.Nil(t, err)
		})
	}
}

func TestAsyncCallback_Fails_RevertsPartialStateAndReturnsError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		failType           string
		expectedReturnCode vmcommon.ReturnCode
	}{
		{
			name:               "execution failed trap in callback",
			failType:           "trap",
			expectedReturnCode: vmcommon.ExecutionFailed,
		},
		{
			name:               "signal user error in callback",
			failType:           "user_error",
			expectedReturnCode: vmcommon.UserError,
		},
		{
			name:               "out of gas in callback",
			failType:           "out_of_gas",
			expectedReturnCode: vmcommon.OutOfGas,
		},
	}

	for _, flagEnabled := range []bool{true, false} {
		flagEnabled := flagEnabled
		for _, tc := range testCases {
			tc := tc
			testName := fmt.Sprintf("%s/flag_%v", tc.name, flagEnabled)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()

				testConfig := makeTestConfig()
				testConfig.GasProvided = 100000
				testConfig.GasLockCost = 1000
				testConfig.IsLegacyAsync = true

				exploitKey := []byte("exploit_key")
				exploitVal := []byte("exploit_val")

				parentContract := test.CreateMockContractOnShard(test.ParentAddress, 0).
					WithBalance(100000).
					WithConfig(testConfig).
					WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
						instanceMock.AddMockMethod("callBack", func() *contextmock.InstanceMock {
							host := instanceMock.Host
							// 1. Partial storage mutation
							_, _ = host.Storage().SetStorage(exploitKey, exploitVal)
							// 2. Partial log emission
							host.Output().WriteLog(test.ParentAddress, [][]byte{[]byte("exploit_topic")}, [][]byte{[]byte("exploit_log")})

							// 3. Trigger configured failure
							switch tc.failType {
							case "trap":
								host.Runtime().FailExecution(errors.New("callback trapped"))
							case "user_error":
								host.Runtime().SignalUserError("callback user error")
							case "out_of_gas":
								_ = host.Metering().UseGasBounded(testConfig.GasProvided + 1)
								host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
							}

							return instanceMock
						})
					})

				asyncArguments := &vmcommon.AsyncArguments{
					CallID:                       []byte{1, 2, 3},
					CallerCallID:                 []byte{3, 2, 1},
					CallbackAsyncInitiatorCallID: []byte{4, 5, 6},
					GasAccumulated:               1,
				}
				arguments := [][]byte{[]byte("thirdparty"), []byte("vault"), {0}}

				builder := test.BuildMockInstanceCallTest(t).
					WithContracts(parentContract).
					WithInput(test.CreateTestContractCallInputBuilder().
						WithCallerAddr(test.ChildAddress).
						WithRecipientAddr(test.ParentAddress).
						WithGasProvided(50000).
						WithFunction("callBack").
						WithAsyncArguments(asyncArguments).
						WithArguments(arguments...).
						WithCallType(vm.AsynchronousCallBack).
						Build()).
					WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
						if enableEpochsHandler, ok := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub); ok {
							enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
								if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
									return flagEnabled
								}
								return true
							}
						}

						world.SelfShardID = 0
						if world.CurrentBlockInfo == nil {
							world.CurrentBlockInfo = &worldmock.BlockInfo{}
						}
						world.CurrentBlockInfo.BlockRound = 2

						accountHandler, _ := world.GetUserAccount(test.ParentAddress)
						setZeroCodeCosts(host)
						setAsyncCosts(host, testConfig.GasLockCost)

						err := host.Runtime().StartWasmerInstance(test.ParentAddress, testConfig.GasUsedByParent, false)
						require.Nil(t, err)

						fakeInput := &host.Runtime().GetVMInput().VMInput
						fakeInput.GasProvided = 1000
						host.Metering().InitStateFromContractCallInput(fakeInput)

						err = contracts.RegisterAsyncCallToChild(host, testConfig, arguments)
						require.Nil(t, err)

						host.Async().SetCallID(asyncArguments.CallbackAsyncInitiatorCallID)
						host.Async().SetCallIDForCallInGroup(0, 0, asyncArguments.CallerCallID)
						err = host.Async().Save()
						require.Nil(t, err)

						for _, account := range host.Output().GetVMOutput().OutputAccounts {
							for _, storageUpdate := range account.StorageUpdates {
								(accountHandler.(*worldmock.Account)).Storage[string(storageUpdate.Offset)] = storageUpdate.Data
							}
						}
					})
				world := worldmock.NewMockWorld()
				_, vmOutput, err := builder.RunTest(world, true, test.RunTest)
				if flagEnabled {
					require.Nil(t, err)
					require.Equal(t, tc.expectedReturnCode, vmOutput.ReturnCode)
					if outAcc, exists := vmOutput.OutputAccounts[string(test.ParentAddress)]; exists {
						_, hasExploitKey := outAcc.StorageUpdates[string(exploitKey)]
						require.False(t, hasExploitKey, "exploit key must not be present in storage updates")
					}
					for _, logEntry := range vmOutput.Logs {
						for _, topic := range logEntry.Topics {
							require.NotEqual(t, []byte("exploit_topic"), topic, "callback log must not be present")
						}
					}
				} else if err != nil {
					require.Contains(t, err.Error(), "storage update found on failed execution")
				} else {
					require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
				}
			})
		}
	}
}

func TestHost_NestedDeleteContract_WithParentGasAndCosts(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(1000000).
						WithOwnerAddress(test.ParentAddress).
						WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testDelete", func() *contextmock.InstanceMock {
								host := parentInstance.Host

								// Parent consumes some gas on its Wasmer instance first
								err := host.Metering().UseGasBounded(5000)
								require.Nil(t, err)

								childInput := test.CreateTestContractCallInputBuilder().
									WithCallerAddr(test.ParentAddress).
									WithRecipientAddr(test.ChildAddress).
									WithFunction(vmhost.DeleteFunctionName).
									WithGasProvided(20000).
									Build()
								childInput.AsyncArguments = &vmcommon.AsyncArguments{
									CallID:       []byte{},
									CallerCallID: []byte{},
								}

								vmOutput, _, err := host.ExecuteOnDestContext(childInput)
								if flagEnabled {
									require.Nil(t, err)
									require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
								} else {
									// In legacy mode, instance mismatch caused input and output gas to not match
									require.NotNil(t, err)
								}

								return parentInstance
							})
						}),
					test.CreateMockContract(test.ChildAddress).
						WithBalance(0).
						WithOwnerAddress(test.ParentAddress).
						WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
						WithMethods(),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(100000).
					WithFunction("testDelete").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
					if flagEnabled {
						verify.DeletedAccounts(test.ChildAddress)
					}
				})
			assert.Nil(t, err)
		})
	}
}

func TestHost_BigIntToString_LargeIntGas(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testBigIntToString", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								managedTypes := host.ManagedTypes()

								// Create a large big int (> 64 bytes)
								largeBytes := bytes.Repeat([]byte{0xFF}, 200)
								largeInt := big.NewInt(0).SetBytes(largeBytes)
								srcHandle := managedTypes.NewBigInt(largeInt)
								dstHandle := managedTypes.NewManagedBuffer()

								gasBefore := host.Metering().GasLeft()
								vmhooks.BigIntToStringWithHost(host, srcHandle, dstHandle)
								gasAfter := host.Metering().GasLeft()

								gasUsed := gasBefore - gasAfter
								copyPerByte := host.Metering().GasSchedule().BigIntAPICost.CopyPerByteForTooBig
								expectedExtraGas := uint64(200) * copyPerByte
								if flagEnabled {
									require.Equal(t, uint64(483)+expectedExtraGas, gasUsed)
								} else {
									require.Equal(t, uint64(483), gasUsed)
								}

								return parentInstance
							})
						}),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(1000000).
					WithFunction("testBigIntToString").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
				})
			assert.Nil(t, err)
		})
	}
}

func TestHost_BigIntSqrt_BoundsCheckUnderFlag(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testBigIntSqrt", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								managedTypes := host.ManagedTypes()

								// 8008 bits > 8000 maxSqrtLen
								largeBytes := bytes.Repeat([]byte{0xFF}, 1001)
								largeInt := big.NewInt(0).SetBytes(largeBytes)
								srcHandle := managedTypes.NewBigInt(largeInt)
								dstHandle := managedTypes.NewBigInt(big.NewInt(0))

								vmHooks := vmhooks.NewVMHooksImpl(host)
								vmHooks.BigIntSqrt(dstHandle, srcHandle)

								return parentInstance
							})
						}),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(1000000).
					WithFunction("testBigIntSqrt").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.AttributeExtraGasUsageFlag {
							return false // Ensure AttributeExtraGasUsageFlag is false to test FixCorrectSmallIntAndGasLimits
						}
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					if flagEnabled {
						verify.ReturnCode(vmcommon.ExecutionFailed)
					} else {
						verify.Ok()
					}
				})
			assert.Nil(t, err)
		})
	}
}

func TestHost_NonPayable_AsyncCalleeCheck(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			world := worldmock.NewMockWorld()
			world.AcctMap.PutAccount(&worldmock.Account{
				Address: test.ChildAddress,
				Balance: big.NewInt(100000),
			})
			world.AcctMap.PutAccount(&worldmock.Account{
				Address:         test.ParentAddress,
				Code:            []byte("some code"),
				IsSmartContract: true,
			})
			enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
				IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
					if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
						return flagEnabled
					}
					return true
				},
			}
			host := test.NewTestHostBuilder(t).
				WithBlockchainHook(world).
				WithEnableEpochsHandler(enableEpochsHandler).
				Build()

			vmInput := &vmcommon.ContractCallInput{
				VMInput: vmcommon.VMInput{
					CallType:   vm.AsynchronousCall,
					CallerAddr: test.ChildAddress,
				},
				RecipientAddr: test.ParentAddress,
			}
			host.Runtime().InitStateFromContractCallInput(vmInput)
			host.Output().InitState()

			// Destination non-payable (parent contract)
			err := host.Output().TransferValueOnly(test.ParentAddress, test.ChildAddress, big.NewInt(100), true)
			if flagEnabled {
				require.Equal(t, vmhost.ErrAccountNotPayable, err)
			} else {
				// Legacy allowed async callee to bypass payability
				require.Nil(t, err)
			}
		})
	}
}

func TestHost_OtherVM_FailedCallDoesNotPoisonOutput(t *testing.T) {
	t.Parallel()

	otherVMAddr := make([]byte, 32)
	otherVMAddr[8] = 0x01
	otherVMAddr[9] = 0x01

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(100000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testOtherVM", func() *contextmock.InstanceMock {
						host := parentInstance.Host

						resHandle := host.ManagedTypes().NewManagedBuffer()
						ret := vmhooks.ManagedExecuteOnDestContextWithErrorReturnWithHost(
							host,
							50000,
							big.NewInt(0),
							"callOther",
							otherVMAddr,
							nil,
							resHandle,
						)
						require.NotEqual(t, int32(0), ret)

						require.Equal(t, vmcommon.Ok, host.Output().ReturnCode())

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(200000).
			WithFunction("testOtherVM").
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
			enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
				return true
			}

			world.AcctMap.PutAccount(&worldmock.Account{
				Address:         otherVMAddr,
				Code:            []byte("some code"),
				IsSmartContract: true,
			})
			if world.OtherVMOutputMap == nil {
				world.OtherVMOutputMap = make(map[string]*vmcommon.VMOutput)
			}
			world.OtherVMOutputMap[string([]byte{0x01, 0x01})] = &vmcommon.VMOutput{
				ReturnCode:    vmcommon.ExecutionFailed,
				ReturnMessage: "failed in other VM",
				GasRemaining:  20000,
			}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
		})
	assert.Nil(t, err)
}

func TestHost_IndirectDeployment_GasDeduction_Finding19(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			childCode := test.GetTestSCCode("init-correct", "../../")
			childAddress := test.MakeTestSCAddressWithDefaultVM("newAddress")
			parentAddress := test.ParentAddress

			test.BuildInstanceCallTest(t).
				WithContracts(
					test.CreateInstanceContract(parentAddress).
						WithCode(test.GetTestSCCode("deployer", "../../")).
						WithBalance(1000),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(parentAddress).
					WithFunction("deployChildContract").
					WithGasProvided(1_000_000).
					WithArguments([]byte{'A'}, []byte{0}).
					WithCurrentTxHash([]byte("txhash")).
					Build()).
				WithSetup(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
					stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, uint32, error) {
						if bytes.Equal(address, parentAddress) && bytes.Equal(key, []byte{'A'}) {
							return childCode, 0, nil
						}
						return nil, 0, nil
					}
					stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
						return childAddress, nil
					}
				}).
				AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
					verify.Ok().
						Balance(parentAddress, 1000).
						BalanceDelta(childAddress, 42).
						Code(childAddress, childCode)
					if flagEnabled {
						// Under fix: parent's GetCode (5000) and AoT (652) initial costs are preserved
						verify.GasUsed(parentAddress, 22047).
							GasRemaining(977363)
					} else {
						// Without fix: parent's initial costs were wiped out / overwritten
						verify.GasUsed(parentAddress, 16395).
							GasRemaining(983015)
					}
				})
		})
	}
}

func TestHost_Wasmer2Instance_DeterminismAndSafety(t *testing.T) {
	t.Parallel()

	// 1. Monotonic unique instance IDs across allocations
	inst1 := wasmer2.NewEmptyInstanceForTest()
	inst2 := wasmer2.NewEmptyInstanceForTest()
	require.NotNil(t, inst1)
	require.NotNil(t, inst2)
	require.NotEmpty(t, inst1.ID())
	require.NotEmpty(t, inst2.ID())
	require.NotEqual(t, inst1.ID(), inst2.ID(), "Each newly allocated instance must have a unique monotonic ID")

	// 2. Thread safety and clean lifecycle without crash
	require.False(t, inst1.IsAlreadyCleaned())
	cleaned := inst1.Clean()
	require.False(t, cleaned, "Cleaning empty instance with nil cgoInstance should return false")
	// Second clean is idempotent
	cleanedAgain := inst1.Clean()
	require.False(t, cleanedAgain)

	// Setting and getting breakpoint on cleaned/empty instance must not panic
	inst1.SetBreakpointValue(uint64(vmhost.BreakpointExecutionFailed))
	val := inst1.GetBreakpointValue()
	require.Equal(t, uint64(0), val)
}

func TestHost_CrossShardForwardedGasAccounting_Finding1(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			crossShardDest := make([]byte, 32)
			crossShardDest[0] = 0x99

			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("parentEntry", func() *contextmock.InstanceMock {
								host := parentInstance.Host

								// Step 1: Parent forwards 50,000 gas cross-shard
								err := host.Metering().UseGasBounded(50000)
								require.Nil(t, err)

								err = host.Output().Transfer(
									crossShardDest,
									test.ParentAddress,
									50000,
									0,
									big.NewInt(0),
									nil,
									[]byte("targetFunc"),
									vm.DirectCall,
								)
								require.Nil(t, err)

								// Step 2: Parent executes sibling child call with provided gas = 20,000 (< 50,000)
								childInput := test.CreateTestContractCallInputBuilder().
									WithCallerAddr(test.ParentAddress).
									WithRecipientAddr(test.ChildAddress).
									WithFunction("childEndpoint").
									WithGasProvided(20000).
									Build()
								childInput.AsyncArguments = &vmcommon.AsyncArguments{
									CallID:       []byte{},
									CallerCallID: []byte{},
								}

								childOutput, _, childErr := host.ExecuteOnDestContext(childInput)
								if flagEnabled {
									require.Nil(t, childErr)
									require.Equal(t, vmcommon.Ok, childOutput.ReturnCode)
									require.Equal(t, uint64(14967), childOutput.GasRemaining)
								} else {
									// In legacy mode without fix, child inherits parent's 50,000 transfer,
									// causing checkGas to fail with ErrInputAndOutputGasDoesNotMatch.
									require.NotNil(t, childErr)
									require.Equal(t, vmcommon.ExecutionFailed, childOutput.ReturnCode)
								}

								return parentInstance
							})
						}),
					test.CreateMockContract(test.ChildAddress).
						WithBalance(100000).
						WithMethods(func(childInstance *contextmock.InstanceMock, config interface{}) {
							childInstance.AddMockMethod("childEndpoint", func() *contextmock.InstanceMock {
								host := childInstance.Host
								_ = host.Metering().UseGasBounded(5000)
								return childInstance
							})
						}),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(200000).
					WithFunction("parentEntry").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
					world.AcctMap.PutAccount(&worldmock.Account{
						Address: crossShardDest,
					})
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
					if flagEnabled {
						// Parent's cross-shard transfer is preserved in the merged output
						transfers, ok := verify.VmOutput.OutputAccounts[string(crossShardDest)]
						require.True(t, ok)
						require.Equal(t, 1, len(transfers.OutputTransfers))
						require.Equal(t, uint64(50000), transfers.OutputTransfers[0].GasLimit)
					}
				})
			assert.Nil(t, err)
		})
	}
}

func TestHost_ESDTBalanceHooks_NilSafe(t *testing.T) {
	t.Parallel()

	gasProvided := uint64(100000)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(100000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testNilESDT", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						addrHandle := managedTypes.NewManagedBufferFromBytes(test.ParentAddress)
						tokenHandle := managedTypes.NewManagedBufferFromBytes(test.ESDTTestTokenName)
						valHandle := managedTypes.NewBigIntFromInt64(999)

						// 1. ManagedGetESDTBalance
						hooks.ManagedGetESDTBalance(addrHandle, tokenHandle, 0, valHandle)
						require.Nil(t, host.Runtime().GetAllErrors())
						resBigInt, err := managedTypes.GetBigInt(valHandle)
						require.Nil(t, err)
						require.Equal(t, int64(100), resBigInt.Int64())

						// 2. GetESDTBalance
						err = hooks.MemStore(0, test.ParentAddress)
						require.Nil(t, err)
						err = hooks.MemStore(32, test.ESDTTestTokenName)
						require.Nil(t, err)

						retLen := hooks.GetESDTBalance(0, 32, 10, 0, 64)
						require.True(t, retLen > 0)
						require.Nil(t, host.Runtime().GetAllErrors())

						return instanceMock
					})
				})).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
			require.Nil(t, err)

			parentAccount := world.AcctMap.GetAccount(test.ParentAddress)
			if parentAccount != nil {
				_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, 100)
			}
		}).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasProvided).
			WithFunction("testNilESDT").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
		})
	assert.Nil(t, err)
}

func TestHost_BUG285_ManagedGetCodeMetadata_MissingAccount(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			stubBlockchainHook := &contextmock.BlockchainHookStub{
				GetUserAccountCalled: func(address []byte) (vmcommon.UserAccountHandler, error) {
					return nil, nil // account is nil, err is nil
				},
			}
			enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
				IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
					if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
						return flagEnabled
					}
					return true
				},
			}
			host := test.NewTestHostBuilder(t).
				WithBlockchainHook(stubBlockchainHook).
				WithEnableEpochsHandler(enableEpochsHandler).
				Build()
			defer host.Reset()

			vmInput := &vmcommon.ContractCallInput{
				VMInput:       vmcommon.VMInput{CallerAddr: test.ParentAddress, GasProvided: 100000},
				RecipientAddr: test.ParentAddress,
			}
			host.Runtime().InitStateFromContractCallInput(vmInput)
			host.Metering().InitStateFromContractCallInput(&vmInput.VMInput)
			host.Output().InitState()

			managedTypes := host.ManagedTypes()
			hooks := vmhooks.NewVMHooksImpl(host)

			addrHandle := managedTypes.NewManagedBufferFromBytes(make([]byte, 32))
			respHandle := managedTypes.NewManagedBuffer()

			hooks.ManagedGetCodeMetadata(addrHandle, respHandle)

			if flagEnabled {
				require.ErrorIs(t, host.Runtime().GetAllErrors(), vmhost.ErrNilContract)
			} else {
				require.Nil(t, host.Runtime().GetAllErrors())
			}
		})
	}
}

func TestHost_BUG240_CryptoEI_AddEC_InvalidHandle_NoPanic(t *testing.T) {
	t.Parallel()

	host := test.NewTestHostBuilder(t).
		WithBlockchainHook(&contextmock.BlockchainHookStub{}).
		Build()
	defer host.Reset()

	vmInput := &vmcommon.ContractCallInput{
		VMInput:       vmcommon.VMInput{CallerAddr: test.ParentAddress, GasProvided: 100000},
		RecipientAddr: test.ParentAddress,
	}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	host.Metering().InitStateFromContractCallInput(&vmInput.VMInput)
	host.Output().InitState()

	managedTypes := host.ManagedTypes()
	hooks := vmhooks.NewVMHooksImpl(host)

	xRes := managedTypes.NewBigIntFromInt64(0)
	yRes := managedTypes.NewBigIntFromInt64(0)
	x1 := managedTypes.NewBigIntFromInt64(1)
	y1 := managedTypes.NewBigIntFromInt64(2)
	x2 := managedTypes.NewBigIntFromInt64(3)
	y2 := managedTypes.NewBigIntFromInt64(4)

	invalidCurveHandle := int32(9999)
	hooks.AddEC(xRes, yRes, invalidCurveHandle, x1, y1, x2, y2)
	require.ErrorIs(t, host.Runtime().GetAllErrors(), vmhost.ErrNoEllipticCurveUnderThisHandle)
}

func TestHost_BUG242_CryptoEI_ManagedUnmarshalEC_InvalidDataHandle(t *testing.T) {
	t.Parallel()

	var capturedErr error
	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(100000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testUnmarshalEC", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						ecHandle := managedTypes.PutEllipticCurve(elliptic.P256().Params())
						invalidDataHandle := int32(9999)
						xRes := managedTypes.NewBigIntFromInt64(0)
						yRes := managedTypes.NewBigIntFromInt64(0)

						res := hooks.ManagedUnmarshalEC(xRes, yRes, ecHandle, invalidDataHandle)
						require.Equal(t, int32(1), res)
						capturedErr = host.Runtime().GetAllErrors()

						return instanceMock
					})
				})).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testUnmarshalEC").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed()
			require.ErrorIs(t, capturedErr, vmhost.ErrNoManagedBufferUnderThisHandle)
		})
	assert.Nil(t, err)
}

func TestHost_BUG297_BigInt_GetArgument_OutOfBounds(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			var capturedErr error

			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
							instanceMock.AddMockMethod("testGetArg", func() *contextmock.InstanceMock {
								host := instanceMock.Host
								managedTypes := host.ManagedTypes()
								hooks := vmhooks.NewVMHooksImpl(host)

								destHandle := managedTypes.NewBigIntFromInt64(0)
								hooks.BigIntGetUnsignedArgument(999, destHandle)
								capturedErr = host.Runtime().GetAllErrors()

								return instanceMock
							})
						})).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(100000).
					WithFunction("testGetArg").
					Build()).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					if flagEnabled {
						verify.ExecutionFailed()
						require.ErrorIs(t, capturedErr, vmhost.ErrArgOutOfRange)
					} else {
						verify.Ok()
					}
				})
			assert.Nil(t, err)
		})
	}
}

func TestHost_ESDT_NilTokenData_Guards(t *testing.T) {
	t.Parallel()

	newTestHost := func() (vmhost.VMHost, *vmhooks.VMHooksImpl) {
		stubHook := &contextmock.BlockchainHookStub{
			GetESDTTokenCalled: func(address []byte, tokenID []byte, nonce uint64) (*esdt.ESDigitalToken, error) {
				return nil, nil // esdtToken is nil, error is nil
			},
		}
		host := test.NewTestHostBuilder(t).
			WithBlockchainHook(stubHook).
			Build()

		vmInput := &vmcommon.ContractCallInput{
			VMInput:       vmcommon.VMInput{CallerAddr: test.ParentAddress, GasProvided: 10000000},
			RecipientAddr: test.ParentAddress,
		}
		host.Runtime().InitStateFromContractCallInput(vmInput)
		host.Metering().InitStateFromContractCallInput(&vmInput.VMInput)
		host.Output().InitState()
		rawCode := test.GetTestSCCode("counter", "../../")
		err := host.Runtime().StartWasmerInstance(rawCode, 10000000, false)
		require.Nil(t, err)
		return host, vmhooks.NewVMHooksImpl(host)
	}

	t.Run("BigIntGetESDTExternalBalance sets 0 without error", func(t *testing.T) {
		host, hooks := newTestHost()
		defer host.Reset()
		managedTypes := host.ManagedTypes()

		valHandle := managedTypes.NewBigIntFromInt64(0)
		_ = hooks.MemStore(0, test.ParentAddress)
		_ = hooks.MemStore(32, []byte("TEST-123456"))
		hooks.BigIntGetESDTExternalBalance(0, 32, 11, 0, valHandle)
		valBigInt, _ := managedTypes.GetBigInt(valHandle)
		require.Equal(t, int64(0), valBigInt.Int64())
		require.Nil(t, host.Runtime().GetAllErrors())
	})

	t.Run("ManagedGetESDTTokenData returns 0 and empty metadata without error", func(t *testing.T) {
		host, hooks := newTestHost()
		defer host.Reset()
		managedTypes := host.ManagedTypes()

		addrHandle := managedTypes.NewManagedBufferFromBytes(test.ParentAddress)
		tokenHandle := managedTypes.NewManagedBufferFromBytes([]byte("TEST-123456"))
		valHandle := managedTypes.NewBigIntFromInt64(0)
		propHandle := managedTypes.NewManagedBuffer()
		hashHandle := managedTypes.NewManagedBuffer()
		nameHandle := managedTypes.NewManagedBuffer()
		attrHandle := managedTypes.NewManagedBuffer()
		creatorHandle := managedTypes.NewManagedBuffer()
		royaltiesHandle := managedTypes.NewBigIntFromInt64(0)
		urisHandle := managedTypes.NewManagedBuffer()

		hooks.ManagedGetESDTTokenData(addrHandle, tokenHandle, 0, valHandle, propHandle, hashHandle, nameHandle, attrHandle, creatorHandle, royaltiesHandle, urisHandle)
		valBigInt, _ := managedTypes.GetBigInt(valHandle)
		require.Equal(t, int64(0), valBigInt.Int64())
		require.Nil(t, host.Runtime().GetAllErrors())
	})

	t.Run("ManagedGetESDTTokenType returns 0 without error", func(t *testing.T) {
		host, hooks := newTestHost()
		defer host.Reset()
		managedTypes := host.ManagedTypes()

		addrHandle := managedTypes.NewManagedBufferFromBytes(test.ParentAddress)
		tokenHandle := managedTypes.NewManagedBufferFromBytes([]byte("TEST-123456"))
		typeHandle := managedTypes.NewBigIntFromInt64(0)

		hooks.ManagedGetESDTTokenType(addrHandle, tokenHandle, 0, typeHandle)
		typeBigInt, _ := managedTypes.GetBigInt(typeHandle)
		require.Equal(t, int64(0), typeBigInt.Int64())
		require.Nil(t, host.Runtime().GetAllErrors())
	})

	t.Run("ManagedIsESDTFrozen returns 0 without error", func(t *testing.T) {
		host, hooks := newTestHost()
		defer host.Reset()
		managedTypes := host.ManagedTypes()

		addrHandle := managedTypes.NewManagedBufferFromBytes(test.ParentAddress)
		tokenHandle := managedTypes.NewManagedBufferFromBytes([]byte("TEST-123456"))

		resFrozen := hooks.ManagedIsESDTFrozen(addrHandle, tokenHandle, 0)
		require.Equal(t, int32(0), resFrozen)
		require.Nil(t, host.Runtime().GetAllErrors())
	})

	t.Run("GetESDTTokenData returns 0 length and sets 0 value without error", func(t *testing.T) {
		host, hooks := newTestHost()
		defer host.Reset()
		managedTypes := host.ManagedTypes()

		_ = hooks.MemStore(0, test.ParentAddress)
		_ = hooks.MemStore(32, []byte("TEST-123456"))
		valHandle := managedTypes.NewBigIntFromInt64(0)
		royaltiesHandle := managedTypes.NewBigIntFromInt64(0)

		resTokenData := hooks.GetESDTTokenData(0, 32, 11, 0, valHandle, 100, 200, 300, 400, 500, royaltiesHandle, 600)
		require.Equal(t, int32(0), resTokenData)
		valBigInt, _ := managedTypes.GetBigInt(valHandle)
		require.Equal(t, int64(0), valBigInt.Int64())
		require.Nil(t, host.Runtime().GetAllErrors())
	})
}

func TestHost_BigFloatFloorCeil_RoundingNearZero(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
				IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
					if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
						return flagEnabled
					}
					return true
				},
			}
			host := test.NewTestHostBuilder(t).
				WithBlockchainHook(&contextmock.BlockchainHookStub{}).
				WithEnableEpochsHandler(enableEpochsHandler).
				Build()
			defer host.Reset()

			vmInput := &vmcommon.ContractCallInput{
				VMInput:       vmcommon.VMInput{CallerAddr: test.ParentAddress, GasProvided: 10000000},
				RecipientAddr: test.ParentAddress,
			}
			host.Runtime().InitStateFromContractCallInput(vmInput)
			host.Metering().InitStateFromContractCallInput(&vmInput.VMInput)
			host.Output().InitState()
			rawCode := test.GetTestSCCode("counter", "../../")
			err := host.Runtime().StartWasmerInstance(rawCode, 10000000, false)
			require.Nil(t, err)

			hooks := vmhooks.NewVMHooksImpl(host)
			managedTypes := host.ManagedTypes()

			negHalfHandle, _ := managedTypes.PutBigFloat(big.NewFloat(-0.5))
			posHalfHandle, _ := managedTypes.PutBigFloat(big.NewFloat(0.5))

			floorDestHandle := managedTypes.NewBigIntFromInt64(0)
			ceilDestHandle := managedTypes.NewBigIntFromInt64(0)

			// Test Floor(-0.5)
			hooks.BigFloatFloor(floorDestHandle, negHalfHandle)
			floorResult, _ := managedTypes.GetBigInt(floorDestHandle)
			if flagEnabled {
				require.Equal(t, int64(-1), floorResult.Int64())
			} else {
				require.Equal(t, int64(0), floorResult.Int64())
			}

			// Test Ceil(+0.5)
			hooks.BigFloatCeil(ceilDestHandle, posHalfHandle)
			ceilResult, _ := managedTypes.GetBigInt(ceilDestHandle)
			if flagEnabled {
				require.Equal(t, int64(1), ceilResult.Int64())
			} else {
				require.Equal(t, int64(0), ceilResult.Int64())
			}
		})
	}
}

func TestHost_BigFloatPow_NegativeExponent(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
				IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
					if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
						return flagEnabled
					}
					return true
				},
			}
			host := test.NewTestHostBuilder(t).
				WithBlockchainHook(&contextmock.BlockchainHookStub{}).
				WithEnableEpochsHandler(enableEpochsHandler).
				Build()
			defer host.Reset()

			vmInput := &vmcommon.ContractCallInput{
				VMInput:       vmcommon.VMInput{CallerAddr: test.ParentAddress, GasProvided: 10000000},
				RecipientAddr: test.ParentAddress,
			}
			host.Runtime().InitStateFromContractCallInput(vmInput)
			host.Metering().InitStateFromContractCallInput(&vmInput.VMInput)
			host.Output().InitState()
			rawCode := test.GetTestSCCode("counter", "../../")
			err := host.Runtime().StartWasmerInstance(rawCode, 10000000, false)
			require.Nil(t, err)

			hooks := vmhooks.NewVMHooksImpl(host)
			managedTypes := host.ManagedTypes()

			baseHandle, _ := managedTypes.PutBigFloat(big.NewFloat(2.0))
			destHandle, _ := managedTypes.PutBigFloat(big.NewFloat(0))

			hooks.BigFloatPow(destHandle, baseHandle, -1)
			err = host.Runtime().GetAllErrors()
			if flagEnabled {
				require.ErrorIs(t, err, vmhost.ErrExponentTooBigOrTooSmall)
			}
		})
	}
}

func TestHost_IsReservedFunctionName_ExtraGasUsage(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
							instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
								host := instanceMock.Host
								managedTypes := host.ManagedTypes()
								hooks := vmhooks.NewVMHooksImpl(host)

								// 200 bytes name (> safeAttributeLength 100)
								longName := bytes.Repeat([]byte("a"), 200)
								nameHandle := managedTypes.NewManagedBufferFromBytes(longName)

								baseCost := host.Metering().GasSchedule().BaseOpsAPICost.IsReservedFunctionName
								dataCopyPerByte := host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte
								expectedExtraGas := uint64(len(longName)) * dataCopyPerByte

								gasBefore := host.Metering().GasLeft()
								res := hooks.IsReservedFunctionName(nameHandle)
								require.Equal(t, int32(0), res)
								gasUsed := gasBefore - host.Metering().GasLeft()

								if flagEnabled {
									require.Equal(t, baseCost+expectedExtraGas, gasUsed)
								} else {
									require.Equal(t, baseCost, gasUsed)
								}

								return contextmock.GetMockInstance(host)
							})
						})).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(1000000).
					WithFunction("test").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.AttributeExtraGasUsageFlag {
							return false
						}
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
				})
			require.Nil(t, err)
		})
	}
}

func TestHost_ManagedIsESDTLimitedTransfer_ExtraGasUsage(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(100000).
						WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
							instanceMock.AddMockMethod("test", func() *contextmock.InstanceMock {
								host := instanceMock.Host
								managedTypes := host.ManagedTypes()
								hooks := vmhooks.NewVMHooksImpl(host)

								// 200 bytes token ID (> safeAttributeLength 100)
								longTokenID := bytes.Repeat([]byte("T"), 200)
								tokenIDHandle := managedTypes.NewManagedBufferFromBytes(longTokenID)

								baseCost := host.Metering().GasSchedule().BaseOpsAPICost.GetExternalBalance
								dataCopyPerByte := host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte
								expectedExtraGas := uint64(len(longTokenID)) * dataCopyPerByte

								gasBefore := host.Metering().GasLeft()
								res := hooks.ManagedIsESDTLimitedTransfer(tokenIDHandle)
								require.Equal(t, int32(0), res)
								gasUsed := gasBefore - host.Metering().GasLeft()

								if flagEnabled {
									require.Equal(t, baseCost+expectedExtraGas, gasUsed)
								} else {
									require.Equal(t, baseCost, gasUsed)
								}

								return contextmock.GetMockInstance(host)
							})
						})).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(1000000).
					WithFunction("test").
					Build()).
				WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
					enableEpochsHandler := host.EnableEpochsHandler().(*worldmock.EnableEpochsHandlerStub)
					enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
						if flag == vmhost.AttributeExtraGasUsageFlag {
							return false
						}
						if flag == vmhost.FixCorrectSmallIntAndGasLimitsFlag {
							return flagEnabled
						}
						return true
					}
				}).
				AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
					verify.Ok()
				})
			require.Nil(t, err)
		})
	}
}


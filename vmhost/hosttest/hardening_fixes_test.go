package hostCoretest

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-vm-go/executor"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHost_ExecuteOnDestContextWithErrorReturn_CallsCounterCleanup(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("flag="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			_, err := test.BuildMockInstanceCallTest(t).
				WithContracts(
					test.CreateMockContract(test.ParentAddress).
						WithBalance(1000000).
						WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
							parentInstance.AddMockMethod("testCallFailingChild", func() *contextmock.InstanceMock {
								host := parentInstance.Host
								initialIsComplete := host.Async().IsComplete()
								require.True(t, initialIsComplete)

								// Call child contract with non-existent method using error return hook
								retCode, err := vmhooks.ManagedExecuteOnDestContextWithErrorReturnWithHost(
									host,
									100000,
									big.NewInt(0),
									"nonExistentMethod",
									test.ChildAddress,
									nil,
									0,
								)
								handleHookError(host, err)
								// Returns 1 indicating error in child execution
								require.Equal(t, int32(1), retCode)

								// Under FixCorrectSmallIntAndGasLimits, callsCounter must be balanced back to 0
								if flagEnabled {
									require.True(t, host.Async().IsComplete(), "callsCounter must be balanced when flag is enabled")
								} else {
									require.False(t, host.Async().IsComplete(), "callsCounter is leaked when flag is disabled")
								}

								return parentInstance
							})
						}),
					test.CreateMockContract(test.ChildAddress).
						WithBalance(0).
						WithMethods(),
				).
				WithInput(test.CreateTestContractCallInputBuilder().
					WithRecipientAddr(test.ParentAddress).
					WithGasProvided(1000000).
					WithFunction("testCallFailingChild").
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

func TestHost_VerifyCustomSecp256k1_MemoryBoundaryGracefulFailure(t *testing.T) {
	t.Parallel()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testSecpBoundary", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						hooks := vmhooks.NewVMHooksLegacy(host)

						memLen := parentInstance.Memory.Length()
						// Place key at offset 0 (33 bytes compressed key)
						keyBytes := make([]byte, 33)
						keyBytes[0] = 0x02 // compressed format prefix
						copy(parentInstance.Memory.Data()[:33], keyBytes)

						// Message offset at 50, len 32
						messageOffset := executor.MemPtr(50)
						messageLength := executor.MemLength(32)

						// Place sigOffset at memory boundary: memLen - 1
						// MemLoad(sigOffset, 2) will only load 1 byte, triggering the bounds check
						sigOffset := executor.MemPtr(memLen - 1)

						var ret int32
						require.NotPanics(t, func() {
							ret = hooks.VerifyCustomSecp256k1(
								0,
								33,
								messageOffset,
								messageLength,
								sigOffset,
								0, // sha256
							)
						})

						require.Equal(t, int32(1), ret)
						require.ErrorIs(t, host.Runtime().GetAllErrors(), executor.ErrMemoryBadBounds)

						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(1000000).
			WithFunction("testSecpBoundary").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {})

	assert.Nil(t, err)
}

func TestHost_CallFunctionOnOtherVM_ReadOnly(t *testing.T) {
	t.Parallel()

	otherVMAddress := []byte("other-vm-address-32-bytes-pad000")

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000000).
				WithMethods(func(parentInstance *contextmock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testOtherVMReadOnly", func() *contextmock.InstanceMock {
						host := parentInstance.Host
						host.Runtime().SetReadOnly(true)

						otherVMInput := test.CreateTestContractCallInputBuilder().
							WithCallerAddr(test.ParentAddress).
							WithRecipientAddr(otherVMAddress).
							WithFunction("someFunction").
							WithGasProvided(10000).
							Build()

						_, _, errCall := host.ExecuteOnDestContext(otherVMInput)
						require.Equal(t, vmhost.ErrInvalidCallOnReadOnlyMode, errCall)

						return parentInstance
					})
				}),
			test.CreateMockContract(otherVMAddress).
				WithVMType([]byte("otherVM")).
				WithBalance(0).
				WithMethods(),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testOtherVMReadOnly").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
		})

	assert.Nil(t, err)
}

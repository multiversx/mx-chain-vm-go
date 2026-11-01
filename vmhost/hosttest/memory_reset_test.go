package hostCoretest

import (
	"testing"
	"unsafe"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

const memResetPageSize = 65536

type memResetWarmGetter interface {
	GetWarmInstance(codeHash []byte) (executor.InstanceNew, bool)
}

// memResetMemoryObserver is what this test needs of an instance beyond executor.InstanceNew:
// its memory, from outside a call.
//
// Not on any interface, because no production code wants it - memory belongs to the per-call
// instance state. The backend offers it for exactly this.
type memResetMemoryObserver interface {
	MemLength() uint32
	MemDump() []byte
}

// memResetInstance returns the instance the host holds for the parent
// contract: the active one right after a call, or the one parked in the warm
// instance cache.
func memResetInstance(t *testing.T, host vmhost.VMHost) memResetMemoryObserver {
	if instance := host.Runtime().GetInstance(); instance != nil && !instance.IsInterfaceNil() {
		return asMemoryObserver(t, instance)
	}
	tracker, ok := host.Runtime().GetInstanceTracker().(memResetWarmGetter)
	require.True(t, ok, "instance tracker does not expose the warm cache")
	codeHash := host.Blockchain().GetCodeHash(test.ParentAddress)
	instance, ok := tracker.GetWarmInstance(codeHash)
	require.True(t, ok, "no warm instance for the parent contract")
	return asMemoryObserver(t, instance)
}

// asMemoryObserver narrows an instance to the memory view this test needs.
func asMemoryObserver(t *testing.T, instance executor.InstanceNew) memResetMemoryObserver {
	observer, ok := instance.(memResetMemoryObserver)
	require.True(t, ok, "instance does not expose its memory")

	return observer
}

func memResetBase(t *testing.T, instance memResetMemoryObserver) uintptr {
	data := instance.MemDump()
	require.NotEmpty(t, data)
	return uintptr(unsafe.Pointer(&data[0]))
}

// A warm instance is Reset() before every reuse. The reset must shrink the
// memory back to its declared minimum without moving it: compiled code for
// static memories has no bounds checks and relies on the guard region that
// was reserved when the instance was created.
func TestExecution_WarmInstance_ResetKeepsMemoryBase(t *testing.T) {
	testCase := test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("memory-reset", "../../")))

	makeInput := func(function string) *vmcommon.ContractCallInput {
		return test.CreateTestContractCallInputBuilder().
			WithGasProvided(10_000_000).
			WithFunction(function).
			Build()
	}

	var base uintptr

	// Cold instance: grow to four pages and leave the instance warm.
	testCase.WithInput(makeInput("growAndDirty")).
		AndAssertResultsWithoutReset(func(host vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok()
			instance := memResetInstance(t, host)
			require.Equal(t, uint32(4*memResetPageSize), instance.MemLength())
			base = memResetBase(t, instance)
		})

	// Warm reuse: Reset() ran before this call. Memory is back at one page and
	// its base address has not moved.
	testCase.WithInput(makeInput("touchMinimum")).
		AndAssertResultsWithoutReset(func(host vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
			verify.Ok()
			instance := memResetInstance(t, host)
			require.Equal(t, uint32(memResetPageSize), instance.MemLength())
			require.Equal(t, base, memResetBase(t, instance), "reset must not move a static memory")
		})

	// After growing and being reset again, the page that was grown, an
	// address deep inside the 4 GiB bound and the last word must all trap.
	// A failed call evicts the warm instance, so grow again before each one.
	for _, function := range []string{"loadGrownPage", "loadHigh", "loadLast"} {
		testCase.WithInput(makeInput("growAndDirty")).
			AndAssertResultsWithoutReset(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
				verify.Ok()
			})
		testCase.WithInput(makeInput(function)).
			AndAssertResultsWithoutReset(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
				verify.ExecutionFailed()
			})
	}
}

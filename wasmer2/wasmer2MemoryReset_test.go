package wasmer2

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

const wasmPage = 65536

// memoryBase returns the address of the first byte of the instance memory, as
// the host sees it through vm_exec_instance_memory_data.
func memoryBase(t *testing.T, instance executor.Instance) uintptr {
	data := instance.MemDump()
	require.NotEmpty(t, data)
	return uintptr(unsafe.Pointer(&data[0]))
}

func newMemoryResetInstance(t *testing.T) executor.Instance {
	code, err := os.ReadFile(filepath.Join("..", "test", "contracts", "memory-reset", "output", "memory-reset.wasm"))
	require.Nil(t, err)

	exec, err := ExecutorFactory().CreateExecutor(executor.ExecutorFactoryArgs{})
	require.Nil(t, err)

	instance, err := exec.NewInstanceWithOptions(code, executor.CompilationOptions{
		GasLimit:           1_000_000_000,
		MaxMemoryGrow:      20,
		MaxMemoryGrowDelta: 10,
	})
	require.Nil(t, err)
	return instance
}

// Warm-instance reuse calls Reset() on a cached instance. The executor uses
// static memories, whose compiled code has no bounds checks and relies on the
// 4 GiB + 2 GiB guard reservation made at instantiation. Reset must keep that
// reservation and the memory base address.
func TestWasmer2Instance_ResetKeepsMemoryBaseAndGuard(t *testing.T) {
	instance := newMemoryResetInstance(t)
	defer instance.Clean()

	base := memoryBase(t, instance)
	require.Equal(t, uint32(wasmPage), instance.MemLength())

	// Grow from inside wasm and dirty the last new page.
	require.Nil(t, instance.CallFunction("growAndDirty"))
	require.Equal(t, uint32(4*wasmPage), instance.MemLength())
	marker, err := instance.MemLoad(3*wasmPage, 4)
	require.Nil(t, err)
	require.Equal(t, []byte{0xed, 0x5e, 0xed, 0x5e}, marker)
	require.Nil(t, instance.CallFunction("loadGrownPage"))

	// This is what instanceTracker.UseWarmInstance does before every reuse.
	require.True(t, instance.Reset())

	require.Equal(t, uint32(wasmPage), instance.MemLength())
	require.Equal(t, base, memoryBase(t, instance), "reset must not move a static memory")
	head, err := instance.MemLoad(0, 4)
	require.Nil(t, err)
	require.Equal(t, []byte{1, 2, 3, 4}, head, "data segments are re-applied")

	// Everything above the minimum must trap again: the page that was grown,
	// an address deep inside the 4 GiB bound, and the last word.
	for _, endpoint := range []string{"loadGrownPage", "loadHigh", "loadLast"} {
		err = instance.CallFunction(endpoint)
		require.ErrorContains(t, err, "out of bounds memory access", "%s must trap after reset", endpoint)
	}

	// Growing again reuses the same mapping and yields zeroed pages.
	require.Nil(t, instance.MemGrow(3))
	require.Equal(t, base, memoryBase(t, instance))
	regrown, err := instance.MemLoad(3*wasmPage, 4)
	require.Nil(t, err)
	require.Equal(t, []byte{0, 0, 0, 0}, regrown)
	require.Nil(t, instance.CallFunction("loadGrownPage"))
}

func TestWasmer2Instance_ResetIsIdempotentAtMinimum(t *testing.T) {
	instance := newMemoryResetInstance(t)
	defer instance.Clean()

	base := memoryBase(t, instance)
	require.Nil(t, instance.MemStore(100, []byte{0xAA, 0xBB}))

	require.True(t, instance.Reset())
	require.True(t, instance.Reset())

	require.Equal(t, base, memoryBase(t, instance))
	require.Equal(t, uint32(wasmPage), instance.MemLength())
	cleared, err := instance.MemLoad(100, 2)
	require.Nil(t, err)
	require.Equal(t, []byte{0, 0}, cleared)
}

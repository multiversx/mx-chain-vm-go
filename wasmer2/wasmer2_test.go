package wasmer2

import (
	"testing"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/assert"
)

func TestCompilationOptions_FFILayout(t *testing.T) {
	t.Parallel()

	opts := executor.CompilationOptions{}

	// CompilationOptions is cast to (*cWasmerCompilationOptions)(unsafe.Pointer(&options))
	// and passed across Cgo to the C/Rust FFI layer (vm_exec_compilation_options_t).
	// Its struct layout, field sizes, and offsets must strictly match the Rust #[repr(C)] struct:
	// pub struct CompilationOptions {
	//     pub gas_limit: u64,
	//     pub unmetered_locals: u64,
	//     pub max_memory_grow: u64,
	//     pub max_memory_grow_delta: u64,
	//     pub opcode_trace: bool,
	//     pub metering: bool,
	//     pub runtime_breakpoints: bool,
	// }
	assert.Equal(t, uintptr(40), unsafe.Sizeof(opts), "CompilationOptions struct size must be exactly 40 bytes")
	assert.Equal(t, uintptr(0), unsafe.Offsetof(opts.GasLimit), "GasLimit offset must be 0")
	assert.Equal(t, uintptr(8), unsafe.Offsetof(opts.UnmeteredLocals), "UnmeteredLocals offset must be 8")
	assert.Equal(t, uintptr(16), unsafe.Offsetof(opts.MaxMemoryGrow), "MaxMemoryGrow offset must be 16")
	assert.Equal(t, uintptr(24), unsafe.Offsetof(opts.MaxMemoryGrowDelta), "MaxMemoryGrowDelta offset must be 24")
	assert.Equal(t, uintptr(32), unsafe.Offsetof(opts.OpcodeTrace), "OpcodeTrace offset must be 32")
	assert.Equal(t, uintptr(33), unsafe.Offsetof(opts.Metering), "Metering offset must be 33")
	assert.Equal(t, uintptr(34), unsafe.Offsetof(opts.RuntimeBreakpoints), "RuntimeBreakpoints offset must be 34")

	// Field sizes
	assert.Equal(t, uintptr(8), unsafe.Sizeof(opts.GasLimit))
	assert.Equal(t, uintptr(8), unsafe.Sizeof(opts.UnmeteredLocals))
	assert.Equal(t, uintptr(8), unsafe.Sizeof(opts.MaxMemoryGrow))
	assert.Equal(t, uintptr(8), unsafe.Sizeof(opts.MaxMemoryGrowDelta))
	assert.Equal(t, uintptr(1), unsafe.Sizeof(opts.OpcodeTrace))
	assert.Equal(t, uintptr(1), unsafe.Sizeof(opts.Metering))
	assert.Equal(t, uintptr(1), unsafe.Sizeof(opts.RuntimeBreakpoints))
}

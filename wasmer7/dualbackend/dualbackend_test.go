//go:build dualbackend

// Package dualbackend answers gate G5 of PLAN-go-ffi.md: can one process hold both
// executor backends at once?
//
// The plan's preferred shape for serving Wasmer 2.2 and Wasmer 7 is a single C library
// carrying both, chosen by a flag, so that the differential run is a flag rather than two
// builds. The one thing that could rule it out is not size or symbols but **trap handling**.
// Both runtimes install process-global signal handlers for the faults a wasm trap raises —
// wasmer 2.2 forcibly, through vm_force_sighandler_reinstall, which wasmer2legacy.CreateExecutor
// calls on every executor it builds; wasmer 7 lazily, when it first builds a store. Whoever
// installs last wins unless the handlers chain, and the failure mode is not a red test but a
// mishandled fault, so it has to be observed rather than assumed.
//
// This package is behind a build tag because it is the one place that deliberately links
// both libraries into one binary, which the plan otherwise forbids: libvmexeccapi and
// libvmexeccapi_w7 each contain their own copy of wasmer 2.2, so a shipped binary must link
// exactly one.
//
// Handler installation is global and guarded by a Once, so the two initialisation orders
// cannot both be exercised in one process. Run them separately:
//
//	go test -tags dualbackend ./wasmer7/dualbackend/ -run TestWasmer7TrapsWithWasmer2First -v
//	go test -tags dualbackend ./wasmer7/dualbackend/ -run TestWasmer2TrapsWithWasmer7First -v
package dualbackend

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2legacy"
	"github.com/multiversx/mx-chain-vm-go/wasmer7"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
	"github.com/stretchr/testify/require"
)

// trapContractHex is a wasm binary, given as bytes because wasmer 2.2 and wasmer 7 need not
// agree about accepting WAT text, and this experiment must hand both backends exactly the
// same input. Its source:
//
//	(module
//	  (memory (export "memory") 1)
//	  (func (export "trapUnreachable") (unreachable))
//	  (func (export "trapMemory") (drop (i32.load (i32.const 2147483632))))
//	  (func (export "noTrap"))
//	)
//
// Two trap kinds, because they take different routes to the handler: `unreachable` is an
// illegal instruction, while the out-of-range load faults on a guard page or fails a bounds
// check. A handler that is installed but wrong can catch one and miss the other.
const trapContractHex = "0061736d010000000104016000000304030000000503010001073204066d656d6f72790" +
	"2000f74726170556e726561636861626c6500000a747261704d656d6f72790001066e6f5472617000020a1" +
	"5030300000b0c0041f0ffffff072802001a0b02000b"

// stubVMHooks satisfies the interface without implementing anything: the trap contract
// imports no hooks, so any call would be a bug in the fixture rather than a gap here.
type stubVMHooks struct {
	executor.VMHooks
}

func (stubVMHooks) VMHooksForCall(_ executor.InstanceState) executor.VMHooks {
	return stubVMHooks{}
}

func (stubVMHooks) VMHooksForInstanceLifetime() executor.VMHooks {
	return stubVMHooks{}
}

func trapContract(t *testing.T) []byte {
	code, err := hex.DecodeString(trapContractHex)
	require.Nil(t, err)

	return code
}

func uniformOpcodeCost(cost uint32) *executor.WASMOpcodeCost {
	opcodeCost := &executor.WASMOpcodeCost{}

	fields := reflect.ValueOf(opcodeCost).Elem()
	for i := 0; i < fields.NumField(); i++ {
		fields.Field(i).SetUint(uint64(cost))
	}

	return opcodeCost
}

// --- Backend setup -----------------------------------------------------------

// newWasmer2Instance builds a legacy executor, which is what forcibly reinstalls the
// wasmer 2.2 trap handlers.
func newWasmer2Instance(t *testing.T) executor.Instance {
	wasmerExecutor, err := wasmer2legacy.ExecutorFactory().CreateExecutor(executor.ExecutorFactoryArgs{
		VMHooksProvider: stubVMHooks{},
		OpcodeVersion:   executor.OpcodeVersionV2,
		OpcodeCosts:     uniformOpcodeCost(1),
	})
	require.Nil(t, err)

	instance, err := wasmerExecutor.NewInstanceWithOptions(trapContract(t), executor.CompilationOptions{
		GasLimit:           10_000_000,
		MaxMemoryGrow:      10,
		MaxMemoryGrowDelta: 10,
		Metering:           true,
		RuntimeBreakpoints: true,
	})
	require.Nil(t, err)
	t.Cleanup(func() { instance.Clean() })

	return instance
}

// newWasmer7Instance builds a new-FFI executor, which is what makes wasmer 7 install its
// own handlers.
func newWasmer7Instance(t *testing.T) *wasmer7.Wasmer7Instance {
	wasmerExecutor, err := wasmer7.NewExecutor(wasmercommon.ExecutorConfig{
		OpcodeVersion:  executor.OpcodeVersionV2,
		WASMOpcodeCost: uniformOpcodeCost(1),
	})
	require.Nil(t, err)
	t.Cleanup(wasmerExecutor.Close)

	instance, err := wasmerExecutor.NewInstance(
		trapContract(t),
		wasmercommon.CompilationOptions{MaxMemoryGrow: 10, MaxMemoryGrowDelta: 10},
		stubVMHooks{},
	)
	require.Nil(t, err)
	t.Cleanup(instance.Close)

	return instance
}

// --- The experiment ----------------------------------------------------------

// TestWasmer7TrapsWithWasmer2First is the order that matters most.
//
// wasmer2legacy.CreateExecutor calls vm_force_sighandler_reinstall unconditionally, so it is the
// one that can overwrite a handler wasmer 7 has already installed. If it does so without
// chaining, the traps below take down the process instead of coming back as results.
func TestWasmer7TrapsWithWasmer2First(t *testing.T) {
	legacy := newWasmer2Instance(t)
	require.Equal(t, executor.CallOk, legacy.CallFunction("noTrap").Kind, "legacy backend is alive")

	instance := newWasmer7Instance(t)

	result, err := instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallOk, result.Kind, "new backend is alive")

	for _, funcName := range []string{"trapUnreachable", "trapMemory"} {
		t.Run(funcName, func(t *testing.T) {
			result, err := instance.Call(funcName, 10_000_000)
			require.Nil(t, err)
			require.Equal(t, wasmercommon.CallRuntimeError, result.Kind,
				"the trap was caught and reported: %+v", result)
			t.Logf("%s reported: %s", funcName, result.Message)
		})
	}

	// Still usable after the traps, which is what distinguishes a handled trap from a
	// handler that merely did not crash.
	result, err = instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallOk, result.Kind, "new backend survived its traps")

	require.Equal(t, executor.CallOk, legacy.CallFunction("noTrap").Kind, "legacy backend survived them too")
}

// TestWasmer2TrapsWithWasmer7First is the mirror: wasmer 7 installs first, and the legacy
// backend has to keep working through its own traps afterwards.
func TestWasmer2TrapsWithWasmer7First(t *testing.T) {
	instance := newWasmer7Instance(t)

	result, err := instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallOk, result.Kind, "new backend is alive")

	legacy := newWasmer2Instance(t)
	require.Equal(t, executor.CallOk, legacy.CallFunction("noTrap").Kind, "legacy backend is alive")

	for _, funcName := range []string{"trapUnreachable", "trapMemory"} {
		t.Run(funcName, func(t *testing.T) {
			result := legacy.CallFunction(funcName)
			require.Equal(t, executor.CallRuntimeError, result.Kind,
				"the trap was caught and reported")
			t.Logf("%s reported: %v", funcName, result.Err)
		})
	}

	require.Equal(t, executor.CallOk, legacy.CallFunction("noTrap").Kind, "legacy backend survived its traps")

	result, err = instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallOk, result.Kind, "new backend survived them too")
}

//go:build dualbackend

// Package dualbackend answers gate G5 of PLAN-go-ffi.md: can one process hold both
// executor backends at once?
//
// The plan's preferred shape for serving Wasmer 2.2 and Wasmer 6 is a single C library
// carrying both, chosen by a flag, so that the differential run is a flag rather than two
// builds. The one thing that could rule it out is not size or symbols but **trap handling**.
// Both runtimes install process-global signal handlers for the faults a wasm trap raises —
// wasmer 2.2 forcibly, through vm_force_sighandler_reinstall, which wasmer2.CreateExecutor
// calls on every executor it builds; wasmer 6 lazily, when it first builds a store. Whoever
// installs last wins unless the handlers chain, and the failure mode is not a red test but a
// mishandled fault, so it has to be observed rather than assumed.
//
// This package is behind a build tag because it is the one place that deliberately links
// both libraries into one binary, which the plan otherwise forbids: libvmexeccapi and
// libvmexeccapi2 each contain their own copy of wasmer 2.2, so a shipped binary must link
// exactly one.
//
// Handler installation is global and guarded by a Once, so the two initialisation orders
// cannot both be exercised in one process. Run them separately:
//
//	go test -tags dualbackend ./wasmer3/dualbackend/ -run TestWasmer3TrapsWithWasmer2First -v
//	go test -tags dualbackend ./wasmer3/dualbackend/ -run TestWasmer2TrapsWithWasmer3First -v
package dualbackend

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/multiversx/mx-chain-vm-go/wasmer3"
	"github.com/stretchr/testify/require"
)

// trapContractHex is a wasm binary, given as bytes because wasmer 2.2 and wasmer 6 need not
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

func (stubVMHooks) VMHooksForCall(_ *wasmer3.InstanceState) executor.VMHooks {
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
	wasmerExecutor, err := wasmer2.ExecutorFactory().CreateExecutor(executor.ExecutorFactoryArgs{
		VMHooks:       stubVMHooks{},
		OpcodeVersion: executor.OpcodeVersionV2,
		OpcodeCosts:   uniformOpcodeCost(1),
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

// newWasmer3Instance builds a new-FFI executor, which is what makes wasmer 6 install its
// own handlers.
func newWasmer3Instance(t *testing.T) *wasmer3.Wasmer3Instance {
	wasmerExecutor, err := wasmer3.NewExecutor(wasmer3.ExecutorConfig{
		OpcodeVersion:  executor.OpcodeVersionV2,
		WASMOpcodeCost: uniformOpcodeCost(1),
	})
	require.Nil(t, err)
	t.Cleanup(wasmerExecutor.Close)

	instance, err := wasmerExecutor.NewInstance(
		trapContract(t),
		wasmer3.CompilationOptions{MaxMemoryGrow: 10, MaxMemoryGrowDelta: 10},
		stubVMHooks{},
	)
	require.Nil(t, err)
	t.Cleanup(instance.Close)

	return instance
}

// --- The experiment ----------------------------------------------------------

// TestWasmer3TrapsWithWasmer2First is the order that matters most.
//
// wasmer2.CreateExecutor calls vm_force_sighandler_reinstall unconditionally, so it is the
// one that can overwrite a handler wasmer 6 has already installed. If it does so without
// chaining, the traps below take down the process instead of coming back as results.
func TestWasmer3TrapsWithWasmer2First(t *testing.T) {
	legacy := newWasmer2Instance(t)
	require.Nil(t, legacy.CallFunction("noTrap"), "legacy backend is alive")

	instance := newWasmer3Instance(t)

	result, err := instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmer3.CallOk, result.Kind, "new backend is alive")

	for _, funcName := range []string{"trapUnreachable", "trapMemory"} {
		t.Run(funcName, func(t *testing.T) {
			result, err := instance.Call(funcName, 10_000_000)
			require.Nil(t, err)
			require.Equal(t, wasmer3.CallRuntimeError, result.Kind,
				"the trap was caught and reported: %+v", result)
			t.Logf("%s reported: %s", funcName, result.Message)
		})
	}

	// Still usable after the traps, which is what distinguishes a handled trap from a
	// handler that merely did not crash.
	result, err = instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmer3.CallOk, result.Kind, "new backend survived its traps")

	require.Nil(t, legacy.CallFunction("noTrap"), "legacy backend survived them too")
}

// TestWasmer2TrapsWithWasmer3First is the mirror: wasmer 6 installs first, and the legacy
// backend has to keep working through its own traps afterwards.
func TestWasmer2TrapsWithWasmer3First(t *testing.T) {
	instance := newWasmer3Instance(t)

	result, err := instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmer3.CallOk, result.Kind, "new backend is alive")

	legacy := newWasmer2Instance(t)
	require.Nil(t, legacy.CallFunction("noTrap"), "legacy backend is alive")

	for _, funcName := range []string{"trapUnreachable", "trapMemory"} {
		t.Run(funcName, func(t *testing.T) {
			err := legacy.CallFunction(funcName)
			require.NotNil(t, err, "the trap was caught and reported")
			t.Logf("%s reported: %s", funcName, err)
		})
	}

	require.Nil(t, legacy.CallFunction("noTrap"), "legacy backend survived its traps")

	result, err = instance.Call("noTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmer3.CallOk, result.Kind, "new backend survived them too")
}

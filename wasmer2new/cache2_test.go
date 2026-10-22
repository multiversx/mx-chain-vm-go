package wasmer2new

// The compiled-code cache, from the Go side.
//
// The Rust tests already pin what the backend does with an artifact; what only a Go test can
// pin is the marshalling. Reading an artifact out is a two-call protocol — ask for the length,
// then hand over a buffer of exactly that size — and it is the one place in this package where
// a variable-length *binary* result crosses cgo. Everything else that crosses is a struct, a
// scalar, or a NUL-terminated string, so an off-by-one here would have nothing else to catch
// it.

import (
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
	"github.com/stretchr/testify/require"
)

func cacheTestOptions() wasmercommon.CompilationOptions {
	return wasmercommon.CompilationOptions{MaxMemoryGrow: 10, MaxMemoryGrowDelta: 10}
}

func newCachedInstance(t *testing.T, artifact []byte) (*Wasmer2Instance, error) {
	instance, err := newTestExecutor(t).NewInstanceFromCache(
		artifact,
		cacheTestOptions(),
		&testHooksProvider{log: &hookLog{}},
	)
	if instance != nil {
		t.Cleanup(instance.Close)
	}

	return instance, err
}

func TestCacheRoundTrip(t *testing.T) {
	code := []byte(hooksContract)
	instance, _ := newTestInstance(t, code)

	artifact, err := instance.Cache()
	require.Nil(t, err)
	require.NotEmpty(t, artifact)
	require.NotEqual(t, code, artifact, "the artifact must be the compiled module, not the source")

	restored, err := newCachedInstance(t, artifact)
	require.Nil(t, err)

	// The restored module is a working module, not merely a loaded one: same endpoints, and
	// the hook imports still resolve, which a module that lost its import section could not do.
	originalNames, err := instance.ExportedFunctionNames()
	require.Nil(t, err)
	restoredNames, err := restored.ExportedFunctionNames()
	require.Nil(t, err)
	require.Equal(t, originalNames, restoredNames)
	require.True(t, restored.CheckSignatures())
}

// The length call and the copy call have to agree exactly. A second read must also work:
// the artifact is held on the instance between the two halves, and nothing about reading it
// consumes it.
func TestCacheIsRepeatable(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))

	first, err := instance.Cache()
	require.Nil(t, err)

	second, err := instance.Cache()
	require.Nil(t, err)

	require.Equal(t, first, second)
}

// An artifact produced by one instance loads into another, which is the whole point: the VM
// stores it under the code hash and hands it back on a later, unrelated execution.
func TestCacheArtifactIsPortableBetweenInstances(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))
	artifact, err := instance.Cache()
	require.Nil(t, err)

	restored, err := newCachedInstance(t, artifact)
	require.Nil(t, err)

	again, err := restored.Cache()
	require.Nil(t, err)
	require.Equal(t, artifact, again, "re-serializing a restored module must reproduce it")
}

func TestCacheRejectsEmptyArtifact(t *testing.T) {
	_, err := newCachedInstance(t, []byte{})
	require.ErrorIs(t, err, wasmercommon.ErrEmptyCacheArtifact)
}

// Not an artifact, so it must be refused with a message rather than mapped as code. This is
// the guard that keeps a corrupt cache entry a recoverable miss instead of undefined
// behaviour; the VM's caller falls back to compiling the bytecode.
func TestCacheRejectsNonArtifactBytes(t *testing.T) {
	_, err := newCachedInstance(t, []byte("not a wasmer artifact"))
	require.NotNil(t, err)
}

// The wasm and the artifact reach the executor through different entry points, so feeding one
// to the other has to fail rather than quietly work.
func TestCacheRejectsPlainWasm(t *testing.T) {
	_, err := newCachedInstance(t, []byte(hooksContract))
	require.NotNil(t, err)
}

// A real contract, through the legacy interface the VM actually calls, so that the gas limit
// handling on the cached path is covered too.
func TestCacheRoundTripThroughLegacyInterface(t *testing.T) {
	code, err := os.ReadFile("../test/contracts/counter/output/counter.wasm")
	require.Nil(t, err, "fixture missing")

	wasmerExecutor := newTestExecutor(t)
	wasmerExecutor.provider = &testHooksProvider{log: &hookLog{}}

	options := executor.CompilationOptions{
		GasLimit:           100000,
		MaxMemoryGrow:      10,
		MaxMemoryGrowDelta: 10,
	}

	compiled, err := wasmerExecutor.NewInstanceWithOptions(code, options)
	require.Nil(t, err)
	t.Cleanup(func() { compiled.Clean() })

	artifact, err := compiled.Cache()
	require.Nil(t, err)
	require.NotEmpty(t, artifact)

	restored, err := wasmerExecutor.NewInstanceFromCompiledCodeWithOptions(artifact, options)
	require.Nil(t, err)
	t.Cleanup(func() { restored.Clean() })

	require.True(t, restored.IsFunctionImported("int64storageStore"))
	require.True(t, restored.HasFunction("increment"))
	require.Nil(t, restored.ValidateFunctionArities())
	require.Equal(t, compiled.GetFunctionNames(), restored.GetFunctionNames())
}

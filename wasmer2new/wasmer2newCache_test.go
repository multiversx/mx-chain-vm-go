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
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
	"github.com/stretchr/testify/require"
)

func newCachedInstance(t *testing.T, artifact []byte) (*Wasmer2Instance, error) {
	instance, err := newTestExecutor(t).NewInstanceFromCache(artifact, testCompilationOptions())
	if instance != nil {
		t.Cleanup(instance.(*Wasmer2Instance).Close)
		return instance.(*Wasmer2Instance), err
	}

	return nil, err
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
	require.Equal(t, instance.GetExportedFunctionNames(), restored.GetExportedFunctionNames())
	require.Equal(t, instance.GetImportedFunctionNames(), restored.GetImportedFunctionNames())
	require.True(t, restored.CheckSignatures())

	// And it runs: a restored instance reaches the hooks exactly as a compiled one does.
	require.Equal(t, executor.CallOk, restored.Call("checkI64Result", 10_000_000).Kind)
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

package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// stubVMHooks implements only the hooks a test calls. The embedded interface supplies the
// rest of the method set and is nil, so an unexpected call panics rather than passing.
type stubVMHooks struct {
	VMHooks

	gasLeft int64
	err     error

	scAddressCalls int
}

func (stub *stubVMHooks) GetGasLeft() (int64, error) {
	return stub.gasLeft, stub.err
}

func (stub *stubVMHooks) GetSCAddress(_ MemPtr) error {
	stub.scAddressCalls++
	return stub.err
}

type recordingErrorHandler struct {
	handled []error
}

func (handler *recordingErrorHandler) HandleHookError(err error) {
	handler.handled = append(handler.handled, err)
}

func TestVMHooksLegacyAdapter_PassesThroughOnSuccess(t *testing.T) {
	t.Parallel()

	handler := &recordingErrorHandler{}
	adapter := NewVMHooksLegacyAdapter(&stubVMHooks{gasLeft: 42}, handler)

	require.Equal(t, int64(42), adapter.GetGasLeft())
	require.Empty(t, handler.handled, "a hook that succeeded must not reach the error handler")
}

// A failing hook still returns the sentinel it produced alongside the error: contracts
// branch on those values, and before the migration the hook returned them directly.
func TestVMHooksLegacyAdapter_KeepsTheResultOfAFailingHook(t *testing.T) {
	t.Parallel()

	hookErr := errors.New("hook failed")
	handler := &recordingErrorHandler{}
	adapter := NewVMHooksLegacyAdapter(&stubVMHooks{gasLeft: -1, err: hookErr}, handler)

	require.Equal(t, int64(-1), adapter.GetGasLeft())
	require.Equal(t, []error{hookErr}, handler.handled)
}

func TestVMHooksLegacyAdapter_HandlesTheErrorOfAHookWithoutAResult(t *testing.T) {
	t.Parallel()

	hookErr := errors.New("hook failed")
	handler := &recordingErrorHandler{}
	hooks := &stubVMHooks{err: hookErr}
	adapter := NewVMHooksLegacyAdapter(hooks, handler)

	adapter.GetSCAddress(0)

	require.Equal(t, 1, hooks.scAddressCalls)
	require.Equal(t, []error{hookErr}, handler.handled)
}

func TestVMHooksLegacyAdapter_ImplementsVMHooksLegacy(t *testing.T) {
	t.Parallel()

	var _ VMHooksLegacy = (*VMHooksLegacyAdapter)(nil)
}

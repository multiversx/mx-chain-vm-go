package executorwrapper

import (
	"errors"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// ErrVMHooksCannotHandleErrors signals that the wrapped VM hooks cannot absorb an error
// returned by a hook, and so cannot be given to a backend that needs them to.
var ErrVMHooksCannotHandleErrors = errors.New("wrapped vm hooks do not handle hook errors")

// vmHooksWithErrorHandler is what the generated wrapper has to wrap: the hooks, plus the
// error handler a legacy backend will look for on them.
//
// The wrapper sits between the VM and the backend, so whatever the backend expects to find on
// the hooks has to survive the wrapping. Stated here, as a requirement of this package, and
// not on executor.VMHooks, which stays free of it.
type vmHooksWithErrorHandler interface {
	executor.VMHooks
	executor.VMHooksErrorHandler
}

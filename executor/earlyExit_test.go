package executor_test

import (
	"errors"
	"fmt"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

func TestEarlyExitCodeOf_PlainErrorGetsTheDefault(t *testing.T) {
	// Most hook failures are plain errors, and they all mean the same thing.
	require.Equal(t,
		uint64(vmcommon.ExecutionFailed),
		executor.EarlyExitCodeOf(errors.New("something went wrong")))
}

func TestEarlyExitCodeOf_UsesTheExitsOwnCode(t *testing.T) {
	userError := &executor.EarlyExit{Code: vmcommon.UserError, Message: "nope"}
	require.Equal(t, uint64(vmcommon.UserError), executor.EarlyExitCodeOf(userError))

	outOfGas := &executor.EarlyExit{Code: vmcommon.OutOfGas}
	require.Equal(t, uint64(vmcommon.OutOfGas), executor.EarlyExitCodeOf(outOfGas))
}

func TestEarlyExitCodeOf_SurvivesWrapping(t *testing.T) {
	// Hook helpers add context with %w; the code must still be found.
	wrapped := fmt.Errorf("while storing: %w", &executor.EarlyExit{Code: vmcommon.UserError})
	require.Equal(t, uint64(vmcommon.UserError), executor.EarlyExitCodeOf(wrapped))
}

func TestEarlyExitCodeOf_AgreesWithTheDefaultExit(t *testing.T) {
	// The code an executor is told and the exit the VM applies have to describe the same
	// termination, so both come from one place.
	cause := errors.New("boom")
	exit := executor.ExitExecutionFailed(cause)

	require.Equal(t, uint64(exit.Code), executor.EarlyExitCodeOf(cause))
	require.True(t, exit.ConsumeAllGas, "the default failure burns the remaining gas")
	require.Equal(t, "boom", exit.Message)
	require.ErrorIs(t, exit, cause, "the cause stays reachable through the exit")
}

func TestEarlyExitCodeOf_ZeroStaysAvailableAsTheAsyncCallSentinel(t *testing.T) {
	// Code 0 is reserved on the wire for the legacy async call, so nothing in the default
	// path may produce it. vmcommon.Ok is 0, and is never a failure.
	require.NotEqual(t, uint64(0), executor.EarlyExitCodeOf(errors.New("boom")))
	require.Equal(t, uint64(vmcommon.Ok), uint64(0))
}

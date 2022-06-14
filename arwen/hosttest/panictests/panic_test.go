package panictests

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/stretchr/testify/require"
)

const increment = "increment"

func TestExecution_PanicInGoWithSilentWasmer_SIGSEGV(t *testing.T) {
	code := test.GetTestSCCode("counter", "../../../")
	host, blockchain := test.DefaultTestArwenForCallSigSegv(t, code, big.NewInt(1))

	defer func() {
		host.Reset()
	}()

	blockchain.GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, error) {
		var i *int
		i = nil

		// dereference a nil pointer
		*i = *i + 1
		return nil, nil
	}

	input := test.CreateTestContractCallInputBuilder().
		WithGasProvided(1000000).
		WithFunction(increment).
		Build()

	// Ensure that host.RunSmartContractCall() still panics, but the panic is a
	// wrapped error.
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		require.True(t, errors.Is(err, arwen.ErrExecutionPanicked))
	}()

	_, _ = host.RunSmartContractCall(input)
}

func TestExecution_PanicInGoWithSilentWasmer_SIGFPE(t *testing.T) {
	code := test.GetTestSCCode("counter", "../../../")
	host, blockchain := test.DefaultTestArwenForCallSigSegv(t, code, big.NewInt(1))
	defer func() {
		host.Reset()
	}()

	blockchain.GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, error) {
		i := 5
		j := 4
		i = i / (j - 4)
		return nil, nil
	}

	input := test.CreateTestContractCallInputBuilder().
		WithGasProvided(1000000).
		WithFunction(increment).
		Build()

	// Ensure that host.RunSmartContractCall() still panics, but the panic is a
	// wrapped error.
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		require.True(t, errors.Is(err, arwen.ErrExecutionPanicked))
	}()

	_, _ = host.RunSmartContractCall(input)
}

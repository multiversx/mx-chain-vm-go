package panictests

import (
	"math/big"
	"sync"
	"testing"
	"time"

	mock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/stretchr/testify/require"
)

const increment = "increment"

func TestExecution_PanicInGoWithSilentWasmer_SIGSEGV(t *testing.T) {
	code := test.GetTestSCCode("counter", "../../../")
	host, blockchain := test.DefaultTestVMForCallSigSegv(t, code, big.NewInt(1))

	defer func() {
		host.Reset()
	}()

	blockchain.GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, uint32, error) {
		var i *int

		// dereference a nil pointer
		*i = *i + 1
		return nil, 0, nil
	}

	input := test.CreateTestContractCallInputBuilder().
		WithGasProvided(1000000).
		WithFunction(increment).
		Build()

	// Ensure that no more panic
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	_, err := host.RunSmartContractCall(input)
	require.Equal(t, vmhost.ErrExecutionPanicked, err)
}

func TestExecution_PanicInGoWithSilentWasmer_SIGFPE(t *testing.T) {
	code := test.GetTestSCCode("counter", "../../../")
	host, blockchain := test.DefaultTestVMForCallSigSegv(t, code, big.NewInt(1))
	defer func() {
		host.Reset()
	}()

	blockchain.GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, uint32, error) {
		i := 5
		j := 4
		// trick the linter with this for loop. It will eventually cause a division by 0 exception
		for counter := 0; counter <= j; counter++ {
			i = i / (counter - 4)
		}
		return nil, 0, nil
	}

	input := test.CreateTestContractCallInputBuilder().
		WithGasProvided(1000000).
		WithFunction(increment).
		Build()

	// Ensure that host.RunSmartContractCall() still panics, but the panic is a
	// wrapped error.
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	_, err := host.RunSmartContractCall(input)
	require.Equal(t, vmhost.ErrExecutionPanicked, err)
}

func TestExecution_PanicInGoWithSilentWasmer_Timeout(t *testing.T) {
	code := test.GetTestSCCode("counter", "../../../")
	host, blockchain := test.DefaultTestVMForCallSigSegv(t, code, big.NewInt(1))
	defer func() {
		host.Reset()
	}()

	blockchain.GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, uint32, error) {
		time.Sleep(2 * time.Second)
		return nil, 0, nil
	}

	input := test.CreateTestContractCallInputBuilder().
		WithGasProvided(1000000).
		WithFunction(increment).
		Build()

	// Ensure that panics are not thrown
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	_, err := host.RunSmartContractCall(input)
	require.Equal(t, vmhost.ErrExecutionFailedWithTimeout, err)
}

func TestExecution_PanicInGoWithSilentWasmer_TimeoutAndSIGSEGV(t *testing.T) {
	code := test.GetTestSCCode("counter", "../../../")
	host, blockchain := test.DefaultTestVMForCallSigSegv(t, code, big.NewInt(1))

	defer func() {
		host.Reset()
	}()

	blockchain.GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, uint32, error) {
		var i *int

		// dereference a nil pointer
		time.Sleep(time.Second)
		*i = *i + 1
		return nil, 0, nil
	}

	input := test.CreateTestContractCallInputBuilder().
		WithGasProvided(1000000).
		WithFunction(increment).
		Build()

	// Ensure that panics are not thrown
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	_, err := host.RunSmartContractCall(input)
	require.NotNil(t, err)
}

func TestExecution_MultipleHostsPanicInGoWithSilentWasmer_TimeoutAndSIGSEGV(t *testing.T) {
	numParallel := 100
	hosts := make([]vmhost.VMHost, numParallel)
	blockchains := make([]*mock.BlockchainHookStub, numParallel)
	for k := 0; k < numParallel; k++ {
		code := test.GetTestSCCode("counter", "../../../")
		hosts[k], blockchains[k] = test.DefaultTestVMForCallSigSegv(t, code, big.NewInt(1))
		blockchains[k].GetStorageDataCalled = func(_ []byte, _ []byte) ([]byte, uint32, error) {
			var i *int

			// dereference a nil pointer
			time.Sleep(time.Second)
			*i = *i + 1
			return nil, 0, nil
		}
	}

	defer func() {
		for _, vm := range hosts {
			vm.Reset()
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(numParallel)
	for k := 0; k < numParallel; k++ {
		go func(idx int) {
			input := test.CreateTestContractCallInputBuilder().
				WithGasProvided(1000000).
				WithFunction(increment).
				Build()
			// Ensure that no more panic
			defer func() {
				r := recover()
				require.Nil(t, r)
			}()

			_, err := hosts[idx].RunSmartContractCall(input)
			wg.Done()
			require.NotNil(t, err)
		}(k)
	}

	wg.Wait()
}

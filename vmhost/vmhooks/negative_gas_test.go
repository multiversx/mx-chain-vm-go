package vmhooks

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

type asyncContextStub struct {
	vmhost.AsyncContext
	called bool
}

func (a *asyncContextStub) RegisterLegacyAsyncCall(_ []byte, _ []byte, _ []byte) error {
	a.called = true
	return nil
}

func TestDeleteContractRejectsNegativeGasLimit(t *testing.T) {
	host := &mock.VMHostMock{}
	host.MeteringContext = &mock.MeteringContextMock{
		GasCost: &config.GasCost{
			BaseOpsAPICost: config.BaseOpsAPICost{
				AsyncCallStep:        5,
				AsyncCallbackGasLock: 5,
			},
		},
	}

	host.RuntimeContext = &mock.RuntimeContextMock{}

	async := &asyncContextStub{}
	host.AsyncContext = async

	exit := executor.EarlyExitOf(deleteContract(host, []byte("dest"), nil, -1))
	require.False(t, async.called)
	require.Equal(t, vmcommon.OutOfGas, exit.Code)

	async.called = false
	exit = executor.EarlyExitOf(deleteContract(host, []byte("dest"), nil, 0))
	require.False(t, async.called)
	require.Equal(t, vmcommon.OutOfGas, exit.Code)
}

func TestUpgradeContractRejectsNegativeGasLimit(t *testing.T) {
	host := &mock.VMHostMock{}
	host.MeteringContext = &mock.MeteringContextMock{
		GasCost: &config.GasCost{
			BaseOpsAPICost: config.BaseOpsAPICost{
				AsyncCallStep:        5,
				AsyncCallbackGasLock: 5,
			},
		},
	}

	host.RuntimeContext = &mock.RuntimeContextMock{}

	async := &asyncContextStub{}
	host.AsyncContext = async

	exit := executor.EarlyExitOf(upgradeContract(host, []byte("dest"), []byte("code"), []byte("meta"), nil, nil, -1))
	require.False(t, async.called)
	require.Equal(t, vmcommon.OutOfGas, exit.Code)

	async.called = false
	exit = executor.EarlyExitOf(upgradeContract(host, []byte("dest"), []byte("code"), []byte("meta"), nil, nil, 0))
	require.False(t, async.called)
	require.Equal(t, vmcommon.OutOfGas, exit.Code)
}

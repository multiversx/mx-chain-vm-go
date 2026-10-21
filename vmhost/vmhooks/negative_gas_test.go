package vmhooks

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/config"
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

type runtimeWithBreakpointStub struct {
	mock.RuntimeContextMock
	breakpoint vmhost.BreakpointValue
}

func (a *runtimeWithBreakpointStub) SetRuntimeBreakpointValue(value vmhost.BreakpointValue) {
	a.breakpoint = value
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

	runtime := &runtimeWithBreakpointStub{}
	host.RuntimeContext = runtime

	async := &asyncContextStub{}
	host.AsyncContext = async

	deleteContract(host, []byte("dest"), nil, -1)
	require.False(t, async.called)
	require.Equal(t, vmhost.BreakpointOutOfGas, runtime.breakpoint)

	async.called = false
	runtime.breakpoint = 0
	deleteContract(host, []byte("dest"), nil, 0)
	require.False(t, async.called)
	require.Equal(t, vmhost.BreakpointOutOfGas, runtime.breakpoint)
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

	runtime := &runtimeWithBreakpointStub{}
	host.RuntimeContext = runtime

	async := &asyncContextStub{}
	host.AsyncContext = async

	upgradeContract(host, []byte("dest"), []byte("code"), []byte("meta"), nil, nil, -1)
	require.False(t, async.called)
	require.Equal(t, vmhost.BreakpointOutOfGas, runtime.breakpoint)

	async.called = false
	runtime.breakpoint = 0
	upgradeContract(host, []byte("dest"), []byte("code"), []byte("meta"), nil, nil, 0)
	require.False(t, async.called)
	require.Equal(t, vmhost.BreakpointOutOfGas, runtime.breakpoint)
}

package host

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

type VMOutputVerifier struct {
	vmOutput *vmcommon.VMOutput

	T testing.TB
}

func NewVMOutputVerifier(t testing.TB, vmOutput *vmcommon.VMOutput, err error) *VMOutputVerifier {
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	return &VMOutputVerifier{
		vmOutput: vmOutput,
		T:        t,
	}
}

func (a *VMOutputVerifier) GasUsed(address []byte, gas uint64) {
	account := a.vmOutput.OutputAccounts[string(address)]
	require.Equal(a.T, gas, account.GasUsed)
}

func (a *VMOutputVerifier) GasRemaining(gas uint64) {
	require.Equal(a.T, gas, a.vmOutput.GasRemaining)
}

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

func (v *VMOutputVerifier) GasUsed(address []byte, gas uint64) {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account)
	require.Equal(v.T, int(gas), int(account.GasUsed))
}

func (v *VMOutputVerifier) GasRemaining(gas uint64) {
	require.Equal(v.T, int(gas), int(v.vmOutput.GasRemaining))
}

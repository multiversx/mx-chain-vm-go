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

func (v *VMOutputVerifier) Ok() *VMOutputVerifier {
	return v.RetCode(vmcommon.Ok)
}

func (v *VMOutputVerifier) RetCode(code vmcommon.ReturnCode) *VMOutputVerifier {
	require.Equal(v.T, code, v.vmOutput.ReturnCode)
	return v
}

func (v *VMOutputVerifier) NoMsg() *VMOutputVerifier {
	require.Equal(v.T, "", v.vmOutput.ReturnMessage)
	return v
}
func (v *VMOutputVerifier) Msg(message string) *VMOutputVerifier {
	require.Equal(v.T, message, v.vmOutput.ReturnMessage)
	return v
}

func (v *VMOutputVerifier) GasUsed(address []byte, gas uint64) *VMOutputVerifier {
	account := v.vmOutput.OutputAccounts[string(address)]
	require.NotNil(v.T, account)
	require.Equal(v.T, int(gas), int(account.GasUsed))
	return v
}

func (v *VMOutputVerifier) GasRemaining(gas uint64) *VMOutputVerifier {
	require.Equal(v.T, int(gas), int(v.vmOutput.GasRemaining))
	return v
}

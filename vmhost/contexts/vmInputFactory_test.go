package contexts

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/stretchr/testify/require"
)

func TestNewVMInputFactory(t *testing.T) {
	t.Parallel()

	f := NewVMInputFactory()
	require.False(t, f.IsInterfaceNil())
	require.Implements(t, new(VMInputFactory), f)
}

func TestVMInputFactory_CreateContractCallInput(t *testing.T) {
	t.Parallel()

	function := "function"
	address := []byte("address")
	vmInput := &vmcommon.VMInput{
		CallerAddr: address,
	}
	contractCallInput := &vmcommon.ContractCallInput{
		Function:      function,
		RecipientAddr: address,
	}
	f := NewVMInputFactory()

	input := f.CreateContractCallInput(*vmInput, contractCallInput)
	require.Equal(t, input.GetVMInput(), vmInput)
	require.Equal(t, input.GetFunction(), function)
	require.Equal(t, input.GetRecipientAddr(), address)
}

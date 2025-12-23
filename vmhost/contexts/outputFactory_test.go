package contexts

import (
	"testing"

	"github.com/stretchr/testify/require"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
)

func TestNewOutputContextFactory(t *testing.T) {
	t.Parallel()

	ocf := NewOutputContextFactory()
	require.False(t, ocf.IsInterfaceNil())
	require.Implements(t, new(OutputContextCreator), ocf)
}

func TestNewOutputContextFactory_CreateOutputContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}

	ocf := NewOutputContextFactory()
	outputContext, err := ocf.CreateOutputContext(host)
	require.NoError(t, err)
	require.NotNil(t, outputContext)
}

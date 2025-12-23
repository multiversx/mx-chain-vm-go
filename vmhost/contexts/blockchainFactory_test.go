package contexts

import (
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/stretchr/testify/require"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
)

func TestNewBlockchainContextFactory(t *testing.T) {
	t.Parallel()

	bcf := NewBlockchainContextFactory()
	require.False(t, bcf.IsInterfaceNil())
	require.Implements(t, new(BlockchainContextCreator), bcf)
}

func TestNewBlockchainContextFactory_CreateBlockchainContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}
	mockWorld := worldmock.NewMockWorld()

	bcf := NewBlockchainContextFactory()
	blockchainContext, err := bcf.CreateBlockchainContext(host, mockWorld)
	require.NoError(t, err)
	require.NotNil(t, blockchainContext)
}

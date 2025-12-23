package contexts

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-vm-go/config"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
)

func TestNewMeteringContextFactory(t *testing.T) {
	t.Parallel()

	mcf := NewMeteringContextFactory()
	require.False(t, mcf.IsInterfaceNil())
	require.Implements(t, new(MeteringContextCreator), mcf)
}

func TestNewMeteringContextFactory_CreateMeteringContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{}
	gasMap := config.MakeGasMapForTests()
	const blockGasLimit = uint64(15000)

	mcf := NewMeteringContextFactory()
	meteringContext, err := mcf.CreateMeteringContext(host, gasMap, blockGasLimit, config.NewGasScheduleFactory())
	require.NoError(t, err)
	require.NotNil(t, meteringContext)
}

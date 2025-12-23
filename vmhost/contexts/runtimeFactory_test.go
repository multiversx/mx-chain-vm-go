package contexts

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/testcommon/testexecutor"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

func TestNewRuntimeContextFactory(t *testing.T) {
	t.Parallel()

	rcf := NewRuntimeContextFactory()
	require.False(t, rcf.IsInterfaceNil())
	require.Implements(t, new(RuntimeContextCreator), rcf)
}

func TestNewRuntimeContextFactory_CreateRuntimeContext(t *testing.T) {
	t.Parallel()

	host := InitializeVMAndWasmer()
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	hasher := defaultHasher
	execFactory := testexecutor.NewDefaultTestExecutorFactory(t)
	exec, err := execFactory.CreateExecutor(executor.ExecutorFactoryArgs{
		VMHooks: vmhooks.NewVMHooksImpl(host),
	})
	require.Nil(t, err)

	rcf := NewRuntimeContextFactory()
	runtimeContext, err := rcf.CreateRuntimeContext(host, vmType, bfc, exec, hasher)
	require.NoError(t, err)
	require.NotNil(t, runtimeContext)
}

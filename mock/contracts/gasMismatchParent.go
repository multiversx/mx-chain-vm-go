package contracts

import (
	"math/big"

	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
)

// GasMismatchAsyncCallParentMock is an exposed mock contract method
func GasMismatchAsyncCallParentMock(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("gasMismatchParent", func() *mock.InstanceMock {
		host := instanceMock.Host
		managedTypes := host.ManagedTypes()
		instance := mock.GetMockInstance(host)

		destHandle := managedTypes.NewManagedBufferFromBytes(test.ChildAddress)
		valueHandle := managedTypes.NewBigIntFromInt64(0)
		functionHandle := managedTypes.NewManagedBufferFromBytes([]byte("gasMismatchChild"))
		argumentsHandle := managedTypes.NewManagedBuffer()
		_ = managedTypes.WriteManagedVecOfManagedBuffers([][]byte{}, argumentsHandle)

		err := vmhooks.ManagedAsyncCallWithHost(
			host,
			destHandle,
			valueHandle,
			functionHandle,
			argumentsHandle,
		)
		if err != nil {
			// ManagedAsyncCallWithHost may return an *executor.EarlyExit (out-of-gas
			// breakpoint) rather than a plain error; route through HandleHookError
			// (as the generated wrapper would) instead of vmhooks.FailExecution,
			// which only knows how to apply the plain-error failure path.
			vmhooks.NewVMHooksImpl(host).HandleHookError(err)
		}

		return instance

	})
}

// GasMismatchCallBackParentMock is an exposed mock contract method
func GasMismatchCallBackParentMock(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("callBack", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		output := host.Output()
		arguments := host.Runtime().Arguments()

		output.Finish(big.NewInt(0xCA11BAC3).Bytes())

		for _, arg := range arguments {
			output.Finish(arg)
		}

		output.Finish(big.NewInt(0xCA11BAC3).Bytes())
		return instance
	})
}

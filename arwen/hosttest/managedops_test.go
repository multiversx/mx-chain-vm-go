package hosttest

import (
	"testing"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/arwen/elrondapi"
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
)

// wasm memory ~~~> managed buffer
func TestManaged_SetByteSlice(t *testing.T) {
	prefix := "ABCD"
	slice := "EFGHIJKLMN"
	sufix := "OPR"
	data := prefix + slice + sufix
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(parentInstance *mock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testFunction", func() *mock.InstanceMock {
						host := parentInstance.Host
						managedType := host.ManagedTypes()
						mBuffer := managedType.NewManagedBufferFromBytes(
							make([]byte, len(data)))
						result := elrondapi.ManagedBufferSetByteSliceWithTypedArgs(
							host, mBuffer, int32(len(prefix)), int32(len(slice)), []byte(data))
						if result != 0 {
							arwen.WithFaultAndHost(host, arwen.ErrSignalError, true)
						}
						bufferBytes, err := managedType.GetBytes(mBuffer)
						if err != nil {
							arwen.WithFaultAndHost(host, err, true)
						}
						host.Output().Finish(bufferBytes)
						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(1000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				// |....ABCDEFGHIJ...|
				ReturnData(append(make([]byte, len(prefix)),
					append([]byte(data[0:len(slice)]), make([]byte, len(sufix))...)...))
		})
}

// managed buffer ~~~> managed buffer
func TestManaged_CopyByteSlice_DifferentBuffer(t *testing.T) {
	prefix := "ABCD"
	slice := "EFGHIJKLMN"
	sufix := "OPR"
	sourceData := prefix + slice + sufix
	destinationData := "01234567890123456789"
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(parentInstance *mock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testFunction", func() *mock.InstanceMock {
						host := parentInstance.Host
						managedType := host.ManagedTypes()
						sourceMBuffer := managedType.NewManagedBufferFromBytes(
							[]byte(sourceData))
						destMBuffer := managedType.NewManagedBufferFromBytes(
							[]byte(destinationData))
						result := elrondapi.ManagedBufferCopyByteSliceWithHost(
							host, sourceMBuffer, int32(len(prefix)), int32(len(slice)), destMBuffer)
						if result != 0 {
							arwen.WithFaultAndHost(host, arwen.ErrSignalError, true)
						}
						destBytes, err := managedType.GetBytes(destMBuffer)
						if err != nil {
							arwen.WithFaultAndHost(host, err, true)
						}
						host.Output().Finish(destBytes)
						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(1000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData([]byte(slice))
		})
}

func TestManaged_CopyByteSlice_SameBuffer(t *testing.T) {
	prefix := "ABCD"
	slice := "EFGHIJKLMN"
	sufix := "OPR"
	sourceData := prefix + slice + sufix
	deltaForSlice := int32(2)
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(parentInstance *mock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("testFunction", func() *mock.InstanceMock {
						host := parentInstance.Host
						managedType := host.ManagedTypes()
						sourceMBuffer := managedType.NewManagedBufferFromBytes(
							[]byte(sourceData))
						result := elrondapi.ManagedBufferCopyByteSliceWithHost(
							host, sourceMBuffer, int32(len(prefix))-deltaForSlice, int32(len(slice)), sourceMBuffer)
						if result != 0 {
							arwen.WithFaultAndHost(host, arwen.ErrSignalError, true)
						}
						destBytes, err := managedType.GetBytes(sourceMBuffer)
						if err != nil {
							arwen.WithFaultAndHost(host, err, true)
						}
						host.Output().Finish(destBytes)
						return parentInstance
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(1000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			prefixLen := int32(len(prefix))
			sliceLen := int32(len(slice))
			verify.Ok().
				// |CDEFGHIJKL|
				ReturnData(
					append([]byte(prefix)[prefixLen-deltaForSlice:prefixLen],
						[]byte(slice)[:sliceLen-deltaForSlice]...))
		})
}

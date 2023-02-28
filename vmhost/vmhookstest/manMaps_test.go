package vmhookstest

import (
	"testing"

	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	worldmock "github.com/multiversx/mx-chain-vm-go/mock/world"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/assert"
)

func TestManagedMap(t *testing.T) {
	key := []byte("key")
	value := []byte("value")

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(instance *mock.InstanceMock, config interface{}) {
					instance.AddMockMethod("testFunction", func() *mock.InstanceMock {
						host := instance.Host
						managedType := host.ManagedTypes()
						mMap := managedType.NewManagedMap()

						keyBuff := managedType.NewManagedBufferFromBytes(key)
						valueBuff := managedType.NewManagedBufferFromBytes(value)
						err := managedType.ManagedMapPut(mMap, keyBuff, valueBuff)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}

						outValueBuf := managedType.NewManagedBufferFromBytes(
							make([]byte, len(value)))
						err = managedType.ManagedMapGet(mMap, keyBuff, outValueBuf)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}
						outValueBytes, err := managedType.GetBytes(outValueBuf)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}
						host.Output().Finish(outValueBytes)

						contains, err := managedType.ManagedMapContains(mMap, keyBuff)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}
						if contains {
							host.Output().Finish([]byte{1})
						} else {
							host.Output().Finish([]byte{0})
						}

						outValueBuf = managedType.NewManagedBufferFromBytes(
							make([]byte, len(value)))
						err = managedType.ManagedMapRemove(mMap, keyBuff, outValueBuf)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}
						host.Output().Finish(outValueBytes)

						containsAfterRemove, err := managedType.ManagedMapContains(mMap, keyBuff)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}
						if containsAfterRemove {
							host.Output().Finish([]byte{1})
						} else {
							host.Output().Finish([]byte{0})
						}

						outValueBuf = managedType.NewManagedBufferFromBytes(
							make([]byte, len(value)))
						err = managedType.ManagedMapGet(mMap, keyBuff, outValueBuf)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}
						outValueBytes, err = managedType.GetBytes(outValueBuf)
						if err != nil {
							vmhooks.WithFaultAndHost(host, err, true)
							return instance
						}
						host.Output().Finish(outValueBytes)

						return instance
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
				ReturnData(value, []byte{1}, value, []byte{0}, []byte{})
		})
	assert.Nil(t, err)
}

package contracts

import (
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen/elrondapi"
	mock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
)

// LoadStore is an exposed mock contract method
func LoadStore(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("loadStore", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 1 {
			host.Runtime().SignalUserError("needs 1 argument")
			return instance
		}

		key := arguments[0]
		elrondapi.StorageLoadWithWithTypedArgs(host, key)
		value := elrondapi.StorageLoadWithWithTypedArgs(host, key)

		host.Output().Finish(value)
		return instance
	})
}

// LoadStoreFromAddress is an exposed mock contract method
func LoadStoreFromAddress(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("loadStoreFromAddress", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 2 {
			host.Runtime().SignalUserError("need 2 arguments")
			return instance
		}

		address := arguments[0]
		key := arguments[1]

		elrondapi.StorageLoadFromAddressWithTypedArgs(host, address, key)
		value := elrondapi.StorageLoadFromAddressWithTypedArgs(host, address, key)

		host.Output().Finish(value)
		return instance
	})
}

// SetStore is an exposed mock contract method
func SetStore(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("setStore", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		arguments := host.Runtime().Arguments()
		if len(arguments) != 2 {
			host.Runtime().SignalUserError("needs 2 arguments")
			return instance
		}

		key := arguments[0]
		value := arguments[1]

		elrondapi.StorageStoreWithTypedArgs(host, key, value)
		elrondapi.StorageStoreWithTypedArgs(host, key, value)

		return instance
	})
}

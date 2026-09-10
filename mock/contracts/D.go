package contracts

import (
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type D struct {
	Host vmhost.VMHost
}

func (d *D) DoSomething() *contextmock.InstanceMock {
	instance := contextmock.GetMockInstance(d.Host)
	host := instance.Host
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("event:D_doSomething")})
	return instance
}

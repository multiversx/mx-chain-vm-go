package contracts

import (
	"math/big"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type C struct {
	Host vmhost.VMHost
}

func (c *C) DoSomething() *contextmock.InstanceMock {
	instance := contextmock.GetMockInstance(c.Host)
	host := instance.Host
	dAddress := host.Runtime().Arguments()[0]
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("event:C_doSomething")})
	err := host.Async().RegisterAsyncCall("", &vmhost.AsyncCall{
		Destination:     dAddress,
		Data:            []byte("doSomething"),
		ValueBytes:      big.NewInt(0).Bytes(),
		GasLimit:        50000000,
		SuccessCallback: "cCallback",
	})
	if err != nil {
		host.Runtime().FailExecution(err)
	}
	return instance
}

func (c *C) CCallback() *contextmock.InstanceMock {
	instance := contextmock.GetMockInstance(c.Host)
	host := instance.Host
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("event:C_callback")})
	return instance
}

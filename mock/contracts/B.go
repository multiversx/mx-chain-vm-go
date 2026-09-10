package contracts

import (
	"math/big"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type B struct {
	Host vmhost.VMHost
}

func (b *B) CallC() *contextmock.InstanceMock {
	instance := contextmock.GetMockInstance(b.Host)
	host := instance.Host
	cAddress := host.Runtime().Arguments()[0]
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("event:B_callC")})
	err := host.Async().RegisterAsyncCall("", &vmhost.AsyncCall{
		Destination:     cAddress,
		Data:            []byte("doSomething"),
		ValueBytes:      big.NewInt(0).Bytes(),
		GasLimit:        50000000,
		SuccessCallback: "bCallback",
	})
	if err != nil {
		host.Runtime().FailExecution(err)
	}
	return instance
}

func (b *B) BCallback() *contextmock.InstanceMock {
	instance := contextmock.GetMockInstance(b.Host)
	host := instance.Host
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("event:B_callback")})
	return instance
}

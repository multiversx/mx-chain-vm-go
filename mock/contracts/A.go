package contracts

import (
	"math/big"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

type A struct {
	Host vmhost.VMHost
}

func (a *A) CallB() *contextmock.InstanceMock {
	instance := contextmock.GetMockInstance(a.Host)
	host := instance.Host
	bAddress := host.Runtime().Arguments()[0]
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("A.callB called by"), host.Runtime().GetCaller()})
	err := host.Async().RegisterAsyncCall("", &vmhost.AsyncCall{
		Destination:     bAddress,
		Data:            []byte("callC@sc:C"),
		ValueBytes:      big.NewInt(0).Bytes(),
		GasLimit:        50000000,
		SuccessCallback: "aCallback",
	})
	if err != nil {
		host.Runtime().FailExecution(err)
	}
	return instance
}

func (a *A) ACallback() *contextmock.InstanceMock {
	instance := contextmock.GetMockInstance(a.Host)
	host := instance.Host
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("A.aCallback called by"), host.Runtime().GetCaller()})
	host.Output().WriteLog(host.Runtime().GetContextAddress(), nil, [][]byte{[]byte("A.aCallback called with data"), host.Runtime().Arguments()[0]})
	return instance
}

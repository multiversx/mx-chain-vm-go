package subcontexts

import (
	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Runtime struct {
}

func (r *Runtime) GetSCAddress() []byte {
	panic("not implemented")
}

func (r *Runtime) Function() string {
	panic("not implemented")
}

func (r *Runtime) Arguments() [][]byte {
	panic("not implemented")
}

func (r *Runtime) SignalUserError() {
	panic("not implemented")
}

func (r *Runtime) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	panic("not implemented")
}

func (r *Runtime) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	panic("not implemented")
}

func (r *Runtime) CallData() []byte {
	panic("not implemented")
}

func (r *Runtime) SetReadOnly(readOnly bool) {
	panic("not implemented")
}

func (r *Runtime) ExecuteOnSameContext(input *vmcommon.ContractCallInput) error {
	panic("not implemented")
}

func (r *Runtime) ExecuteOnDestContext(input *vmcommon.ContractCallInput) error {
	panic("not implemented")
}

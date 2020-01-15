package contexts

import (
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type outputState struct {
	ReturnData      [][]byte
	ReturnCode      vmcommon.ReturnCode
	ReturnMessage   string
	GasRemaining    uint64
	GasRefund       *big.Int
	OutputAccounts  map[string]*outputAccount
	DeletedAccounts [][]byte
	TouchedAccounts [][]byte
	Logs            map[string]*vmcommon.LogEntry
}

func newOutputState() *outputState {
	return &outputState{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       nil,
		OutputAccounts:  make(map[string]*outputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make(map[string]*vmcommon.LogEntry),
	}
}

func (state *outputState) update(otherState *outputState) {
	for key, log := range otherState.Logs {
		state.Logs[key] = log
	}

	for address, otherAccount := range otherState.OutputAccounts {
		if _, ok := state.OutputAccounts[address]; !ok {
			state.OutputAccounts[address] = &outputAccount{}
		}
		state.OutputAccounts[address].update(otherAccount)
	}

	// TODO merge DeletedAccounts and TouchedAccounts as well?

	state.ReturnData = append(state.ReturnData, otherState.ReturnData...)
	state.GasRemaining = otherState.GasRemaining
	state.GasRefund = otherState.GasRefund
	state.ReturnCode = otherState.ReturnCode
	state.ReturnMessage = otherState.ReturnMessage
}

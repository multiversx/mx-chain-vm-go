package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
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

type outputAccount struct {
	Address        []byte
	Nonce          uint64
	Balance        *big.Int
	BalanceDelta   *big.Int
	StorageUpdates map[string]*vmcommon.StorageUpdate
	Code           []byte
	Data           []byte
	GasLimit       uint64
}

func (output *Output) popStateWithoutMerge() (*Output, error) {
	stateStackLen := len(output.stateStack)
	if stateStackLen < 1 {
		return nil, arwen.StateStackUnderflow
	}

	state := output.stateStack[stateStackLen-1]
	output.stateStack = output.stateStack[:stateStackLen-1]

	return state, nil
}

func (output *Output) mergeState(state *Output) error {
	return nil
}

func mergeAccountMaps(left AccountsMap, right AccountsMap) AccountsMap {
	return nil
}

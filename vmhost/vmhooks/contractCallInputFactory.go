package vmhooks

import (
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// ContractCallInputFactory holds the args for contract call input factory
type ContractCallInputFactory struct{}

// PrepareContractCallInput will return the contract call input
func (f *ContractCallInputFactory) PrepareContractCallInput(
	host vmhost.VMHost,
	sender []byte,
	value *big.Int,
	gasLimit int64,
	destination []byte,
	function []byte,
	data [][]byte,
	_ uint64,
	syncExecutionRequired bool,
) (vmcommon.ContractCallInputHandler, error) {
	return prepareIndirectContractCallInput(host, sender, value, gasLimit, destination, function, data, 0, syncExecutionRequired)
}

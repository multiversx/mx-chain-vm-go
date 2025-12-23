package vmhooks

import (
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// ContractCreateInputFactory holds the args for contract create input factory
type ContractCreateInputFactory struct{}

// PrepareContractCreateInput will return the evm contract create input
func (f *ContractCreateInputFactory) PrepareContractCreateInput(
	host vmhost.VMHost,
	sender []byte,
	data [][]byte,
	value *big.Int,
	gasLimit int64,
	code []byte,
	codeMetadata []byte,
) vmcommon.ContractCreateInputHandler {
	return prepareCreateContractInput(host, sender, data, value, gasLimit, code, codeMetadata)
}

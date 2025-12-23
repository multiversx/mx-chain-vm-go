package vmhooks

import (
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// ContractCallInputHandlerFactory defines the behavior for contract call input handler factory
type ContractCallInputHandlerFactory interface {
	PrepareContractCallInput(
		host vmhost.VMHost,
		sender []byte,
		value *big.Int,
		gasLimit int64,
		destination []byte,
		function []byte,
		data [][]byte,
		_ uint64,
		syncExecutionRequired bool,
	) (vmcommon.ContractCallInputHandler, error)
}

// ContractCreateInputHandlerFactory defines the behavior for contract create input handler factory
type ContractCreateInputHandlerFactory interface {
	PrepareContractCreateInput(
		host vmhost.VMHost,
		sender []byte,
		data [][]byte,
		value *big.Int,
		gasLimit int64,
		code []byte,
		codeMetadata []byte,
	) vmcommon.ContractCreateInputHandler
}

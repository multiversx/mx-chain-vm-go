package arwenpart

import vmcommon "github.com/ElrondNetwork/elrond-vm-common"

// VMHost is
type VMHost interface {
	RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error)
	RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error)
}
